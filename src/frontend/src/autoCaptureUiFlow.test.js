import assert from 'node:assert/strict'
import test from 'node:test'

import { modeAfterAutoCaptureResult } from './autoCaptureUiFlow.js'

test('設定画面を開いている間は自動撮影結果で画面を切り替えない', () => {
  assert.equal(modeAfterAutoCaptureResult('settings'), 'settings')
})

test('設定画面以外では自動撮影結果を表示する', () => {
  assert.equal(modeAfterAutoCaptureResult('results'), 'results')
  assert.equal(modeAfterAutoCaptureResult('error'), 'results')
})
