#!/usr/bin/env node
import fs from 'node:fs'
import path from 'node:path'

const root = process.cwd()
const frontendPath = path.join(root, 'src', 'frontend', 'src', 'main.js')
const source = fs.readFileSync(frontendPath, 'utf8')

const problems = []

function requireSource(pattern, message) {
  if (!pattern.test(source)) problems.push(message)
}

requireSource(
  /if \(output\.uploadDiscord && autoPhoto\.enabled\) \{\s*const target = this\.effectiveWebhookURL\(autoPhoto\.webhookUrl, primaryWebhook\)\s*items\.push\(/s,
  'autoPhoto confirmation must be created whenever Discord upload and auto photo are enabled'
)

requireSource(
  /if \(output\.uploadDiscord && screenshot\.enabled\) \{\s*const target = this\.effectiveWebhookURL\(screenshot\.webhookUrl, primaryWebhook\)\s*items\.push\(/s,
  'screenshot confirmation must be created whenever Discord upload and screenshot auto processing are enabled'
)

requireSource(
  /autoProcessingWatchDirectoryWarning\(autoPhoto\.photoDirectory\)/,
  'autoPhoto confirmation must include watch directory risk warning'
)

requireSource(
  /autoProcessingWatchDirectoryWarning\(screenshot\.screenshotDirectory\)/,
  'screenshot confirmation must include watch directory risk warning'
)

requireSource(
  /v-model="state\.config\.autoPhoto\.photoDirectory"/,
  'settings UI must expose autoPhoto.photoDirectory'
)

requireSource(
  /v-model="state\.config\.screenshotAutoPost\.screenshotDirectory"/,
  'settings UI must expose screenshotAutoPost.screenshotDirectory'
)

requireSource(
  /v-if="item\.warning" class="warning"/,
  'confirmation dialog must render warning messages'
)

requireSource(
  /parts\[1\] === 'users' && parts\.length <= 3/,
  'watch directory warning must treat Windows user profile roots as broad folders'
)

if (problems.length > 0) {
  console.error(problems.join('\n'))
  process.exit(1)
}

console.log('Auto processing confirmation check OK')
