#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const [, , versionArg = 'develop', wailsPathArg = 'src/wails.json'] = process.argv
const wailsPath = path.resolve(process.cwd(), wailsPathArg)

function windowsVersion(value) {
  const raw = String(value || '').trim().replace(/^[vV]/, '').split(/[+-]/, 1)[0]
  const parts = raw.split('.').map((part) => Number.parseInt(part, 10))
  const nums = []
  for (let i = 0; i < 3; i += 1) {
    nums.push(Number.isInteger(parts[i]) && parts[i] >= 0 ? parts[i] : 0)
  }
  return nums.join('.')
}

const data = JSON.parse(fs.readFileSync(wailsPath, 'utf8'))
data.info = {
  ...(data.info || {}),
  companyName: 'hatolife',
  productName: 'ClipForVRChat',
  productVersion: windowsVersion(versionArg),
  copyright: 'Copyright (c) 2026 hatolife',
  comments: 'ClipForVRChat official build',
}

fs.writeFileSync(wailsPath, `${JSON.stringify(data, null, 2)}\n`)
