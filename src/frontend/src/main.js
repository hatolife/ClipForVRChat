import { createApp } from 'vue/dist/vue.esm-bundler.js'
import './style.css'

const api = window.go?.main?.App

createApp({
  data() {
    return {
      info: { name: 'ClipForVRChat', version: 'dev', github: 'https://github.com/hatolife/ClipForVRChat' },
      state: { mode: 'results', message: '', configPath: '', config: null, results: [] },
      licenses: [],
      webhookGuideUrl: 'https://support.discord.com/hc/ja/articles/228383668-%E3%82%A6%E3%82%A7%E3%83%96%E3%83%95%E3%83%83%E3%82%AF%E3%81%AE%E3%81%94%E7%B4%B9%E4%BB%8B',
      issuesUrl: 'https://github.com/hatolife/ClipForVRChat/issues',
      authorTwitterUrl: 'https://x.com/hato_poppo_life',
      feedbackTweetUrl: 'https://x.com/hato_poppo_life/status/2068611307830710667',
      releasesUrl: 'https://github.com/hatolife/ClipForVRChat/releases',
      latestReleaseUrl: 'https://github.com/hatolife/ClipForVRChat/releases/latest',
      boothUrl: 'https://hatolife.booth.pm/items/8531663',
      updateInfo: { available: false, currentVersion: '', currentReleaseTime: '', latestVersion: '', latestReleasePublished: '', url: '' },
      updateBannerDismissed: false,
      view: 'main',
      processing: false,
      dragging: false,
      selectedHistoryIds: [],
      lastSelectedHistoryIndex: -1,
      toast: '',
      saving: false,
      saved: false,
      error: '',
      settingsBaseline: '',
      settingsTab: 'feature',
      pendingSettingsLeave: null,
      pendingAutoPostConfirmation: null,
      pendingDropPaths: [],
      historyDragSelecting: false,
      historySelectionAdditive: false,
      historyDragStart: { x: 0, y: 0 },
      historyDragCurrent: { x: 0, y: 0 },
      historyDragBaseIds: [],
      diagnosticGenerating: false,
      autoCaptureTestResults: {},
      spoutStatus: null,
      spoutChecking: false,
      spoutSendersLoading: false
    }
  },
  computed: {
    hasResults() {
      return this.state.results && this.state.results.length > 0
    },
    visibleHistory() {
      return (this.state.history || []).filter((item) => !item.cleared)
    },
    resultItems() {
      if (this.hasResults) return this.state.results
      return this.visibleHistory.map((item) => ({
        ...item,
        historyId: item.id,
        processing: false,
        fromHistory: true
      }))
    },
    hasResultItems() {
      return this.resultItems.length > 0
    },
    hasAnyHistory() {
      return (this.state.history || []).length > 0
    },
    selectedHistoryEntries() {
      const selected = new Set(this.selectedHistoryIds)
      return (this.state.history || []).filter((item) => selected.has(item.id))
    },
    hasDiscordDeletableSelection() {
      return this.selectedHistoryEntries.some((item) => this.canDeleteHistoryDiscord(item))
    },
    hasLocalDeletableSelection() {
      return this.selectedHistoryEntries.some((item) => this.canDeleteHistoryLocal(item))
    },
    hasHistoryDeletableSelection() {
      return this.selectedHistoryEntries.some((item) => !item.pinned)
    },
    isSettings() {
      return this.state.mode === 'settings'
    },
    isError() {
      return this.state.mode === 'error'
    },
    isJpegOutput() {
      return this.state.config?.image?.outputFormat === 'jpg'
    },
    processedCount() {
      return (this.state.results || []).filter((item) => !item.processing).length
    },
    totalProcessingCount() {
      return (this.state.results || []).length
    },
    overallProgress() {
      if (!this.processing || !this.totalProcessingCount) return 0
      return Math.round((this.processedCount / this.totalProcessingCount) * 100)
    },
    outputExample() {
      const suffix = this.state.config?.image?.suffix || '_2048'
      const ext = this.state.config?.image?.outputFormat === 'jpg' ? 'jpg' : 'png'
      return `image.png -> image${suffix}.${ext}`
    },
    activeView() {
      if (this.isSettings) return 'settings'
      return this.view
    },
    hasUnsavedSettings() {
      if (!this.isSettings || !this.state.config || !this.settingsBaseline) return false
      return this.serializeSettings(this.state.config) !== this.settingsBaseline
    },
    historySelectionRectStyle() {
      if (!this.historyDragSelecting) return {}
      const left = Math.min(this.historyDragStart.x, this.historyDragCurrent.x)
      const top = Math.min(this.historyDragStart.y, this.historyDragCurrent.y)
      const width = Math.abs(this.historyDragCurrent.x - this.historyDragStart.x)
      const height = Math.abs(this.historyDragCurrent.y - this.historyDragStart.y)
      return {
        left: `${left}px`,
        top: `${top}px`,
        width: `${width}px`,
        height: `${height}px`
      }
    },
    updateSettings() {
      return this.state.config?.update || { checkEnabled: true, notificationEnabled: true }
    },
    autoCaptureSettings() {
      if (!this.state.config) return {}
      if (!this.state.config.autoCapture) {
        this.state.config.autoCapture = {}
      }
      const autoCapture = this.state.config.autoCapture
      autoCapture.schedule ||= {}
      autoCapture.osc ||= {}
      autoCapture.playerLocal ||= {}
      autoCapture.capture ||= {}
      autoCapture.stream ||= {}
      autoCapture.output ||= {}
      autoCapture.presence ||= {}
      autoCapture.discord ||= {}
      autoCapture.views ||= []
      return autoCapture
    },
    autoCaptureViews() {
      const views = this.autoCaptureSettings.views || []
      views.forEach((view, index) => {
        view.sortOrder = Number.isFinite(Number(view.sortOrder)) ? Number(view.sortOrder) : index
        view.pose ||= {}
        view.pose.position ||= {}
        view.pose.rotation ||= {}
      })
      return views.sort((a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0))
    },
    settingsTabs() {
      return [
        { id: 'feature', label: '機能' },
        { id: 'autoCapture', label: '自動撮影' },
        { id: 'process', label: '処理' },
        { id: 'webhook', label: 'Discord投稿' },
        { id: 'update', label: '更新' }
      ]
    },
    shouldShowUpdateBanner() {
      return Boolean(
        this.updateInfo.available &&
        this.updateSettings.notificationEnabled !== false &&
        !this.updateBannerDismissed
      )
    },
    shouldShowDiscordWebhookWarning() {
      return this.shouldWarnMissingPrimaryWebhook()
    },
    autoPostConfirmationItems() {
      const config = this.state.config || {}
      const items = []
      const autoPhoto = config.autoPhoto || {}
      const screenshot = config.screenshotAutoPost || {}
      const output = config.output || {}
      const discord = config.discord || {}
      if (autoPhoto.enabled) {
        const target = String(autoPhoto.webhookUrl || discord.webhookUrl || '').trim()
        items.push({
          label: 'VRChat写真自動処理',
          detail: `監視フォルダ: ${autoPhoto.photoDirectory || '(未設定)'}`,
          discord: output.uploadDiscord ? `Discord投稿: ON / 送信先: ${target ? this.maskWebhook(target) : '(未設定)'}` : 'Discord投稿: OFF'
        })
      }
      if (screenshot.enabled) {
        const target = String(screenshot.webhookUrl || discord.webhookUrl || '').trim()
        items.push({
          label: 'スクリーンショット自動処理',
          detail: `監視フォルダ: ${screenshot.screenshotDirectory || '(未設定)'}`,
          discord: output.uploadDiscord ? `Discord投稿: ON / 送信先: ${target ? this.maskWebhook(target) : '(未設定)'}` : 'Discord投稿: OFF'
        })
      }
      return items
    },
    shouldConfirmAutoPostSettings() {
      return this.autoPostConfirmationItems.length > 0
    }
  },
  async mounted() {
    this.info = await api.GetAppInfo()
    this.state = await api.GetInitialState()
    this.licenses = await api.GetOSSLicenses()
    window.runtime?.EventsOn?.('process:progress', (event) => {
      this.applyProgress(event)
    })
    window.runtime?.EventsOn?.('auto-photo:result', (event) => {
      this.applyAutoPhotoResult(event)
    })
    window.runtime?.OnFileDrop?.(async (_x, _y, paths) => {
      this.dragging = false
      await this.handleDrop(paths || [])
    }, false)
    window.addEventListener('dragenter', this.showDropOverlay)
    window.addEventListener('dragover', this.showDropOverlay)
    window.addEventListener('dragleave', this.hideDropOverlay)
    window.addEventListener('drop', () => {
      this.dragging = false
    })
    window.addEventListener('keydown', this.handleKeyDown)
    void this.checkForUpdate()
  },
  unmounted() {
    window.removeEventListener('keydown', this.handleKeyDown)
    window.removeEventListener('mousemove', this.updateHistoryDragSelect)
    window.removeEventListener('mouseup', this.finishHistoryDragSelect)
  },
  methods: {
    logUserAction(action, detail = '') {
      if (!api?.LogUserAction) return
      void api.LogUserAction(String(action || ''), String(detail || '')).catch(() => {})
    },
    setView(view, detail = '') {
      if (this.view !== view) {
        this.logUserAction('screen_transition', `${this.view}->${view}${detail ? ` ${detail}` : ''}`)
      }
      this.view = view
    },
    handleKeyDown(event) {
      if (this.view !== 'history' || !(event.ctrlKey || event.metaKey) || event.key.toLowerCase() !== 'a') return
      event.preventDefault()
      this.selectAllHistory()
    },
    showDropOverlay(event) {
      event.preventDefault()
      this.dragging = true
    },
    hideDropOverlay(event) {
      if (event.clientX <= 0 || event.clientY <= 0 || event.clientX >= window.innerWidth || event.clientY >= window.innerHeight) {
        this.dragging = false
      }
    },
    serializeSettings(config) {
      return JSON.stringify(config || {})
    },
    rememberSettingsBaseline() {
      this.settingsBaseline = this.serializeSettings(this.state.config)
      this.settingsTab = 'feature'
    },
    resetSettingsBaseline() {
      this.settingsBaseline = ''
    },
    async leaveSettings(action) {
      if (this.hasUnsavedSettings) {
        this.pendingSettingsLeave = action
        return false
      }
      await this.performSettingsLeave(action)
      return true
    },
    async performSettingsLeave(action) {
      this.error = ''
      if (this.isSettings) {
        this.state = await api.CloseSettings()
        this.resetSettingsBaseline()
      }
      if (action === 'help') {
        this.setView('help', 'after_settings')
      } else if (action === 'about') {
        this.setView('about', 'after_settings')
      } else if (action === 'history') {
        this.setView('history', 'after_settings')
        this.state.history = await api.GetHistory()
      } else if (action === 'drop') {
        const paths = [...this.pendingDropPaths]
        this.pendingDropPaths = []
        this.setView('main', 'drop_after_settings')
        await this.handleDrop(paths, true)
      } else {
        this.setView('main', 'after_settings')
      }
    },
    async confirmSaveAndLeaveSettings(skipAutoPostConfirmation = false, overrideAction = '') {
      const action = overrideAction || this.pendingSettingsLeave || 'home'
      if (!skipAutoPostConfirmation && this.shouldConfirmAutoPostSettings) {
        this.pendingSettingsLeave = null
        this.requestAutoPostConfirmation(`leave:${action}`)
        return
      }
      this.pendingSettingsLeave = null
      this.saving = true
      this.saved = false
      this.error = ''
      try {
        this.sanitizeOutputDirectory()
        this.sanitizePhotoDirectory()
        this.sanitizeScreenshotDirectory()
        if (this.state.processOnSave) {
          this.logUserAction('button_click', 'settings_save_and_process')
          this.state = await api.SaveConfigAndProcess(this.state.config, this.state.pendingPaths || [])
          if (this.shouldWarnMissingPrimaryWebhook()) {
            this.logUserAction('settings_warning', 'missing_primary_discord_webhook')
          }
          this.resetSettingsBaseline()
          this.setView('main', 'save_config_and_process')
          return
        }
        this.logUserAction('button_click', 'settings_save')
        await api.SaveConfig(this.state.config)
        if (this.shouldWarnMissingPrimaryWebhook()) {
          this.logUserAction('settings_warning', 'missing_primary_discord_webhook')
        }
        this.resetSettingsBaseline()
        await this.performSettingsLeave(action)
      } catch (err) {
        this.error = String(err)
        this.pendingSettingsLeave = action
      } finally {
        this.saving = false
      }
    },
    async discardSettingsAndLeave() {
      const action = this.pendingSettingsLeave || 'home'
      this.pendingSettingsLeave = null
      await this.performSettingsLeave(action)
    },
    cancelSettingsLeave() {
      this.logUserAction('button_click', 'cancel_settings_leave')
      this.pendingSettingsLeave = null
      this.pendingDropPaths = []
    },
    async goHome() {
      if (this.isSettings) {
        await this.leaveSettings('home')
        return
      }
      this.setView('main', 'go_home')
    },
    async toggleHelp() {
      this.logUserAction('button_click', 'header_help')
      if (this.activeView === 'help') {
        await this.goHome()
        return
      }
      if (this.isSettings) {
        await this.leaveSettings('help')
        return
      }
      this.setView('help', 'header_help')
    },
    async toggleAbout() {
      this.logUserAction('button_click', 'header_about')
      if (this.activeView === 'about' || this.activeView === 'licenses') {
        await this.goHome()
        return
      }
      if (this.isSettings) {
        await this.leaveSettings('about')
        return
      }
      this.setView('about', 'header_about')
    },
    async toggleSettings() {
      this.logUserAction('button_click', 'header_settings')
      if (this.activeView === 'settings') {
        await this.leaveSettings('home')
        return
      }
      await this.openSettings()
    },
    async openSettings() {
      this.error = ''
      this.setView('main', 'open_settings')
      try {
        this.state = await api.OpenSettings('')
        this.rememberSettingsBaseline()
      } catch (err) {
        this.error = String(err)
      }
    },
    async closeSettings() {
      this.logUserAction('button_click', 'settings_close')
      await this.leaveSettings('home')
    },
    selectSettingsTab(tabId) {
      this.logUserAction('button_click', `settings_tab ${tabId}`)
      this.settingsTab = tabId
    },
    shouldWarnMissingPrimaryWebhook(config = this.state.config) {
      return Boolean(config?.output?.uploadDiscord && !String(config?.discord?.webhookUrl || '').trim())
    },
    maskWebhook(value) {
      const text = String(value || '').trim()
      if (!text) return ''
      const parts = text.split('/')
      if (parts.length < 2) return '(設定済み)'
      const id = parts[parts.length - 2] || ''
      const token = parts[parts.length - 1] || ''
      const shortToken = token.length > 8 ? `${token.slice(0, 4)}...${token.slice(-4)}` : '****'
      return `${id}/${shortToken}`
    },
    requestAutoPostConfirmation(action) {
      this.pendingAutoPostConfirmation = action
      this.logUserAction('settings_confirmation_required', `auto_post ${action}`)
    },
    async confirmAutoPostSettings() {
      const action = this.pendingAutoPostConfirmation || 'save'
      this.pendingAutoPostConfirmation = null
      if (action.startsWith('leave:')) {
        await this.confirmSaveAndLeaveSettings(true, action.slice('leave:'.length) || 'home')
        return
      }
      await this.saveSettings(true)
    },
    cancelAutoPostConfirmation() {
      this.logUserAction('button_click', 'cancel_auto_post_confirmation')
      this.pendingAutoPostConfirmation = null
    },
    openDiscordWebhookSettings() {
      this.logUserAction('button_click', 'open_discord_webhook_settings_from_banner')
      if (this.isSettings) {
        this.settingsTab = 'webhook'
        return
      }
      void this.openSettings().then(() => {
        this.settingsTab = 'webhook'
      })
    },
    resultPlaceholder(path, index, total) {
      const normalized = path || 'clipboard.png'
      const parts = normalized.split(/[\\/]/)
      return {
        sourcePath: normalized,
        name: parts[parts.length - 1] || normalized,
        outputPath: '',
        url: '',
        thumbnail: '',
        error: '',
        processing: true,
        progress: total > 0 ? Math.max(5, Math.round((index / total) * 100)) : 20
      }
    },
    applyProgress(event) {
      if (!event || event.index < 0) return
      const results = [...(this.state.results || [])]
      const current = results[event.index] || this.resultPlaceholder(event.result?.sourcePath, event.index, event.total)
      if (event.stage === 'done') {
        results[event.index] = { ...event.result, processing: false, progress: 100 }
      } else {
        results[event.index] = { ...current, ...event.result, processing: true, progress: 35 }
      }
      const done = results.filter((item) => !item.processing).length
      const message = event.total > 1 ? `画像を処理しています。${done} / ${event.total}` : this.state.message
      this.state = { ...this.state, results, message }
    },
    async applyAutoPhotoResult(event) {
      if (!event?.result) return
      const result = event.result
      const hasVisibleWork = this.hasResultVisibleWork(result)
      const results = hasVisibleWork ? [result, ...(this.state.results || [])] : (this.state.results || [])
      const message = event.error ? '自動処理でエラーが発生しました。' : this.resultSummaryMessage([result])
      this.state = { ...this.state, mode: 'results', results, message }
      try {
        this.state.history = await api.GetHistory()
      } catch {
        // Result display should not fail just because history refresh failed.
      }
      this.toast = event.error ? '自動処理に失敗しました' : (message || '自動処理を確認しました')
      setTimeout(() => {
        this.toast = ''
      }, 2200)
    },
    hasResultVisibleWork(item) {
      return !!(item && (item.error || item.url || item.outputPath || (item.qrUrls && item.qrUrls.length)))
    },
    resultSummaryMessage(items) {
      const results = (items || []).filter((item) => this.hasResultVisibleWork(item))
      if (!results.length) return '実行された処理はありません。'
      const actions = []
      if (results.some((item) => item.url)) actions.push('Discord投稿')
      if (results.some((item) => item.outputPath)) actions.push('ローカル保存')
      if (results.some((item) => item.qrUrls && item.qrUrls.length)) actions.push('QRコードURL取得')
      if (!actions.length) return '処理しました。'
      return `${actions.join('、')}を行いました。`
    },
    async handleDrop(paths, skipSettingsGuard = false) {
      this.error = ''
      this.saved = false
      this.logUserAction('input', `drop count=${(paths || []).length}`)
      if (!skipSettingsGuard && this.isSettings && this.hasUnsavedSettings) {
        this.pendingDropPaths = [...paths]
        this.pendingSettingsLeave = 'drop'
        return
      }
      this.setView('main', 'drop')
      if (!paths.length) return
      const jsonPaths = paths.filter((path) => path.toLowerCase().endsWith('.json'))
      try {
        if (jsonPaths.length === 1 && paths.length === 1) {
          this.state = await api.OpenSettings(jsonPaths[0])
          this.rememberSettingsBaseline()
          return
        }
        if (jsonPaths.length) {
          this.state = {
            ...this.state,
            mode: 'error',
            message: '画像ファイルと設定ファイルが混在しています。設定編集と画像処理は別々に実行してください。',
            results: []
          }
          return
        }
        this.processing = true
        this.state = {
          ...this.state,
          mode: 'results',
          message: `画像を処理しています。0 / ${paths.length}`,
          results: paths.map((path, index) => this.resultPlaceholder(path, index, paths.length))
        }
        this.state = await api.ProcessToState(paths)
      } catch (err) {
        this.error = String(err)
      } finally {
        this.processing = false
      }
    },
    async processClipboard() {
      this.error = ''
      this.logUserAction('button_click', 'process_clipboard')
      if (this.isSettings && !(await this.leaveSettings('home'))) return
      this.setView('main', 'process_clipboard')
      try {
        this.processing = true
        this.state = {
          ...this.state,
          mode: 'results',
          message: 'クリップボード画像を処理しています。',
          results: [this.resultPlaceholder('clipboard.png', 0, 1)]
        }
        this.state = await api.ProcessToState([])
      } catch (err) {
        this.error = String(err)
      } finally {
        this.processing = false
      }
    },
    async clearResults() {
      this.error = ''
      this.logUserAction('button_click', this.hasResults ? 'clear_results' : 'clear_visible_history')
      if (this.hasResults) {
        this.state = await api.ClearResults()
      } else {
        const ids = this.visibleHistory.map((item) => item.id)
        this.state.history = await api.MarkHistoryCleared(ids)
      }
      this.setView('main', 'clear_results')
    },
    async openHistory() {
      this.logUserAction('button_click', 'open_history')
      if (this.isSettings) {
        await this.leaveSettings('history')
        return
      }
      this.setView('history', 'open_history')
      this.state.history = await api.GetHistory()
    },
    selectAllHistory() {
      this.selectedHistoryIds = (this.state.history || []).map((item) => item.id).filter(Boolean)
      this.lastSelectedHistoryIndex = this.selectedHistoryIds.length ? 0 : -1
      this.logUserAction('button_click', `history_select_all count=${this.selectedHistoryIds.length}`)
    },
    selectHistory(event, index, item) {
      if (this.historyDragSelecting) return
      if (!item?.id) return
      const ids = [...this.selectedHistoryIds]
      if (event.shiftKey && this.lastSelectedHistoryIndex >= 0) {
        const start = Math.min(this.lastSelectedHistoryIndex, index)
        const end = Math.max(this.lastSelectedHistoryIndex, index)
        const range = (this.state.history || []).slice(start, end + 1).map((entry) => entry.id)
        this.selectedHistoryIds = Array.from(new Set([...ids, ...range]))
        this.logUserAction('selection', `history_range selected=${this.selectedHistoryIds.length}`)
        return
      }
      if (event.ctrlKey || event.metaKey) {
        this.selectedHistoryIds = ids.includes(item.id) ? ids.filter((id) => id !== item.id) : [...ids, item.id]
        this.lastSelectedHistoryIndex = index
        this.logUserAction('selection', `history_toggle id=${item.id} selected=${this.selectedHistoryIds.length}`)
        return
      }
      this.selectedHistoryIds = [item.id]
      this.lastSelectedHistoryIndex = index
      this.logUserAction('selection', `history_single id=${item.id}`)
    },
    startHistoryDragSelect(event) {
      if (event.button !== 0 || event.target.closest('.history-card')) return
      const grid = event.currentTarget
      const rect = grid.getBoundingClientRect()
      const point = {
        x: event.clientX - rect.left + grid.scrollLeft,
        y: event.clientY - rect.top + grid.scrollTop
      }
      this.historyDragSelecting = true
      this.historySelectionAdditive = event.ctrlKey || event.metaKey
      this.historyDragStart = point
      this.historyDragCurrent = point
      this.historyDragBaseIds = this.historySelectionAdditive ? [...this.selectedHistoryIds] : []
      if (!this.historySelectionAdditive) {
        this.selectedHistoryIds = []
      }
      window.addEventListener('mousemove', this.updateHistoryDragSelect)
      window.addEventListener('mouseup', this.finishHistoryDragSelect)
      event.preventDefault()
    },
    updateHistoryDragSelect(event) {
      if (!this.historyDragSelecting) return
      const grid = this.$refs.historyGrid
      if (!grid) return
      const rect = grid.getBoundingClientRect()
      this.historyDragCurrent = {
        x: event.clientX - rect.left + grid.scrollLeft,
        y: event.clientY - rect.top + grid.scrollTop
      }
      const selectionRect = this.normalizedHistorySelectionRect()
      const selected = this.historySelectionAdditive ? new Set(this.historyDragBaseIds) : new Set()
      grid.querySelectorAll('.history-card').forEach((card) => {
        const id = card.dataset.historyId
        if (!id) return
        if (this.rectsIntersect(selectionRect, this.elementRectInGrid(card, grid))) {
          selected.add(id)
        }
      })
      this.selectedHistoryIds = Array.from(selected)
    },
    finishHistoryDragSelect() {
      if (!this.historyDragSelecting) return
      window.removeEventListener('mousemove', this.updateHistoryDragSelect)
      window.removeEventListener('mouseup', this.finishHistoryDragSelect)
      this.historyDragSelecting = false
      this.historyDragBaseIds = []
      this.logUserAction('selection', `history_drag selected=${this.selectedHistoryIds.length}`)
    },
    normalizedHistorySelectionRect() {
      const left = Math.min(this.historyDragStart.x, this.historyDragCurrent.x)
      const top = Math.min(this.historyDragStart.y, this.historyDragCurrent.y)
      const right = Math.max(this.historyDragStart.x, this.historyDragCurrent.x)
      const bottom = Math.max(this.historyDragStart.y, this.historyDragCurrent.y)
      return { left, top, right, bottom }
    },
    elementRectInGrid(element, grid) {
      const elementRect = element.getBoundingClientRect()
      const gridRect = grid.getBoundingClientRect()
      return {
        left: elementRect.left - gridRect.left + grid.scrollLeft,
        top: elementRect.top - gridRect.top + grid.scrollTop,
        right: elementRect.right - gridRect.left + grid.scrollLeft,
        bottom: elementRect.bottom - gridRect.top + grid.scrollTop
      }
    },
    rectsIntersect(a, b) {
      return a.left <= b.right && a.right >= b.left && a.top <= b.bottom && a.bottom >= b.top
    },
    isHistorySelected(id) {
      return this.selectedHistoryIds.includes(id)
    },
    async toggleHistoryPinned(item) {
      if (!item?.id) return
      this.error = ''
      this.logUserAction('button_click', `history_pin id=${item.id} next=${!item.pinned}`)
      try {
        this.state.history = await api.SetHistoryPinned(item.id, !item.pinned)
      } catch (err) {
        this.error = String(err)
      }
    },
    isTrustedDiscordImageURL(url) {
      try {
        const parsed = new URL(String(url || ''))
        return parsed.protocol === 'https:' &&
          ['cdn.discordapp.com', 'media.discordapp.net'].includes(parsed.hostname.toLowerCase()) &&
          parsed.pathname.startsWith('/attachments/')
      } catch {
        return false
      }
    },
    hasHistoryDiscordRecord(item) {
      return !!(item && this.isTrustedDiscordImageURL(item.url))
    },
    canDeleteHistoryDiscord(item) {
      return !!(item && !item.pinned && !item.discordDeleted && this.hasHistoryDiscordRecord(item) && item.discordMessageId && item.discordWebhookId && item.discordToken)
    },
    historyDiscordLabel(item) {
      if (!this.hasHistoryDiscordRecord(item)) return 'Discord: なし'
      if (item.discordDeleted) return 'Discord: 削除済み'
      return 'Discord: あり'
    },
    historyLocalLabel(item) {
      if (!item?.outputPath) return 'ローカル: なし'
      if (item.localDeleted) return 'ローカル: 削除済み'
      return item.localExists ? 'ローカル: あり' : 'ローカル: なし'
    },
    canDeleteHistoryLocal(item) {
      return !!(item && !item.pinned && item.outputPath && !item.localDeleted && item.localExists)
    },
    async deleteSelectedFromDiscord() {
      if (!this.hasDiscordDeletableSelection) return
      this.error = ''
      this.logUserAction('button_click', `history_delete_discord selected=${this.selectedHistoryIds.length}`)
      try {
        this.state.history = await api.DeleteDiscordHistoryEntries(this.selectedHistoryIds)
        this.toast = 'Discordから削除しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async deleteSelectedLocalFiles() {
      if (!this.hasLocalDeletableSelection) return
      this.error = ''
      this.logUserAction('button_click', `history_delete_local selected=${this.selectedHistoryIds.length}`)
      try {
        this.state.history = await api.DeleteLocalHistoryFiles(this.selectedHistoryIds)
        this.toast = 'ローカルファイルを削除しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async deleteSelectedHistoryEntries() {
      if (!this.hasHistoryDeletableSelection) return
      this.error = ''
      this.logUserAction('button_click', `history_delete_entries selected=${this.selectedHistoryIds.length}`)
      try {
        this.state.history = await api.DeleteHistoryEntries(this.selectedHistoryIds)
        this.selectedHistoryIds = this.selectedHistoryIds.filter((id) => (this.state.history || []).some((item) => item.id === id))
        this.toast = '履歴から削除しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async purgeDeletedHistory() {
      this.error = ''
      try {
        this.state.history = await api.PurgeDeletedHistoryEntries()
        this.selectedHistoryIds = this.selectedHistoryIds.filter((id) => (this.state.history || []).some((item) => item.id === id))
        this.toast = '削除済みURLの履歴を整理しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async copy(url) {
      if (!url) return
      this.logUserAction('button_click', 'copy_text')
      await api.CopyText(url)
      this.toast = 'コピーしました'
      setTimeout(() => {
        this.toast = ''
      }, 1800)
    },
    async copyQRURL(event, url) {
      event.stopPropagation()
      await this.copy(url)
    },
    canCopyResultURL(item) {
      return !!(item && item.url && !item.processing)
    },
    canRevealResultFile(item) {
      return !!(item && item.outputPath && !item.processing)
    },
    hasResultImageAction(item) {
      return this.canCopyResultURL(item) || this.canRevealResultFile(item)
    },
    resultImageActionLabel(item) {
      if (this.canCopyResultURL(item) && this.canRevealResultFile(item)) return '上: URLをコピー / 下: 保存先で表示'
      if (this.canCopyResultURL(item)) return 'URLをコピー'
      if (this.canRevealResultFile(item)) return '保存先で表示'
      return ''
    },
    async copyResultURL(event, item) {
      event.stopPropagation()
      if (!this.canCopyResultURL(item)) return
      this.logUserAction('button_click', `result_copy_url name=${item.name || ''}`)
      await this.copy(item.url)
    },
    async revealResultFile(event, item) {
      event.stopPropagation()
      if (!this.canRevealResultFile(item)) return
      this.logUserAction('button_click', `result_reveal_file name=${item.name || ''}`)
      try {
        await api.RevealFileInExplorer(item.outputPath)
      } catch (err) {
        this.error = String(err)
      }
    },
    async openURL(url) {
      this.logUserAction('button_click', `open_url ${url}`)
      await api.OpenURL(url)
    },
    async createDiagnosticPackage() {
      if (this.diagnosticGenerating) return
      this.error = ''
      this.diagnosticGenerating = true
      this.logUserAction('button_click', 'create_diagnostic_package')
      try {
        const path = await api.CreateEncryptedDiagnosticPackage()
        await api.RevealFileInExplorer(path)
        this.toast = `不具合報告用データを作成しました: ${path}`
        setTimeout(() => {
          this.toast = ''
        }, 5000)
      } catch (err) {
        this.error = String(err)
      } finally {
        this.diagnosticGenerating = false
      }
    },
    normalizeAutoCaptureViewOrders() {
      this.autoCaptureViews.forEach((view, index) => {
        view.sortOrder = index
      })
    },
    newAutoCaptureView(overrides = {}) {
      const id = overrides.id || this.uniqueAutoCaptureViewId('view')
      return {
        id,
        name: overrides.name || `構図 ${this.autoCaptureViews.length + 1}`,
        enabled: overrides.enabled ?? true,
        sortOrder: overrides.sortOrder ?? this.autoCaptureViews.length,
        coordinateSpace: overrides.coordinateSpace || 'world',
        pose: {
          position: { x: 0, y: 0, z: 0, ...(overrides.pose?.position || {}) },
          rotation: { x: 0, y: 0, z: 0, ...(overrides.pose?.rotation || {}) }
        },
        zoom: overrides.zoom ?? null,
        exposure: overrides.exposure ?? null,
        focalDistance: overrides.focalDistance ?? null,
        aperture: overrides.aperture ?? null,
        lookAtMe: overrides.lookAtMe ?? null,
        showUiInCamera: overrides.showUiInCamera ?? null,
        localPlayer: overrides.localPlayer ?? null,
        remotePlayer: overrides.remotePlayer ?? null,
        environment: overrides.environment ?? null,
        settleDelayMs: overrides.settleDelayMs ?? 1500,
        captureDelayMs: overrides.captureDelayMs ?? 0,
        calibrated: overrides.calibrated ?? false
      }
    },
    uniqueAutoCaptureViewId(prefix) {
      const used = new Set(this.autoCaptureViews.map((view) => view.id))
      const base = String(prefix || 'view').trim().replace(/[^a-zA-Z0-9_-]+/g, '-').replace(/^-+|-+$/g, '') || 'view'
      let id = base
      let index = 2
      while (used.has(id)) {
        id = `${base}-${index}`
        index += 1
      }
      return id
    },
    addAutoCaptureView() {
      this.autoCaptureSettings.views.push(this.newAutoCaptureView({ calibrated: false, coordinateSpace: 'player_local' }))
      this.normalizeAutoCaptureViewOrders()
    },
    duplicateAutoCaptureView(view) {
      const copy = JSON.parse(JSON.stringify(view || {}))
      copy.id = this.uniqueAutoCaptureViewId(`${copy.id || 'view'}-copy`)
      copy.name = `${copy.name || '構図'} コピー`
      copy.sortOrder = this.autoCaptureViews.indexOf(view) + 1
      const index = this.autoCaptureSettings.views.indexOf(view)
      this.autoCaptureSettings.views.splice(index + 1, 0, this.newAutoCaptureView(copy))
      this.normalizeAutoCaptureViewOrders()
    },
    deleteAutoCaptureView(view) {
      const index = this.autoCaptureSettings.views.indexOf(view)
      if (index >= 0) {
        this.autoCaptureSettings.views.splice(index, 1)
        this.normalizeAutoCaptureViewOrders()
      }
    },
    moveAutoCaptureView(view, direction) {
      const views = this.autoCaptureViews
      const index = views.indexOf(view)
      const nextIndex = index + direction
      if (index < 0 || nextIndex < 0 || nextIndex >= views.length) return
      const currentOrder = views[index].sortOrder
      views[index].sortOrder = views[nextIndex].sortOrder
      views[nextIndex].sortOrder = currentOrder
      this.normalizeAutoCaptureViewOrders()
    },
    viewPoseValue(view, part, axis) {
      const value = view?.pose?.[part]?.[axis]
      return Number.isFinite(Number(value)) ? Number(value).toFixed(3) : '-'
    },
    coordinateSpaceLabel(value) {
      if (value === 'world') return 'ワールド'
      if (value === 'player_local') return 'プレイヤー基準'
      return '未保存'
    },
    syncAutoCaptureConfig(updated) {
      if (updated?.autoCapture) {
        this.state.config.autoCapture = updated.autoCapture
      } else if (updated?.views) {
        this.state.config.autoCapture.views = updated.views
      }
      this.normalizeAutoCaptureViewOrders()
    },
    async saveCurrentCameraPoseToView(view) {
      if (!view?.id) return
      if (!api?.SaveCurrentCameraPoseToView) {
        this.error = '現在Pose保存APIが利用できません。'
        return
      }
      this.error = ''
      try {
        const updated = await api.SaveCurrentCameraPoseToView(view.id)
        this.syncAutoCaptureConfig(updated)
        if (updated?.id) {
          Object.assign(view, updated)
        }
        this.toast = '現在Poseを保存しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async addCurrentCameraPoseAsView(view) {
      if (!api?.AddCurrentCameraPoseAsView) {
        this.error = '現在Pose追加APIが利用できません。'
        return
      }
      this.error = ''
      try {
        const updated = await api.AddCurrentCameraPoseAsView(view?.id || '')
        this.syncAutoCaptureConfig(updated)
        if (updated?.id) {
          this.autoCaptureSettings.views.push(this.newAutoCaptureView(updated))
        }
        this.toast = '現在Poseから構図を追加しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async saveCurrentCameraPoseAsPlayerBasis() {
      if (!api?.SaveCurrentCameraPoseAsPlayerBasis) {
        this.error = 'プレイヤー基準Pose保存APIが利用できません。'
        return
      }
      this.error = ''
      try {
        const updated = await api.SaveCurrentCameraPoseAsPlayerBasis()
        this.syncAutoCaptureConfig(updated)
        this.toast = 'プレイヤー基準Poseを保存しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async moveCameraToView(view) {
      if (!view?.id) return
      if (!api?.MoveCameraToView) {
        this.error = 'カメラ移動APIが利用できません。'
        return
      }
      this.error = ''
      try {
        await api.MoveCameraToView(view.id)
        this.toast = 'カメラを構図のPoseへ移動しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async resetCameraOSC() {
      if (!api?.ResetCameraOSC) {
        this.error = 'カメラOSCリセットAPIが利用できません。'
        return
      }
      this.error = ''
      try {
        await api.ResetCameraOSC()
        this.toast = 'カメラOSCをリセットしました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async resetCameraPoseToDefault(view) {
      if (!view?.id) return
      if (!api?.ResetCameraPoseToDefault) {
        this.error = '初期PoseリセットAPIが利用できません。'
        return
      }
      this.error = ''
      try {
        const updated = await api.ResetCameraPoseToDefault(view.id)
        this.syncAutoCaptureConfig(updated)
        this.toast = '初期Poseへ戻しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async resetCameraViewsToDefaults() {
      if (!api?.ResetCameraViewsToDefaults) {
        this.error = '初期構図リセットAPIが利用できません。'
        return
      }
      this.error = ''
      try {
        const updated = await api.ResetCameraViewsToDefaults()
        this.syncAutoCaptureConfig(updated)
        this.toast = '初期構図へ戻しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
      }
    },
    async testAutoCaptureView(view) {
      if (!view?.id) return
      if (!api?.TestAutoCaptureView) {
        this.error = 'テスト撮影APIが利用できません。'
        return
      }
      this.error = ''
      this.toast = 'テスト撮影中です'
      try {
        const results = await api.TestAutoCaptureView(view.id)
        const entry = this.formatAutoCaptureTestResult(results)
        this.autoCaptureTestResults = { ...this.autoCaptureTestResults, [view.id]: entry }
        const firstError = entry.results.find((result) => result?.error)?.error
        if (firstError) {
          this.error = firstError
        }
        this.toast = 'テスト撮影を実行しました'
        setTimeout(() => {
          this.toast = ''
        }, 1800)
      } catch (err) {
        this.error = String(err)
        this.autoCaptureTestResults = {
          ...this.autoCaptureTestResults,
          [view.id]: { ok: false, message: String(err), results: [], updatedAt: new Date().toLocaleString() }
        }
        this.toast = ''
      }
    },
    formatAutoCaptureTestResult(results) {
      const normalized = (results || []).map((result) => ({
        ...result,
        sidecarPath: result?.outputPath ? `${result.outputPath}.json` : ''
      }))
      const firstError = normalized.find((result) => result?.error)?.error
      const successCount = normalized.filter((result) => result?.outputPath || result?.url || result?.sourcePath).length
      let message = 'テスト撮影結果が返りませんでした。'
      if (firstError) {
        message = firstError
      } else if (successCount > 0) {
        message = `テスト撮影に成功しました。${successCount}件の結果があります。`
      }
      return {
        ok: !firstError && successCount > 0,
        message,
        results: normalized,
        updatedAt: new Date().toLocaleString()
      }
    },
    async revealAutoCaptureResult(result) {
      const path = result?.outputPath || result?.sourcePath
      if (!path || !api?.RevealFileInExplorer) return
      this.error = ''
      try {
        await api.RevealFileInExplorer(path)
      } catch (err) {
        this.error = String(err)
      }
    },
    async checkSpoutHelper() {
      if (!api?.CheckSpoutHelper) {
        this.error = 'Spout helper確認APIが利用できません。'
        return
      }
      this.error = ''
      this.spoutChecking = true
      try {
        this.spoutStatus = await api.CheckSpoutHelper()
      } catch (err) {
        this.error = String(err)
      } finally {
        this.spoutChecking = false
      }
    },
    async refreshSpoutSenders() {
      if (!api?.ListSpoutSenders) {
        this.error = 'Spout sender一覧APIが利用できません。'
        return
      }
      this.error = ''
      this.spoutSendersLoading = true
      try {
        this.spoutStatus = await api.ListSpoutSenders()
      } catch (err) {
        this.error = String(err)
      } finally {
        this.spoutSendersLoading = false
      }
    },
    async checkForUpdate() {
      if (!api?.CheckForUpdate || this.updateSettings.checkEnabled === false) {
        this.updateInfo = { available: false, currentVersion: '', currentReleaseTime: '', latestVersion: '', latestReleasePublished: '', url: '' }
        return
      }
      try {
        this.updateInfo = await api.CheckForUpdate()
      } catch (_err) {
        this.updateInfo = { available: false, currentVersion: '', currentReleaseTime: '', latestVersion: '', latestReleasePublished: '', url: '' }
      }
    },
    dismissUpdateBanner() {
      this.logUserAction('button_click', 'dismiss_update_banner')
      this.updateBannerDismissed = true
    },
    sanitizeOutputDirectory() {
      if (!this.state.config?.image) return
      this.state.config.image.outputDirectory = String(this.state.config.image.outputDirectory || '').trim().replace(/^"+|"+$/g, '')
    },
    sanitizePhotoDirectory() {
      if (!this.state.config?.autoPhoto) return
      this.state.config.autoPhoto.photoDirectory = String(this.state.config.autoPhoto.photoDirectory || '').trim().replace(/^"+|"+$/g, '')
    },
    sanitizeScreenshotDirectory() {
      if (!this.state.config?.screenshotAutoPost) return
      this.state.config.screenshotAutoPost.screenshotDirectory = String(this.state.config.screenshotAutoPost.screenshotDirectory || '').trim().replace(/^"+|"+$/g, '')
    },
    async chooseOutputDirectory() {
      this.sanitizeOutputDirectory()
      try {
        const selected = await api.SelectOutputDirectory(this.state.config.image.outputDirectory)
        if (selected) {
          this.state.config.image.outputDirectory = selected
        }
      } catch (err) {
        this.error = String(err)
      }
    },
    async choosePhotoDirectory() {
      this.sanitizePhotoDirectory()
      try {
        const selected = await api.SelectPhotoDirectory(this.state.config.autoPhoto.photoDirectory)
        if (selected) {
          this.state.config.autoPhoto.photoDirectory = selected
        }
      } catch (err) {
        this.error = String(err)
      }
    },
    async chooseScreenshotDirectory() {
      this.sanitizeScreenshotDirectory()
      try {
        const selected = await api.SelectScreenshotDirectory(this.state.config.screenshotAutoPost.screenshotDirectory)
        if (selected) {
          this.state.config.screenshotAutoPost.screenshotDirectory = selected
        }
      } catch (err) {
        this.error = String(err)
      }
    },
    async saveSettings(skipAutoPostConfirmation = false) {
      this.saving = true
      this.saved = false
      this.error = ''
      try {
        this.sanitizeOutputDirectory()
        this.sanitizePhotoDirectory()
        this.sanitizeScreenshotDirectory()
        if (!skipAutoPostConfirmation && this.shouldConfirmAutoPostSettings) {
          this.saving = false
          this.requestAutoPostConfirmation('save')
          return
        }
        if (this.state.processOnSave) {
          this.state = await api.SaveConfigAndProcess(this.state.config, this.state.pendingPaths || [])
          if (this.shouldWarnMissingPrimaryWebhook()) {
            this.logUserAction('settings_warning', 'missing_primary_discord_webhook')
          }
        } else {
          await api.SaveConfig(this.state.config)
          if (this.shouldWarnMissingPrimaryWebhook()) {
            this.logUserAction('settings_warning', 'missing_primary_discord_webhook')
          }
          this.state = await api.CloseSettings()
        }
        this.resetSettingsBaseline()
        this.saved = true
        this.updateBannerDismissed = false
        void this.checkForUpdate()
        setTimeout(() => {
          this.saved = false
        }, 1800)
      } catch (err) {
        this.error = String(err)
      } finally {
        this.saving = false
      }
    }
  },
  template: `
    <main class="shell">
      <header>
        <div>
          <h1>{{ info.name }}</h1>
        </div>
        <nav>
          <button :class="{ active: activeView === 'settings' }" @click="toggleSettings">設定</button>
          <button :class="{ active: activeView === 'help' }" @click="toggleHelp">使い方</button>
          <button :class="{ active: activeView === 'about' || activeView === 'licenses' }" @click="toggleAbout">情報</button>
        </nav>
      </header>

      <div v-if="shouldShowUpdateBanner || shouldShowDiscordWebhookWarning" class="banner-stack">
        <div v-if="shouldShowDiscordWebhookWarning" class="update-banner warning-banner">
          <span>Discord投稿がONですが、通常投稿用Webhook URLが空欄です。空の時は投稿できません。</span>
          <div class="update-actions">
            <button @click="openDiscordWebhookSettings">設定する</button>
          </div>
        </div>
        <div v-if="shouldShowUpdateBanner" class="update-banner">
          <span>新しいバージョン {{ updateInfo.latestVersion }} があります。</span>
          <div class="update-actions">
            <button @click="openURL(updateInfo.url || latestReleaseUrl)">GitHub Releases</button>
            <button class="secondary" @click="openURL(boothUrl)">BOOTH</button>
            <button class="icon-button" @click="dismissUpdateBanner" aria-label="更新通知を閉じる" title="更新通知を閉じる">×</button>
          </div>
        </div>
      </div>

      <section v-if="view === 'help'" class="panel help">
        <div class="section-title">
          <h2>使い方</h2>
          <p class="subtle">VRChatで使う画像URLを作るための基本操作です。</p>
        </div>

        <div class="help-grid">
          <article class="help-card">
            <h3>1. 画像を用意する</h3>
            <p>画像ファイルをウィンドウへドラッグ&ドロップします。複数画像もまとめて処理できます。</p>
            <p>スクリーンショットなど、クリップボードに入っている画像は「クリップボード画像を処理」ボタンで処理できます。</p>
          </article>

          <article class="help-card">
            <h3>2. 必要なら設定する</h3>
            <p>設定画面では、VRChat写真やスクリーンショットの自動処理、Discord投稿、QRコードURL検出、ローカル保存を変更できます。</p>
            <p>ローカル保存を使う場合は、出力形式、出力先フォルダ、ファイル名サフィックス、JPEG品質を変更できます。</p>
          </article>

          <article class="help-card">
            <h3>3. Discordへ投稿する</h3>
            <p>Discord投稿を使う場合は、投稿先チャンネルのWebhook URLを設定します。</p>
            <p>Webhook URLの発行方法はDiscord公式ヘルプで確認できます。</p>
            <button class="link-button" @click="openURL(webhookGuideUrl)">Discord公式: ウェブフックのご紹介</button>
          </article>

          <article class="help-card">
            <h3>4. 結果を使う</h3>
            <p>1枚だけ処理した場合は、設定がONなら画像URLを自動でクリップボードへコピーします。</p>
            <p>結果画面では、サムネイルの上側でURLコピー、下側でローカル保存先の表示ができます。URLやローカル保存先がない場合、その操作は表示されません。</p>
          </article>
        </div>

        <div class="button-row">
          <button class="secondary" @click="setView('main', 'help_close')">閉じる</button>
        </div>
      </section>

      <section v-else-if="view === 'about'" class="panel about">
        <h2>このアプリについて</h2>
        <dl>
          <div><dt>プログラム名</dt><dd>{{ info.name }}</dd></div>
          <div><dt>バージョン</dt><dd>{{ info.version }}</dd></div>
          <div><dt>ライセンス</dt><dd>MIT License / Copyright (c) 2026 hatolife</dd></div>
          <div><dt>GitHub</dt><dd><button class="link-button" @click="openURL(info.github)">{{ info.github }}</button></dd></div>
          <div><dt>作者</dt><dd><button class="link-button" @click="openURL(authorTwitterUrl)">@hato_poppo_life</button></dd></div>
          <div><dt>バグ報告</dt><dd><button class="link-button" @click="openURL(issuesUrl)">{{ issuesUrl }}</button></dd></div>
        </dl>
        <section class="about-note">
          <h3>公式の配布場所</h3>
          <p>公式の配布場所は下記のみです。他で取得したファイルについては、内容や安全性に責任を取れません。</p>
          <ul>
            <li><button class="link-button inline" @click="openURL(latestReleaseUrl)">GitHub - https://github.com/hatolife/ClipForVRChat/releases/latest</button></li>
            <li><button class="link-button inline" @click="openURL(boothUrl)">BOOTH - https://hatolife.booth.pm/items/8531663</button></li>
          </ul>
        </section>
        <section class="about-note">
          <h3>PGPで改竄確認できます</h3>
          <p>
            GitHub Releasesでは、通常版exeとは別に <code>ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc</code> 署名ファイルも配布しています。
          </p>
          <p>
            信頼済みfingerprintの公開鍵で検証した場合に限り、ダウンロードした通常版exeと署名の組み合わせを確認できます。
          </p>
          <ol>
            <li><button class="link-button inline" @click="openURL(releasesUrl)">GitHub Releases</button> で使いたいバージョンを開きます。</li>
            <li>確認したい <code>ClipForVRChat-vX.Y.Z-windows-amd64.exe</code> を用意します。</li>
            <li>同じReleaseに別添付されている <code>.exe.asc</code> を、exeと同じフォルダに保存します。</li>
            <li><code>release-signing@hato.life</code> の公開鍵を取り込み、fingerprintが <code>BE40 AA8D 082F 493F 613B C072 21DC 3486 1B40 E77D</code> と一致することを、このアプリのREADMEや公式配布ページなどRelease assetとは別の信頼経路で確認します。</li>
            <li>コマンドプロンプトやPowerShellで <code>gpg --verify ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc ClipForVRChat-vX.Y.Z-windows-amd64.exe</code> を実行します。</li>
            <li>信頼済みfingerprintの公開鍵で <code>Good signature</code> と表示された場合に限り、その鍵で署名されたexeとして確認できます。</li>
          </ol>
          <p>同じReleaseに同梱されたURLや公開鍵だけでは、公開鍵自体の真正性は確認できません。</p>
          <p>PGPがよく分からない場合は、公式の配布場所から直接ダウンロードしてください。</p>
          <ul>
            <li><button class="link-button inline" @click="openURL(latestReleaseUrl)">GitHub - https://github.com/hatolife/ClipForVRChat/releases/latest</button></li>
            <li><button class="link-button inline" @click="openURL(boothUrl)">BOOTH - https://hatolife.booth.pm/items/8531663</button></li>
          </ul>
        </section>
        <section class="about-note">
          <h3>連絡・要望</h3>
          <p>
            不具合や使いにくい点、要望などがありましたら、Twitterの <button class="link-button inline" @click="openURL(authorTwitterUrl)">@hato_poppo_life</button> までお気軽にご連絡ください。
          </p>
          <p>
            「こんな機能が欲しい」といったご意見でも問題ありません。全てのご要望にお応えすることは難しいですが、できる限り改善していきます。
          </p>
          <p>
            できれば
            <button class="link-button inline" @click="openURL(feedbackTweetUrl)">このツイートへの返信</button>
            でご連絡いただけると、管理しやすく助かります。
          </p>
          <p>
            GitHubのIssueからの報告でも問題ありません。日本語で気軽に投稿してください。
          </p>
          <p>
            恐らく反応はTwitterのほうが早いです。GitHub Issueは返信が遅くなる場合がありますが、ご容赦ください。
          </p>
        </section>
        <section class="about-note">
          <h3>不具合報告について</h3>
          <p>
            不具合報告の際は、下記のボタンから生成できる不具合報告用データを添付していただけると助かります。
          </p>
          <p>
            不具合報告用データ作成ボタンを押して作成されるデータはzip形式のデータとzipを暗号化したgpgファイルになります。
          </p>
          <p>
            正式リリースされたプログラムで作成される不具合報告用データは、作者の <button class="link-button inline" @click="openURL(authorTwitterUrl)">@hato_poppo_life</button> のみ復号できます。
          </p>
          <p>
            不具合報告用データには、設定ファイル、履歴、ログ、実行ファイル本体、画像保持フォルダの情報が含まれます。
          </p>
          <p>
            ログや設定などのテキストに含まれるユーザーフォルダのパスは、可能な範囲で <code>%USERPROFILE%</code> などの環境変数表記へ置き換えてから入れます。
          </p>
          <p>
            作成時には時刻付きの作業フォルダを作り、一時的な <code>data</code> フォルダへデータを用意してから、暗号化前zipと暗号化後の <code>.zip.gpg</code> を同じ場所に作成します。<code>data</code> フォルダはzip作成後に削除します。
          </p>
          <p>
            zipは確認用のものです。何が暗号化されたzipに含まれるか確認したいときにご使用ください。gpgファイルが暗号化されたzipファイルで、これを不具合報告で使用してください。
          </p>
          <p>
            誤ってzipを添付された場合、中身が公開されてしまいます。とは言え、含まれるデータに致命的な情報はないと思われますので、基本的に問題はないと思われます。
          </p>
          <p>
            お送りいただいたデータは、不具合の調査および原因解析の目的にのみ使用します。
          </p>
          <div class="button-row">
            <button @click="createDiagnosticPackage" :disabled="diagnosticGenerating">不具合報告用データ生成</button>
          </div>
        </section>
        <section class="about-note">
          <h3>その他</h3>
          <div class="button-row">
            <button @click="setView('licenses', 'about_licenses')">OSSライセンス</button>
            <button class="secondary" @click="setView('main', 'about_close')">閉じる</button>
          </div>
        </section>
      </section>

      <section v-else-if="view === 'licenses'" class="panel licenses">
        <h2>OSSライセンス</h2>
        <p class="subtle">このアプリで使用しているOSSです。</p>
        <div class="license-list">
          <article v-for="license in licenses" :key="license.name" class="license-card">
            <h3>{{ license.name }}</h3>
            <p>{{ license.license }}</p>
            <p>{{ license.copyright }}</p>
            <button class="link-button" @click="openURL(license.url)">{{ license.url }}</button>
          </article>
        </div>
        <button class="secondary" @click="setView('about', 'licenses_back')">戻る</button>
      </section>

      <section v-else-if="view === 'history'" class="panel history-page">
        <div class="history-toolbar">
          <div>
            <h2>画像履歴</h2>
          </div>
          <div class="button-row">
            <span class="tooltip-action">
              <button class="secondary" @click="selectAllHistory" :disabled="!(state.history && state.history.length)">全選択</button>
              <span class="tooltip">履歴に表示されている画像をすべて選択します。Ctrl+Aでも同じ操作ができます。</span>
            </span>
            <span class="tooltip-action">
              <button @click="deleteSelectedFromDiscord" :disabled="!hasDiscordDeletableSelection">Discordから削除</button>
              <span class="tooltip">選択されている履歴のうち、Discord上にあり削除可能な投稿だけを削除します。ピン止めした履歴は対象外です。</span>
            </span>
            <span class="tooltip-action">
              <button class="secondary" @click="deleteSelectedLocalFiles" :disabled="!hasLocalDeletableSelection">ローカルから削除</button>
              <span class="tooltip">選択されている履歴のうち、ローカル保存ファイルがある履歴だけをPCから削除します。ピン止めした履歴は対象外です。</span>
            </span>
            <span class="tooltip-action">
              <button class="secondary danger-button" @click="deleteSelectedHistoryEntries" :disabled="!hasHistoryDeletableSelection">履歴から削除</button>
              <span class="tooltip">選択されている履歴だけを履歴から削除します。Discord投稿やローカルファイルは削除しません。ピン止めした履歴は対象外です。</span>
            </span>
            <span class="tooltip-action">
              <button class="secondary" @click="goHome">閉じる</button>
              <span class="tooltip">履歴画面を閉じて結果画面へ戻ります。</span>
            </span>
          </div>
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <div v-if="state.history && state.history.length" ref="historyGrid" class="history-grid" :class="{ selecting: historyDragSelecting }" @mousedown="startHistoryDragSelect">
          <button v-for="(item, index) in state.history" :key="item.id" class="history-card" :data-history-id="item.id" :class="{ selected: isHistorySelected(item.id), discordDeleted: item.discordDeleted, localDeleted: item.localDeleted, pinned: item.pinned }" @click="selectHistory($event, index, item)">
            <span class="pin-action">
              <span class="pin-button" :class="{ active: item.pinned }" @click.stop="toggleHistoryPinned(item)" title="ピン止め" aria-label="ピン止め"></span>
            </span>
            <div class="thumb-media">
              <img v-if="item.thumbnail || isTrustedDiscordImageURL(item.url)" :src="item.thumbnail || item.url" alt="" />
              <div v-else class="thumb-placeholder"></div>
              <div class="history-badges">
                <span v-if="item.pinned">ピン止め</span>
                <span :class="{ ok: hasHistoryDiscordRecord(item) && !item.discordDeleted, muted: !hasHistoryDiscordRecord(item), deleted: item.discordDeleted }">{{ historyDiscordLabel(item) }}</span>
                <span :class="{ ok: item.localExists, muted: !item.outputPath, deleted: item.localDeleted }">{{ historyLocalLabel(item) }}</span>
                <span v-if="item.qrUrls && item.qrUrls.length" class="ok">QR: {{ item.qrUrls.length }}</span>
              </div>
            </div>
            <span class="history-name">{{ item.name || '画像' }}</span>
            <small>{{ item.createdAt }}</small>
            <div v-if="item.url || item.outputPath" class="history-paths">
              <button v-if="item.url" class="link-button inline" @click.stop="copy(item.url)" :disabled="item.discordDeleted">{{ item.url }}</button>
              <button v-if="item.outputPath && item.localExists && !item.localDeleted" class="link-button inline" @click.stop="revealResultFile($event, item)">{{ item.outputPath }}</button>
            </div>
            <div v-if="item.qrUrls && item.qrUrls.length" class="qr-url-list">
              <strong>QRコードURL</strong>
              <button v-for="qrUrl in item.qrUrls" :key="qrUrl" class="qr-url-chip" @click="copyQRURL($event, qrUrl)" :title="qrUrl">{{ qrUrl }}</button>
            </div>
          </button>
          <div v-if="historyDragSelecting" class="selection-rect" :style="historySelectionRectStyle"></div>
        </div>
        <p v-else class="empty">履歴はありません。</p>
      </section>

      <section v-else-if="isSettings" class="panel settings-page">
        <div class="section-title settings-titlebar">
          <div>
            <h2>設定</h2>
            <p v-if="state.message" class="message" :class="{ warning: isError }">{{ state.message }}</p>
          </div>
          <div v-if="state.config" class="settings-title-actions">
            <button @click="saveSettings()" :disabled="saving">{{ saving ? '保存中' : '保存' }}</button>
            <button class="secondary" @click="closeSettings">閉じる</button>
            <span v-if="saved" class="saved">保存しました</span>
          </div>
        </div>
        <div v-if="state.config" class="settings-layout">
          <div class="settings-topbar">
            <div class="settings-tabs" role="tablist" aria-label="設定カテゴリ">
              <button
                v-for="tab in settingsTabs"
                :key="tab.id"
                type="button"
                role="tab"
                :aria-selected="settingsTab === tab.id"
                :class="{ active: settingsTab === tab.id }"
                @click="selectSettingsTab(tab.id)"
              >{{ tab.label }}</button>
            </div>
          </div>
          <p v-if="error" class="error settings-error">{{ error }}</p>

          <section v-if="settingsTab === 'feature'" class="settings-group" role="tabpanel">
            <h3>機能</h3>
            <div class="setting-row">
              <div><strong>VRChat写真自動処理</strong><p>VRChat上で撮影されたときに処理します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.autoPhoto.enabled" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>スクリーンショット自動処理</strong><p>Win + Shift + Sでスクリーンショットが撮られたときに処理します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.screenshotAutoPost.enabled" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>QRコードURL検出</strong><p>画像内のQRコードからURLを取得します。取得したURLはDiscord本文と結果画面に表示します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.detectQrCodeUrls" /><span></span></label>
            </div>
          </section>

          <section v-if="settingsTab === 'autoCapture'" class="settings-group" role="tabpanel">
            <h3>自動撮影</h3>
            <div class="settings-explainer">
              <strong>VRChatのUser CameraをOSCで操作し、指定間隔で写真を撮影する機能です。</strong>
              <p>VRChat側でOSCを有効にし、Stream方式ではVRChatのStream Camera(Spout)映像を直接受信して静止画として保存します。Photo方式はVRChat標準写真を使うフォールバックです。</p>
              <p>正面、背後、斜めの初期構図にはプレーヤーを写す想定のPoseと拡大率が入っています。構図ごとのテスト撮影で見え方を確認できます。</p>
              <p>v0.1.8では撮影時の解像度変更は行わず、VRChat側の現在のカメラ解像度設定で保存します。</p>
            </div>
            <section class="auto-capture-views" aria-label="構図プリセット">
              <div class="auto-capture-views-header">
                <div>
                  <h4>構図プリセット</h4>
                  <p>「撮影する」がONの構図を上から順番に撮影します。</p>
                  <p>プレイヤー基準構図を使う場合は、基準にしたい位置でカメラPoseを受信してから基準Poseを保存してください。</p>
                </div>
                <div class="button-row">
                  <button type="button" class="secondary" @click="saveCurrentCameraPoseAsPlayerBasis">現在Poseをプレイヤー基準に保存</button>
                  <button type="button" class="secondary" @click="resetCameraViewsToDefaults">初期3構図に戻す</button>
                  <button type="button" class="secondary" @click="resetCameraOSC">カメラOSCをリセット</button>
                </div>
              </div>
              <p class="setting-note" :class="autoCaptureSettings.playerLocal.calibrated ? 'ok' : 'warning'">
                プレイヤー基準Pose: {{ autoCaptureSettings.playerLocal.calibrated ? '保存済み' : '未設定' }}
                <span v-if="autoCaptureSettings.playerLocal.updatedAt"> / {{ autoCaptureSettings.playerLocal.updatedAt }}</span>
              </p>
              <div v-if="autoCaptureViews.length" class="view-list">
                <article v-for="(cameraView, index) in autoCaptureViews" :key="cameraView.id" class="view-card">
                  <div class="view-card-main">
                    <label class="switch compact-switch" :title="cameraView.enabled ? '撮影する' : '撮影しない'">
                      <input type="checkbox" v-model="cameraView.enabled" />
                      <span></span>
                    </label>
                    <div class="view-fields">
                      <label>
                        <small>構図名</small>
                        <input v-model.trim="cameraView.name" placeholder="構図名" />
                      </label>
                      <label>
                        <small>座標系</small>
                        <select v-model="cameraView.coordinateSpace">
                          <option value="world">ワールド</option>
                          <option value="player_local">プレイヤー基準</option>
                        </select>
                      </label>
                      <div class="view-meta">
                        <span :class="['status-pill', cameraView.enabled ? 'ok' : 'muted']">{{ cameraView.enabled ? '撮影する' : '撮影しない' }}</span>
                        <span :class="['status-pill', cameraView.calibrated ? 'ok' : 'muted']">{{ cameraView.calibrated ? 'キャリブレーション済み' : '未キャリブレーション' }}</span>
                        <span class="status-pill">{{ coordinateSpaceLabel(cameraView.coordinateSpace) }}</span>
                      </div>
                    </div>
                  </div>
                  <div class="pose-grid">
                    <label><small>位置 X</small><input type="number" step="0.001" v-model.number="cameraView.pose.position.x" /></label>
                    <label><small>位置 Y</small><input type="number" step="0.001" v-model.number="cameraView.pose.position.y" /></label>
                    <label><small>位置 Z</small><input type="number" step="0.001" v-model.number="cameraView.pose.position.z" /></label>
                    <label><small>回転 X</small><input type="number" step="0.001" v-model.number="cameraView.pose.rotation.x" /></label>
                    <label><small>回転 Y</small><input type="number" step="0.001" v-model.number="cameraView.pose.rotation.y" /></label>
                    <label><small>回転 Z</small><input type="number" step="0.001" v-model.number="cameraView.pose.rotation.z" /></label>
                    <label><small>拡大率</small><input type="number" min="0.1" max="10" step="0.01" v-model.number="cameraView.zoom" /></label>
                  </div>
                  <div class="view-actions">
                    <button type="button" class="secondary" @click="moveAutoCaptureView(cameraView, -1)" :disabled="index === 0">↑</button>
                    <button type="button" class="secondary" @click="moveAutoCaptureView(cameraView, 1)" :disabled="index === autoCaptureViews.length - 1">↓</button>
                    <button type="button" class="secondary" @click="addCurrentCameraPoseAsView(cameraView)">現在Poseから追加</button>
                    <button type="button" class="secondary" @click="saveCurrentCameraPoseToView(cameraView)">現在Poseを保存</button>
                    <button type="button" class="secondary" @click="moveCameraToView(cameraView)">このPoseへカメラ移動</button>
                    <button type="button" class="secondary" @click="testAutoCaptureView(cameraView)">テスト撮影</button>
                    <button type="button" class="secondary" @click="resetCameraPoseToDefault(cameraView)">初期Poseへ戻す</button>
                    <button type="button" class="secondary" @click="duplicateAutoCaptureView(cameraView)">複製</button>
                    <button type="button" class="secondary danger-button" @click="deleteAutoCaptureView(cameraView)">削除</button>
                  </div>
                  <p v-if="autoCaptureTestResults[cameraView.id]" :class="['setting-note', autoCaptureTestResults[cameraView.id].ok ? 'ok' : 'warning']">
                    {{ autoCaptureTestResults[cameraView.id].message }}
                  </p>
                  <div v-if="autoCaptureTestResults[cameraView.id]" class="auto-capture-test-result">
                    <small>更新: {{ autoCaptureTestResults[cameraView.id].updatedAt }}</small>
                    <div v-if="autoCaptureTestResults[cameraView.id].results.length" class="test-result-list">
                      <div v-for="(result, resultIndex) in autoCaptureTestResults[cameraView.id].results" :key="cameraView.id + '-test-' + resultIndex" class="test-result-item">
                        <img v-if="result.thumbnail" :src="result.thumbnail" alt="" class="test-result-thumb" />
                        <dl class="test-result-meta">
                          <div><dt>保存先</dt><dd>{{ result.outputPath || result.sourcePath || 'なし' }}</dd></div>
                          <div><dt>sidecar</dt><dd>{{ result.sidecarPath || 'なし' }}</dd></div>
                          <div><dt>Discord URL</dt><dd>{{ result.url || 'なし' }}</dd></div>
                          <div v-if="result.error"><dt>エラー</dt><dd>{{ result.error }}</dd></div>
                        </dl>
                        <button type="button" class="secondary" @click="revealAutoCaptureResult(result)" :disabled="!(result.outputPath || result.sourcePath)">表示</button>
                      </div>
                    </div>
                  </div>
                </article>
              </div>
              <p v-else class="empty">構図プリセットがありません。</p>
            </section>
            <div class="setting-row">
              <div><strong>自動撮影スケジュール</strong><p>一定間隔でVRChatカメラ撮影を実行します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.schedule.enabled" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.schedule.enabled }">
              <div><strong>撮影間隔</strong><p>1回の撮影開始から次の撮影開始までの秒数です。最小10秒です。</p></div>
              <label>
                <input type="number" min="10" step="1" v-model.number="autoCaptureSettings.schedule.captureIntervalSec" :disabled="!autoCaptureSettings.schedule.enabled" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.schedule.enabled }">
              <div><strong>初回待機時間</strong><p>自動撮影開始後、最初の撮影まで待つ秒数です。</p></div>
              <label>
                <input type="number" min="0" step="1" v-model.number="autoCaptureSettings.schedule.initialDelaySec" :disabled="!autoCaptureSettings.schedule.enabled" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.schedule.enabled }">
              <div><strong>開始時に撮影</strong><p>自動撮影を開始した直後に1回撮影します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.schedule.captureOnStart" :disabled="!autoCaptureSettings.schedule.enabled" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>OSCホスト</strong><p>VRChat OSCへ接続するホストです。</p></div>
              <label>
                <input v-model="autoCaptureSettings.osc.vrcHost" placeholder="127.0.0.1" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>OSC送信ポート</strong><p>外部アプリからVRChatへ送るUDPポートです。</p></div>
              <label>
                <input type="number" min="1" max="65535" step="1" v-model.number="autoCaptureSettings.osc.vrcInPort" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>OSC受信ポート</strong><p>VRChatから外部アプリへ届くUDPポートです。</p></div>
              <label>
                <input type="number" min="1" max="65535" step="1" v-model.number="autoCaptureSettings.osc.appOutPort" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>現在Pose保存の有効秒数</strong><p>「現在Poseを保存」を押した時に、何秒以内にVRChatから受信したPoseなら保存に使うかです。撮影間隔ではありません。</p></div>
              <label>
                <input type="number" min="1" step="1" v-model.number="autoCaptureSettings.osc.poseFreshnessSec" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>撮影方式</strong><p>StreamはStream Camera(Spout)映像を保存します。PhotoはVRChat標準写真を使うためシャッター音が出ます。</p></div>
              <label>
                <select v-model="autoCaptureSettings.capture.mode">
                  <option value="stream">Stream</option>
                  <option value="photo">Photo</option>
                </select>
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: autoCaptureSettings.capture.mode !== 'stream' }">
              <div><strong>Spout helper</strong><p>Stream Camera(Spout)映像を受信する同梱ヘルパーです。通常は初期値のまま使います。</p></div>
              <div class="settings-control-stack">
                <input v-model="autoCaptureSettings.stream.spoutHelperPath" :disabled="autoCaptureSettings.capture.mode !== 'stream' || spoutChecking || spoutSendersLoading" placeholder="spout-capture.exe" />
                <div class="inline-actions">
                  <button type="button" class="secondary" @click="checkSpoutHelper" :disabled="autoCaptureSettings.capture.mode !== 'stream' || spoutChecking || spoutSendersLoading">{{ spoutChecking ? '確認中' : 'helper確認' }}</button>
                  <button type="button" class="secondary" @click="refreshSpoutSenders" :disabled="autoCaptureSettings.capture.mode !== 'stream' || spoutChecking || spoutSendersLoading">{{ spoutSendersLoading ? '取得中' : 'sender一覧更新' }}</button>
                </div>
                <p v-if="spoutStatus" :class="['setting-note', spoutStatus.available ? 'ok' : 'warning']">
                  {{ spoutStatus.message }}
                  <span v-if="spoutStatus.path"> / {{ spoutStatus.path }}</span>
                </p>
              </div>
            </div>
            <div class="setting-row" :class="{ disabled: autoCaptureSettings.capture.mode !== 'stream' }">
              <div><strong>Spout sender自動選択</strong><p>sender未指定時に単一senderまたはVRChatらしいsenderを自動選択します。複数候補で判断不能な場合は一覧から選択してください。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.stream.spoutAutoSelect" :disabled="autoCaptureSettings.capture.mode !== 'stream'" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: autoCaptureSettings.capture.mode !== 'stream' || autoCaptureSettings.stream.spoutAutoSelect }">
              <div><strong>Spout sender名</strong><p>自動選択をOFFにした場合に使うsender名です。VRChatでStream Cameraを起動してから一覧更新してください。</p></div>
              <div class="settings-control-stack">
                <input v-model="autoCaptureSettings.stream.spoutSenderName" :disabled="autoCaptureSettings.capture.mode !== 'stream' || autoCaptureSettings.stream.spoutAutoSelect" list="spout-senders" placeholder="VRChat Stream Camera" />
                <datalist id="spout-senders">
                  <option v-for="sender in spoutStatus?.senders || []" :key="sender.name" :value="sender.name">{{ sender.name }}</option>
                </datalist>
                <p v-if="spoutStatus?.senders?.length" class="setting-note">
                  検出: <span v-for="sender in spoutStatus.senders" :key="sender.name">{{ sender.name }} {{ sender.width && sender.height ? '(' + sender.width + 'x' + sender.height + ')' : '' }} </span>
                </p>
              </div>
            </div>
            <div class="setting-row" :class="{ disabled: autoCaptureSettings.capture.mode !== 'stream' }">
              <div><strong>Stream取得タイムアウト</strong><p>Spout helperが1枚のStream Cameraフレームを保存するまで待つ最大ミリ秒です。</p></div>
              <label>
                <input type="number" min="1000" max="60000" step="500" v-model.number="autoCaptureSettings.stream.captureTimeoutMs" :disabled="autoCaptureSettings.capture.mode !== 'stream'" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: autoCaptureSettings.capture.mode !== 'stream' }">
              <div><strong>Stream起動後待機</strong><p>VRChat Stream CameraをONにしてからSpoutフレーム取得を始めるまで待つミリ秒です。</p></div>
              <label>
                <input type="number" min="0" max="10000" step="100" v-model.number="autoCaptureSettings.stream.startDelayMs" :disabled="autoCaptureSettings.capture.mode !== 'stream'" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>自動撮影保存先</strong><p>Stream方式で保存した画像と関連ファイルの保存先です。Photo方式ではVRChat写真フォルダに保存された画像を使います。</p></div>
              <label>
                <input v-model="autoCaptureSettings.output.directory" placeholder="%USERPROFILE%/Pictures/VRChat/VRC-AutoCapture" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>Stream保存形式</strong><p>Stream方式で保存する画像形式です。Spout helperはPNGで取得し、必要に応じて後続処理で扱います。</p></div>
              <label>
                <select v-model="autoCaptureSettings.output.imageFormat">
                  <option value="png">PNG</option>
                  <option value="jpg">JPG</option>
                </select>
              </label>
            </div>
            <div class="setting-row">
              <div><strong>ファイル名テンプレート</strong><p>{timestamp_local}、{batch_id}、{shot_index}、{view_name}、{mode}、{ext} を使用できます。</p></div>
              <label>
                <input v-model="autoCaptureSettings.output.filenameTemplate" placeholder="{timestamp_local}_{batch_id}_{shot_index}_{view_name}_{mode}.{ext}" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>サイドカーJSON</strong><p>撮影情報を画像と同じ場所にJSONとして保存します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.output.writeSidecarJson" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>画像埋め込みメタデータ</strong><p>PNG/JPEG画像内に自動撮影情報を埋め込みます。sidecar JSONと合わせて画像との紐づけ確認に使います。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.output.writeExif" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.output.writeExif }">
              <div><strong>同席ユーザー一覧を画像に埋め込む</strong><p>取得できた同席ユーザー表示名を画像メタデータに含めます。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.output.writeUserListToExif" :disabled="!autoCaptureSettings.output.writeExif" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.output.writeExif || !autoCaptureSettings.output.writeUserListToExif }">
              <div><strong>ユーザーIDを画像に埋め込む</strong><p>ユーザーIDは識別性が高いため、必要な場合だけONにしてください。sidecar/Discordの設定とは独立しています。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.output.writeUserIdsToExif" :disabled="!autoCaptureSettings.output.writeExif || !autoCaptureSettings.output.writeUserListToExif" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>VRChat output log監視</strong><p>output_logから同じインスタンスのユーザー情報を取得します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.presence.watchOutputLog" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.output.writeSidecarJson }">
              <div><strong>サイドカーJSONにユーザーIDを含める</strong><p>取得できたユーザーIDを撮影情報に保存します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.presence.includeUserIdsInSidecar" :disabled="!autoCaptureSettings.output.writeSidecarJson" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>Discordに表示名を含める</strong><p>自動撮影のDiscord本文に参加者の表示名を含めます。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.presence.includeDisplayNamesInDiscord" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>DiscordにユーザーIDを含める</strong><p>自動撮影のDiscord本文に参加者のユーザーIDを含めます。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.presence.includeUserIdsInDiscord" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>Discord自動投稿</strong><p>自動撮影した画像をDiscord Webhookへ投稿します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.discord.enabled" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.discord.enabled }">
              <div><strong>自動撮影用Webhook URL</strong><p>自動撮影の投稿先Webhook URLです。</p></div>
              <label>
                <input type="password" v-model="autoCaptureSettings.discord.webhookUrl" placeholder="https://discord.com/api/webhooks/..." :disabled="!autoCaptureSettings.discord.enabled" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !autoCaptureSettings.discord.enabled }">
              <div><strong>Discordに画像を添付する</strong><p>OFFの場合は撮影情報の本文だけを投稿します。画像はローカル保存先とsidecar JSONで保持します。</p></div>
              <label class="switch"><input type="checkbox" v-model="autoCaptureSettings.discord.includeImages" :disabled="!autoCaptureSettings.discord.enabled" /><span></span></label>
            </div>
          </section>

          <section v-if="settingsTab === 'process'" class="settings-group" role="tabpanel">
            <h3>処理</h3>
            <div class="setting-row">
              <div><strong>ローカル保存</strong><p>処理した画像をローカルに保存します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.saveLocal" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.saveLocal }">
              <div><strong>出力先フォルダ</strong><p>ローカル保存時の保存先です。初期値はアプリと同じ場所にある output フォルダです。</p></div>
              <div class="input-with-button">
                <input v-model="state.config.image.outputDirectory" @blur="sanitizeOutputDirectory" placeholder="./output" :disabled="!state.config.output.saveLocal" />
                <button class="secondary" @click="chooseOutputDirectory" :disabled="!state.config.output.saveLocal">選択</button>
              </div>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.saveLocal }">
              <div>
                <strong>ファイル名サフィックス</strong>
                <p>ローカル保存時のファイル名末尾に付ける文字です。</p>
                <p>例: {{ outputExample }}</p>
              </div>
              <label>
                <input v-model="state.config.image.suffix" :disabled="!state.config.output.saveLocal" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.saveLocal }">
              <div><strong>出力形式</strong><p>保存または投稿に使う画像形式です。PNGは画質を保ちやすく、JPGは写真向きです。</p></div>
              <label>
                <select v-model="state.config.image.outputFormat" :disabled="!state.config.output.saveLocal">
                  <option value="png">PNG</option>
                  <option value="jpg">JPG</option>
                </select>
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.saveLocal || !isJpegOutput }">
              <div><strong>JPEG品質</strong><p>{{ state.config.output.saveLocal && isJpegOutput ? 'JPG出力時の画質です。数字が大きいほど高画質です。' : 'ローカル保存OFFまたはPNG出力では使用しません。' }}</p></div>
              <label>
                <input type="number" min="1" max="100" v-model.number="state.config.image.jpegQuality" :disabled="!state.config.output.saveLocal || !isJpegOutput" />
              </label>
            </div>
          </section>

          <section v-if="settingsTab === 'update'" class="settings-group" role="tabpanel">
            <h3>更新</h3>
            <div class="setting-row">
              <div><strong>更新確認</strong><p>起動時にGitHub Releasesを確認し、新しいバージョンがあるか調べます。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.update.checkEnabled" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.update.checkEnabled }">
              <div><strong>更新通知</strong><p>新しいバージョンが見つかったとき、画面上部に通知を表示します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.update.notificationEnabled" :disabled="!state.config.update.checkEnabled" /><span></span></label>
            </div>
          </section>

          <section v-if="settingsTab === 'webhook'" class="settings-group" role="tabpanel">
            <h3>Discord投稿</h3>
            <div class="setting-row">
              <div><strong>Discord投稿</strong><p>縮小した画像をDiscord Webhookへ投稿し、VRChatで使うURLを取得します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.uploadDiscord" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.uploadDiscord }">
              <div><strong>投稿URLの自動コピー</strong><p>Discordに投稿したURLをクリップボードに保存します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.copySingleUrlToClipboard" :disabled="!state.config.output.uploadDiscord" /><span></span></label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.uploadDiscord }">
              <div>
                <strong>通常投稿用Webhook URL</strong>
                <p>Discordの投稿先チャンネルでWebhookを作成し、そのURLを貼り付けます。空の時は投稿できません。</p>
                <button class="link-button" @click="openURL(webhookGuideUrl)" :disabled="!state.config.output.uploadDiscord">Discord公式: Webhookの作成方法</button>
              </div>
              <label>
              <input type="password" v-model="state.config.discord.webhookUrl" placeholder="https://discord.com/api/webhooks/..." :disabled="!state.config.output.uploadDiscord" :class="{ 'attention-input': shouldWarnMissingPrimaryWebhook() }" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.uploadDiscord || !state.config.autoPhoto.enabled }">
              <div><strong>VRChat写真用Webhook URL</strong><p>通常投稿とは別の投稿先にしたい場合だけ入力します。空の場合は通常投稿用Webhook URLへ投稿します。</p></div>
              <label>
                <input type="password" v-model="state.config.autoPhoto.webhookUrl" placeholder="空なら通常投稿用Webhook URLを使用" :disabled="!state.config.output.uploadDiscord || !state.config.autoPhoto.enabled" />
              </label>
            </div>
            <div class="setting-row" :class="{ disabled: !state.config.output.uploadDiscord || !state.config.screenshotAutoPost.enabled }">
              <div><strong>スクリーンショット用Webhook URL</strong><p>通常投稿とは別の投稿先にしたい場合だけ入力します。空の場合は通常投稿用Webhook URLへ投稿します。</p></div>
              <label>
                <input type="password" v-model="state.config.screenshotAutoPost.webhookUrl" placeholder="空なら通常投稿用Webhook URLを使用" :disabled="!state.config.output.uploadDiscord || !state.config.screenshotAutoPost.enabled" />
              </label>
            </div>
          </section>
        </div>
      </section>

      <section v-else class="workspace">
        <div class="drop-zone">
          <div class="drop-card">
            <div class="drop-icon">画像</div>
            <h2>画像をここにドラッグ&ドロップ</h2>
            <p>複数画像をまとめて処理できます。config.json をドロップすると設定画面を開きます。</p>
            <span class="tooltip-action">
              <button @click="processClipboard">クリップボード画像を処理</button>
              <span class="tooltip">クリップボードにある画像を縮小し、設定に応じて保存やDiscord投稿を行います。</span>
            </span>
          </div>
        </div>

        <div class="result-panel">
          <div class="result-heading">
            <div>
              <h2>{{ isError ? '確認が必要です' : '結果' }}</h2>
              <p class="subtle">サムネイル上の表示からURLコピーやローカル保存先の表示ができます。</p>
            </div>
            <div class="result-actions">
              <span class="tooltip-action">
                <button class="secondary clear-button" @click="clearResults" :disabled="processing">クリア</button>
                <span class="tooltip">結果一覧から非表示にします。Discord上の画像や履歴データは削除しません。</span>
              </span>
              <span class="tooltip-action">
                <button class="secondary" @click="openHistory" :disabled="processing">履歴</button>
                <span class="tooltip">過去の処理履歴を表示します。Discord、ローカル保存、QRコードURLの状態確認や削除操作ができます。</span>
              </span>
            </div>
          </div>
          <p v-if="state.message" class="message" :class="{ warning: isError }">{{ state.message }}</p>
          <div v-if="processing && totalProcessingCount" class="overall-progress">
            <div><span>全体の進捗</span><strong>{{ processedCount }} / {{ totalProcessingCount }}</strong></div>
            <div class="progress-track"><span :style="{ width: overallProgress + '%' }"></span></div>
          </div>
          <p v-if="error" class="error">{{ error }}</p>

          <div v-if="hasResultItems" class="thumb-grid">
            <article v-for="(item, index) in resultItems" :key="item.name + item.outputPath + item.url + item.error + index" class="thumb-card">
              <div class="thumb-media" :class="{ actionable: hasResultImageAction(item) }" :title="resultImageActionLabel(item)">
                <img v-if="item.thumbnail" :src="item.thumbnail" alt="" />
                <div v-else class="thumb-placeholder">
                  <span class="progress-ring" :style="{ '--progress': (item.progress || 35) + '%' }"></span>
                </div>
                <div v-if="item.processing" class="processing-overlay">
                  <span class="progress-ring" :style="{ '--progress': (item.progress || 35) + '%' }"></span>
                  <strong>処理中</strong>
                </div>
                <div v-else-if="hasResultImageAction(item)" class="result-action-overlay" aria-hidden="true">
                  <div v-if="canCopyResultURL(item)" class="result-action-zone copy-zone">
                    <strong>URLをコピー</strong>
                    <span>Discord投稿URL</span>
                  </div>
                  <div v-if="canRevealResultFile(item)" class="result-action-zone reveal-zone">
                    <strong>保存先で表示</strong>
                    <span>Explorerで画像を選択</span>
                  </div>
                </div>
                <button v-if="canCopyResultURL(item)" class="result-action-button result-action-copy" @click="copyResultURL($event, item)" aria-label="URLをコピー"></button>
                <button v-if="canRevealResultFile(item)" class="result-action-button result-action-reveal" @click="revealResultFile($event, item)" aria-label="保存先で表示"></button>
              </div>
              <span>{{ item.name }}</span>
              <div v-if="item.qrUrls && item.qrUrls.length" class="qr-url-list">
                <strong>QRコードURL</strong>
                <button v-for="qrUrl in item.qrUrls" :key="qrUrl" class="qr-url-chip" @click="copyQRURL($event, qrUrl)" :title="qrUrl">{{ qrUrl }}</button>
              </div>
              <small v-if="item.error" class="error">{{ item.error }}</small>
            </article>
          </div>
          <p v-else class="empty">まだ処理結果はありません。</p>
        </div>
      </section>

      <div v-if="dragging" class="drop-overlay">
        <div>
          <strong>ここにドロップ</strong>
          <span>画像ファイルまたは config.json を処理できます</span>
        </div>
      </div>
      <div v-if="pendingSettingsLeave" class="modal-backdrop" role="dialog" aria-modal="true">
        <div class="confirm-dialog">
          <h2>設定が保存されていません</h2>
          <p>変更した設定を保存してから移動しますか。保存しない場合、変更前の設定が維持されます。</p>
          <p v-if="error" class="error">{{ error }}</p>
          <div class="button-row dialog-actions">
            <button @click="confirmSaveAndLeaveSettings()" :disabled="saving">{{ saving ? '保存中' : '保存して移動' }}</button>
            <button class="secondary" @click="discardSettingsAndLeave" :disabled="saving">保存せずに移動</button>
            <button class="secondary" @click="cancelSettingsLeave" :disabled="saving">キャンセル</button>
          </div>
        </div>
      </div>
      <div v-if="pendingAutoPostConfirmation" class="modal-backdrop" role="dialog" aria-modal="true">
        <div class="confirm-dialog">
          <h2>自動処理の設定を確認してください</h2>
          <p>保存すると、下記の自動処理が有効になり、条件に一致した画像が自動で処理されます。</p>
          <ul class="confirmation-list">
            <li v-for="item in autoPostConfirmationItems" :key="item.label">
              <strong>{{ item.label }}</strong>
              <span>{{ item.detail }}</span>
              <span>{{ item.discord }}</span>
            </li>
          </ul>
          <p>意図しない送信を避けるため、Discord投稿タブでWebhook URLと送信先を確認してから保存してください。</p>
          <p v-if="error" class="error">{{ error }}</p>
          <div class="button-row dialog-actions">
            <button @click="confirmAutoPostSettings" :disabled="saving">{{ saving ? '保存中' : '確認して保存' }}</button>
            <button class="secondary" @click="openDiscordWebhookSettings" :disabled="saving">Discord投稿設定を開く</button>
            <button class="secondary" @click="cancelAutoPostConfirmation" :disabled="saving">キャンセル</button>
          </div>
        </div>
      </div>
      <div v-if="diagnosticGenerating" class="modal-backdrop busy-backdrop" role="dialog" aria-modal="true" aria-live="polite">
        <div class="busy-dialog">
          <h2>不具合報告用データを作成しています</h2>
          <p>ログ、設定、履歴、実行ファイルをzip化し、添付用に暗号化しています。</p>
          <div class="indeterminate-progress"><span></span></div>
        </div>
      </div>
      <div v-if="toast" class="toast">{{ toast }}</div>
    </main>
  `
}).mount('#app')
