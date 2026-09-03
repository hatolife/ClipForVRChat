import assert from 'node:assert/strict'
import test from 'node:test'

import { resumeAfterAutoPostConfirmation } from './settingsConfirmationFlow.js'

test('自動投稿確認後の保存では両方の確認を完了済みとして扱う', () => {
  assert.deepEqual(resumeAfterAutoPostConfirmation('save'), {
    target: 'save',
    leaveAction: '',
    skipAutoPostConfirmation: true,
    skipSensitiveSettingsConfirmation: true
  })
})

test('自動投稿確認後の画面移動でも両方の確認を完了済みとして扱う', () => {
  assert.deepEqual(resumeAfterAutoPostConfirmation('leave:history'), {
    target: 'leave',
    leaveAction: 'history',
    skipAutoPostConfirmation: true,
    skipSensitiveSettingsConfirmation: true
  })
})

test('移動先が空の場合はホームへ戻る', () => {
  assert.deepEqual(resumeAfterAutoPostConfirmation('leave:'), {
    target: 'leave',
    leaveAction: 'home',
    skipAutoPostConfirmation: true,
    skipSensitiveSettingsConfirmation: true
  })
})
