#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const [, , versionArg, releaseNotesArg = 'RELEASE_NOTES.md', outPathArg = 'dist/release-body.md'] = process.argv

if (!versionArg) {
  console.error('usage: extract-release-notes.mjs <version> [release-notes.md] [out.md]')
  process.exit(2)
}

const version = String(versionArg).trim()
const releaseNotesPath = path.resolve(process.cwd(), releaseNotesArg)
const outPath = path.resolve(process.cwd(), outPathArg)
const content = fs.readFileSync(releaseNotesPath, 'utf8')
const escapedVersion = version.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
const pattern = new RegExp(`^## ${escapedVersion}\\s*$([\\s\\S]*?)(?=^##\\s+|(?![\\s\\S]))`, 'm')
const match = content.match(pattern)

if (!match) {
  console.error(`release notes for ${version} were not found in ${releaseNotesArg}`)
  process.exit(1)
}

const body = match[1].trim()
if (!body) {
  console.error(`release notes for ${version} are empty`)
  process.exit(1)
}

fs.mkdirSync(path.dirname(outPath), { recursive: true })
fs.writeFileSync(outPath, `${body}\n`)
