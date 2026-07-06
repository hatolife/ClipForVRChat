#!/usr/bin/env node
import { execFileSync } from 'node:child_process'

const [, , rangeArg = 'HEAD~1..HEAD'] = process.argv
const subjectPattern = /^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9-]+\))?!?: .+/

let output = ''
try {
  output = execFileSync('git', ['log', '--format=%H%x09%s', rangeArg], { encoding: 'utf8' })
} catch (error) {
  console.error(`failed to read git log for range ${rangeArg}`)
  if (error.stderr) console.error(String(error.stderr).trim())
  process.exit(2)
}

const problems = []
for (const line of output.trim().split(/\r?\n/).filter(Boolean)) {
  const [sha, subject] = line.split('\t')
  if (!subjectPattern.test(subject)) {
    problems.push(`${sha.slice(0, 12)} ${subject}`)
  }
}

if (problems.length > 0) {
  console.error('Commit subjects must follow Conventional Commits, for example "fix(ui): handle empty state".')
  console.error(problems.join('\n'))
  process.exit(1)
}

console.log(`Commit subject check OK: ${rangeArg}`)
