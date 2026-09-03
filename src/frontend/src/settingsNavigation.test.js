import assert from 'node:assert/strict'
import test from 'node:test'

import { autoCaptureSettingsCategory, settingsNavigation } from './settingsNavigation.js'

test('設定ナビゲーションは主要カテゴリと詳細カテゴリを縦に整理する', () => {
  assert.deepEqual(settingsNavigation.map(({ id, label, group }) => ({ id, label, group })), [
    { id: 'feature', label: '機能 ON/OFF', group: 'primary' },
    { id: 'process', label: '縮小処理', group: 'primary' },
    { id: 'webhook', label: '投稿処理', group: 'primary' },
    { id: 'autoCapture', label: '撮影処理', group: 'primary' },
    { id: 'osc', label: 'OSC', group: 'advanced' },
    { id: 'other', label: 'その他', group: 'advanced' }
  ])
})

test('自動定期撮影の未保存変更は機能カテゴリへ割り当てる', () => {
  assert.equal(autoCaptureSettingsCategory('autoCapture.schedule.enabled'), 'feature')
  assert.equal(autoCaptureSettingsCategory('autoCapture.stream.startDelayMs'), 'autoCapture')
  assert.equal(autoCaptureSettingsCategory('autoCapture.osc.vrcHost'), 'osc')
  assert.equal(autoCaptureSettingsCategory('autoCapture.playerLocal.basisSource'), 'osc')
  assert.equal(autoCaptureSettingsCategory('discord.webhookUrl'), '')
})
