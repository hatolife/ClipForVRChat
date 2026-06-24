#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const [, , versionArg = 'develop', revisionArg = 'unknown', outPathArg = 'src/build/windows/info.json'] = process.argv
const outPath = path.resolve(process.cwd(), outPathArg)

function numericVersion(value) {
  const raw = String(value || '').trim().replace(/^[vV]/, '').split(/[+-]/, 1)[0]
  const parts = raw.split('.').map((part) => Number.parseInt(part, 10))
  const nums = []
  for (let i = 0; i < 3; i += 1) {
    nums.push(Number.isInteger(parts[i]) && parts[i] >= 0 ? parts[i] : 0)
  }
  nums.push(0)
  return nums.join('.')
}

function displayVersion(version, revision) {
  const cleanVersion = String(version || 'develop').trim() || 'develop'
  const cleanRevision = String(revision || 'unknown').trim() || 'unknown'
  return `${cleanVersion}.${cleanRevision}`
}

const windowsVersion = numericVersion(versionArg)
const fullName = `ClipForVRChat ${displayVersion(versionArg, revisionArg)}`

const data = {
  fixed: {
    file_version: windowsVersion,
    product_version: windowsVersion,
  },
  info: {
    '0000': {
      ProductVersion: fullName,
      FileVersion: fullName,
      CompanyName: 'hatolife',
      FileDescription: fullName,
      LegalCopyright: 'Copyright (c) 2026 hatolife',
      ProductName: 'ClipForVRChat',
      InternalName: 'ClipForVRChat.exe',
      OriginalFilename: 'ClipForVRChat.exe',
      Comments: 'ClipForVRChat official build',
    },
  },
}

fs.writeFileSync(outPath, `${JSON.stringify(data, null, 2)}\n`)
process.stdout.write(windowsVersion)
