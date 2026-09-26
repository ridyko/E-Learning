<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DosenProfile from '../components/DosenProfile.vue'
import { 
  BookOpen, 
  Code, 
  Megaphone, 
  Calendar, 
  Clock, 
  ArrowRight, 
  CheckCircle,
  Users,
  LogIn,
  UserPlus,
  Sparkles,
  ShieldCheck,
  CheckSquare,
  Zap,
  Terminal
} from 'lucide-vue-next'

const router = useRouter()
const profile = ref(null)
const courses = ref([])
const announcements = ref([])
const landingSettings = ref({
  hero_title: 'Portal Pembelajaran Digital Terpadu',
  hero_subtitle: 'Selamat datang di Portal E-Learning Resmi Dosen Rio Widyatmoko, S.Kom, M.M.S.I. Fasilitas pembelajaran interaktif khusus mata kuliah Rekayasa Perangkat Lunak dan Pemrograman Web di ITB Swadharma.',
  institution_badge: 'ITB SWADHARMA JAKARTA',
  semester_badge: 'SEMESTER GANJIL 2026/2027'
})
const loading = ref(true)
const currentUser = ref(null)

const loadUserSession = () => {
  const userJson = localStorage.getItem('user')
  if (userJson) {
    try {
      currentUser.value = JSON.parse(userJson)
    } catch (e) {
      currentUser.value = null
    }
  } else {
    currentUser.value = null
  }
}

const fetchDashboardData = async () => {
  try {
    const [resProf, resCourse, resAnn, resSet] = await Promise.all([
      fetch('/api/profile'),
      fetch('/api/courses'),
      fetch('/api/announcements'),
      fetch('/api/settings/landing')
    ])

    profile.value = await resProf.json()
    courses.value = await resCourse.json()
    announcements.value = await resAnn.json()
    const set = await resSet.json()
    if (set && set.hero_title) {
      landingSettings.value = set
    }

    if (currentUser.value && currentUser.value.role === 'dosen' && profile.value) {
      const degreeStr = (profile.value.degree && profile.value.degree.trim() !== '' && profile.value.degree.trim() !== '-')
        ? ', ' + profile.value.degree.trim()
        : ''
      const fullName = (profile.value.name || '').trim() + degreeStr
      if (fullName && currentUser.value.name !== fullName) {
        currentUser.value.name = fullName
        localStorage.setItem('user', JSON.stringify(currentUser.value))
        window.dispatchEvent(new Event('auth-changed'))
      }
    }
  } catch (err) {
    console.warn('Backend fetch error:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUserSession()
  fetchDashboardData()
  window.addEventListener('auth-changed', loadUserSession)
})
</script>

<template>
  <div class="home-container animate-fade-in">
    <!-- ========================================================================= -->
    <!-- VIEW A: TAMPILAN BERANDA UTAMA SEBELUM LOGIN (PUBLIC LANDING PAGE)       -->
    <!-- ========================================================================= -->
    <div v-if="!currentUser" class="landing-public-wrapper">
      
      <!-- HERO BANNER (CUSTOMIZABLE VIA DOSEN PANEL) -->
      <section class="hero-section glass-card">
        <div class="hero-content">
          <div class="badge-row">
            <span class="badge badge-gold">{{ landingSettings.institution_badge || 'ITB SWADHARMA JAKARTA' }}</span>
            <span class="badge badge-blue">{{ landingSettings.semester_badge || 'SEMESTER GANJIL 2026/2027' }}</span>
          </div>

          <h1 class="hero-title">
            {{ landingSettings.hero_title || 'Portal Pembelajaran Digital Terpadu' }}
          </h1>

          <p class="hero-subtitle">
            {{ landingSettings.hero_subtitle }}
          </p>

          <!-- Quick Action CTAs -->
          <div class="hero-actions">
            <router-link to="/login" class="btn btn-gold btn-lg shadow-gold">
              <LogIn class="btn-icon-sm" />
              <span>Masuk Portal (Login)</span>
            </router-link>

            <router-link to="/register" class="btn btn-secondary btn-lg">
              <UserPlus class="btn-icon-sm" />
              <span>Pendaftaran Mahasiswa Baru</span>
            </router-link>
          </div>

          <!-- Highlight Quick Stats -->
          <div class="hero-stats-grid">
            <div class="stat-card">
              <span class="stat-num">2</span>
              <span class="stat-label">Mata Kuliah Utama</span>
            </div>
            <div class="stat-card">
              <span class="stat-num">28</span>
              <span class="stat-label">Modul Pertemuan</span>
            </div>
            <div class="stat-card">
              <span class="stat-num">100%</span>
              <span class="stat-label">Live Code Sandbox</span>
            </div>
          </div>
        </div>
      </section>

      <!-- PORTAL CHOICE CARDS (Akses Cepat Dosen vs Mahasiswa) -->
      <section class="section-block">
        <div class="section-header text-center">
          <span class="badge badge-blue">AKSES PORTAL</span>
          <h2>Pilih Akses Pengguna Anda</h2>
          <p class="section-subtitle">Gunakan NIP Dosen atau NIM Mahasiswa untuk masuk ke sistem.</p>
        </div>

        <div class="portal-choice-grid">
          <!-- Card 1: Portal Dosen -->
          <div class="glass-card choice-card dosen-choice">
            <div class="choice-icon-box gold-box">
              <ShieldCheck class="choice-icon text-gold" />
            </div>
            <h3>Portal Dosen Pengampu</h3>
            <p>Khusus Dosen Pengampu. Kelola tampilan beranda, pengumuman, modul pertemuan, rekapitulasi tugas masuk, dan absensi mahasiswa.</p>

            <router-link to="/login" class="btn btn-gold w-full mt-auto">
              <LogIn class="btn-icon-xs" />
              <span>Login Dosen</span>
            </router-link>
          </div>

          <!-- Card 2: Portal Mahasiswa -->
          <div class="glass-card choice-card student-choice">
            <div class="choice-icon-box blue-box">
              <Users class="choice-icon text-blue" />
            </div>
            <h3>Portal Mahasiswa ITB Swadharma</h3>
            <p>Khusus Mahasiswa terdaftar. Akses silabus, download slide, kuis online, upload tugas, dan jalankan Live Code Editor.</p>

            <div class="button-group-row mt-auto">
              <router-link to="/login" class="btn btn-primary w-full">
                <LogIn class="btn-icon-xs" />
                <span>Login Mahasiswa</span>
              </router-link>
              <router-link to="/register" class="btn btn-secondary w-full">
                <UserPlus class="btn-icon-xs" />
                <span>Daftar NIM Mandiri</span>
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Dosen Profile Overview Section -->
      <DosenProfile :profile="profile" />

      <!-- Features Showcase Section -->
      <section class="section-block">
        <div class="section-header text-center">
          <span class="badge badge-emerald">FITUR DIGITAL</span>
          <h2>Keunggulan Portal E-Learning</h2>
          <p class="section-subtitle">Fasilitas modern untuk mempermudah pengalaman belajar mengajar.</p>
        </div>

        <div class="features-grid">
          <div class="glass-card feature-card">
            <div class="feature-icon-box gold-box">
              <Terminal class="feature-icon text-gold" />
            </div>
            <h4>Live Code Sandbox</h4>
            <p>Editor HTML, CSS & JavaScript langsung di browser untuk latihan praktikum Pemrograman Web secara realtime.</p>
          </div>

          <div class="glass-card feature-card">
            <div class="feature-icon-box blue-box">
              <CheckSquare class="feature-icon text-blue" />
            </div>
            <h4>Simulator Kuis Online</h4>
            <p>Kuis pilihan ganda dengan timer countdown, auto-grading instan, serta pembahasan jawaban langsung dari Dosen.</p>
          </div>

          <div class="glass-card feature-card">
            <div class="feature-icon-box emerald-box">
              <Zap class="feature-icon text-emerald" />
            </div>
            <h4>Presensi & Tugas Digital</h4>
            <p>Sistem absensi otomatis dan portal pengumpulan tugas terintegrasi repository GitHub/Google Drive.</p>
          </div>
        </div>
      </section>

    </div>

    <!-- ========================================================================= -->
    <!-- VIEW B: TAMPILAN DASHBOARD SETELAH USER LOGIN (PERSONAL DASHBOARD)      -->
    <!-- ========================================================================= -->
    <div v-else class="dashboard-user-wrapper">
      <!-- Welcome Logged In User Header -->
      <div class="glass-card welcome-banner">
        <div class="welcome-text">
          <span class="badge badge-gold">
            {{ currentUser.role === 'dosen' ? 'MODE DOSEN PENGAMPU' : 'MODE MAHASISWA TERDAFTAR' }}
          </span>
          <h2>Selamat Datang Kembali, {{ currentUser.name }}!</h2>
          <p v-if="currentUser.role === 'dosen'">
            Anda dapat mengedit tampilan Beranda, mengelola pengumuman, melihat rekapitulasi tugas masuk, dan absensi di Panel Dosen.
          </p>
          <p v-else>
            NIM: <strong>{{ currentUser.username }}</strong> | Prodi: <strong>{{ currentUser.prodi || 'Teknik Informatika' }}</strong>. Selamat melanjutkan perkuliahan!
          </p>
        </div>
        
        <div class="welcome-actions">
          <router-link v-if="currentUser.role === 'dosen'" to="/admin" class="btn btn-gold btn-lg">
            <ShieldCheck class="btn-icon-sm" />
            <span>Buka Panel Dosen</span>
          </router-link>
          <router-link v-else to="/quizzes" class="btn btn-gold btn-lg">
            <CheckSquare class="btn-icon-sm" />
            <span>Mulai Kuis Online</span>
          </router-link>
        </div>
      </div>

      <!-- Dosen Profile Card (Khusus Dosen) -->
      <DosenProfile v-if="currentUser?.role === 'dosen'" :profile="profile" />

      <!-- Matakuliah Grid -->
      <section class="section-block">
        <div class="section-header">
          <h2>Mata Kuliah Semester Ini</h2>
          <p class="section-subtitle">Silakan pilih mata kuliah untuk mengakses modul pertemuan 1 - 14.</p>
        </div>

        <div class="course-grid">
          <div 
            v-for="course in courses" 
            :key="course.id" 
            class="glass-card course-card"
          >
            <div class="card-badge-row">
              <span class="badge badge-gold">{{ course.code }}</span>
              <span class="badge badge-blue">{{ course.sks }} SKS</span>
            </div>

            <h3 class="course-title">{{ course.name }}</h3>
            <p class="course-desc">{{ course.description }}</p>

            <div class="course-meta">
              <div class="meta-row">
                <Calendar class="meta-icon" />
                <span>{{ course.class_time }}</span>
              </div>
            </div>

            <div class="card-footer">
              <router-link :to="'/course/' + course.id" class="btn btn-primary w-full">
                <span>Buka Modul Pertemuan (1 - 14)</span>
                <ArrowRight class="btn-icon-xs" />
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <!-- Pengumuman Akademik -->
      <section class="section-block grid-2-col">
        <div class="announcement-section">
          <div class="section-header">
            <div class="title-with-icon">
              <Megaphone class="header-icon text-gold" />
              <h3>Pengumuman Akademik Terbaru</h3>
            </div>
          </div>

          <div class="announcement-list">
            <div 
              v-for="ann in announcements" 
              :key="ann.id" 
              class="glass-card ann-card"
            >
              <div class="ann-header">
                <span 
                  :class="['badge', ann.category === 'Penting' ? 'badge-rose' : ann.category === 'Tugas' ? 'badge-gold' : 'badge-blue']"
                >
                  {{ ann.category }}
                </span>
                <span class="ann-date">{{ new Date(ann.created_at).toLocaleDateString('id-ID') }}</span>
              </div>
              <h4 class="ann-title">{{ ann.title }}</h4>
              <p class="ann-content">{{ ann.content }}</p>
              <div class="ann-author">
                <span>Oleh: {{ ann.author }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="info-sidebar">
          <div class="glass-card info-card">
            <h3>📌 Panduan Perkuliahan</h3>
            <ul class="rule-list">
              <li>
                <CheckCircle class="rule-icon text-emerald" />
                <span>Kehadiran minimal 75% untuk dapat mengikuti UAS.</span>
              </li>
              <li>
                <CheckCircle class="rule-icon text-emerald" />
                <span>Tugas dikumpulkan melalui portal ini sebelum deadline.</span>
              </li>
              <li>
                <CheckCircle class="rule-icon text-emerald" />
                <span>Gunakan <strong>Live Code Sandbox</strong> untuk latihan Pemrograman Web.</span>
              </li>
            </ul>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  max-width: 1350px;
  margin: 0 auto;
  padding: 1.5rem;
}

.landing-public-wrapper {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.hero-section {
  padding: 3.5rem 2.5rem;
  text-align: center;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  border-radius: var(--radius-lg);
  box-shadow: 0 10px 30px -5px rgba(37, 99, 235, 0.08);
}

.hero-content {
  max-width: 880px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.badge-row {
  display: flex;
  gap: 0.6rem;
  margin-bottom: 1.25rem;
}

.hero-title {
  font-size: 2.75rem;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 1rem;
  color: #0f172a;
}

.hero-subtitle {
  color: #334155;
  font-size: 1.15rem;
  line-height: 1.6;
  margin-bottom: 2rem;
  max-width: 780px;
  font-weight: 400;
}

.highlight-text {
  color: #1e3a8a;
  font-weight: 700;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  margin-bottom: 2.5rem;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-lg {
  padding: 0.85rem 1.75rem;
  font-size: 1rem;
}

.shadow-gold {
  box-shadow: 0 6px 20px rgba(217, 119, 6, 0.35);
}

.hero-stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  width: 100%;
  max-width: 720px;
  padding-top: 1.5rem;
  border-top: 1.5px solid #e2e8f0;
}

.stat-card {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 1rem;
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stat-num {
  font-family: var(--font-heading);
  font-size: 2.2rem;
  font-weight: 800;
  color: #d97706;
  line-height: 1;
}

.stat-label {
  font-size: 0.82rem;
  font-weight: 700;
  color: #0f172a;
  text-transform: uppercase;
  margin-top: 0.3rem;
  letter-spacing: 0.03em;
}

.text-center {
  text-align: center;
}

.portal-choice-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-top: 1.5rem;
}

.choice-card {
  padding: 2.25rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  background: #ffffff;
}

.dosen-choice {
  border-top: 4px solid #d97706;
}

.student-choice {
  border-top: 4px solid #2563eb;
}

.choice-icon-box {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.25rem;
}

.gold-box {
  background: #fffbeb;
  border: 1px solid #fde68a;
}

.blue-box {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
}

.emerald-box {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
}

.choice-icon {
  width: 30px;
  height: 30px;
}

.choice-card h3 {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 0.6rem;
}

.choice-card p {
  color: #475569;
  font-size: 0.92rem;
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.mt-auto {
  margin-top: auto;
}

.button-group-row {
  display: flex;
  gap: 0.75rem;
  width: 100%;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-top: 1.5rem;
}

.feature-card {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  background: #ffffff;
}

.feature-icon-box {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.feature-icon {
  width: 28px;
  height: 28px;
}

.feature-card h4 {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.feature-card p {
  color: #475569;
  font-size: 0.9rem;
}

.welcome-banner {
  padding: 2rem;
  margin-bottom: 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  background: #ffffff;
}

.welcome-text h2 {
  font-size: 1.8rem;
  font-weight: 800;
  margin-top: 0.4rem;
  margin-bottom: 0.3rem;
}

.welcome-text p {
  color: #475569;
  font-size: 0.95rem;
}

.section-block {
  margin-bottom: 2.5rem;
}

.section-header {
  margin-bottom: 1.5rem;
}

.section-subtitle {
  color: #475569;
  font-size: 0.92rem;
  margin-top: 0.2rem;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
}

.course-card {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: #ffffff;
}

.card-badge-row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.course-title {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 0.6rem;
}

.course-desc {
  color: #475569;
  font-size: 0.92rem;
  margin-bottom: 1.25rem;
}

.course-meta {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  margin-bottom: 1.5rem;
  padding: 0.85rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: #475569;
}

.meta-icon {
  width: 16px;
  height: 16px;
  color: #d97706;
}

.w-full {
  width: 100%;
}

.grid-2-col {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.header-icon {
  width: 24px;
  height: 24px;
}

.announcement-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.ann-card {
  padding: 1.25rem;
  background: #ffffff;
}

.ann-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.6rem;
}

.ann-date {
  font-size: 0.78rem;
  color: #64748b;
}

.ann-title {
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
}

.ann-content {
  color: #334155;
  font-size: 0.9rem;
  margin-bottom: 0.75rem;
}

.ann-author {
  font-size: 0.78rem;
  color: #b45309;
  font-weight: 600;
}

.info-card {
  padding: 1.5rem;
  background: #ffffff;
}

.info-card h3 {
  font-size: 1.1rem;
  margin-bottom: 1rem;
}

.rule-list {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.rule-list li {
  display: flex;
  align-items: flex-start;
  gap: 0.65rem;
  font-size: 0.88rem;
  color: #334155;
}

.rule-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  margin-top: 2px;
}

.text-gold {
  color: #d97706;
}

.text-blue {
  color: #2563eb;
}

.text-emerald {
  color: #059669;
}

.btn-icon-sm {
  width: 18px;
  height: 18px;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}

@media (max-width: 900px) {
  .hero-title {
    font-size: 2.1rem;
  }
  .hero-stats-grid {
    grid-template-columns: 1fr;
  }
  .portal-choice-grid {
    grid-template-columns: 1fr;
  }
  .grid-2-col {
    grid-template-columns: 1fr;
  }
  .welcome-banner {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
