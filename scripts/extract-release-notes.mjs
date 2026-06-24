#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const [, , versionArg, releaseNotesArg = 'RELEASE_NOTES.md', outPathArg = 'dist/release-body.md', ...flags] = process.argv

if (!versionArg) {
  console.error('usage: extract-release-notes.mjs <version> [release-notes.md] [out.md]')
  process.exit(2)
}

const version = String(versionArg).trim()
const releaseNotesPath = path.resolve(process.cwd(), releaseNotesArg)
const outPath = path.resolve(process.cwd(), outPathArg)
const content = fs.readFileSync(releaseNotesPath, 'utf8')
const findSection = (sectionVersion) => {
  const escapedVersion = sectionVersion.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const pattern = new RegExp(`^## ${escapedVersion}\\s*$([\\s\\S]*?)(?=^##\\s+|(?![\\s\\S]))`, 'm')
  return content.match(pattern)
}

let match = findSection(version)
let sourceVersion = version

if (!match) {
  const rcMatch = version.match(/^(v\d+\.\d+\.\d+)-rc\d+$/)
  if (flags.includes('--fallback-rc-notes') && rcMatch) {
    match = findSection(rcMatch[1])
    if (match) {
      sourceVersion = rcMatch[1]
    }
  }
}

if (!match) {
  if (flags.includes('--fallback-draft-notes')) {
    fs.mkdirSync(path.dirname(outPath), { recursive: true })
    fs.writeFileSync(outPath, `Draft release for ${version}.\n\nThis tag does not have a dedicated RELEASE_NOTES.md entry.\n`)
    process.exit(0)
  }
  console.error(`release notes for ${version} were not found in ${releaseNotesArg}`)
  process.exit(1)
}

let body = match[1].trim()
if (!body) {
  console.error(`release notes for ${version} are empty`)
  process.exit(1)
}

if (sourceVersion !== version) {
  const escapedSourceVersion = sourceVersion.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  body = body.replace(new RegExp(escapedSourceVersion, 'g'), version)
}

fs.mkdirSync(path.dirname(outPath), { recursive: true })
fs.writeFileSync(outPath, `${body}\n`)
