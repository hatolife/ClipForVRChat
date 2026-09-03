param(
    [Parameter(Mandatory = $true)]
    [string]$ExecutablePath,

    [int]$Port = 9222,

    [int]$TimeoutSeconds = 45
)

$ErrorActionPreference = "Stop"

if ($Port -lt 1 -or $Port -gt 65535) {
    throw "Port must be between 1 and 65535: $Port"
}
if ($TimeoutSeconds -lt 1) {
    throw "TimeoutSeconds must be positive: $TimeoutSeconds"
}

$resolvedExecutable = (Resolve-Path -LiteralPath $ExecutablePath).Path
$workingDirectory = Split-Path -Parent $resolvedExecutable
$previousArguments = $env:WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS
$previousUserDataFolder = $env:WEBVIEW2_USER_DATA_FOLDER
$debugArgument = "--remote-debugging-port=$Port"
if ([string]::IsNullOrWhiteSpace($previousArguments)) {
    $env:WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS = $debugArgument
} else {
    $env:WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS = "$previousArguments $debugArgument"
}
$env:WEBVIEW2_USER_DATA_FOLDER = Join-Path ([System.IO.Path]::GetTempPath()) "ClipForVRChat-WebView2-$([Guid]::NewGuid().ToString('N'))"

$startedAtUtc = (Get-Date).ToUniversalTime()
$logDirectory = Join-Path $workingDirectory "logs"

function Get-LifecycleLogText {
    param(
        [string]$Directory,
        [datetime]$SinceUtc
    )

    if (!(Test-Path -LiteralPath $Directory)) {
        return ""
    }

    $parts = Get-ChildItem -LiteralPath $Directory -Filter "*.log" -File -ErrorAction SilentlyContinue |
        Where-Object { $_.LastWriteTimeUtc -ge $SinceUtc.AddSeconds(-2) } |
        ForEach-Object { Get-Content -LiteralPath $_.FullName -Raw -ErrorAction SilentlyContinue }
    return ($parts -join "`n")
}

$process = $null
try {
    $process = Start-Process -FilePath $resolvedExecutable -WorkingDirectory $workingDirectory -PassThru
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    $endpoint = "http://127.0.0.1:$Port/json/list"
    $devToolsTarget = $null
    $logText = ""

    while ((Get-Date) -lt $deadline) {
        if ($process.HasExited) {
            throw "ClipForVRChat exited before WebView2 became ready. ExitCode=$($process.ExitCode)"
        }

        if (!$devToolsTarget) {
            try {
                $targets = Invoke-RestMethod -Uri $endpoint -TimeoutSec 2
                $devToolsTarget = @($targets) |
                    Where-Object { $_.type -eq "page" -and ![string]::IsNullOrWhiteSpace($_.webSocketDebuggerUrl) } |
                    Select-Object -First 1
            } catch {
                if ($process.HasExited) {
                    throw
                }
            }
        }

        $logText = [string](Get-LifecycleLogText -Directory $logDirectory -SinceUtc $startedAtUtc)
        $domReady = $logText.Contains("ui lifecycle dom_ready")
        $frontendScriptLoaded = $logText.Contains('ui action="frontend_script_loaded"')
        $initialStateReady = $logText.Contains("api GetInitialState complete")

        if ($domReady -and $frontendScriptLoaded -and $initialStateReady) {
            Write-Output "Wails frontend startup detected: dom_ready=true frontend_script_loaded=true initial_state=true"
            if ($devToolsTarget) {
                Write-Output "Optional WebView2 DevTools target detected: title=$($devToolsTarget.title) url=$($devToolsTarget.url)"
            } else {
                Write-Output "Optional WebView2 DevTools endpoint was unavailable; lifecycle diagnostics verified the packaged frontend instead."
            }
            return
        }

        if ($logText -match 'ui action="frontend_(error|mount_error|unhandledrejection)"') {
            $logTail = (($logText -split "`r?`n") | Select-Object -Last 30) -join "`n"
            throw "Frontend startup error was recorded in the lifecycle log:`n$logTail"
        }

        Start-Sleep -Milliseconds 500
    }

    $logTail = (($logText -split "`r?`n") | Select-Object -Last 30) -join "`n"
    throw "Wails frontend lifecycle was not ready within $TimeoutSeconds seconds. Expected dom_ready, frontend_script_loaded, and GetInitialState complete. DevToolsEndpoint=$endpoint`n$logTail"
} finally {
    $env:WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS = $previousArguments
    $env:WEBVIEW2_USER_DATA_FOLDER = $previousUserDataFolder
    if ($process -and !$process.HasExited) {
        Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
        $process.WaitForExit(5000) | Out-Null
    }
}
