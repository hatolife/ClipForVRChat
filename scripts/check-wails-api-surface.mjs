#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const frontendPath = path.join(root, 'src', 'frontend', 'src', 'main.js')
const appPath = path.join(root, 'src', 'app.go')

const frontend = fs.readFileSync(frontendPath, 'utf8')
const app = fs.readFileSync(appPath, 'utf8')

const calls = new Set()
for (const match of frontend.matchAll(/\bapi\??\.([A-Z][A-Za-z0-9_]*)/g)) {
  calls.add(match[1])
}

const methods = new Set()
for (const match of app.matchAll(/func\s*\(\s*a\s+\*App\s*\)\s+([A-Z][A-Za-z0-9_]*)\s*\(/g)) {
  methods.add(match[1])
}

const missing = [...calls].filter((name) => !methods.has(name)).sort()
if (missing.length > 0) {
  console.error(`frontend calls missing App methods: ${missing.join(', ')}`)
  process.exit(1)
}

console.log(`Wails API surface OK: ${calls.size} frontend calls matched App methods.`)
