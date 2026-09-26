<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { BookOpen, Plus, Trash2, RotateCcw, ShieldCheck, CheckCircle2, Edit3, Save, X, FolderOpen, Calendar, Clock, MapPin } from 'lucide-vue-next'
import { showSuccess, showError, showWarning, showConfirm } from '../utils/swal.js'

const router = useRouter()
const courses = ref([])

const goToCourseMateri = (courseId) => {
  router.push(`/course/${courseId}`)
}

// Form Fields for New Course
const newCourseCode = ref('')
const newCourseName = ref('')
const newCourseSks = ref(3)
const newCourseSemester = ref('Ganjil 2026/2027')
const newCourseClassTime = ref('Senin, 08.00 - 10.30 WIB')
const newCourseRoom = ref('Lab Komputer / Hybrid Zoom')
const newCourseDescription = ref('')
const newCourseStatus = ref('Aktif')

// Modal & Form Fields for Editing Course
const showEditCourseModal = ref(false)
const isUpdating = ref(false)
const editCourseForm = ref({
  id: '',
  code: '',
  name: '',
  sks: 3,
  semester: '',
  class_time: '',
  room: '',
  description: '',
  status: 'Aktif'
})

const fetchCourses = async () => {
  try {
    const res = await fetch('/api/courses')
    if (res.ok) {
      courses.value = await res.json()
    }
  } catch (err) {
    console.warn('Fetch courses error:', err)
  }
}

onMounted(() => {
  fetchCourses()
})

const toggleCourseStatus = async (c) => {
  const newStatus = c.status === 'Non Aktif' ? 'Aktif' : 'Non Aktif'
  try {
    const res = await fetch('/api/courses', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...c, status: newStatus })
    })

    if (res.ok) {
      showSuccess(
        `Status Diubah! ${newStatus === 'Aktif' ? '🟢' : '🔴'}`,
        `Mata kuliah '${c.name}' sekarang berstatus ${newStatus}.`
      )
      window.dispatchEvent(new Event('course-changed'))
      fetchCourses()
    } else {
      showError('Gagal!', 'Gagal memperbarui status mata kuliah!')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan saat mengubah status.')
  }
}

const createCourse = async () => {
  if (!newCourseCode.value || !newCourseName.value) {
    showWarning('Perhatian', 'Kode Mata Kuliah dan Nama Mata Kuliah wajib diisi!')
    return
  }

  try {
    const res = await fetch('/api/courses', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code: newCourseCode.value,
        name: newCourseName.value,
        sks: parseInt(newCourseSks.value) || 3,
        semester: newCourseSemester.value,
        class_time: newCourseClassTime.value,
        room: newCourseRoom.value,
        description: newCourseDescription.value,
        status: newCourseStatus.value
      })
    })

    if (res.ok) {
      showSuccess('Mata Kuliah Ditambahkan! 📚', `Mata Kuliah '${newCourseName.value}' berhasil ditambahkan!`)
      newCourseCode.value = ''
      newCourseName.value = ''
      newCourseDescription.value = ''
      newCourseStatus.value = 'Aktif'
      window.dispatchEvent(new Event('course-changed'))
      fetchCourses()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menambahkan mata kuliah baru!')
  }
}

// Open Edit Course Modal
const openEditCourseModal = (c) => {
  editCourseForm.value = {
    id: c.id,
    code: c.code || '',
    name: c.name || '',
    sks: c.sks || 3,
    semester: c.semester || 'Ganjil 2026/2027',
    class_time: c.class_time || '',
    room: c.room || '',
    description: c.description || '',
    status: c.status || 'Aktif'
  }
  showEditCourseModal.value = true
}

// Update Existing Course
const updateCourse = async () => {
  if (!editCourseForm.value.code || !editCourseForm.value.name) {
    showWarning('Perhatian', 'Kode dan Nama Mata Kuliah Wajib Diisi!')
    return
  }

  isUpdating.value = true
  try {
    const res = await fetch('/api/courses', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(editCourseForm.value)
    })

    if (res.ok) {
      showSuccess('Mata Kuliah Diperbarui! ✏️', `Data mata kuliah '${editCourseForm.value.name}' telah berhasil diubah.`)
      showEditCourseModal.value = false
      window.dispatchEvent(new Event('course-changed'))
      fetchCourses()
    } else {
      showError('Gagal!', 'Gagal memperbarui data mata kuliah!')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan saat menyimpan perubahan.')
  } finally {
    isUpdating.value = false
  }
}

const deleteCourse = async (id, name) => {
  const confirmed = await showConfirm(
    'Hapus Mata Kuliah?',
    `Apakah Anda yakin ingin menghapus mata kuliah '${name}'?`,
    'Ya, Hapus'
  )
  if (!confirmed) return

  try {
    const res = await fetch(`/api/courses?id=${id}`, { method: 'DELETE' })
    if (res.ok) {
      showSuccess('Berhasil Dihapus! 🗑️', `Mata kuliah '${name}' berhasil dihapus!`)
      window.dispatchEvent(new Event('course-changed'))
      fetchCourses()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menghapus mata kuliah!')
  }
}

const resetCourses = async () => {
  const confirmed = await showConfirm(
    'Pulihkan Mata Kuliah Bawaan?',
    'Apakah Anda yakin ingin memulihkan kembali semua mata kuliah bawaan (Rekayasa Perangkat Lunak & Pemrograman Web)?',
    'Ya, Pulihkan'
  )
  if (!confirmed) return

  try {
    const res = await fetch('/api/courses/reset', { method: 'POST' })
    if (res.ok) {
      showSuccess('Berhasil Dipulihkan! 🔄', 'Mata kuliah bawaan berhasil dipulihkan!')
      window.dispatchEvent(new Event('course-changed'))
      fetchCourses()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal memulihkan mata kuliah!')
  }
}
</script>

<template>
  <div class="manage-courses-container animate-fade-in">
    <!-- Header -->
    <div class="page-header">
      <div class="title-row">
        <BookOpen class="header-icon text-gold" />
        <div>
          <h2>Kelola, Edit & Tambah Mata Kuliah</h2>
          <p class="subtitle">Tambahkan mata kuliah baru atau edit informasi mata kuliah (Kode, SKS, Jadwal, Ruang, & Deskripsi) yang diampu Dosen.</p>
        </div>
      </div>
    </div>

    <!-- Active Courses Grid -->
    <section class="section-block">
      <div class="glass-card manage-card">
        <div class="card-title-row">
          <BookOpen class="title-icon text-gold" />
          <div>
            <h3>Daftar Mata Kuliah Diampu Saat Ini</h3>
            <p class="card-sub">Klik <strong>Kelola Materi</strong> untuk mengunggah modul atau <strong>Edit</strong> untuk mengubah detail mata kuliah.</p>
          </div>

          <button @click="resetCourses" class="btn-restore" title="Kembalikan Mata Kuliah Bawaan">
            <RotateCcw class="icon-sm" />
            <span>Pulihkan Mata Kuliah Bawaan</span>
          </button>
        </div>

        <div class="active-courses-grid">
          <div v-for="c in courses" :key="c.id" class="course-manage-card">
            <!-- Header: Badges & Status -->
            <div class="crs-manage-header">
              <div class="badge-group">
                <span class="badge badge-gold">{{ c.code }}</span>
                <span class="badge badge-blue">{{ c.sks }} SKS</span>
              </div>
              <div class="status-action-box">
                <div :class="['status-pill-full', c.status === 'Non Aktif' ? 'pill-inactive' : 'pill-active']">
                  <span class="status-dot"></span>
                  <span class="status-text-label">Status: <strong>{{ c.status === 'Non Aktif' ? 'NON-AKTIF' : 'AKTIF' }}</strong></span>
                </div>
                <button 
                  @click="toggleCourseStatus(c)" 
                  class="btn-toggle-quick" 
                  :class="c.status === 'Non Aktif' ? 'btn-quick-activate' : 'btn-quick-deactivate'"
                  :title="c.status === 'Non Aktif' ? 'Klik untuk Mengaktifkan Mata Kuliah' : 'Klik untuk Menonaktifkan Mata Kuliah'"
                >
                  <span>{{ c.status === 'Non Aktif' ? '⚡ Ubah ke Aktif' : '🔒 Ubah ke Non-Aktif' }}</span>
                </button>
              </div>
            </div>

            <!-- Visibility Notice Box -->
            <div :class="['visibility-notice-box', c.status === 'Non Aktif' ? 'notice-off' : 'notice-on']">
              <span v-if="c.status === 'Non Aktif'">🔴 <strong>Status Non-Aktif:</strong> Sembunyi dari Mahasiswa (Kuis & Modul Terunci)</span>
              <span v-else>🟢 <strong>Status Aktif:</strong> Tampak di Mahasiswa (Kuis & Modul Dapat Diakses)</span>
            </div>

            <!-- Content Body -->
            <div class="crs-body">
              <h4 class="crs-title">{{ c.name }}</h4>
              <p class="crs-desc">{{ c.description }}</p>
            </div>

            <!-- Meta details -->
            <div class="crs-meta-list">
              <div class="meta-item">
                <Calendar class="meta-icon text-gold" />
                <span>{{ c.semester }}</span>
              </div>
              <div class="meta-item">
                <Clock class="meta-icon text-blue" />
                <span>{{ c.class_time }}</span>
              </div>
              <div v-if="c.room" class="meta-item">
                <MapPin class="meta-icon text-emerald" />
                <span>{{ c.room }}</span>
              </div>
            </div>

            <!-- Full-width Action Footer -->
            <div class="crs-action-footer">
              <button @click="goToCourseMateri(c.id)" class="btn-action-primary" title="Kelola Modul & Materi Pertemuan">
                <FolderOpen class="icon-sm" />
                <span>Kelola Materi</span>
              </button>
              <div class="footer-btn-secondary">
                <button @click="openEditCourseModal(c)" class="btn-action-edit" title="Edit Data Mata Kuliah">
                  <Edit3 class="icon-sm" />
                  <span>Edit</span>
                </button>
                <button @click="deleteCourse(c.id, c.name)" class="btn-action-delete" title="Hapus Mata Kuliah">
                  <Trash2 class="icon-sm" />
                  <span>Hapus</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Add New Course Form -->
    <section class="section-block">
      <div class="glass-card manage-card">
        <div class="card-title-row">
          <Plus class="title-icon text-gold" />
          <div>
            <h3>➕ Form Tambah Mata Kuliah Baru</h3>
            <p class="card-sub">Isi detail di bawah untuk menerbitkan mata kuliah baru di portal e-learning.</p>
          </div>
        </div>

        <form @submit.prevent="createCourse" class="settings-form">
          <div class="form-row">
            <div>
              <label class="input-label">Kode Mata Kuliah *</label>
              <input v-model="newCourseCode" class="glass-input" required placeholder="Contoh: TIF-303 / SIV-201" />
            </div>
            <div>
              <label class="input-label">Nama Mata Kuliah *</label>
              <input v-model="newCourseName" class="glass-input" required placeholder="Contoh: Kecerdasan Buatan / Database" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Jumlah SKS</label>
              <input v-model.number="newCourseSks" type="number" min="1" max="6" class="glass-input" required />
            </div>
            <div>
              <label class="input-label">Semester Akademik</label>
              <input v-model="newCourseSemester" class="glass-input" placeholder="Ganjil/Genap 2026/2027" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Jadwal Kuliah</label>
              <input v-model="newCourseClassTime" class="glass-input" placeholder="Senin, 08.00 - 10.30 WIB" />
            </div>
            <div>
              <label class="input-label">Ruang Kelas / Lab</label>
              <input v-model="newCourseRoom" class="glass-input" placeholder="Lab Komputer / Hybrid Zoom" />
            </div>
          </div>

          <div>
            <label class="input-label">Status Publikasi Mata Kuliah *</label>
            <select v-model="newCourseStatus" class="glass-input">
              <option value="Aktif">🟢 Aktif (Tampil di Mahasiswa & Kuis)</option>
              <option value="Non Aktif">🔴 Non-Aktif (Sembunyikan dari Mahasiswa)</option>
            </select>
          </div>

          <div>
            <label class="input-label">Deskripsi & Silabus Mata Kuliah</label>
            <textarea v-model="newCourseDescription" class="glass-input textarea" placeholder="Jelaskan fokus materi dan silabus pembelajaran..."></textarea>
          </div>

          <button type="submit" class="btn btn-gold btn-lg">
            <Plus class="btn-icon-sm" />
            <span>Tambah Mata Kuliah Sekarang</span>
          </button>
        </form>
      </div>
    </section>

    <!-- MODAL EDIT MATA KULIAH -->
    <Teleport to="body">
      <div v-if="showEditCourseModal" class="modal-overlay" @click.self="showEditCourseModal = false">
        <div class="glass-card modal-card">
          <div class="modal-header">
            <div class="m-title-box">
              <Edit3 class="modal-title-icon text-gold" />
              <div>
                <h3>Edit Data Mata Kuliah</h3>
                <p class="modal-sub">Ubah Kode, Nama, SKS, Status, Jadwal, Ruangan, dan Deskripsi</p>
              </div>
            </div>
            <button @click="showEditCourseModal = false" class="btn-close-modal">✕</button>
          </div>

          <form @submit.prevent="updateCourse" class="modal-form">
            <div class="form-row">
              <div>
                <label class="input-label">Kode Mata Kuliah *</label>
                <input v-model="editCourseForm.code" class="glass-input" required placeholder="Contoh: TIF-301" />
              </div>
              <div>
                <label class="input-label">Jumlah SKS *</label>
                <input v-model.number="editCourseForm.sks" type="number" min="1" max="6" class="glass-input" required />
              </div>
            </div>

            <div class="form-row">
              <div>
                <label class="input-label">Nama Mata Kuliah *</label>
                <input v-model="editCourseForm.name" class="glass-input" required placeholder="Contoh: Rekayasa Perangkat Lunak" />
              </div>
              <div>
                <label class="input-label">Status Publikasi *</label>
                <select v-model="editCourseForm.status" class="glass-input">
                  <option value="Aktif">🟢 Aktif (Tampil di Mahasiswa & Kuis)</option>
                  <option value="Non Aktif">🔴 Non-Aktif (Sembunyikan dari Mahasiswa)</option>
                </select>
              </div>
            </div>

            <div class="form-row">
              <div>
                <label class="input-label">Semester Akademik</label>
                <input v-model="editCourseForm.semester" class="glass-input" placeholder="Ganjil 2026/2027" />
              </div>
              <div>
                <label class="input-label">Jadwal Kuliah</label>
                <input v-model="editCourseForm.class_time" class="glass-input" placeholder="Senin, 08.00 - 10.30 WIB" />
              </div>
            </div>

            <div>
              <label class="input-label">Ruang Kelas / Lab</label>
              <input v-model="editCourseForm.room" class="glass-input" placeholder="Lab Komputer 2 / Hybrid Zoom" />
            </div>

            <div>
              <label class="input-label">Deskripsi Mata Kuliah</label>
              <textarea v-model="editCourseForm.description" class="glass-input textarea" placeholder="Jelaskan cakupan materi perkuliahan..."></textarea>
            </div>

            <div class="modal-actions">
              <button type="button" @click="showEditCourseModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" :disabled="isUpdating" class="btn btn-gold">
                <Save class="btn-icon-sm" />
                {{ isUpdating ? 'Menyimpan...' : 'Simpan Perubahan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.manage-courses-container {
  max-width: 1200px;
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

.section-block {
  margin-bottom: 2rem;
}

.manage-card {
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.card-title-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
  position: relative;
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

.btn-restore {
  position: absolute;
  right: 0;
  top: 0;
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  padding: 0.4rem 0.8rem;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
}

.btn-restore:hover {
  background: #e2e8f0;
}

.active-courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.course-manage-card {
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  transition: all 0.25s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.02);
}

.course-manage-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.06);
  transform: translateY(-2px);
}

.crs-manage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f1f5f9;
}

.badge-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: nowrap;
}

.status-action-box {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  flex-wrap: wrap;
}

.status-pill-full {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.75rem;
  padding: 0.25rem 0.7rem;
  border-radius: 99px;
  white-space: nowrap;
}

.status-pill-full.pill-active {
  color: #15803d;
  background: #dcfce7;
  border: 1.5px solid #86efac;
}

.status-pill-full.pill-inactive {
  color: #b91c1c;
  background: #fee2e2;
  border: 1.5px solid #fca5a5;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.pill-active .status-dot {
  background: #22c55e;
  box-shadow: 0 0 6px #22c55e;
}

.pill-inactive .status-dot {
  background: #ef4444;
  box-shadow: 0 0 6px #ef4444;
}

.visibility-notice-box {
  margin: 0.5rem 0 1rem 0;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  font-size: 0.78rem;
  line-height: 1.4;
}

.visibility-notice-box.notice-on {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #166534;
}

.visibility-notice-box.notice-off {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

.btn-toggle-quick {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.3rem 0.65rem;
  border-radius: 7px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

.btn-quick-activate {
  background: #16a34a;
  color: white;
}

.btn-quick-activate:hover {
  background: #15803d;
}

.btn-quick-deactivate {
  background: #ffffff;
  color: #475569;
  border: 1px solid #cbd5e1;
}

.btn-quick-deactivate:hover {
  background: #fee2e2;
  color: #dc2626;
  border-color: #fca5a5;
}

.crs-body {
  margin-bottom: 1rem;
  flex: 1;
}

.crs-title {
  font-size: 1.2rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.4rem;
  line-height: 1.35;
}

.crs-desc {
  font-size: 0.88rem;
  color: #64748b;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.crs-meta-list {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  border-radius: var(--radius-sm);
  padding: 0.75rem 0.9rem;
  margin-bottom: 1.25rem;
  font-size: 0.82rem;
  color: #475569;
  font-weight: 500;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.55rem;
}

.meta-icon {
  width: 15px;
  height: 15px;
  flex-shrink: 0;
}

.crs-action-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  padding-top: 1rem;
  border-top: 1.5px solid #f1f5f9;
}

.btn-action-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  background: #059669;
  color: #ffffff;
  border: none;
  padding: 0.5rem 0.9rem;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
}

.btn-action-primary:hover {
  background: #047857;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.35);
}

.footer-btn-secondary {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.btn-action-edit {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  background: #eff6ff;
  border: 1.5px solid #bfdbfe;
  color: #1d4ed8;
  padding: 0.45rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action-edit:hover {
  background: #dbeafe;
  border-color: #93c5fd;
  color: #1e40af;
}

.btn-action-delete {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  background: #fff1f2;
  border: 1.5px solid #fecdd3;
  color: #be123c;
  padding: 0.45rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action-delete:hover {
  background: #ffe4e6;
  border-color: #fda4af;
  color: #9f1239;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  font-weight: 700;
  color: #334155;
  margin-bottom: 0.4rem;
}

.textarea {
  min-height: 90px;
}

.btn-lg {
  padding: 0.75rem 1.5rem;
  font-size: 0.95rem;
  align-self: flex-start;
}

.icon-sm {
  width: 14px;
  height: 14px;
}

.btn-icon-sm {
  width: 16px;
  height: 16px;
}

.text-gold {
  color: #b45309;
}

/* Modal Styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(6px);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-card {
  width: 100%;
  max-width: 580px;
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.m-title-box {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.modal-title-icon {
  width: 24px;
  height: 24px;
  margin-top: 0.2rem;
}

.btn-close-modal {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #64748b;
  cursor: pointer;
  padding: 0.2rem 0.5rem;
}

.btn-close-modal:hover {
  color: #0f172a;
}

.modal-sub {
  font-size: 0.85rem;
  color: #64748b;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
