import assert from 'node:assert/strict'
import test from 'node:test'

import { buildAutoPostConfirmationItems } from './autoPostConfirmationDiffPatch.js'

function makeConfig() {
  return {
    output: {
      uploadDiscord: true
    },
    discord: {
      webhookUrl: 'https://discord.com/api/webhooks/1000/primary-token'
    },
    autoPhoto: {
      enabled: true,
      photoDirectory: 'C:\\Users\\tester\\Pictures\\VRChat',
      webhookUrl: ''
    },
    screenshotAutoPost: {
      enabled: true,
      screenshotDirectory: 'C:\\Users\\tester\\Pictures\\Screenshots',
      webhookUrl: ''
    },
    autoCapture: {
      schedule: {
        enabled: false
      },
      discord: {
        enabled: false,
        webhookUrl: ''
      },
      output: {
        directory: 'output'
      }
    },
    update: {
      checkEnabled: true
    }
  }
}

function makeVM(before, after, hasUnsavedSettings = true) {
  return {
    hasUnsavedSettings,
    settingsBaseline: JSON.stringify(before),
    state: {
      config: after
    }
  }
}

test('未保存変更がなければ確認項目を返さない', () => {
  const config = makeConfig()
  const vm = makeVM(config, structuredClone(config), false)

  assert.deepEqual(buildAutoPostConfirmationItems(vm), [])
})

test('自動投稿に無関係な変更では確認項目を返さない', () => {
  const before = makeConfig()
  const after = structuredClone(before)
  after.update.checkEnabled = false

  assert.deepEqual(buildAutoPostConfirmationItems(makeVM(before, after)), [])
})

test('VRChat写真の監視フォルダ変更を確認対象にする', () => {
  const before = makeConfig()
  const after = structuredClone(before)
  after.autoPhoto.photoDirectory = 'D:\\VRChat'

  const items = buildAutoPostConfirmationItems(makeVM(before, after))

  assert.equal(items.length, 1)
  assert.equal(items[0].label, 'VRChat写真自動処理')
  assert.equal(items[0].detail, '監視フォルダ: D:\\VRChat')
})

test('通常Webhook変更をfallback利用中の自動投稿へ反映する', () => {
  const before = makeConfig()
  const after = structuredClone(before)
  after.discord.webhookUrl = 'https://discord.com/api/webhooks/2000/changed-token'

  const items = buildAutoPostConfirmationItems(makeVM(before, after))

  assert.deepEqual(items.map((item) => item.label), [
    'VRChat写真自動処理',
    'スクリーンショット自動処理'
  ])
})

test('無効化後の自動投稿は確認対象にしない', () => {
  const before = makeConfig()
  const after = structuredClone(before)
  after.autoPhoto.enabled = false

  assert.deepEqual(buildAutoPostConfirmationItems(makeVM(before, after)), [])
})

test('送信先未設定の自動撮影変更を確認対象にする', () => {
  const before = makeConfig()
  before.output.uploadDiscord = false
  before.discord.webhookUrl = ''
  const after = structuredClone(before)
  after.autoCapture.schedule.enabled = true
  after.autoCapture.discord.enabled = true

  const items = buildAutoPostConfirmationItems(makeVM(before, after))

  assert.equal(items.length, 1)
  assert.equal(items[0].label, '自動撮影')
  assert.equal(items[0].discord, 'Discord投稿: ON / 送信先: 未設定')
})
