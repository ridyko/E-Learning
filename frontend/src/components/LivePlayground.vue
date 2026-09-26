<script setup>
import { ref, computed } from 'vue'
import { Code, Play, Layers, Sparkles, Copy, Check } from 'lucide-vue-next'

const htmlCode = ref(`<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: #f8fafc;
      color: #0f172a;
      padding: 2rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-height: 80vh;
    }
    .card {
      background: #ffffff;
      border: 1px solid #cbd5e1;
      border-radius: 16px;
      padding: 2rem;
      max-width: 450px;
      box-shadow: 0 10px 30px rgba(0,0,0,0.08);
      text-align: center;
    }
    h2 { color: #1d4ed8; margin-bottom: 0.5rem; }
    p { color: #475569; font-size: 0.95rem; }
    .badge {
      background: #fffbeb;
      color: #b45309;
      border: 1px solid #fde68a;
      padding: 0.3rem 0.8rem;
      border-radius: 99px;
      font-weight: bold;
      font-size: 0.8rem;
    }
    .btn {
      margin-top: 1.5rem;
      background: #2563eb;
      color: white;
      border: none;
      padding: 0.75rem 1.5rem;
      border-radius: 8px;
      font-weight: 600;
      cursor: pointer;
      transition: 0.2s;
    }
    .btn:hover { transform: scale(1.05); }
  </style>
</head>
<body>

  <div class="card">
    <span class="badge">ITB SWADHARMA</span>
    <h2>Pemrograman Web Live Sandbox</h2>
    <p>Selamat mencoba Live Playground Pemrograman Web!</p>
    <button class="btn" onclick="tambahAngka()">Klik Saya: <span id="counter">0</span></button>
  </div>

  <script>
    let count = 0;
    function tambahAngka() {
      count++;
      document.getElementById('counter').textContent = count;
    }
  <\/script>
</body>
</html>`)

const copied = ref(false)

const templates = [
  {
    name: 'Flexbox & CSS Card',
    code: `<!DOCTYPE html>
<html>
<head>
<style>
  body { background: #f8fafc; color: #0f172a; font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
  .card-container { display: flex; gap: 20px; flex-wrap: wrap; }
  .card { background: #ffffff; border: 1px solid #cbd5e1; padding: 20px; border-radius: 12px; width: 200px; text-align: center; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
  .card h3 { color: #d97706; }
</style>
</head>
<body>
  <div class="card-container">
    <div class="card">
      <h3>Rekayasa Perangkat Lunak</h3>
      <p>Modul SDLC & UML</p>
    </div>
    <div class="card">
      <h3>Pemrograman Web</h3>
      <p>Vue.js & Golang API</p>
    </div>
  </div>
</body>
</html>`
  },
  {
    name: 'Form Validasi JS',
    code: `<!DOCTYPE html>
<html>
<head>
<style>
  body { background: #f8fafc; color: #0f172a; font-family: sans-serif; padding: 30px; }
  input { display: block; margin: 10px 0; padding: 10px; width: 100%; border-radius: 6px; border: 1px solid #cbd5e1; background: #ffffff; color: #0f172a; }
  button { background: #059669; color: white; border: none; padding: 10px 20px; border-radius: 6px; cursor: pointer; font-weight: bold; }
</style>
</head>
<body>
  <h2>Form Registrasi Mahasiswa</h2>
  <input type="text" id="nim" placeholder="NIM Mahasiswa (cth: 20260801)">
  <input type="text" id="nama" placeholder="Nama Lengkap">
  <button onclick="validasi()">Daftar Praktikum</button>
  <p id="msg"></p>

  <script>
    function validasi() {
      const nim = document.getElementById('nim').value;
      const nama = document.getElementById('nama').value;
      const msg = document.getElementById('msg');
      if (!nim || !nama) {
        msg.textContent = '❌ Harap isi NIM dan Nama!';
        msg.style.color = '#e11d48';
      } else {
        msg.textContent = '✅ Berhasil mendaftar, ' + nama + ' (' + nim + ')!';
        msg.style.color = '#059669';
      }
    }
  <\/script>
</body>
</html>`
  }
]

const loadTemplate = (tmpl) => {
  htmlCode.value = tmpl.code
}

const srcDoc = computed(() => htmlCode.value)

const copyCode = () => {
  navigator.clipboard.writeText(htmlCode.value)
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}
</script>

<template>
  <div class="playground-wrapper">
    <div class="playground-header">
      <div class="title-group">
        <Sparkles class="header-icon text-gold" />
        <div>
          <h3>Live Code Sandbox — Pemrograman Web</h3>
          <p>Praktikum HTML, CSS & JavaScript</p>
        </div>
      </div>

      <!-- Quick Template Selector -->
      <div class="template-selector">
        <span class="tmpl-label">Template Praktikum:</span>
        <button 
          v-for="tmpl in templates" 
          :key="tmpl.name"
          @click="loadTemplate(tmpl)"
          class="btn btn-secondary btn-sm"
        >
          <Layers class="btn-icon-xs" />
          {{ tmpl.name }}
        </button>
      </div>
    </div>

    <!-- Playground Workspace -->
    <div class="editor-grid">
      <!-- Left: Code Editor -->
      <div class="code-panel glass-card">
        <div class="panel-header">
          <div class="tab-group">
            <button class="tab-btn active">
              <Code class="tab-icon" />
              HTML / CSS / JS Editor
            </button>
          </div>
          <button @click="copyCode" class="btn-action" title="Salin Kode">
            <Check v-if="copied" class="action-icon text-emerald" />
            <Copy v-else class="action-icon" />
          </button>
        </div>
        <textarea 
          v-model="htmlCode" 
          class="code-textarea" 
          placeholder="Tulis kode HTML/CSS/JS di sini..."
          spellcheck="false"
        ></textarea>
      </div>

      <!-- Right: Live Output Preview -->
      <div class="preview-panel glass-card">
        <div class="panel-header">
          <div class="preview-status">
            <Play class="preview-icon text-emerald" />
            <span>Live Output Preview</span>
          </div>
        </div>
        <iframe 
          :srcdoc="srcDoc" 
          class="preview-iframe"
          sandbox="allow-scripts allow-modals"
        ></iframe>
      </div>
    </div>
  </div>
</template>

<style scoped>
.playground-wrapper {
  margin-bottom: 2rem;
}

.playground-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.25rem;
  flex-wrap: wrap;
}

.title-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-icon {
  width: 32px;
  height: 32px;
}

.title-group h3 {
  font-size: 1.3rem;
  font-weight: 800;
  color: #0f172a;
}

.title-group p {
  font-size: 0.85rem;
  color: #475569;
}

.template-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tmpl-label {
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
}

.btn-sm {
  padding: 0.4rem 0.75rem;
  font-size: 0.8rem;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}

.editor-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
  height: 580px;
}

.code-panel, .preview-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: transparent;
  border: none;
  color: #0f172a;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 0.85rem;
}

.tab-icon {
  width: 16px;
  height: 16px;
  color: #1d4ed8;
}

.preview-status {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
}

.preview-icon {
  width: 16px;
  height: 16px;
}

.btn-action {
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 0.3rem;
  border-radius: var(--radius-sm);
}

.btn-action:hover {
  color: #0f172a;
  background: #e2e8f0;
}

.action-icon {
  width: 16px;
  height: 16px;
}

.code-textarea {
  flex: 1;
  width: 100%;
  background: #0f172a;
  color: #f8fafc;
  font-family: var(--font-code);
  font-size: 0.9rem;
  padding: 1rem;
  border: none;
  resize: none;
  outline: none;
  line-height: 1.5;
}

.preview-iframe {
  flex: 1;
  width: 100%;
  border: none;
  background: #ffffff;
}

.text-gold {
  color: #d97706;
}

.text-emerald {
  color: #059669;
}

@media (max-width: 900px) {
  .editor-grid {
    grid-template-columns: 1fr;
    height: auto;
  }
  .code-textarea, .preview-iframe {
    height: 350px;
  }
}
</style>
