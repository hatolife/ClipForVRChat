const patchedInstances = new WeakSet()

function serializeComparable(value) {
  return JSON.stringify(value || {})
}

function valueAtPath(source, path) {
  return String(path || '').split('.').reduce((current, key) => {
    if (current === undefined || current === null) return undefined
    return current[key]
  }, source)
}

function parseBaseline(vm) {
  try {
    return JSON.parse(vm.settingsBaseline || '{}')
  } catch {
    return null
  }
}

function effectiveWebhookURL(vm, specific, fallback) {
  if (typeof vm.effectiveWebhookURL === 'function') {
    return vm.effectiveWebhookURL(specific, fallback)
  }
  const own = String(specific || '').trim()
  if (own) return own
  return String(fallback || '').trim()
}

function maskWebhook(vm, value) {
  if (!value) return '未設定'
  if (typeof vm.maskWebhook === 'function') return vm.maskWebhook(value)
  const text = String(value)
  const match = text.match(/\/webhooks\/([^/?#]+)\/([^/?#]+)/)
  if (!match) return '設定済み'
  const token = match[2]
  return `${match[1]}/${token.slice(0, 4)}...${token.slice(-4)}`
}

function watchDirectoryWarning(vm, directory) {
  if (typeof vm.autoProcessingWatchDirectoryWarning === 'function') {
    return vm.autoProcessingWatchDirectoryWarning(directory)
  }
  return ''
}

function primaryWebhook(config) {
  return String(config?.discord?.webhookUrl || '').trim()
}

function autoPhotoSnapshot(vm, config) {
  const autoPhoto = config?.autoPhoto || {}
  const output = config?.output || {}
  const target = effectiveWebhookURL(vm, autoPhoto.webhookUrl, primaryWebhook(config))
  return {
    enabled: Boolean(output.uploadDiscord && autoPhoto.enabled),
    directory: String(autoPhoto.photoDirectory || ''),
    target: String(target || '')
  }
}

function screenshotSnapshot(vm, config) {
  const screenshot = config?.screenshotAutoPost || {}
  const output = config?.output || {}
  const target = effectiveWebhookURL(vm, screenshot.webhookUrl, primaryWebhook(config))
  return {
    enabled: Boolean(output.uploadDiscord && screenshot.enabled),
    directory: String(screenshot.screenshotDirectory || ''),
    target: String(target || '')
  }
}

function autoCaptureSnapshot(vm, config) {
  const autoCapture = config?.autoCapture || {}
  const schedule = autoCapture.schedule || {}
  const discord = autoCapture.discord || {}
  const target = effectiveWebhookURL(vm, discord.webhookUrl, primaryWebhook(config))
  return {
    scheduleEnabled: Boolean(schedule.enabled),
    discordEnabled: Boolean(discord.enabled),
    directory: String(autoCapture.output?.directory || ''),
    target: String(target || '')
  }
}

function snapshotChanged(before, after) {
  return serializeComparable(before) !== serializeComparable(after)
}

function buildCurrentAutoPostConfirmationItems(vm, config) {
  const items = []
  const autoPhoto = config?.autoPhoto || {}
  const screenshot = config?.screenshotAutoPost || {}
  const autoCapture = config?.autoCapture || {}
  const autoCaptureSchedule = autoCapture.schedule || {}
  const autoCaptureDiscord = autoCapture.discord || {}
  const output = config?.output || {}
  const primary = primaryWebhook(config)

  if (output.uploadDiscord && autoPhoto.enabled) {
    const target = effectiveWebhookURL(vm, autoPhoto.webhookUrl, primary)
    items.push({
      label: 'VRChat写真自動処理',
      detail: `監視フォルダ: ${autoPhoto.photoDirectory || '(未設定)'}`,
      discord: target ? `Discord投稿: ON / 送信先: ${maskWebhook(vm, target)}` : 'Discord投稿: ON / 送信先: 未設定',
      warning: watchDirectoryWarning(vm, autoPhoto.photoDirectory)
    })
  }
  if (output.uploadDiscord && screenshot.enabled) {
    const target = effectiveWebhookURL(vm, screenshot.webhookUrl, primary)
    items.push({
      label: 'スクリーンショット自動処理',
      detail: `監視フォルダ: ${screenshot.screenshotDirectory || '(未設定)'}`,
      discord: target ? `Discord投稿: ON / 送信先: ${maskWebhook(vm, target)}` : 'Discord投稿: ON / 送信先: 未設定',
      warning: watchDirectoryWarning(vm, screenshot.screenshotDirectory)
    })
  }
  if (autoCaptureSchedule.enabled && autoCaptureDiscord.enabled) {
    const target = effectiveWebhookURL(vm, autoCaptureDiscord.webhookUrl, primary)
    if (!target) {
      items.push({
        label: '自動撮影',
        detail: `保存先: ${autoCapture.output?.directory || '(未設定)'}`,
        discord: 'Discord投稿: ON / 送信先: 未設定'
      })
    }
  }
  return items
}

function buildAutoPostConfirmationItems(vm) {
  const after = vm.state?.config || {}
  if (!vm.hasUnsavedSettings) return []

  const before = parseBaseline(vm)
  if (!before) return buildCurrentAutoPostConfirmationItems(vm, after)

  const items = []
  const afterPhoto = autoPhotoSnapshot(vm, after)
  if (afterPhoto.enabled && snapshotChanged(autoPhotoSnapshot(vm, before), afterPhoto)) {
    const autoPhoto = after.autoPhoto || {}
    items.push({
      label: 'VRChat写真自動処理',
      detail: `監視フォルダ: ${autoPhoto.photoDirectory || '(未設定)'}`,
      discord: afterPhoto.target ? `Discord投稿: ON / 送信先: ${maskWebhook(vm, afterPhoto.target)}` : 'Discord投稿: ON / 送信先: 未設定',
      warning: watchDirectoryWarning(vm, autoPhoto.photoDirectory)
    })
  }

  const afterScreenshot = screenshotSnapshot(vm, after)
  if (afterScreenshot.enabled && snapshotChanged(screenshotSnapshot(vm, before), afterScreenshot)) {
    const screenshot = after.screenshotAutoPost || {}
    items.push({
      label: 'スクリーンショット自動処理',
      detail: `監視フォルダ: ${screenshot.screenshotDirectory || '(未設定)'}`,
      discord: afterScreenshot.target ? `Discord投稿: ON / 送信先: ${maskWebhook(vm, afterScreenshot.target)}` : 'Discord投稿: ON / 送信先: 未設定',
      warning: watchDirectoryWarning(vm, screenshot.screenshotDirectory)
    })
  }

  const afterAutoCapture = autoCaptureSnapshot(vm, after)
  if (
    afterAutoCapture.scheduleEnabled &&
    afterAutoCapture.discordEnabled &&
    !afterAutoCapture.target &&
    snapshotChanged(autoCaptureSnapshot(vm, before), afterAutoCapture)
  ) {
    items.push({
      label: '自動撮影',
      detail: `保存先: ${after.autoCapture?.output?.directory || '(未設定)'}`,
      discord: 'Discord投稿: ON / 送信先: 未設定'
    })
  }

  return items
}

function installComputedGetter(ctx, key, getter) {
  const descriptor = Object.getOwnPropertyDescriptor(ctx, key)
  if (descriptor && descriptor.configurable === false) return false
  Object.defineProperty(ctx, key, {
    configurable: true,
    enumerable: true,
    get: getter
  })
  return true
}

function patchMethod(ctx, key, wrapper) {
  const original = ctx[key]
  if (typeof original !== 'function') return
  Object.defineProperty(ctx, key, {
    configurable: true,
    enumerable: true,
    writable: true,
    value: wrapper(original)
  })
}

function installPatch() {
  const root = document.getElementById('app')
  const instance = root?.__vue_app__?._instance
  const proxy = instance?.proxy
  const ctx = instance?.ctx
  if (!proxy || !ctx || patchedInstances.has(instance)) return Boolean(proxy && ctx)

  installComputedGetter(ctx, 'autoPostConfirmationItems', () => buildAutoPostConfirmationItems(proxy))
  installComputedGetter(ctx, 'shouldConfirmAutoPostSettings', () => buildAutoPostConfirmationItems(proxy).length > 0)

  patchMethod(ctx, 'saveSettings', (original) => async function patchedSaveSettings(skipAutoPostConfirmation = false, skipSensitiveSettingsConfirmation = false) {
    const skipAutoPost = skipAutoPostConfirmation || buildAutoPostConfirmationItems(proxy).length === 0
    return original.call(this, skipAutoPost, skipSensitiveSettingsConfirmation)
  })

  patchMethod(ctx, 'confirmSaveAndLeaveSettings', (original) => async function patchedConfirmSaveAndLeaveSettings(skipAutoPostConfirmation = false, overrideAction = '', skipSensitiveSettingsConfirmation = false) {
    const skipAutoPost = skipAutoPostConfirmation || buildAutoPostConfirmationItems(proxy).length === 0
    return original.call(this, skipAutoPost, overrideAction, skipSensitiveSettingsConfirmation)
  })

  patchedInstances.add(instance)
  return true
}

let attempts = 0
function waitForVueApp() {
  if (installPatch()) return
  attempts += 1
  if (attempts < 200) window.setTimeout(waitForVueApp, 25)
}

waitForVueApp()
