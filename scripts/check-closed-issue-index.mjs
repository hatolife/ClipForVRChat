#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const target = path.join(root, 'issues', 'closed', 'README.md')
const text = fs.readFileSync(target, 'utf8')
const rows = text.split(/\r?\n/).filter((line) => /^\|\s*\d+\s*\|/.test(line))

let previous = 0
for (const row of rows) {
  const match = row.match(/^\|\s*(\d+)\s*\|/)
  if (!match) continue
  const current = Number(match[1])
  if (current < previous) {
    console.error(`Closed issue index is not sorted: ${current} appears after ${previous}`)
    process.exit(1)
  }
  previous = current
}

console.log(`Closed issue index OK: ${rows.length} rows sorted`)
