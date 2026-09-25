<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showSuccess, showError, showWarning, showConfirm } from '../utils/swal.js'
import { 
  CheckSquare, 
  HelpCircle, 
  Plus, 
  Trash2,
  Search, 
  Award, 
  CheckCircle2, 
  XCircle, 
  Clock, 
  BookOpen, 
  Users, 
  BarChart3,
  FileText,
  Code,
  Layers,
  ChevronRight,
  PlusCircle,
  X,
  Radio
} from 'lucide-vue-next'

const router = useRouter()

const courses = ref([])
const quizzes = ref([])
const quizSubmissions = ref([])

const selectedCourseId = ref('')
const selectedQuizId = ref('')
const activeMainTab = ref('soal') // 'soal' | 'rekap'
const showAddQuestionForm = ref(false)
const showCreateQuizModal = ref(false)
const isCreatingQuiz = ref(false)

const searchQuery = ref('')

// Form State for New Quiz Package
const newQuizTitle = ref('')
const newQuizDescription = ref('')
const newQuizTimeLimit = ref(15)

// Form State for New Question
const newQType = ref('mc') // 'mc' | 'essay'
const newQQuestion = ref('')
const newQOpt0 = ref('')
const newQOpt1 = ref('')
const newQOpt2 = ref('')
const newQOpt3 = ref('')
const newQCorrect = ref(0)
const newQExplanation = ref('')

const fetchData = async () => {
  try {
    const [resCourses, resQuizzes, resSubs] = await Promise.all([
      fetch('/api/courses'),
      fetch('/api/quizzes'),
      fetch('/api/quizzes/submissions')
    ])
    courses.value = await resCourses.json()
    quizzes.value = await resQuizzes.json()
    quizSubmissions.value = await resSubs.json()

    // Default select first course if not selected
    if (courses.value.length > 0 && !selectedCourseId.value) {
      selectedCourseId.value = courses.value[0].id
    }
  } catch (err) {
    console.warn('BankSoal fetch error:', err)
  }
}

// Check current user
const currentUser = ref(null)
const checkUser = () => {
  const u = localStorage.getItem('user')
  if (u) {
    try { currentUser.value = JSON.parse(u) } catch { currentUser.value = null }
  }
}

onMounted(() => {
  checkUser()
  fetchData()
  window.addEventListener('auth-changed', checkUser)
  window.addEventListener('course-changed', fetchData)
})

const isDosen = computed(() => currentUser.value?.role === 'dosen')

// Active course object
const selectedCourse = computed(() => {
  return courses.value.find(c => c.id === selectedCourseId.value) || null
})

// Quizzes belonging to selected course
const courseQuizzes = computed(() => {
  if (!selectedCourseId.value) return []
  return quizzes.value.filter(q => q.course_id === selectedCourseId.value)
})

// Auto update selectedQuizId when selectedCourseId changes or when quizzes load
watch(selectedCourseId, (newCourseId) => {
  const firstQuiz = quizzes.value.find(q => q.course_id === newCourseId)
  if (firstQuiz) {
    selectedQuizId.value = firstQuiz.id
  } else {
    selectedQuizId.value = ''
  }
}, { immediate: true })

watch(quizzes, (newQuizzes) => {
  if (selectedCourseId.value && !selectedQuizId.value) {
    const firstQuiz = newQuizzes.find(q => q.course_id === selectedCourseId.value)
    if (firstQuiz) {
      selectedQuizId.value = firstQuiz.id
    }
  }
})

// Active selected quiz object
const selectedQuiz = computed(() => {
  return quizzes.value.find(x => x.id === selectedQuizId.value) || null
})

// Questions of selected quiz
const activeQuestions = computed(() => {
  return selectedQuiz.value?.questions || []
})

// Count quizzes per course
const getCourseQuizCount = (courseId) => {
  return quizzes.value.filter(q => q.course_id === courseId).length
}

// Stats overview
const totalQuestionsCount = computed(() => {
  return quizzes.value.reduce((acc, q) => acc + (q.questions?.length || 0), 0)
})

const averageScore = computed(() => {
  if (quizSubmissions.value.length === 0) return 0
  const total = quizSubmissions.value.reduce((acc, s) => acc + s.score, 0)
  return Math.round(total / quizSubmissions.value.length)
})

const passRate = computed(() => {
  if (quizSubmissions.value.length === 0) return 0
  const passed = quizSubmissions.value.filter(s => s.score >= 70).length
  return Math.round((passed / quizSubmissions.value.length) * 100)
})

// Submissions filtered for active course & search query
const filteredSubmissions = computed(() => {
  return quizSubmissions.value.filter(s => {
    // If student (not dosen), only show their own submissions
    if (!isDosen.value && currentUser.value) {
      if (s.student_nim !== currentUser.value.username) return false
    }

    // Filter by selected quiz if present
    if (selectedQuizId.value && s.quiz_id !== selectedQuizId.value) {
      const subQuiz = quizzes.value.find(q => q.id === s.quiz_id)
      if (subQuiz && subQuiz.course_id !== selectedCourseId.value) {
        return false
      }
    }
    
    const q = searchQuery.value.toLowerCase()
    return (
      s.student_name?.toLowerCase().includes(q) ||
      s.student_nim?.toLowerCase().includes(q) ||
      s.quiz_id?.toLowerCase().includes(q)
    )
  })
})

// Create New Quiz Package
const createQuizPackage = async () => {
  if (!newQuizTitle.value) {
    showWarning('Perhatian', 'Judul Paket Kuis wajib diisi!')
    return
  }

  isCreatingQuiz.value = true
  try {
    const res = await fetch('/api/quizzes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: selectedCourseId.value,
        title: newQuizTitle.value,
        description: newQuizDescription.value || 'Paket Kuis / Ujian untuk ' + (selectedCourse.value?.name || ''),
        time_limit: parseInt(newQuizTimeLimit.value) || 15
      })
    })

    if (res.ok) {
      const created = await res.json()
      showSuccess('Paket Kuis Dibuat! 📝', `Paket kuis '${created.title}' berhasil ditambahkan ke mata kuliah ini.`)
      newQuizTitle.value = ''
      newQuizDescription.value = ''
      showCreateQuizModal.value = false
      await fetchData()
      selectedQuizId.value = created.id
    } else {
      showError('Gagal!', 'Gagal membuat paket kuis baru.')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan sistem saat membuat paket kuis.')
  } finally {
    isCreatingQuiz.value = false
  }
}

// Add Question Handler
const addQuizQuestion = async () => {
  if (!selectedQuizId.value) {
    showWarning('Perhatian', 'Harap pilih atau buat paket kuis terlebih dahulu!')
    return
  }
  if (newQType.value === 'mc') {
    if (!newQQuestion.value || !newQOpt0.value || !newQOpt1.value || !newQOpt2.value || !newQOpt3.value) {
      showWarning('Perhatian', 'Harap lengkapi Teks Pertanyaan dan ke-4 opsi jawaban (A, B, C, D)!')
      return
    }
  } else {
    if (!newQQuestion.value) {
      showWarning('Perhatian', 'Harap isi Teks Pertanyaan Essay!')
      return
    }
  }

  try {
    const res = await fetch('/api/quizzes/questions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        quiz_id: selectedQuizId.value,
        question: {
          type: newQType.value,
          question: newQQuestion.value,
          options: newQType.value === 'mc' ? [newQOpt0.value, newQOpt1.value, newQOpt2.value, newQOpt3.value] : [],
          correct_answer: newQType.value === 'mc' ? parseInt(newQCorrect.value) : 0,
          explanation: newQExplanation.value || 'Pembahasan materi oleh Pak Rio Widyatmoko.'
        }
      })
    })

    if (res.ok) {
      showSuccess('Soal Berhasil Ditambahkan! 📝', 'Soal kuis baru telah tersimpan di Bank Soal.')
      newQQuestion.value = ''
      newQOpt0.value = ''
      newQOpt1.value = ''
      newQOpt2.value = ''
      newQOpt3.value = ''
      newQExplanation.value = ''
      showAddQuestionForm.value = false
      fetchData()
    } else {
      showError('Gagal!', 'Gagal menambahkan soal kuis baru.')
    }
  } catch (err) {
    showError('Gagal!', 'Terjadi kesalahan sistem saat menambahkan soal!')
  }
}

// Launch Live Quiz Interaktif
const startLiveQuiz = async (quiz) => {
  if (!quiz) return
  try {
    const res = await fetch('/api/live-quiz/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        quiz_id: quiz.id,
        time_per_question: quiz.time_limit || 20
      })
    })

    if (!res.ok) {
      showError('Gagal Membuat Sesi', 'Tidak dapat membuat sesi kuis live.')
      return
    }

    const session = await res.json()
    showSuccess('Sesi Live Quiz Dibuat! 🚀', `PIN: ${session.pin}. Mengalihkan ke Layar Proyektor...`)
    router.push(`/live-quiz/${session.pin}?host=true`)
  } catch (err) {
    showError('Terjadi Kesalahan', 'Gagal membuat sesi live kuis.')
  }
}

// Delete entire Quiz Package
const deleteQuizPackage = async (quizId, quizTitle) => {
  const confirmed = await showConfirm(
    'Hapus Paket Kuis?',
    `Apakah Anda yakin ingin menghapus paket kuis '${quizTitle}' beserta seluruh soal di dalamnya?`,
    'Ya, Hapus Paket'
  )
  if (!confirmed) return

  try {
    const res = await fetch(`/api/quizzes?id=${quizId}`, { method: 'DELETE' })
    if (res.ok) {
      showSuccess('Paket Kuis Dihapus! 🗑️', `Paket kuis '${quizTitle}' telah berhasil dihapus.`)
      await fetchData()
      if (courseQuizzes.value.length > 0) {
        selectedQuizId.value = courseQuizzes.value[0].id
      } else {
        selectedQuizId.value = ''
      }
    } else {
      showError('Gagal!', 'Gagal menghapus paket kuis.')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan sistem saat menghapus paket kuis.')
  }
}

// Delete individual question from active quiz
const deleteQuestion = async (qId, index) => {
  const confirmed = await showConfirm(
    'Hapus Soal Kuis?',
    `Apakah Anda yakin ingin menghapus Soal #${index + 1}?`,
    'Ya, Hapus Soal'
  )
  if (!confirmed) return

  try {
    const res = await fetch(`/api/quizzes/questions?quiz_id=${selectedQuizId.value}&question_id=${qId}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      showSuccess('Soal Dihapus! 🗑️', `Soal #${index + 1} berhasil dihapus dari paket kuis.`)
      fetchData()
    } else {
      showError('Gagal!', 'Gagal menghapus soal.')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan sistem saat menghapus soal.')
  }
}
</script>

<template>
  <div class="bank-soal-container animate-fade-in">
    <!-- Header -->
    <div class="page-header">
      <div class="title-row">
        <HelpCircle class="header-icon text-gold" />
        <div>
          <h2>Bank Soal Kuis & Rekapitulasi Nilai Mahasiswa</h2>
          <p class="subtitle">Kelola soal kuis CBT dan rekapitulasi nilai mahasiswa terstruktur berdasarkan mata kuliah.</p>
        </div>
      </div>
    </div>

    <!-- Stats Overview Cards Grid -->
    <div class="stats-grid">
      <div class="glass-card stat-card border-gold">
        <div class="stat-icon-box bg-gold-subtle text-gold">
          <HelpCircle class="icon-md" />
        </div>
        <div class="stat-info">
          <span class="stat-num">{{ totalQuestionsCount }}</span>
          <span class="stat-label">Total Soal di Bank Kuis</span>
        </div>
      </div>

      <div class="glass-card stat-card border-blue">
        <div class="stat-icon-box bg-blue-subtle text-blue">
          <Users class="icon-md" />
        </div>
        <div class="stat-info">
          <span class="stat-num">{{ quizSubmissions.length }}</span>
          <span class="stat-label">Kuis Dikerjakan Mahasiswa</span>
        </div>
      </div>

      <div class="glass-card stat-card border-emerald">
        <div class="stat-icon-box bg-emerald-subtle text-emerald">
          <Award class="icon-md" />
        </div>
        <div class="stat-info">
          <span class="stat-num">{{ averageScore }} / 100</span>
          <span class="stat-label">Rata-Rata Nilai Kuis</span>
        </div>
      </div>

      <div class="glass-card stat-card border-purple">
        <div class="stat-icon-box bg-purple-subtle text-purple">
          <BarChart3 class="icon-md" />
        </div>
        <div class="stat-info">
          <span class="stat-num">{{ passRate }}%</span>
          <span class="stat-label">Tingkat Kelulusan KKM</span>
        </div>
      </div>
    </div>

    <!-- WORKFLOW STEP 1: PILIH MATA KULIAH -->
    <section class="step-section">
      <div class="step-header">
        <div class="step-number">1</div>
        <div>
          <h3>Langkah 1: Pilih Mata Kuliah</h3>
          <p class="step-sub">Pilih mata kuliah di bawah untuk mengelola kuis, bank soal, dan melihat rekap nilai mahasiswa.</p>
        </div>
      </div>

      <div class="course-cards-grid">
        <button 
          v-for="c in courses" 
          :key="c.id"
          @click="selectedCourseId = c.id"
          :class="['course-select-card', selectedCourseId === c.id ? 'active' : '']"
        >
          <div class="crs-card-header">
            <component :is="(c.name.toLowerCase().includes('web') || c.id.includes('web')) ? Code : BookOpen" class="crs-card-icon" />
            <span class="code-badge">{{ c.code }}</span>
          </div>
          <h4 class="crs-card-title">{{ c.name }}</h4>
          <div class="crs-card-meta">
            <span>{{ c.sks }} SKS</span>
            <span>•</span>
            <span>{{ getCourseQuizCount(c.id) }} Paket Kuis</span>
          </div>
          <div v-if="selectedCourseId === c.id" class="selected-indicator">
            <CheckCircle2 class="check-icon" /> Terpilih
          </div>
        </button>
      </div>
    </section>

    <!-- WORKFLOW STEP 2: PILIH / BUAT PAKET KUIS (ACTIVE WHEN COURSE SELECTED) -->
    <section v-if="selectedCourse" class="step-section animate-fade-in">
      <div class="step-header">
        <div class="step-number">2</div>
        <div>
          <h3>Langkah 2: Pilih Paket Kuis — <span>{{ selectedCourse.name }}</span></h3>
          <p class="step-sub">Pilih paket kuis yang ingin dikelola atau tambahkan paket kuis baru untuk mata kuliah ini.</p>
        </div>
      </div>

      <div class="quiz-packages-row">
        <!-- Quiz Package Pills/Tabs -->
        <div 
          v-for="q in courseQuizzes" 
          :key="q.id"
          :class="['quiz-pkg-card-item', selectedQuizId === q.id ? 'active' : '']"
        >
          <button 
            @click="selectedQuizId = q.id"
            class="quiz-package-btn"
          >
            <div class="qp-info">
              <CheckSquare class="qp-icon" />
              <div class="qp-text">
                <span class="qp-title">{{ q.title }}</span>
                <span class="qp-meta">{{ q.questions?.length || 0 }} Soal • {{ q.time_limit || 15 }} Menit Durasi</span>
              </div>
            </div>
            <span v-if="selectedQuizId === q.id" class="qp-active-badge">Aktif</span>
          </button>
          <button 
            v-if="isDosen"
            @click.stop="deleteQuizPackage(q.id, q.title)" 
            class="btn-delete-quiz-pkg" 
            title="Hapus Paket Kuis Ini"
          >
            <Trash2 class="del-pkg-icon" />
          </button>
        </div>

        <!-- Button Create New Quiz Package -->
        <button v-if="isDosen" @click="showCreateQuizModal = true" class="btn-add-quiz-package">
          <PlusCircle class="btn-icon-sm" />
          <span>+ Buat Paket Kuis Baru</span>
        </button>
      </div>

      <div v-if="courseQuizzes.length === 0" class="empty-quiz-notice glass-card">
        <HelpCircle class="empty-icon text-amber" />
        <div>
          <h4>Belum Ada Paket Kuis untuk {{ selectedCourse.name }}</h4>
          <p>Klik tombol <strong>"+ Buat Paket Kuis Baru"</strong> untuk menerbitkan kuis pertama pada mata kuliah ini.</p>
        </div>
        <button @click="showCreateQuizModal = true" class="btn btn-gold btn-sm">
          <Plus class="btn-icon-xs" /> Buat Paket Kuis Baru
        </button>
      </div>
    </section>

    <!-- WORKFLOW STEP 3: KELOLA SOAL & REKAP NILAI / REVIEW MAHASISWA -->
    <section v-if="selectedQuiz" class="step-section animate-fade-in">
      <div class="step-header step-header-space">
        <div class="step-title-group">
          <div class="step-number">3</div>
          <div>
            <h3>{{ isDosen ? 'Langkah 3: Kelola Soal & Rekap Nilai' : 'Review Bank Soal & Pembahasan Kuis' }} — <span>{{ selectedQuiz.title }}</span></h3>
            <p class="step-sub">{{ isDosen ? 'Tambah pertanyaan pilihan ganda, atur kunci jawaban, atau pantau rekapitulasi nilai kuis mahasiswa.' : 'Silakan pelajari kumpulan soal kuis, kunci jawaban, dan pembahasan materi untuk pendalaman perkuliahan.' }}</p>
          </div>
        </div>

        <div v-if="isDosen" class="step-actions-group" style="display: flex; gap: 0.75rem; align-items: center;">
          <button 
            @click="startLiveQuiz(selectedQuiz)" 
            class="btn btn-gold btn-start-live-quiz"
            style="background: #d97706; color: white; border: none; padding: 0.65rem 1.2rem; border-radius: 0.75rem; font-weight: 800; cursor: pointer; display: flex; align-items: center; gap: 0.5rem; box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25);"
            title="Mulai Sesi Presenter Live Quiz Interaktif pada Layar Proyektor"
          >
            <Radio class="btn-icon-xs animate-pulse" />
            <span>Mulai Live Quiz Interaktif 🚀</span>
          </button>

          <button 
            @click="deleteQuizPackage(selectedQuiz.id, selectedQuiz.title)" 
            class="btn-del-active-quiz"
            title="Hapus Paket Kuis Terpilih Ini"
          >
            <Trash2 class="btn-icon-xs" />
            <span>Hapus Paket Kuis Ini</span>
          </button>
        </div>
      </div>

      <!-- Tabs Navigation for Step 3 -->
      <div class="step3-tab-bar">
        <button 
          @click="activeMainTab = 'soal'" 
          :class="['step3-tab-btn', activeMainTab === 'soal' ? 'active' : '']"
        >
          <HelpCircle class="tab-icon-sm" />
          <span>Kelola Bank Soal ({{ activeQuestions.length }} Soal)</span>
        </button>
        <button 
          @click="activeMainTab = 'rekap'" 
          :class="['step3-tab-btn', activeMainTab === 'rekap' ? 'active' : '']"
        >
          <Award class="tab-icon-sm" />
          <span>Rekapitulasi Nilai Mahasiswa ({{ filteredSubmissions.length }} Submisi)</span>
        </button>
      </div>

      <!-- TAB 1: KELOLA SOAL / REVIEW SOAL -->
      <div v-if="activeMainTab === 'soal'" class="glass-card main-panel-card animate-fade-in">
        <div class="panel-header-row">
          <div>
            <h4 class="panel-title">{{ isDosen ? 'Daftar Soal pada' : 'Review Soal & Kunci Jawaban' }} {{ selectedQuiz.title }}</h4>
            <p class="panel-sub">Total {{ activeQuestions.length }} soal pilihan ganda telah tersimpan.</p>
          </div>

          <button 
            v-if="isDosen"
            @click="showAddQuestionForm = !showAddQuestionForm" 
            :class="['btn', showAddQuestionForm ? 'btn-secondary' : 'btn-gold']"
          >
            <component :is="showAddQuestionForm ? X : Plus" class="btn-icon-xs" />
            <span>{{ showAddQuestionForm ? 'Tutup Form Soal' : '+ Tambah Soal Baru ke Paket Ini' }}</span>
          </button>
        </div>

        <!-- Collapsible Form Input Soal Baru -->
        <div v-if="showAddQuestionForm" class="add-question-form-wrapper animate-fade-in">
          <div class="form-banner">
            <Plus class="form-banner-icon text-gold" />
            <span>Form Tambah Soal Baru ({{ selectedQuiz.title }})</span>
          </div>

          <form @submit.prevent="addQuizQuestion" class="bank-form">
            <div class="form-group" style="margin-bottom: 1rem;">
              <label class="input-label">Pilih Jenis Soal *</label>
              <div style="display: flex; gap: 1rem;">
                <label style="cursor: pointer; display: flex; align-items: center; gap: 0.4rem; font-weight: 700;">
                  <input type="radio" v-model="newQType" value="mc" /> Pilihan Ganda (Multiple Choice)
                </label>
                <label style="cursor: pointer; display: flex; align-items: center; gap: 0.4rem; font-weight: 700;">
                  <input type="radio" v-model="newQType" value="essay" /> Soal Essay / Uraian
                </label>
              </div>
            </div>

            <div class="form-group">
              <label class="input-label">Teks Pertanyaan / Soal Kuis *</label>
              <textarea 
                v-model="newQQuestion" 
                class="glass-input textarea" 
                required 
                :placeholder="newQType === 'essay' ? 'Contoh: Jelaskan perbedaan utama antara model pengembangan SDLC Waterfall dan Agile Scrum!' : 'Contoh: Apa kegunaan utama dari metode Waterfall dalam SDLC?'"
              ></textarea>
            </div>

            <div v-if="newQType === 'mc'">
              <div class="form-row">
                <div class="form-group">
                  <label class="input-label">Pilihan Jawaban A *</label>
                  <input v-model="newQOpt0" class="glass-input" :required="newQType === 'mc'" placeholder="Teks Opsi A..." />
                </div>
                <div class="form-group">
                  <label class="input-label">Pilihan Jawaban B *</label>
                  <input v-model="newQOpt1" class="glass-input" :required="newQType === 'mc'" placeholder="Teks Opsi B..." />
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label class="input-label">Pilihan Jawaban C *</label>
                  <input v-model="newQOpt2" class="glass-input" :required="newQType === 'mc'" placeholder="Teks Opsi C..." />
                </div>
                <div class="form-group">
                  <label class="input-label">Pilihan Jawaban D *</label>
                  <input v-model="newQOpt3" class="glass-input" :required="newQType === 'mc'" placeholder="Teks Opsi D..." />
                </div>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="input-label">Kunci Jawaban Benar *</label>
                <select v-model.number="newQCorrect" class="glass-input">
                  <option :value="0">Opsi A (Pilihan Pertama)</option>
                  <option :value="1">Opsi B (Pilihan Kedua)</option>
                  <option :value="2">Opsi C (Pilihan Ketiga)</option>
                  <option :value="3">Opsi D (Pilihan Keempat)</option>
                </select>
              </div>
              <div class="form-group">
                <label class="input-label">Penjelasan / Pembahasan Kunci</label>
                <input v-model="newQExplanation" class="glass-input" placeholder="Pembahasan singkat..." />
              </div>
            </div>

            <div class="form-actions">
              <button type="button" @click="showAddQuestionForm = false" class="btn btn-secondary">Batal</button>
              <button type="submit" class="btn btn-gold btn-lg">
                <Plus class="btn-icon-sm" />
                <span>Simpan Soal ke Bank Kuis</span>
              </button>
            </div>
          </form>
        </div>

        <!-- Questions List Grid -->
        <div v-if="activeQuestions.length === 0" class="empty-state">
          Belum ada soal pada paket kuis ini. Klik tombol <strong>"+ Tambah Soal Baru ke Paket Ini"</strong> untuk menginput soal.
        </div>
        <div v-else class="questions-list">
          <div 
            v-for="(q, idx) in activeQuestions" 
            :key="idx" 
            class="question-item"
          >
            <div class="q-header">
              <div class="q-header-left">
                <span class="q-num">Soal #{{ idx + 1 }}</span>
                <span :class="['badge', q.type === 'essay' ? 'badge-purple' : 'badge-gold']" style="margin-left: 0.5rem;">
                  {{ q.type === 'essay' ? 'ESSAY' : 'PILHAN GANDA' }}
                </span>
                <span v-if="q.type !== 'essay'" class="q-key-badge" style="margin-left: 0.5rem;">Kunci Jawaban: Opsi {{ ['A', 'B', 'C', 'D'][q.correct_answer] }}</span>
              </div>
              <button v-if="isDosen" @click="deleteQuestion(q.id, idx)" class="btn-del-question" title="Hapus Soal Ini">
                <Trash2 class="btn-icon-xs text-rose" />
                <span>Hapus Soal</span>
              </button>
            </div>
            <p class="q-text">{{ q.question }}</p>
            <div v-if="q.type !== 'essay'" class="q-options-grid">
              <div :class="['q-opt', q.correct_answer === 0 ? 'correct' : '']">A. {{ q.options[0] }}</div>
              <div :class="['q-opt', q.correct_answer === 1 ? 'correct' : '']">B. {{ q.options[1] }}</div>
              <div :class="['q-opt', q.correct_answer === 2 ? 'correct' : '']">C. {{ q.options[2] }}</div>
              <div :class="['q-opt', q.correct_answer === 3 ? 'correct' : '']">D. {{ q.options[3] }}</div>
            </div>
            <p v-if="q.explanation" class="q-explanation">
              💡 <strong>Pembahasan:</strong> {{ q.explanation }}
            </p>
          </div>
        </div>
      </div>

      <!-- TAB 2: REKAPITULASI NILAI MAHASISWA -->
      <div v-if="activeMainTab === 'rekap'" class="glass-card main-panel-card animate-fade-in">
        <div class="panel-header-row">
          <div>
            <h4 class="panel-title">Rekapitulasi Nilai Ujian — {{ selectedQuiz.title }}</h4>
            <p class="panel-sub">Daftar mahasiswa ITB Swadharma yang telah menyelesaikan kuis ini.</p>
          </div>

          <!-- Search Bar -->
          <div class="search-box">
            <Search class="search-icon" />
            <input 
              v-model="searchQuery" 
              class="glass-input search-input" 
              placeholder="Cari Nama Mahasiswa / NIM..."
            />
          </div>
        </div>

        <div v-if="filteredSubmissions.length === 0" class="empty-state">
          Belum ada rekapitulasi nilai kuis mahasiswa yang ditemukan untuk paket kuis ini.
        </div>
        <div v-else class="submissions-grid">
          <div 
            v-for="sub in filteredSubmissions" 
            :key="sub.id" 
            class="submission-card"
          >
            <div class="sub-header-row">
              <div class="student-meta">
                <span class="student-name">{{ sub.student_name }}</span>
                <span class="student-nim">NIM: {{ sub.student_nim }}</span>
              </div>
              <span :class="['badge', sub.score >= 70 ? 'badge-emerald' : 'badge-rose']">
                {{ sub.score >= 70 ? 'LULUS (KKM 70)' : 'REMIDIAL' }}
              </span>
            </div>

            <div class="sub-detail-row">
              <span class="quiz-tag">📝 {{ selectedQuiz.title }}</span>
              <span class="time-tag">📅 {{ new Date(sub.submitted_at).toLocaleString('id-ID') }}</span>
            </div>

            <div class="sub-score-row">
              <span class="score-label">Nilai Ujian:</span>
              <span :class="['score-value', sub.score >= 70 ? 'text-emerald' : 'text-rose']">
                {{ sub.score }} <small>/ 100</small>
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- MODAL: BUAT PAKET KUIS BARU -->
    <Teleport to="body">
      <div v-if="showCreateQuizModal" class="modal-overlay" @click.self="showCreateQuizModal = false">
        <div class="glass-card modal-card">
          <div class="modal-header">
            <div class="m-title-box">
              <PlusCircle class="modal-title-icon text-gold" />
              <div>
                <h3>Buat Paket Kuis Baru</h3>
                <p class="modal-sub">Mata Kuliah: {{ selectedCourse?.name }}</p>
              </div>
            </div>
            <button @click="showCreateQuizModal = false" class="btn-close-modal">✕</button>
          </div>

          <form @submit.prevent="createQuizPackage" class="modal-form">
            <div class="form-group">
              <label class="input-label">Judul Paket Kuis *</label>
              <input v-model="newQuizTitle" class="glass-input" required placeholder="Contoh: Kuis 2: UML Diagram & Design Patterns" />
            </div>

            <div class="form-group">
              <label class="input-label">Deskripsi / Petunjuk Kuis</label>
              <textarea v-model="newQuizDescription" class="glass-input textarea" placeholder="Jelaskan cakupan materi kuis..."></textarea>
            </div>

            <div class="form-group">
              <label class="input-label">Durasi Pengerjaan Kuis (Dalam Menit) *</label>
              <input v-model.number="newQuizTimeLimit" type="number" min="5" max="180" class="glass-input" required />
            </div>

            <div class="modal-actions">
              <button type="button" @click="showCreateQuizModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" :disabled="isCreatingQuiz" class="btn btn-gold">
                <Plus class="btn-icon-xs" />
                {{ isCreatingQuiz ? 'Membuat...' : 'Buat Paket Kuis' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.bank-soal-container {
  max-width: 1350px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  margin-bottom: 2rem;
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
  margin: 0;
}

.subtitle {
  color: #475569;
  font-size: 0.95rem;
  margin-top: 0.3rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.25rem;
  margin-bottom: 2rem;
}

.stat-card {
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1.25rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.stat-icon-box {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-md {
  width: 26px;
  height: 26px;
}

.bg-gold-subtle { background: #fefce8; }
.bg-blue-subtle { background: #eff6ff; }
.bg-emerald-subtle { background: #ecfdf5; }
.bg-purple-subtle { background: #faf5ff; }

.text-gold { color: #d97706; }
.text-blue { color: #2563eb; }
.text-emerald { color: #059669; }
.text-purple { color: #9333ea; }
.text-rose { color: #e11d48; }

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-num {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
}

.stat-label {
  font-size: 0.82rem;
  color: #64748b;
  font-weight: 600;
}

.step-section {
  margin-bottom: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.02);
}

.step-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.25rem;
  padding-bottom: 0.85rem;
  border-bottom: 1.5px solid #f1f5f9;
}

.step-header.step-header-space {
  justify-content: space-between;
  align-items: center;
}

.step-title-group {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.btn-del-active-quiz {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.55rem 0.95rem;
  background: #fff1f2;
  border: 1px solid #fecdd3;
  border-radius: var(--radius-sm);
  color: #e11d48;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.btn-del-active-quiz:hover {
  background: #e11d48;
  border-color: #be123c;
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(225, 29, 72, 0.25);
}

.btn-del-question {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.35rem 0.65rem;
  background: #fff1f2;
  border: 1px solid #fecdd3;
  border-radius: var(--radius-sm);
  color: #e11d48;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-del-question:hover {
  background: #e11d48;
  border-color: #be123c;
  color: #ffffff;
}

.step-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #1d4ed8;
  color: #ffffff;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(29, 78, 216, 0.25);
}

.step-header h3 {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.step-header h3 span {
  color: #1d4ed8;
}

.step-sub {
  font-size: 0.85rem;
  color: #64748b;
  margin-top: 0.2rem;
}

/* Step 1 Course Cards Grid */
.course-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

.course-select-card {
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: var(--radius-sm);
  padding: 1.25rem;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  display: flex;
  flex-direction: column;
}

.course-select-card:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  transform: translateY(-2px);
}

.course-select-card.active {
  background: #ffffff;
  border: 2px solid #1d4ed8;
  box-shadow: 0 6px 20px rgba(29, 78, 216, 0.12);
}

.crs-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.crs-card-icon {
  width: 22px;
  height: 22px;
  color: #b45309;
}

.code-badge {
  font-size: 0.75rem;
  font-weight: 800;
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.crs-card-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.4rem;
}

.crs-card-meta {
  font-size: 0.8rem;
  color: #64748b;
  display: flex;
  gap: 0.4rem;
  align-items: center;
}

.selected-indicator {
  margin-top: 0.75rem;
  font-size: 0.75rem;
  font-weight: 800;
  color: #15803d;
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
}

.check-icon {
  width: 14px;
  height: 14px;
}

/* Step 2 Quiz Package Row */
.quiz-packages-row {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  flex-wrap: wrap;
}

.quiz-pkg-card-item {
  display: inline-flex;
  align-items: stretch;
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: var(--radius-sm);
  overflow: hidden;
  transition: all 0.2s ease;
}

.quiz-pkg-card-item:hover {
  border-color: #93c5fd;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.04);
}

.quiz-pkg-card-item.active {
  background: #1d4ed8;
  border-color: #1e40af;
  box-shadow: 0 4px 12px rgba(29, 78, 216, 0.25);
}

.quiz-package-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.85rem 1.1rem;
  background: transparent;
  border: none;
  color: #334155;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 260px;
}

.quiz-pkg-card-item.active .quiz-package-btn {
  color: #ffffff;
}

.btn-delete-quiz-pkg {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 0.85rem;
  background: transparent;
  border: none;
  border-left: 1px solid #e2e8f0;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quiz-pkg-card-item.active .btn-delete-quiz-pkg {
  border-left-color: rgba(255, 255, 255, 0.25);
  color: rgba(255, 255, 255, 0.8);
}

.btn-delete-quiz-pkg:hover {
  background: #ffe4e6;
  color: #e11d48;
}

.quiz-pkg-card-item.active .btn-delete-quiz-pkg:hover {
  background: #b91c1c;
  color: #ffffff;
}

.del-pkg-icon {
  width: 17px;
  height: 17px;
}

.qp-info {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.qp-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.qp-text {
  display: flex;
  flex-direction: column;
  text-align: left;
}

.qp-title {
  font-size: 0.92rem;
  font-weight: 800;
}

.qp-meta {
  font-size: 0.75rem;
  opacity: 0.85;
}

.qp-active-badge {
  font-size: 0.7rem;
  font-weight: 800;
  background: #10b981;
  color: #ffffff;
  padding: 0.15rem 0.5rem;
  border-radius: 99px;
  text-transform: uppercase;
}

.btn-add-quiz-package {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.85rem 1.25rem;
  background: #fffbeb;
  border: 1.5px dashed #f59e0b;
  border-radius: var(--radius-sm);
  color: #b45309;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-quiz-package:hover {
  background: #fef3c7;
  border-color: #d97706;
}

.empty-quiz-notice {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.25rem 1.5rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  margin-top: 1rem;
}

.empty-quiz-notice h4 {
  font-size: 1rem;
  font-weight: 800;
  color: #92400e;
  margin-bottom: 0.2rem;
}

.empty-quiz-notice p {
  font-size: 0.85rem;
  color: #b45309;
}

.empty-icon {
  width: 28px;
  height: 28px;
  flex-shrink: 0;
}

/* Step 3 Sub Tabs */
.step3-tab-bar {
  display: flex;
  gap: 0.6rem;
  margin-bottom: 1.25rem;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 0.5rem;
}

.step3-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 1.25rem;
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: var(--radius-sm);
  color: #475569;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.step3-tab-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.step3-tab-btn.active {
  background: #1d4ed8;
  border-color: #1e40af;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(29, 78, 216, 0.2);
}

.tab-icon-sm {
  width: 17px;
  height: 17px;
}

.main-panel-card {
  padding: 1.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

/* Panel Header Row */
.panel-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
  padding-bottom: 0.85rem;
  border-bottom: 1px solid #f1f5f9;
  flex-wrap: wrap;
  gap: 1rem;
}

.panel-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.panel-sub {
  font-size: 0.85rem;
  color: #64748b;
  margin-top: 0.2rem;
}

/* Add Question Form Wrapper */
.add-question-form-wrapper {
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: var(--radius-sm);
  padding: 1.25rem;
  margin-bottom: 1.5rem;
}

.form-banner {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: var(--font-heading);
  font-size: 0.95rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.form-banner-icon {
  width: 18px;
  height: 18px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.questions-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.question-item {
  padding: 1.25rem;
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-sm);
}

.q-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.q-num {
  font-weight: 800;
  font-size: 0.85rem;
  color: #2563eb;
}

.q-key-badge {
  font-size: 0.78rem;
  font-weight: 700;
  color: #b45309;
  background: #fefce8;
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-sm);
  border: 1px solid #fef08a;
}

.q-text {
  font-weight: 700;
  color: #0f172a;
  font-size: 0.95rem;
  margin-bottom: 0.75rem;
}

.q-options-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.q-opt {
  font-size: 0.82rem;
  padding: 0.45rem 0.7rem;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: var(--radius-sm);
  color: #475569;
}

.q-opt.correct {
  background: #ecfdf5;
  border-color: #86efac;
  color: #15803d;
  font-weight: 700;
}

.q-explanation {
  font-size: 0.8rem;
  color: #b45309;
  margin-top: 0.4rem;
}

.submissions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1rem;
}

.submission-card {
  padding: 1.1rem;
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.sub-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.student-meta {
  display: flex;
  flex-direction: column;
}

.student-name {
  font-weight: 800;
  font-size: 0.98rem;
  color: #0f172a;
}

.student-nim {
  font-size: 0.82rem;
  color: #2563eb;
  font-weight: 700;
}

.sub-detail-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
  padding-bottom: 0.4rem;
  border-bottom: 1px dashed #cbd5e1;
}

.sub-score-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.2rem;
}

.score-label {
  font-size: 0.85rem;
  font-weight: 700;
  color: #334155;
}

.score-value {
  font-size: 1.3rem;
  font-weight: 800;
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: #64748b;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px dashed #cbd5e1;
  font-size: 0.9rem;
}

.text-amber {
  color: #d97706;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}

.btn-icon-sm {
  width: 16px;
  height: 16px;
}

/* Modal Styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-card {
  width: 100%;
  max-width: 520px;
  padding: 2rem;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  box-shadow: 0 25px 50px -12px rgba(15, 23, 42, 0.25);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1.5px solid #f1f5f9;
}

.m-title-box {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.modal-title-icon {
  width: 24px;
  height: 24px;
  margin-top: 0.15rem;
  flex-shrink: 0;
}

.modal-header h3 {
  font-size: 1.3rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 0.2rem 0;
}

.modal-sub {
  font-size: 0.85rem;
  color: #64748b;
  margin: 0;
}

.btn-close-modal {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.btn-close-modal:hover {
  background: #e2e8f0;
  color: #0f172a;
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
  margin-top: 1.25rem;
  padding-top: 1rem;
  border-top: 1.5px solid #f1f5f9;
}

@media (max-width: 768px) {
  .q-options-grid {
    grid-template-columns: 1fr;
  }
}
</style>
