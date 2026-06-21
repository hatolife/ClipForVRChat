import { createApp } from 'vue/dist/vue.esm-bundler.js'
import './style.css'

const api = window.go?.main?.App

createApp({
  data() {
    return {
      info: { name: 'ClipForVRChat', version: 'dev', github: 'https://github.com/hatolife/ClipForVRChat' },
      state: { mode: 'results', message: '', configPath: '', config: null, results: [] },
      licenses: [],
      webhookGuideUrl: 'https://support.discord.com/hc/ja/articles/228383668-%E3%82%A6%E3%82%A7%E3%83%96%E3%83%95%E3%83%83%E3%82%AF%E3%81%AE%E3%81%94%E7%B4%B9%E4%BB%8B#:~:text=Webhook%E3%81%AE%E4%BD%9C%E6%88%90',
      issuesUrl: 'https://github.com/hatolife/ClipForVRChat/issues',
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
      pendingSettingsLeave: null,
      pendingDropPaths: [],
      historyDragSelecting: false,
      historySelectionAdditive: false,
      historyDragStart: { x: 0, y: 0 },
      historyDragCurrent: { x: 0, y: 0 },
      historyDragBaseIds: []
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
      return `例: image.png -> image${suffix}.${ext}`
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
  },
  methods: {
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
        this.view = 'help'
      } else if (action === 'about') {
        this.view = 'about'
      } else if (action === 'history') {
        this.view = 'history'
        this.state.history = await api.GetHistory()
      } else if (action === 'drop') {
        const paths = [...this.pendingDropPaths]
        this.pendingDropPaths = []
        this.view = 'main'
        await this.handleDrop(paths, true)
      } else {
        this.view = 'main'
      }
    },
    async confirmSaveAndLeaveSettings() {
      const action = this.pendingSettingsLeave || 'home'
      this.pendingSettingsLeave = null
      this.saving = true
      this.saved = false
      this.error = ''
      try {
        this.sanitizeOutputDirectory()
        this.sanitizePhotoDirectory()
        if (this.state.processOnSave) {
          this.state = await api.SaveConfigAndProcess(this.state.config, this.state.pendingPaths || [])
          this.resetSettingsBaseline()
          this.view = 'main'
          return
        }
        await api.SaveConfig(this.state.config)
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
      this.pendingSettingsLeave = null
      this.pendingDropPaths = []
    },
    async goHome() {
      if (this.isSettings) {
        await this.leaveSettings('home')
        return
      }
      this.view = 'main'
    },
    async toggleHelp() {
      if (this.activeView === 'help') {
        await this.goHome()
        return
      }
      if (this.isSettings) {
        await this.leaveSettings('help')
        return
      }
      this.view = 'help'
    },
    async toggleAbout() {
      if (this.activeView === 'about' || this.activeView === 'licenses') {
        await this.goHome()
        return
      }
      if (this.isSettings) {
        await this.leaveSettings('about')
        return
      }
      this.view = 'about'
    },
    async toggleSettings() {
      if (this.activeView === 'settings') {
        await this.leaveSettings('home')
        return
      }
      await this.openSettings()
    },
    async openSettings() {
      this.error = ''
      this.view = 'main'
      try {
        this.state = await api.OpenSettings('')
        this.rememberSettingsBaseline()
      } catch (err) {
        this.error = String(err)
      }
    },
    async closeSettings() {
      await this.leaveSettings('home')
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
    applyAutoPhotoResult(event) {
      if (!event?.result) return
      const results = [event.result, ...(this.state.results || [])]
      this.state = { ...this.state, mode: 'results', results, message: event.error ? 'VRChat写真の自動処理でエラーが発生しました。' : 'VRChat写真を自動投稿しました。' }
      this.toast = event.error ? '自動投稿に失敗しました' : 'VRChat写真を自動投稿しました'
      setTimeout(() => {
        this.toast = ''
      }, 2200)
    },
    async handleDrop(paths, skipSettingsGuard = false) {
      this.error = ''
      this.saved = false
      if (!skipSettingsGuard && this.isSettings && this.hasUnsavedSettings) {
        this.pendingDropPaths = [...paths]
        this.pendingSettingsLeave = 'drop'
        return
      }
      this.view = 'main'
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
      if (this.isSettings && !(await this.leaveSettings('home'))) return
      this.view = 'main'
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
      if (this.hasResults) {
        this.state = await api.ClearResults()
      } else {
        const ids = this.visibleHistory.map((item) => item.id)
        this.state.history = await api.MarkHistoryCleared(ids)
      }
      this.view = 'main'
    },
    async openHistory() {
      if (this.isSettings) {
        await this.leaveSettings('history')
        return
      }
      this.view = 'history'
      this.state.history = await api.GetHistory()
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
        return
      }
      if (event.ctrlKey || event.metaKey) {
        this.selectedHistoryIds = ids.includes(item.id) ? ids.filter((id) => id !== item.id) : [...ids, item.id]
        this.lastSelectedHistoryIndex = index
        return
      }
      this.selectedHistoryIds = [item.id]
      this.lastSelectedHistoryIndex = index
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
    async deleteSelectedFromDiscord() {
      if (!this.selectedHistoryIds.length) return
      this.error = ''
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
      await api.CopyText(url)
      this.toast = 'コピーしました'
      setTimeout(() => {
        this.toast = ''
      }, 1800)
    },
    async openURL(url) {
      await api.OpenURL(url)
    },
    sanitizeOutputDirectory() {
      if (!this.state.config?.image) return
      this.state.config.image.outputDirectory = String(this.state.config.image.outputDirectory || '').trim().replace(/^"+|"+$/g, '')
    },
    sanitizePhotoDirectory() {
      if (!this.state.config?.autoPhoto) return
      this.state.config.autoPhoto.photoDirectory = String(this.state.config.autoPhoto.photoDirectory || '').trim().replace(/^"+|"+$/g, '')
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
    async saveSettings() {
      this.saving = true
      this.saved = false
      this.error = ''
      try {
        this.sanitizeOutputDirectory()
        this.sanitizePhotoDirectory()
        if (this.state.processOnSave) {
          this.state = await api.SaveConfigAndProcess(this.state.config, this.state.pendingPaths || [])
        } else {
          await api.SaveConfig(this.state.config)
          this.state = await api.CloseSettings()
        }
        this.resetSettingsBaseline()
        this.saved = true
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
          <button :class="{ active: activeView === 'help' }" @click="toggleHelp">使い方</button>
          <button :class="{ active: activeView === 'about' || activeView === 'licenses' }" @click="toggleAbout">情報</button>
          <button :class="{ active: activeView === 'settings' }" @click="toggleSettings">設定</button>
        </nav>
      </header>

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
            <p>設定画面では、ローカル保存、Discord投稿、出力形式、出力先フォルダ、ファイル名サフィックスを変更できます。</p>
            <p>出力形式はPNGまたはJPGを選べます。JPEG品質はJPG出力のときだけ使用します。</p>
          </article>

          <article class="help-card">
            <h3>3. Discordへ投稿する</h3>
            <p>Discord投稿を使う場合は、投稿先チャンネルのWebhook URLを設定します。</p>
            <p>Webhook URLの発行方法はDiscord公式ヘルプで確認できます。</p>
            <button class="link-button" @click="openURL(webhookGuideUrl)">Discord公式: ウェブフックのご紹介</button>
          </article>

          <article class="help-card">
            <h3>4. URLを使う</h3>
            <p>1枚だけ処理した場合は、設定がONなら画像URLを自動でクリップボードへコピーします。</p>
            <p>結果画面では、サムネイルをクリックするとその画像URLを再度コピーできます。</p>
          </article>
        </div>

        <div class="button-row">
          <button class="secondary" @click="view = 'main'">閉じる</button>
        </div>
      </section>

      <section v-else-if="view === 'about'" class="panel about">
        <h2>このアプリについて</h2>
        <dl>
          <div><dt>プログラム名</dt><dd>{{ info.name }}</dd></div>
          <div><dt>バージョン</dt><dd>{{ info.version }}</dd></div>
          <div><dt>ライセンス</dt><dd>MIT License / Copyright (c) 2026 hatolife</dd></div>
          <div><dt>GitHub</dt><dd><button class="link-button" @click="openURL(info.github)">{{ info.github }}</button></dd></div>
          <div><dt>バグ報告</dt><dd><button class="link-button" @click="openURL(issuesUrl)">{{ issuesUrl }}</button></dd></div>
        </dl>
        <p class="subtle">不具合や使いにくい点があれば、バグ報告ページから連絡できます。</p>
        <div class="button-row">
          <button @click="view = 'licenses'">OSSライセンス</button>
          <button class="secondary" @click="view = 'main'">閉じる</button>
        </div>
      </section>

      <section v-else-if="view === 'licenses'" class="panel licenses">
        <h2>OSSライセンス</h2>
        <p class="subtle">このアプリで使用している主なOSSです。</p>
        <div class="license-list">
          <article v-for="license in licenses" :key="license.name" class="license-card">
            <h3>{{ license.name }}</h3>
            <p>{{ license.license }}</p>
            <p>{{ license.copyright }}</p>
            <button class="link-button" @click="openURL(license.url)">{{ license.url }}</button>
          </article>
        </div>
        <button class="secondary" @click="view = 'about'">戻る</button>
      </section>

      <section v-else-if="view === 'history'" class="panel history-page">
        <div class="history-toolbar">
          <div>
            <h2>画像履歴</h2>
            <p class="subtle">削除済みを含む履歴です。選択した画像はDiscord上の投稿だけを削除します。</p>
          </div>
          <div class="button-row">
            <span class="tooltip-action">
              <button @click="deleteSelectedFromDiscord" :disabled="!selectedHistoryIds.length">選択をDiscordから削除</button>
              <span class="tooltip">選択した画像のDiscord上の投稿を削除します。履歴データ自体は残ります。</span>
            </span>
            <span class="tooltip-action">
              <button class="secondary" @click="purgeDeletedHistory">削除済みURLの履歴を削除</button>
              <span class="tooltip">URLを確認し、Discord上で削除済みになっている画像だけを履歴から取り除きます。</span>
            </span>
            <span class="tooltip-action">
              <button class="secondary" @click="goHome">閉じる</button>
              <span class="tooltip">履歴画面を閉じて結果画面へ戻ります。</span>
            </span>
          </div>
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <div v-if="state.history && state.history.length" ref="historyGrid" class="history-grid" :class="{ selecting: historyDragSelecting }" @mousedown="startHistoryDragSelect">
          <button v-for="(item, index) in state.history" :key="item.id" class="history-card" :data-history-id="item.id" :class="{ selected: isHistorySelected(item.id), cleared: item.cleared, deleted: item.discordDeleted }" @click="selectHistory($event, index, item)">
            <div class="thumb-media">
              <img v-if="item.thumbnail || isTrustedDiscordImageURL(item.url)" :src="item.thumbnail || item.url" alt="" />
              <div v-else class="thumb-placeholder"></div>
              <div class="history-badges">
                <span v-if="item.cleared">クリア済み</span>
                <span v-if="item.discordDeleted">Discord削除済み</span>
              </div>
            </div>
            <span>{{ item.name }}</span>
            <small>{{ item.createdAt }}</small>
          </button>
          <div v-if="historyDragSelecting" class="selection-rect" :style="historySelectionRectStyle"></div>
        </div>
        <p v-else class="empty">履歴はありません。</p>
      </section>

      <section v-else-if="isSettings" class="panel settings-page">
        <div class="section-title">
          <h2>設定</h2>
          <p v-if="state.message" class="message" :class="{ warning: isError }">{{ state.message }}</p>
        </div>
        <div v-if="state.config" class="settings-layout">
          <section class="settings-group">
            <h3>出力</h3>
            <div class="setting-row">
              <div><strong>ローカル保存</strong><p>縮小した画像ファイルをPCにも保存します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.saveLocal" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>Discord投稿</strong><p>縮小した画像をDiscord Webhookへ投稿し、VRChatで使うURLを取得します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.uploadDiscord" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>1枚時にURLを自動コピー</strong><p>1枚だけ処理したとき、取得したURLを自動でクリップボードへコピーします。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.output.copySingleUrlToClipboard" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>出力形式</strong><p>PNGは画質を保ちやすく、JPGは写真向きです。</p></div>
              <label>
                <select v-model="state.config.image.outputFormat">
                  <option value="png">PNG</option>
                  <option value="jpg">JPG</option>
                </select>
              </label>
            </div>
            <div class="setting-row">
              <div><strong>最大入力サイズ</strong><p>処理する画像ファイルとクリップボード画像の上限です。大きくしすぎると処理が重くなります。</p></div>
              <label>
                <input type="number" min="1" max="512" v-model.number="state.config.image.maxInputMb" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>JPEG品質</strong><p>{{ isJpegOutput ? 'JPG出力時の画質です。数字が大きいほど高画質です。' : 'PNG出力では使用しません。' }}</p></div>
              <label>
                <input type="number" min="1" max="100" v-model.number="state.config.image.jpegQuality" :disabled="!isJpegOutput" />
              </label>
            </div>
            <div class="setting-row">
              <div><strong>UI表示</strong><p>処理後に画面を表示する条件を選びます。通常はautoのままで問題ありません。</p></div>
              <label>
                <select v-model="state.config.output.showUi">
                  <option value="auto">auto</option>
                  <option value="always">always</option>
                  <option value="never">never</option>
                </select>
              </label>
            </div>
          </section>

          <section class="settings-group">
            <h3>保存</h3>
            <div class="setting-row">
              <div><strong>出力先フォルダ</strong><p>保存先です。初期値はアプリと同じ場所にある output フォルダです。</p></div>
              <div class="input-with-button">
                <input v-model="state.config.image.outputDirectory" @blur="sanitizeOutputDirectory" placeholder="./output" />
                <button class="secondary" @click="chooseOutputDirectory">選択</button>
              </div>
            </div>
            <div class="setting-row">
              <div><strong>サフィックス</strong><p>保存するファイル名の末尾に付ける文字です。{{ outputExample }}</p></div>
              <label>
                <input v-model="state.config.image.suffix" />
              </label>
            </div>
          </section>

          <section class="settings-group">
            <h3>Discord</h3>
            <div class="setting-row">
              <div>
                <strong>Webhook URL</strong>
                <p>Discordの投稿先チャンネルでWebhookを作成し、そのURLを貼り付けます。</p>
                <button class="link-button" @click="openURL(webhookGuideUrl)">Discord公式: Webhookの作成方法</button>
              </div>
              <label>
              <input type="password" v-model="state.config.discord.webhookUrl" placeholder="https://discord.com/api/webhooks/..." />
              </label>
            </div>
          </section>

          <section class="settings-group">
            <h3>VRChat写真自動投稿</h3>
            <div class="setting-row">
              <div><strong>自動投稿</strong><p>VRChatの写真フォルダに新しい画像が保存されたら、自動で縮小してDiscordへ投稿します。</p></div>
              <label class="switch"><input type="checkbox" v-model="state.config.autoPhoto.enabled" /><span></span></label>
            </div>
            <div class="setting-row">
              <div><strong>写真フォルダ</strong><p>VRChatが写真を保存するフォルダです。通常は「ピクチャ」内のVRChatフォルダです。</p></div>
              <div class="input-with-button">
                <input v-model="state.config.autoPhoto.photoDirectory" @blur="sanitizePhotoDirectory" placeholder="C:\\Users\\...\\Pictures\\VRChat" />
                <button class="secondary" @click="choosePhotoDirectory">選択</button>
              </div>
            </div>
            <div class="setting-row">
              <div><strong>自動投稿用Webhook URL</strong><p>通常のDiscord投稿とは別の投稿先にしたい場合だけ入力します。空の場合は上のDiscord Webhook URLへ投稿します。</p></div>
              <label>
                <input type="password" v-model="state.config.autoPhoto.webhookUrl" placeholder="空なら通常のWebhook URLを使用" />
              </label>
            </div>
          </section>

          <div class="button-row footer-actions">
            <button @click="saveSettings" :disabled="saving">{{ saving ? '保存中' : '保存' }}</button>
            <button class="secondary" @click="closeSettings">閉じる</button>
            <span v-if="saved" class="saved">保存しました</span>
            <p v-if="error" class="error">{{ error }}</p>
          </div>
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
              <p class="subtle">サムネイルをクリックすると画像URLをコピーできます。</p>
            </div>
            <div class="result-actions">
              <span class="tooltip-action">
                <button class="secondary clear-button" @click="clearResults" :disabled="processing">クリア</button>
                <span class="tooltip">結果一覧から非表示にします。Discord上の画像や履歴データは削除しません。</span>
              </span>
              <span class="tooltip-action">
                <button class="secondary" @click="openHistory" :disabled="processing">履歴</button>
                <span class="tooltip">過去に取得した画像URLを表示します。削除済み画像の確認やDiscord上の削除操作ができます。</span>
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
            <button v-for="(item, index) in resultItems" :key="item.name + item.outputPath + item.url + item.error + index" class="thumb-card" @click="copy(item.url)" :disabled="!item.url || item.processing">
              <div class="thumb-media">
                <img v-if="item.thumbnail" :src="item.thumbnail" alt="" />
                <div v-else class="thumb-placeholder">
                  <span class="progress-ring" :style="{ '--progress': (item.progress || 35) + '%' }"></span>
                </div>
                <div v-if="item.processing" class="processing-overlay">
                  <span class="progress-ring" :style="{ '--progress': (item.progress || 35) + '%' }"></span>
                  <strong>処理中</strong>
                </div>
                <div v-else-if="item.url" class="copy-overlay">クリックでURLをコピー</div>
              </div>
              <span>{{ item.name }}</span>
              <small v-if="item.error" class="error">{{ item.error }}</small>
            </button>
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
            <button @click="confirmSaveAndLeaveSettings" :disabled="saving">{{ saving ? '保存中' : '保存して移動' }}</button>
            <button class="secondary" @click="discardSettingsAndLeave" :disabled="saving">保存せずに移動</button>
            <button class="secondary" @click="cancelSettingsLeave" :disabled="saving">キャンセル</button>
          </div>
        </div>
      </div>
      <div v-if="toast" class="toast">{{ toast }}</div>
    </main>
  `
}).mount('#app')
