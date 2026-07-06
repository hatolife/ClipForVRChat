#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'
import { createRequire } from 'node:module'

const root = process.cwd()
const frontendPath = path.join(root, 'src', 'frontend', 'src', 'main.js')
const frontendPackagePath = path.join(root, 'src', 'frontend', 'package.json')
const source = fs.readFileSync(frontendPath, 'utf8')
const lines = source.split(/\r?\n/)

let compile
try {
  const frontendRequire = createRequire(frontendPackagePath)
  ;({ compile } = frontendRequire('@vue/compiler-dom'))
} catch (error) {
  console.error('Failed to load @vue/compiler-dom from src/frontend dependencies.')
  console.error('Run npm ci in src/frontend before this check.')
  console.error(String(error))
  process.exit(1)
}

const startIndex = lines.findIndex((line) => line.includes('template: `'))
if (startIndex === -1) {
  console.error(`${frontendPath}: Vue runtime template literal was not found`)
  process.exit(1)
}

let endIndex = -1
for (let i = startIndex + 1; i < lines.length; i += 1) {
  if (/^\s*`\s*,?\s*$/.test(lines[i])) {
    endIndex = i
    break
  }
}

if (endIndex === -1) {
  console.error(`${frontendPath}:${startIndex + 1}: Vue runtime template literal was not closed`)
  process.exit(1)
}

const templateStartLine = startIndex + 2
const template = lines.slice(startIndex + 1, endIndex).join('\n')
const errors = []

try {
  compile(template, {
    onError(error) {
      errors.push(error)
    }
  })
} catch (error) {
  errors.push(error)
}

if (errors.length > 0) {
  for (const error of errors) {
    const loc = error.loc?.start
    const line = loc ? templateStartLine + loc.line - 1 : startIndex + 1
    const column = loc?.column ?? 1
    console.error(`${frontendPath}:${line}:${column}: Vue runtime template error ${error.code ?? 'unknown'}: ${error.message}`)
    const from = Math.max(1, line - 2)
    const to = Math.min(lines.length, line + 2)
    for (let i = from; i <= to; i += 1) {
      console.error(`${String(i).padStart(5)} ${lines[i - 1]}`)
    }
  }
  process.exit(1)
}

console.log('Vue runtime template check OK')
