import { createApp } from 'vue/dist/vue.esm-bundler.js'
import './style.css'

const api = window.go?.main?.App

createApp({
  data() {
    return {
      info: { name: 'ClipForVRChat', version: 'dev', github: 'https://github.com/hatolife/ClipForVRChat' },
      state: { mode: 'results', message: '', configPath: '', config: null, results: [] },
      licenses: [],
      view: 'main',
      toast: '',
      saving: false,
      saved: false,
      error: ''
    }
  },
  computed: {
    hasResults() {
      return this.state.results && this.state.results.length > 0
    },
    isSettings() {
      return this.state.mode === 'settings'
    },
    isError() {
      return this.state.mode === 'error'
    },
    isJpegOutput() {
      return this.state.config?.image?.outputFormat === 'jpg'
    }
  },
  async mounted() {
    this.info = await api.GetAppInfo()
    this.state = await api.GetInitialState()
    this.licenses = await api.GetOSSLicenses()
    window.runtime?.OnFileDrop?.(async (_x, _y, paths) => {
      await this.handleDrop(paths || [])
    }, false)
  },
  methods: {
    async openSettings() {
      this.error = ''
      this.view = 'main'
      try {
        this.state = await api.OpenSettings('')
      } catch (err) {
        this.error = String(err)
      }
    },
    async closeSettings() {
      this.error = ''
      this.state = await api.CloseSettings()
    },
    async handleDrop(paths) {
      this.error = ''
      this.saved = false
      this.view = 'main'
      if (!paths.length) return
      const jsonPaths = paths.filter((path) => path.toLowerCase().endsWith('.json'))
      try {
        if (jsonPaths.length === 1 && paths.length === 1) {
          this.state = await api.OpenSettings(jsonPaths[0])
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
        this.state = await api.ProcessToState(paths)
      } catch (err) {
        this.error = String(err)
      }
    },
    async processClipboard() {
      this.error = ''
      this.view = 'main'
      try {
        this.state = await api.ProcessToState([])
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
    async saveSettings() {
      this.saving = true
      this.saved = false
      this.error = ''
      try {
        if (this.state.processOnSave) {
          this.state = await api.SaveConfigAndProcess(this.state.config, this.state.pendingPaths || [])
        } else {
          await api.SaveConfig(this.state.config)
          this.state = await api.CloseSettings()
        }
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
          <p>version {{ info.version }}</p>
        </div>
        <nav>
          <button class="secondary" @click="view = 'about'; state.mode = 'results'">情報</button>
          <button @click="openSettings">設定</button>
        </nav>
      </header>

      <section v-if="view === 'about'" class="panel about">
        <h2>このアプリについて</h2>
        <dl>
          <div><dt>プログラム名</dt><dd>{{ info.name }}</dd></div>
          <div><dt>バージョン</dt><dd>{{ info.version }}</dd></div>
          <div><dt>GitHub</dt><dd><button class="link-button" @click="openURL(info.github)">{{ info.github }}</button></dd></div>
        </dl>
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

      <section v-else-if="isSettings" class="panel settings-page">
        <div class="section-title">
          <h2>設定</h2>
          <p v-if="state.message" class="message">{{ state.message }}</p>
        </div>
        <div v-if="state.config" class="settings-layout">
          <section class="settings-group">
            <h3>出力</h3>
            <div class="toggle-list">
              <label><input type="checkbox" v-model="state.config.output.saveLocal" /><span>ローカル保存</span></label>
              <label><input type="checkbox" v-model="state.config.output.uploadDiscord" /><span>Discord投稿</span></label>
              <label><input type="checkbox" v-model="state.config.output.copySingleUrlToClipboard" /><span>1枚時にURLを自動コピー</span></label>
            </div>
            <div class="field-grid">
              <label>
                出力形式
                <select v-model="state.config.image.outputFormat">
                  <option value="png">PNG</option>
                  <option value="jpg">JPG</option>
                </select>
              </label>
              <label>
                JPEG品質
                <input type="number" min="1" max="100" v-model.number="state.config.image.jpegQuality" :disabled="!isJpegOutput" />
                <small>{{ isJpegOutput ? 'JPG出力時に使用します。' : 'PNG出力では使用しません。' }}</small>
              </label>
              <label>
                UI表示
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
            <div class="field-grid two">
              <label>
                出力先フォルダ
                <input v-model="state.config.image.outputDirectory" placeholder="未指定なら入力画像と同じ場所" />
              </label>
              <label>
                サフィックス
                <input v-model="state.config.image.suffix" />
              </label>
            </div>
          </section>

          <section class="settings-group">
            <h3>Discord</h3>
            <label>
              Webhook URL
              <input type="password" v-model="state.config.discord.webhookUrl" placeholder="https://discord.com/api/webhooks/..." />
            </label>
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
            <button @click="processClipboard">クリップボード画像を処理</button>
          </div>
        </div>

        <div class="result-panel">
          <h2>{{ isError ? '確認が必要です' : '結果' }}</h2>
          <p class="subtle">サムネイルをクリックすると画像URLをコピーできます。</p>
          <p v-if="state.message" class="message">{{ state.message }}</p>
          <p v-if="error" class="error">{{ error }}</p>

          <div v-if="hasResults" class="thumb-grid">
            <button v-for="item in state.results" :key="item.name + item.outputPath + item.url + item.error" class="thumb-card" @click="copy(item.url)" :disabled="!item.url">
              <img v-if="item.thumbnail" :src="item.thumbnail" alt="" />
              <span>{{ item.name }}</span>
              <small v-if="item.error" class="error">{{ item.error }}</small>
            </button>
          </div>
          <p v-else class="empty">まだ処理結果はありません。</p>
        </div>
      </section>

      <div v-if="toast" class="toast">{{ toast }}</div>
    </main>
  `
}).mount('#app')
