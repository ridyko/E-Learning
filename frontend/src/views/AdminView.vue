<script setup>
import { ref, onMounted } from 'vue'
import { showSuccess, showError, showWarning, showConfirm } from '../utils/swal.js'
import { 
  CheckSquare, 
  Zap, 
  ShieldCheck, 
  Layout, 
  CheckCircle2, 
  Save, 
  BookOpen, 
  Layers, 
  Megaphone, 
  ClipboardList, 
  HelpCircle, 
  Users, 
  Sparkles, 
  FolderOpen, 
  Calendar, 
  Clock, 
  MapPin, 
  Edit3,
  ArrowRight 
} from 'lucide-vue-next'

const announcements = ref([])
const assignments = ref([])
const attendances = ref([])
const registeredStudents = ref([])
const courses = ref([])
const quizzes = ref([])
const quizSubmissions = ref([])

const selectedQuizId = ref('quiz-rpl-1')
const newQQuestion = ref('')
const newQOpt0 = ref('')
const newQOpt1 = ref('')
const newQOpt2 = ref('')
const newQOpt3 = ref('')
const newQCorrect = ref(0)
const newQExplanation = ref('')

const newAnnTitle = ref('')
const newAnnCategory = ref('Penting')
const newAnnCourse = ref('all')
const newAnnContent = ref('')

// Dynamic Course Form Fields
const newCourseCode = ref('')
const newCourseName = ref('')
const newCourseSks = ref(3)
const newCourseSemester = ref('Ganjil 2026/2027')
const newCourseClassTime = ref('Senin, 08.00 - 10.30 WIB')
const newCourseRoom = ref('Lab Komputer / Hybrid Zoom')
const newCourseDescription = ref('')

// Beranda Landing Page Settings
const landingSettings = ref({
  hero_title: '',
  hero_subtitle: '',
  institution_badge: '',
  semester_badge: '',
  dosen_bio: '',
  office_hours: '',
  room: '',
  phone: ''
})
const saveSuccess = ref(false)

const fetchData = async () => {
  try {
    const [resAnn, resAsg, resAtt, resStd, resSet, resCrs, resQuizzes, resSubs] = await Promise.all([
      fetch('/api/announcements'),
      fetch('/api/assignments'),
      fetch('/api/attendance'),
      fetch('/api/students'),
      fetch('/api/settings/landing'),
      fetch('/api/courses'),
      fetch('/api/quizzes'),
      fetch('/api/quizzes/submissions')
    ])
    announcements.value = await resAnn.json()
    assignments.value = await resAsg.json()
    attendances.value = await resAtt.json()
    registeredStudents.value = await resStd.json()
    landingSettings.value = await resSet.json()
    courses.value = await resCrs.json()
    quizzes.value = await resQuizzes.json()
    quizSubmissions.value = await resSubs.json()
  } catch (err) {
    console.warn('Admin fetch error:', err)
  }
}

const addQuizQuestion = async () => {
  if (!newQQuestion.value || !newQOpt0.value || !newQOpt1.value || !newQOpt2.value || !newQOpt3.value) {
    showWarning('Perhatian', 'Harap lengkapi Teks Pertanyaan dan ke-4 opsi jawaban (A, B, C, D)!')
    return
  }

  try {
    const res = await fetch('/api/quizzes/questions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        quiz_id: selectedQuizId.value,
        question: {
          question: newQQuestion.value,
          options: [newQOpt0.value, newQOpt1.value, newQOpt2.value, newQOpt3.value],
          correct_answer: parseInt(newQCorrect.value),
          explanation: newQExplanation.value || 'Pembahasan kunci jawaban oleh Dosen Pengampu.'
        }
      })
    })

    if (res.ok) {
      showSuccess('Soal Ditambahkan! 📝', 'Soal kuis baru berhasil ditambahkan ke Bank Soal.')
      newQQuestion.value = ''
      newQOpt0.value = ''
      newQOpt1.value = ''
      newQOpt2.value = ''
      newQOpt3.value = ''
      newQExplanation.value = ''
      fetchData()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menambahkan soal kuis baru!')
  }
}

onMounted(() => {
  fetchData()
})

const saveLandingSettings = async () => {
  try {
    const res = await fetch('/api/settings/landing', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(landingSettings.value)
    })

    if (res.ok) {
      saveSuccess.value = true
      showSuccess('Tampilan Beranda Disimpan! 🎨', 'Perubahan beranda publik telah diperbarui secara realtime.')
      setTimeout(() => saveSuccess.value = false, 3500)
    }
  } catch (err) {
    showError('Gagal Menyimpan!', 'Gagal menyimpan pengaturan beranda!')
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
        description: newCourseDescription.value
      })
    })

    if (res.ok) {
      showSuccess('Mata Kuliah Ditambahkan! 📚', `Mata Kuliah '${newCourseName.value}' berhasil ditambahkan!`)
      newCourseCode.value = ''
      newCourseName.value = ''
      newCourseDescription.value = ''
      window.dispatchEvent(new Event('course-changed'))
      fetchData()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menambahkan mata kuliah baru!')
  }
}

const deleteCourse = async (id, name) => {
  const confirmed = await showConfirm(
    'Hapus Mata Kuliah?',
    `Apakah Anda yakin ingin menghapus/mengarsipkan mata kuliah '${name}'?`,
    'Ya, Hapus'
  )
  if (!confirmed) return

  try {
    const res = await fetch(`/api/courses?id=${id}`, { method: 'DELETE' })
    if (res.ok) {
      showSuccess('Berhasil Dihapus! 🗑️', `Mata kuliah '${name}' berhasil dihapus!`)
      window.dispatchEvent(new Event('course-changed'))
      fetchData()
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
      fetchData()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal memulihkan mata kuliah!')
  }
}

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
      showSuccess('Pengumuman Diterbitkan! 📢', 'Pengumuman baru berhasil diterbitkan!')
      newAnnTitle.value = ''
      newAnnContent.value = ''
      fetchData()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menerbitkan pengumuman!')
  }
}

const deleteAnnouncement = async (id) => {
  const confirmed = await showConfirm('Hapus Pengumuman?', 'Apakah Anda yakin ingin menghapus pengumuman ini?', 'Ya, Hapus')
  if (!confirmed) return

  try {
    const res = await fetch(`/api/announcements?id=${id}`, { method: 'DELETE' })
    if (res.ok) {
      showSuccess('Terhapus!', 'Pengumuman berhasil dihapus.')
      fetchData()
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menghapus pengumuman!')
  }
}
const triggerStartServer = async () => {
  try {
    const res = await fetch('/start_server.php')
    const data = await res.json()
    if (data.status === 'success') {
      showSuccess('Server Diaktifkan! 🚀', data.message)
    } else {
      showSuccess('Status Server 🟢', data.message)
    }
  } catch (err) {
    showError('Gagal!', 'Gagal menghubungi server trigger.')
  }
}
</script>

<template>
  <div class="admin-container animate-fade-in">
    <!-- Admin Header -->
    <div class="page-header">
      <div class="header-flex">
        <div class="title-row">
          <ShieldCheck class="header-icon text-gold" />
          <div>
            <h2>Panel Kontrol Dashboard Dosen</h2>
            <p class="subtitle">Kelola tampilan Beranda, pengumuman, daftar mahasiswa, rekapitulasi tugas, dan presensi ITB Swadharma.</p>
          </div>
        </div>
        <button @click="triggerStartServer" class="btn btn-emerald btn-sm btn-server-trigger" title="1-Klik Aktifkan Server Backend Golang">
          <Zap class="btn-icon-xs" />
          <span>🚀 Aktifkan / Cek Server Backend</span>
        </button>
      </div>
    </div>

    <!-- PUSAT NAVIGASI & KENDALI UTAMA DOSEN (QUICK ACTION HUB) -->
    <section class="section-block">
      <div class="hub-header-title">
        <h3>⚡ Menu Utama Manajemen Perkuliahan</h3>
        <p>Akses cepat seluruh modul pengelolaan untuk Dosen Pengampu.</p>
      </div>

      <div class="quick-hub-grid">
        <!-- 1. Kelola Matakuliah -->
        <router-link to="/mata-kuliah" class="glass-card hub-card hub-gold">
          <div class="hub-icon-box gold-box">
            <Layers class="hub-icon text-gold" />
          </div>
          <div class="hub-content">
            <h4>Kelola Mata Kuliah & Jadwal</h4>
            <p>Atur Kode, SKS (4 SKS), Jadwal Perkuliahan, Ruang Kelas, dan Deskripsi Silabus.</p>
            <span class="hub-action-link">Buka Pengaturan Matkul →</span>
          </div>
        </router-link>

        <!-- 2. Kelola Modul / Materi 1-14 -->
        <router-link :to="courses.length > 0 ? ('/course/' + courses[0].id) : '/mata-kuliah'" class="glass-card hub-card hub-blue">
          <div class="hub-icon-box blue-box">
            <BookOpen class="hub-icon text-blue" />
          </div>
          <div class="hub-content">
            <h4>Upload & Kelola Modul (P1 - P14)</h4>
            <p>Buka/kunci modul mahasiswa, upload file slide PPTX/PDF, video materi, dan deadline tugas.</p>
            <span class="hub-action-link">Akses Modul Perkuliahan →</span>
          </div>
        </router-link>

        <!-- 3. Terbitkan Pengumuman -->
        <router-link to="/pengumuman" class="glass-card hub-card hub-amber">
          <div class="hub-icon-box amber-box">
            <Megaphone class="hub-icon text-amber" />
          </div>
          <div class="hub-content">
            <h4>Terbitkan & Kelola Pengumuman</h4>
            <p>Buat pengumuman penting, info tugas, info kuis, atau hapus pengumuman lama.</p>
            <span class="hub-action-link">Kelola Pengumuman →</span>
          </div>
        </router-link>

        <!-- 4. Presensi Perkuliahan -->
        <router-link to="/presensi" class="glass-card hub-card hub-emerald">
          <div class="hub-icon-box emerald-box">
            <ClipboardList class="hub-icon text-emerald" />
          </div>
          <div class="hub-content">
            <h4>Presensi & Absensi Digital</h4>
            <p>Tampilkan QR code presensi ke layar kelas, rekap hadir, sakit, atau izin mahasiswa.</p>
            <span class="hub-action-link">Buka Sesi Presensi →</span>
          </div>
        </router-link>

        <!-- 5. Bank Soal & Kuis -->
        <router-link to="/bank-soal" class="glass-card hub-card hub-purple">
          <div class="hub-icon-box purple-box">
            <HelpCircle class="hub-icon text-purple" />
          </div>
          <div class="hub-content">
            <h4>Bank Soal & Rekap Nilai Kuis</h4>
            <p>Kelola bank soal ujian, atur live quiz proyektor, dan ekspor nilai mahasiswa.</p>
            <span class="hub-action-link">Kelola Bank Soal →</span>
          </div>
        </router-link>

        <!-- 6. Data Mahasiswa -->
        <router-link to="/mahasiswa" class="glass-card hub-card hub-indigo">
          <div class="hub-icon-box indigo-box">
            <Users class="hub-icon text-indigo" />
          </div>
          <div class="hub-content">
            <h4>Daftar Mahasiswa Terdaftar</h4>
            <p>Pantau mahasiswa yang sudah mendaftar NIM mandiri dan status akunnya.</p>
            <span class="hub-action-link">Lihat Semua Mahasiswa →</span>
          </div>
        </router-link>
      </div>
    </section>

    <!-- RINGKASAN MATA KULIAH AKTIF -->
    <section class="section-block">
      <div class="glass-card admin-card">
        <div class="card-title-row-between">
          <div class="card-title-row">
            <BookOpen class="title-icon text-gold" />
            <div>
              <h3>📚 Status Mata Kuliah Semester Ini</h3>
              <p class="card-sub">Mata kuliah yang aktif dapat diakses mahasiswa untuk mengunduh modul dan kuis.</p>
            </div>
          </div>
          <router-link to="/mata-kuliah" class="btn btn-secondary btn-sm">
            <Edit3 class="btn-icon-xs" />
            <span>Kelola Lengkap di Menu Matkul</span>
          </router-link>
        </div>

        <div class="courses-summary-grid">
          <div v-for="c in courses" :key="c.id" class="course-summary-card">
            <div class="summary-top">
              <span class="badge badge-gold">{{ c.code }}</span>
              <span class="badge badge-blue">{{ c.sks }} SKS</span>
              <span :class="c.status === 'Non Aktif' ? 'badge badge-rose' : 'badge badge-emerald'">
                {{ c.status || 'Aktif' }}
              </span>
            </div>
            <h4 class="summary-title">{{ c.name }}</h4>
            <p class="summary-desc">{{ c.description }}</p>
            <div class="summary-meta">
              <div class="meta-line">
                <Clock class="meta-icon text-blue" />
                <span>{{ c.class_time }}</span>
              </div>
              <div class="meta-line">
                <MapPin class="meta-icon text-emerald" />
                <span>{{ c.room }}</span>
              </div>
            </div>
            <div class="summary-actions">
              <router-link :to="'/course/' + c.id" class="btn btn-primary btn-sm">
                <FolderOpen class="btn-icon-xs" />
                <span>Buka Modul (1 - 14)</span>
              </router-link>
              <router-link to="/mata-kuliah" class="btn btn-secondary btn-sm">
                <Edit3 class="btn-icon-xs" />
                <span>Edit Detail</span>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Edit Beranda Settings Card -->
    <section class="section-block">
      <div class="glass-card admin-card edit-beranda-card">
        <div class="card-title-row">
          <Layout class="title-icon text-gold" />
          <div>
            <h3>🎨 Edit & Kelola Tampilan Beranda</h3>
            <p class="card-sub">Ubah judul, kata sambutan, bio dosen, dan kontak yang tampil di halaman beranda publik.</p>
          </div>
        </div>

        <div v-if="saveSuccess" class="alert-success">
          <CheckCircle2 class="icon-sm" />
          <span>Tampilan Beranda berhasil diperbarui secara realtime!</span>
        </div>

        <form @submit.prevent="saveLandingSettings" class="settings-form">
          <div class="form-row">
            <div>
              <label class="input-label">Judul Utama Hero Banner</label>
              <input v-model="landingSettings.hero_title" class="glass-input" required placeholder="Judul beranda..." />
            </div>
            <div>
              <label class="input-label">Badge Kampus / Institusi</label>
              <input v-model="landingSettings.institution_badge" class="glass-input" required placeholder="Contoh: ITB SWADHARMA JAKARTA" />
            </div>
          </div>

          <div>
            <label class="input-label">Sub Judul / Tagline Beranda</label>
            <textarea v-model="landingSettings.hero_subtitle" class="glass-input textarea" required placeholder="Tuliskan kata sambutan beranda..."></textarea>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Jam Konsultasi Dosen</label>
              <input v-model="landingSettings.office_hours" class="glass-input" placeholder="Senin & Rabu: 13.00 - 16.00 WIB" />
            </div>
            <div>
              <label class="input-label">Lokasi Ruang Dosen</label>
              <input v-model="landingSettings.room" class="glass-input" placeholder="Ruang Dosen Lantai 3" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Nomor WhatsApp Direct Dosen</label>
              <input v-model="landingSettings.phone" class="glass-input" placeholder="+62 812-9876-5432" />
            </div>
            <div>
              <label class="input-label">Badge Semester Akademik</label>
              <input v-model="landingSettings.semester_badge" class="glass-input" placeholder="SEMESTER GANJIL 2026/2027" />
            </div>
          </div>

          <div>
            <label class="input-label">Bio & Profil Singkat Dosen</label>
            <textarea v-model="landingSettings.dosen_bio" class="glass-input textarea" placeholder="Deskripsi profil dosen..."></textarea>
          </div>

          <button type="submit" class="btn btn-gold btn-lg">
            <Save class="btn-icon-sm" />
            <span>Simpan Perubahan Beranda Now</span>
          </button>
        </form>
      </div>
    </section>

    <!-- 2 Column Admin Grid for Overview Stats -->
    <div class="admin-grid">
      <!-- Left Column: Rekapitulasi Tugas & Presensi -->
      <div class="left-column">
        <!-- Rekap Pengumpulan Tugas Mahasiswa -->
        <div class="glass-card admin-card">
          <div class="card-title-row">
            <BookOpen class="title-icon text-gold" />
            <div>
              <h3>📁 Rekapitulasi Tugas Masuk</h3>
              <p class="card-sub">Daftar pengumpulan tugas dari mahasiswa seluruh mata kuliah.</p>
            </div>
          </div>

          <div v-if="assignments.length === 0" class="empty-notice">
            <p>Belum ada pengumpulan tugas dari mahasiswa.</p>
          </div>
          <div v-else class="asg-list">
            <div 
              v-for="asg in assignments" 
              :key="asg.id" 
              class="asg-item"
            >
              <div class="asg-info">
                <span class="asg-student">{{ asg.student_name }} (NIM: {{ asg.student_nim }})</span>
                <span class="asg-title">P{{ asg.meeting_no }}: {{ asg.title }}</span>
                <a :href="asg.repo_link" target="_blank" class="asg-link">🔗 {{ asg.repo_link }}</a>
                <p v-if="asg.notes" class="asg-notes">Catatan: {{ asg.notes }}</p>
              </div>
              <span class="badge badge-emerald">{{ asg.status }}</span>
            </div>
          </div>
        </div>

        <!-- Rekap Kehadiran -->
        <div class="glass-card admin-card margin-top">
          <div class="card-title-row">
            <ShieldCheck class="title-icon text-gold" />
            <div>
              <h3>📊 Rekap Presensi Terkini</h3>
              <p class="card-sub">Log kehadiran mahasiswa secara realtime.</p>
            </div>
          </div>
          <div class="att-list">
            <div v-for="att in attendances" :key="att.id" class="att-item">
              <span class="att-name">{{ att.student_name }} ({{ att.student_nim }})</span>
              <span class="badge badge-blue">{{ att.status }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Registered Students Table -->
      <div class="right-column">
        <div class="glass-card admin-card">
          <div class="card-title-row">
            <Users class="title-icon text-gold" />
            <div>
              <h3>👥 Daftar Mahasiswa Terdaftar</h3>
              <p class="card-sub">Total: {{ registeredStudents.length }} Mahasiswa terdaftar</p>
            </div>
          </div>

          <div class="std-list">
            <div v-for="std in registeredStudents" :key="std.id" class="std-item">
              <div class="std-info">
                <span class="std-name">{{ std.name }}</span>
                <span class="std-nim">NIM: {{ std.username }} | {{ std.prodi || 'Teknik Informatika' }}</span>
                <span class="std-email">📧 {{ std.email }}</span>
              </div>
              <span class="badge badge-blue">MAHASISWA</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin-container {
  max-width: 1350px;
  margin: 0 auto;
  padding: 1.5rem;
}

.hub-header-title {
  margin-bottom: 1.25rem;
}

.hub-header-title h3 {
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.25rem;
}

.hub-header-title p {
  color: #64748b;
  font-size: 0.88rem;
  margin: 0;
}

.quick-hub-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.25rem;
  margin-bottom: 1rem;
}

.hub-card {
  display: flex;
  gap: 1.2rem;
  padding: 1.5rem;
  text-decoration: none;
  border-radius: var(--radius-md);
  border: 1.5px solid #e2e8f0;
  background: #ffffff;
  transition: all 0.25s ease;
  position: relative;
  overflow: hidden;
}

.hub-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 24px -6px rgba(15, 23, 42, 0.1);
  border-color: #cbd5e1;
}

.hub-icon-box {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.hub-icon {
  width: 26px;
  height: 26px;
}

.gold-box { background: #fef3c7; }
.blue-box { background: #dbeafe; }
.amber-box { background: #fef3c7; }
.emerald-box { background: #d1fae5; }
.purple-box { background: #ede9fe; }
.indigo-box { background: #e0e7ff; }

.text-gold { color: #d97706; }
.text-blue { color: #2563eb; }
.text-amber { color: #b45309; }
.text-emerald { color: #059669; }
.text-purple { color: #7c3aed; }
.text-indigo { color: #4f46e5; }

.hub-content h4 {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.4rem;
}

.hub-content p {
  font-size: 0.84rem;
  color: #475569;
  line-height: 1.45;
  margin-bottom: 0.75rem;
}

.hub-action-link {
  font-size: 0.82rem;
  font-weight: 700;
  color: #2563eb;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.hub-card:hover .hub-action-link {
  color: #1d4ed8;
  text-decoration: underline;
}

.card-title-row-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 1.25rem;
}

.courses-summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 1.25rem;
}

.course-summary-card {
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-sm);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.summary-top {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.summary-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.summary-desc {
  font-size: 0.84rem;
  color: #475569;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.summary-meta {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  font-size: 0.82rem;
  color: #64748b;
  font-weight: 600;
  padding: 0.5rem 0;
  border-top: 1px dashed #e2e8f0;
  border-bottom: 1px dashed #e2e8f0;
}

.meta-line {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.meta-icon {
  width: 15px;
  height: 15px;
}

.summary-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
}

.page-header {
  margin-bottom: 2rem;
}

.header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.btn-server-trigger {
  padding: 0.65rem 1.25rem;
  font-weight: 800;
  border-radius: var(--radius-sm);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2);
}

.title-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  width: 38px;
  height: 38px;
}

.title-row h2 {
  font-size: 1.8rem;
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

.edit-beranda-card {
  border-top: 4px solid #d97706;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  margin-top: 1rem;
}

.alert-success {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #047857;
  padding: 0.85rem 1rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  margin-bottom: 1rem;
}

.admin-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.admin-card {
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  margin-bottom: 0.5rem;
}

.title-icon {
  width: 24px;
  height: 24px;
}

.card-sub {
  color: #475569;
  font-size: 0.88rem;
  margin-bottom: 1.25rem;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  color: #334155;
  font-weight: 600;
  margin-bottom: 0.35rem;
}

.textarea {
  min-height: 80px;
  resize: vertical;
}

.w-full {
  width: 100%;
}

.ann-manage-section h4 {
  font-size: 1rem;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 0.4rem;
}

.ann-manage-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.ann-manage-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 0.9rem 1rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.ann-manage-item h5 {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
  margin-top: 0.3rem;
}

.ann-manage-item p {
  font-size: 0.85rem;
  color: #475569;
}

.btn-del {
  background: #fff1f2;
  border: 1px solid #fecdd3;
  color: #e11d48;
  padding: 0.35rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: var(--transition);
}

.btn-del:hover {
  background: #ffe4e6;
}

.icon-sm {
  width: 16px;
  height: 16px;
}

.margin-top {
  margin-top: 1.5rem;
}

.std-list, .asg-list, .att-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.std-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1.5px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.std-info {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.std-name {
  font-weight: 800;
  font-size: 0.95rem;
  color: #0f172a;
}

.std-nim {
  font-size: 0.82rem;
  color: #1d4ed8;
  font-weight: 700;
}

.std-email {
  font-size: 0.8rem;
  color: #475569;
}

.asg-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1.5px solid #e2e8f0;
}

.asg-info {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.asg-student {
  font-weight: 800;
  font-size: 0.95rem;
  color: #0f172a;
}

.asg-title {
  font-size: 0.85rem;
  color: #b45309;
  font-weight: 700;
}

.asg-link {
  font-size: 0.82rem;
  color: #1d4ed8;
  font-weight: 600;
  word-break: break-all;
  text-decoration: none;
}

.asg-link:hover {
  text-decoration: underline;
}

.asg-notes {
  font-size: 0.8rem;
  color: #64748b;
}

.att-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
  font-size: 0.9rem;
}

.att-name {
  font-weight: 700;
  color: #0f172a;
}

.text-gold {
  color: #d97706;
}

.btn-icon-sm {
  width: 18px;
  height: 18px;
}

.sub-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e2e8f0;
}

.sub-title-row .form-sub-title {
  margin: 0;
  padding: 0;
  border: none;
}

.btn-restore {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.4rem 0.85rem;
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-restore:hover {
  background: #1d4ed8;
  color: #ffffff;
}

.form-sub-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #e2e8f0;
}

.active-courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.25rem;
  margin-bottom: 2rem;
}

.course-manage-card {
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-sm);
  padding: 1.25rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.course-manage-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
  border-color: #cbd5e1;
}

.crs-manage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.badge-group {
  display: flex;
  gap: 0.5rem;
}

.badge-blue {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.btn-del-course {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.35rem 0.65rem;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-del-course:hover {
  background: #dc2626;
  color: #ffffff;
}

.crs-title {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.crs-desc {
  font-size: 0.85rem;
  color: #475569;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.crs-meta {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
  padding-top: 0.5rem;
  border-top: 1px dashed #e2e8f0;
}

.add-course-box {
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  padding: 1.5rem;
}

.quiz-opsi-list {
  font-size: 0.82rem;
  color: #334155;
  margin: 0.4rem 0;
  font-weight: 600;
}

.quiz-key {
  font-size: 0.8rem;
  color: #b45309;
  margin: 0.2rem 0 0 0;
}

.sub-date {
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.score-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.3rem;
}

.score-val {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.badge-rose {
  background: #fff1f2;
  color: #e11d48;
  border: 1px solid #fecdd3;
}

.empty-rekap {
  padding: 2rem;
  text-align: center;
  color: #64748b;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px dashed #cbd5e1;
}

@media (max-width: 900px) {
  .admin-grid {
    grid-template-columns: 1fr;
  }
}
</style>
