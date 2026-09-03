import assert from 'node:assert/strict'
import test from 'node:test'

import { settingsNavigation } from './settingsNavigation.js'

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
