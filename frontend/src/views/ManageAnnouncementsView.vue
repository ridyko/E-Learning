<script setup>
import { ref, onMounted } from 'vue'
import { Megaphone, Plus, Trash2, CheckCircle2 } from 'lucide-vue-next'
import { showSuccess, showError, showWarning, showConfirm } from '../utils/swal.js'

const announcements = ref([])

const newAnnTitle = ref('')
const newAnnCategory = ref('Penting')
const newAnnCourse = ref('all')
const newAnnContent = ref('')

const fetchAnnouncements = async () => {
  try {
    const res = await fetch('/api/announcements')
    if (res.ok) {
      announcements.value = await res.json()
    }
  } catch (err) {
    console.warn('Fetch announcements error:', err)
  }
}

onMounted(() => {
  fetchAnnouncements()
})

const createAnnouncement = async () => {
  if (!newAnnTitle.value || !newAnnContent.value) {
    showWarning('Perhatian', 'Harap isi Judul dan Isi Pengumuman!')
    return
  }

  try {
    const res = await fetch('/api/announcements', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: newAnnTitle.value,
        category: newAnnCategory.value,
        course_id: newAnnCourse.value,
        content: newAnnContent.value
      })
    })

    if (res.ok) {
      showSuccess('Pengumuman Diterbitkan! 📢', 'Pengumuman baru berhasil diterbitkan oleh Pak Rio!')
      newAnnTitle.value = ''
      newAnnContent.value = ''
      fetchAnnouncements()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menerbitkan pengumuman!')
  }
}

const deleteAnnouncement = async (id) => {
  const confirmed = await showConfirm(
    'Hapus Pengumuman?',
    'Apakah Anda yakin ingin menghapus pengumuman ini?',
    'Ya, Hapus'
  )
  if (!confirmed) return

  try {
    const res = await fetch(`/api/announcements?id=${id}`, { method: 'DELETE' })
    if (res.ok) {
      showSuccess('Terhapus! 🗑️', 'Pengumuman berhasil dihapus.')
      fetchAnnouncements()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menghapus pengumuman!')
  }
}
</script>

<template>
  <div class="manage-ann-container animate-fade-in">
    <!-- Header -->
    <div class="page-header">
      <div class="title-row">
        <Megaphone class="header-icon text-gold" />
        <div>
          <h2>📢 Terbitkan & Kelola Pengumuman Perkuliahan</h2>
          <p class="subtitle">Terbitkan pengumuman resmi perkuliahan, tugas, info kuis, dan pengumuman kampus untuk mahasiswa ITB Swadharma.</p>
        </div>
      </div>
    </div>

    <div class="ann-grid">
      <!-- Left: Create Announcement Form -->
      <section class="section-block">
        <div class="glass-card manage-card">
          <div class="card-title-row">
            <Plus class="title-icon text-gold" />
            <div>
              <h3>➕ Form Terbitkan Pengumuman Baru</h3>
              <p class="card-sub">Pengumuman akan langsung tampil secara realtime di Beranda Publik & Mahasiswa.</p>
            </div>
          </div>

          <form @submit.prevent="createAnnouncement" class="ann-form">
            <div>
              <label class="input-label">Judul Pengumuman *</label>
              <input v-model="newAnnTitle" class="glass-input" required placeholder="Contoh: Perubahan Jam Kuliah Pemrograman Web" />
            </div>

            <div class="form-row">
              <div>
                <label class="input-label">Kategori Pengumuman</label>
                <select v-model="newAnnCategory" class="glass-input">
                  <option value="Penting">🔴 Penting</option>
                  <option value="Tugas">📝 Tugas Perkuliahan</option>
                  <option value="Kuis">❓ Kuis Online</option>
                  <option value="Info Akademik">🎓 Info Akademik Kampus</option>
                </select>
              </div>
              <div>
                <label class="input-label">Target Mata Kuliah</label>
                <select v-model="newAnnCourse" class="glass-input">
                  <option value="all">🌐 Semua Mata Kuliah</option>
                  <option value="rpl-2026">Rekayasa Perangkat Lunak</option>
                  <option value="webdev-2026">Pemrograman Web</option>
                </select>
              </div>
            </div>

            <div>
              <label class="input-label">Isi Pengumuman Lengkap *</label>
              <textarea v-model="newAnnContent" class="glass-input textarea" required placeholder="Tuliskan detail pesan pengumuman perkuliahan..."></textarea>
            </div>

            <button type="submit" class="btn btn-gold btn-lg w-full">
              <Megaphone class="btn-icon-sm" />
              <span>Terbitkan Pengumuman Sekarang</span>
            </button>
          </form>
        </div>
      </section>

      <!-- Right: List of Published Announcements -->
      <section class="section-block">
        <div class="glass-card manage-card">
          <div class="card-title-row">
            <Megaphone class="title-icon text-gold" />
            <div>
              <h3>📋 Daftar Pengumuman Aktif saat Ini</h3>
              <p class="card-sub">Manajemen pengumuman yang sudah dipublikasikan.</p>
            </div>
          </div>

          <div v-if="announcements.length === 0" class="empty-notice">
            <p>Belum ada pengumuman yang diterbitkan.</p>
          </div>

          <div v-else class="ann-list">
            <div v-for="ann in announcements" :key="ann.id" class="ann-manage-item">
              <div class="ann-item-header">
                <span :class="['badge', ann.category === 'Penting' ? 'badge-rose' : 'badge-gold']">
                  {{ ann.category }}
                </span>
                <button @click="deleteAnnouncement(ann.id)" class="btn-del" title="Hapus Pengumuman">
                  <Trash2 class="icon-xs" />
                  <span>Hapus</span>
                </button>
              </div>

              <h4 class="ann-item-title">{{ ann.title }}</h4>
              <p class="ann-item-content">{{ ann.content }}</p>

              <div class="ann-item-meta">
                <span>👤 {{ ann.author || 'Rio Widyatmoko' }}</span>
                <span>📅 {{ new Date(ann.created_at).toLocaleString('id-ID') }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.manage-ann-container {
  max-width: 1300px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  margin-bottom: 2rem;
  background: #ffffff;
  padding: 1.5rem 2rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  width: 36px;
  height: 36px;
}

.page-header h2 {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
}

.subtitle {
  color: #475569;
  font-size: 0.92rem;
}

.ann-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.section-block {
  height: 100%;
}

.manage-card {
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  height: 100%;
}

.card-title-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.title-icon {
  width: 24px;
  height: 24px;
  margin-top: 0.2rem;
}

.card-title-row h3 {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
}

.card-sub {
  font-size: 0.88rem;
  color: #64748b;
}

.ann-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  font-weight: 700;
  color: #334155;
  margin-bottom: 0.4rem;
}

.textarea {
  min-height: 120px;
}

.btn-lg {
  padding: 0.75rem 1.5rem;
  font-size: 0.95rem;
}

.w-full {
  width: 100%;
}

.empty-notice {
  text-align: center;
  padding: 3rem 1rem;
  color: #94a3b8;
  font-size: 0.9rem;
}

.ann-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: 600px;
  overflow-y: auto;
  padding-right: 0.25rem;
}

.ann-manage-item {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: var(--radius-sm);
  padding: 1.25rem;
}

.ann-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.6rem;
}

.btn-del {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  background: #fff1f2;
  border: 1px solid #fecdd3;
  color: #be123c;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
}

.btn-del:hover {
  background: #ffe4e6;
}

.ann-item-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.4rem;
}

.ann-item-content {
  font-size: 0.88rem;
  color: #475569;
  margin-bottom: 0.75rem;
  white-space: pre-line;
}

.ann-item-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 500;
}

.badge-rose {
  background: #ffe4e6;
  color: #be123c;
  border: 1px solid #fecdd3;
}

.icon-xs {
  width: 12px;
  height: 12px;
}

.btn-icon-sm {
  width: 16px;
  height: 16px;
}

.text-gold {
  color: #b45309;
}

@media (max-width: 900px) {
  .ann-grid {
    grid-template-columns: 1fr;
  }
}
</style>
