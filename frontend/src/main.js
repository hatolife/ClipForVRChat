import { createApp } from 'vue'
import './style.css'

const api = window.go?.main?.App

createApp({
  data() {
    return {
      state: {
        mode: 'results',
        message: '',
        configPath: '',
        config: null,
        results: []
      },
      toast: '',
      saving: false,
      saved: false,
      error: ''
    }
  },
  async mounted() {
    this.state = await api.GetInitialState()
    window.runtime?.OnFileDrop?.(async (_x, _y, paths) => {
      await this.handleDrop(paths || [])
    }, false)
  },
  methods: {
    async openSettings() {
      this.error = ''
      try {
        this.state = await api.OpenSettings('')
      } catch (err) {
        this.error = String(err)
      }
    },
    async handleDrop(paths) {
      this.error = ''
      this.saved = false
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
    async copy(url) {
      if (!url) return
      await api.CopyText(url)
      this.toast = 'コピーしました'
      setTimeout(() => {
        this.toast = ''
      }, 1800)
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
          <h1>ClipForVRChat</h1>
          <p>{{ state.configPath }}</p>
        </div>
        <button class="settings-button" @click="openSettings">設定</button>
      </header>

      <section v-if="state.mode === 'settings'" class="panel drop-target">
        <h2>設定</h2>
        <p v-if="state.message" class="message">{{ state.message }}</p>
        <div v-if="state.config" class="settings">
          <label><input type="checkbox" v-model="state.config.output.saveLocal" /> ローカル保存</label>
          <label><input type="checkbox" v-model="state.config.output.uploadDiscord" /> Discord投稿</label>
          <label><input type="checkbox" v-model="state.config.output.copySingleUrlToClipboard" /> 1枚時にURLを自動コピー</label>

          <label>
            Discord Webhook URL
            <input type="password" v-model="state.config.discord.webhookUrl" placeholder="https://discord.com/api/webhooks/..." />
          </label>

          <label>
            出力先フォルダ
            <input v-model="state.config.image.outputDirectory" placeholder="未指定なら入力画像と同じ場所" />
          </label>

          <div class="grid">
            <label>
              サフィックス
              <input v-model="state.config.image.suffix" />
            </label>
            <label>
              JPEG品質
              <input type="number" min="1" max="100" v-model.number="state.config.image.jpegQuality" />
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

          <button @click="saveSettings" :disabled="saving">{{ saving ? '保存中' : '保存' }}</button>
          <span v-if="saved" class="saved">保存しました</span>
          <p v-if="error" class="error">{{ error }}</p>
        </div>
      </section>

      <section v-else class="panel drop-target">
        <h2>{{ state.mode === 'error' ? '確認が必要です' : '画像URL' }}</h2>
        <p v-if="state.message" class="message">{{ state.message }}</p>

        <div v-if="state.results && state.results.length" class="list">
          <button v-for="item in state.results" :key="item.name + item.outputPath + item.url" class="row" @click="copy(item.url)" :disabled="!item.url">
            <img v-if="item.thumbnail" :src="item.thumbnail" alt="" />
            <div>
              <strong>{{ item.name }}</strong>
              <span v-if="item.outputPath">{{ item.outputPath }}</span>
              <span v-if="item.url">{{ item.url }}</span>
              <span v-if="item.error" class="error">{{ item.error }}</span>
            </div>
          </button>
        </div>
      </section>

      <div v-if="toast" class="toast">{{ toast }}</div>
    </main>
  `
}).mount('#app')
