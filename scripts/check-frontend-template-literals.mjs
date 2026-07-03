#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const frontendPath = path.join(root, 'src', 'frontend', 'src', 'main.js')
const source = fs.readFileSync(frontendPath, 'utf8')
const lines = source.split(/\r?\n/)

let inTemplate = false
let startLine = 0
const problems = []

for (let i = 0; i < lines.length; i += 1) {
  const line = lines[i]
  if (!inTemplate && line.includes('template: `')) {
    inTemplate = true
    startLine = i + 1
    continue
  }
  if (!inTemplate) continue
  if (/^\s*`\s*,?\s*$/.test(line)) {
    inTemplate = false
    continue
  }
  if (line.includes('`')) {
    problems.push(`${frontendPath}:${i + 1}: raw backtick inside Vue template literal started at line ${startLine}`)
  }
}

if (inTemplate) {
  problems.push(`${frontendPath}:${startLine}: Vue template literal was not closed`)
}

if (problems.length > 0) {
  console.error(problems.join('\n'))
  console.error('Use <code>...</code> or escaped text instead of raw backticks inside template: `...`.')
  process.exit(1)
}

console.log('Frontend template literal check OK')
