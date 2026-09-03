export const settingsNavigation = Object.freeze([
  { id: 'feature', label: '機能 ON/OFF', group: 'primary' },
  { id: 'process', label: '縮小処理', group: 'primary' },
  { id: 'webhook', label: '投稿処理', group: 'primary' },
  { id: 'autoCapture', label: '撮影処理', group: 'primary' },
  { id: 'osc', label: 'OSC', group: 'advanced' },
  { id: 'other', label: 'その他', group: 'advanced' }
])

export function autoCaptureSettingsCategory(path) {
  if (path === 'autoCapture.schedule' || path.startsWith('autoCapture.schedule.')) return 'feature'
  if (path === 'autoCapture.osc' || path.startsWith('autoCapture.osc.') || path === 'autoCapture.playerLocal' || path.startsWith('autoCapture.playerLocal.')) return 'osc'
  if (path === 'autoCapture' || path.startsWith('autoCapture.')) return 'autoCapture'
  return ''
}
