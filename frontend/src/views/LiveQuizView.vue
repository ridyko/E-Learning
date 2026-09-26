<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showSuccess, showError, showWarning } from '../utils/swal.js'
import { 
  Play, 
  Users, 
  Award, 
  CheckCircle2, 
  XCircle, 
  Clock, 
  ChevronRight, 
  Trophy, 
  Sparkles, 
  Flame, 
  BarChart2, 
  ArrowLeft, 
  Copy, 
  HelpCircle,
  Radio,
  Zap,
  Crown,
  BookOpen,
  Plus,
  FileText
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

// User & State
const currentUser = ref(null)
const inputPin = ref('')
const inputStudentName = ref('')
const inputStudentNim = ref('')

const activeSession = ref(null)
const currentPin = ref('')
const mode = ref('join') // 'dosen-setup' | 'join' | 'host' | 'player'
const isPollerActive = ref(false)
let syncTimer = null
let questionCountdownTimer = null

// Dosen Setup State
const dosenSetupTab = ref('existing') // 'existing' | 'new_live'
const courses = ref([])
const quizzes = ref([])
const selectedCourseId = ref('')
const selectedQuizId = ref('')
const selectedTimeLimit = ref(20)

// State for Creating Instant New Live Quiz
const newLiveTitle = ref('')
const newLiveQuestions = ref([
  { question: '', options: ['', '', '', ''], correct_answer: 0, explanation: '' }
])

// Quiz Playing State
const timeLeft = ref(20)
const hasVoted = ref(false)
const selectedOpt = ref(null)
const voteTimeSec = ref(0)
const startTimeMs = ref(0)

// Check current user
const checkUser = () => {
  const u = localStorage.getItem('user')
  if (u) {
    try {
      currentUser.value = JSON.parse(u)
      if (currentUser.value.role === 'student' || currentUser.value.role === 'mahasiswa') {
        inputStudentName.value = currentUser.value.name || 'Mahasiswa'
        inputStudentNim.value = currentUser.value.username || ''
      } else if (currentUser.value.role === 'dosen') {
        inputStudentName.value = currentUser.value.name || 'Dosen Pengampu'
        inputStudentNim.value = currentUser.value.username || '21099001'
      }
    } catch {
      currentUser.value = null
    }
  }
}

const isDosen = computed(() => currentUser.value?.role === 'dosen')

const courseQuizzes = computed(() => {
  if (!selectedCourseId.value) return quizzes.value
  return quizzes.value.filter(q => q.course_id === selectedCourseId.value)
})

const fetchDosenSetupData = async () => {
  try {
    const [resC, resQ] = await Promise.all([
      fetch('/api/courses'),
      fetch('/api/quizzes')
    ])
    courses.value = await resC.json()
    quizzes.value = await resQ.json()

    if (courses.value.length > 0) {
      selectedCourseId.value = courses.value[0].id
    }
    if (quizzes.value.length > 0) {
      selectedQuizId.value = quizzes.value[0].id
    }
  } catch (err) {
    console.warn('Fetch setup data error:', err)
  }
}

watch(selectedCourseId, (newId) => {
  const firstQ = quizzes.value.find(q => q.course_id === newId)
  if (firstQ) {
    selectedQuizId.value = firstQ.id
  }
})

// Route & Mode Initialization
onMounted(async () => {
  checkUser()
  await fetchDosenSetupData()

  const pinFromUrl = route.params.pin || route.query.pin
  const isHostQuery = route.query.host === 'true'

  if (pinFromUrl) {
    currentPin.value = pinFromUrl
    if (isHostQuery && isDosen.value) {
      mode.value = 'host'
    } else {
      mode.value = 'player'
    }
    syncSession()
    startPoller()
  } else {
    if (isDosen.value) {
      mode.value = 'dosen-setup'
    } else {
      mode.value = 'join'
    }
  }
})

onUnmounted(() => {
  stopPoller()
  clearInterval(questionCountdownTimer)
})

// Add New Question row in instant live setup
const addInstantQuestionRow = () => {
  newLiveQuestions.value.push({
    question: '',
    options: ['', '', '', ''],
    correct_answer: 0,
    explanation: ''
  })
}

// Remove Question row
const removeInstantQuestionRow = (idx) => {
  if (newLiveQuestions.value.length > 1) {
    newLiveQuestions.value.splice(idx, 1)
  }
}

// Dosen Create Live Session
const handleDosenCreateLive = async () => {
  let targetQuizId = selectedQuizId.value

  if (dosenSetupTab.value === 'new_live') {
    if (!newLiveTitle.value) {
      showWarning('Judul Wajib Diisi', 'Masukkan judul Sesi Live Quiz Baru.')
      return
    }
    // Create new Quiz Package first
    try {
      const resQuiz = await fetch('/api/quizzes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          course_id: selectedCourseId.value,
          title: newLiveTitle.value + ' (Live)',
          description: 'Sesi Live Quiz Interaktif Perkuliahan',
          time_limit: selectedTimeLimit.value || 20
        })
      })
      if (!resQuiz.ok) {
        showError('Gagal', 'Gagal membuat paket kuis baru.')
        return
      }
      const createdQuiz = await resQuiz.json()
      targetQuizId = createdQuiz.id

      // Add questions
      for (const q of newLiveQuestions.value) {
        if (q.question) {
          await fetch('/api/quizzes/questions', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              quiz_id: targetQuizId,
              question: {
                type: 'mc',
                question: q.question,
                options: q.options,
                correct_answer: parseInt(q.correct_answer),
                explanation: q.explanation || 'Pembahasan kunci jawaban oleh Dosen Pengampu.'
              }
            })
          })
        }
      }
    } catch (err) {
      showError('Gagal', 'Gagal menyimpan pertanyaan kuis baru.')
      return
    }
  }

  if (!targetQuizId) {
    showWarning('Pilih Paket Kuis', 'Silakan pilih paket kuis yang ingin diluncurkan.')
    return
  }

  try {
    const res = await fetch('/api/live-quiz/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        quiz_id: targetQuizId,
        time_per_question: parseInt(selectedTimeLimit.value) || 20
      })
    })

    if (!res.ok) {
      const errTxt = await res.text()
      showError('Gagal Peluncuran', errTxt)
      return
    }

    const session = await res.json()
    activeSession.value = session
    currentPin.value = session.pin
    mode.value = 'host'
    showSuccess('Sesi Live Quiz Berhasil Dibuat! 🚀', `PIN Kuis: ${session.pin}. Layar Presenter Siap!`)
    startPoller()
  } catch (err) {
    showError('Terjadi Kesalahan', 'Gagal menghubungkan ke server live kuis.')
  }
}

// Sync function (Polled every 1s)
const syncSession = async () => {
  if (!currentPin.value) return
  try {
    const res = await fetch(`/api/live-quiz/sync?pin=${currentPin.value}`)
    if (res.ok) {
      const data = await res.json()
      const prevStatus = activeSession.value?.status
      const prevQIndex = activeSession.value?.current_q_index

      activeSession.value = data

      if (data.status === 'QUESTION' && (prevStatus !== 'QUESTION' || prevQIndex !== data.current_q_index)) {
        hasVoted.value = false
        selectedOpt.value = null
        startTimeMs.value = Date.now()
        timeLeft.value = data.time_per_question || 20
        startCountdown()
      }
    } else if (res.status === 404) {
      stopPoller()
      showError('Sesi Tidak Ditemukan', 'Sesi kuis live ini telah berakhir atau PIN salah.')
      if (isDosen.value) {
        mode.value = 'dosen-setup'
      } else {
        mode.value = 'join'
      }
    }
  } catch (err) {
    console.warn('Sync live quiz error:', err)
  }
}

const startPoller = () => {
  stopPoller()
  isPollerActive.value = true
  syncTimer = setInterval(syncSession, 1000)
}

const stopPoller = () => {
  if (syncTimer) {
    clearInterval(syncTimer)
    syncTimer = null
  }
  isPollerActive.value = false
}

// Countdown timer
const startCountdown = () => {
  clearInterval(questionCountdownTimer)
  questionCountdownTimer = setInterval(() => {
    if (activeSession.value?.status === 'QUESTION') {
      if (timeLeft.value > 0) {
        timeLeft.value--
      } else {
        clearInterval(questionCountdownTimer)
        if (mode.value === 'host' && activeSession.value?.status === 'QUESTION') {
          updateSessionState('REVEAL', activeSession.value.current_q_index)
        }
      }
    }
  }, 1000)
}

// Join Quiz via PIN (Mahasiswa)
const handleJoinQuiz = async () => {
  if (!inputPin.value || inputPin.value.length < 6) {
    showWarning('PIN Tidak Lengkap', 'Masukkan 6 digit PIN sesi kuis.')
    return
  }

  const nim = inputStudentNim.value || ('GUEST-' + Math.floor(1000 + Math.random() * 9000))
  const name = inputStudentName.value || 'Peserta Mahasiswa'

  try {
    const res = await fetch('/api/live-quiz/join', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        pin: inputPin.value,
        student_nim: nim,
        student_name: name
      })
    })

    if (!res.ok) {
      const errText = await res.text()
      showError('Gagal Bergabung', errText)
      return
    }

    const data = await res.json()
    activeSession.value = data
    currentPin.value = data.pin
    mode.value = 'player'
    showSuccess('Berhasil Bergabung! 🚀', `Selamat datang, ${name}! Anda telah masuk ke lobby kuis.`)
    startPoller()
  } catch (err) {
    showError('Terjadi Kesalahan', 'Gagal menghubungkan ke server kuis.')
  }
}

// Host controls
const updateSessionState = async (nextStatus, qIndex = 0) => {
  if (!currentPin.value) return
  try {
    const res = await fetch('/api/live-quiz/state', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        pin: currentPin.value,
        status: nextStatus,
        current_q_index: qIndex
      })
    })

    if (res.ok) {
      const updated = await res.json()
      activeSession.value = updated
    }
  } catch (err) {
    console.error('Update session status failed:', err)
  }
}

const handleStartQuiz = () => {
  updateSessionState('QUESTION', 0)
}

const handleNextQuestion = () => {
  if (!activeSession.value) return
  const nextIdx = activeSession.value.current_q_index + 1
  if (nextIdx < activeSession.value.quiz.questions.length) {
    updateSessionState('QUESTION', nextIdx)
  } else {
    updateSessionState('FINISHED', activeSession.value.current_q_index)
  }
}

const handleShowReveal = () => {
  updateSessionState('REVEAL', activeSession.value.current_q_index)
}

const handleShowLeaderboard = () => {
  updateSessionState('LEADERBOARD', activeSession.value.current_q_index)
}

// Player submit vote
const handleVoteOption = async (optIdx) => {
  if (hasVoted.value || activeSession.value?.status !== 'QUESTION') return

  const elapsedSec = (Date.now() - startTimeMs.value) / 1000.0
  voteTimeSec.value = elapsedSec.toFixed(1)
  selectedOpt.value = optIdx
  hasVoted.value = true

  const nim = inputStudentNim.value || currentUser.value?.username || 'GUEST-0000'

  try {
    const res = await fetch('/api/live-quiz/vote', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        pin: currentPin.value,
        student_nim: nim,
        question_index: activeSession.value.current_q_index,
        selected_opt: optIdx,
        answered_in_sec: elapsedSec
      })
    })

    if (res.ok) {
      const updated = await res.json()
      activeSession.value = updated
    }
  } catch (err) {
    console.error('Vote submission failed:', err)
  }
}

const copyPin = () => {
  if (!currentPin.value) return
  navigator.clipboard.writeText(currentPin.value)
  showSuccess('PIN Disalin!', `PIN Kuis: ${currentPin.value}`)
}

// Computed metrics
const participantsList = computed(() => {
  if (!activeSession.value?.participants) return []
  return Object.values(activeSession.value.participants).sort((a, b) => b.total_score - a.total_score)
})

const participantCount = computed(() => participantsList.value.length)

const currentQuestion = computed(() => {
  if (!activeSession.value?.quiz?.questions) return null
  return activeSession.value.quiz.questions[activeSession.value.current_q_index] || null
})

const voteDistribution = computed(() => {
  if (!activeSession.value?.votes || !currentQuestion.value) return [0, 0, 0, 0]
  const dist = [0, 0, 0, 0]
  const curQIdx = activeSession.value.current_q_index

  activeSession.value.votes.forEach(v => {
    if (v.question_index === curQIdx && v.selected_opt >= 0 && v.selected_opt <= 3) {
      dist[v.selected_opt]++
    }
  })
  return dist
})

const totalVotesForCurrentQ = computed(() => {
  return voteDistribution.value.reduce((a, b) => a + b, 0)
})

const myParticipantData = computed(() => {
  if (!activeSession.value?.participants || !inputStudentNim.value) return null
  return activeSession.value.participants[inputStudentNim.value] || null
})

const myRank = computed(() => {
  if (!inputStudentNim.value || participantsList.value.length === 0) return 0
  const idx = participantsList.value.findIndex(p => p.student_nim === inputStudentNim.value)
  return idx !== -1 ? idx + 1 : 0
})

// SOLID COLORS INSTEAD OF GRADIENTS
const optionThemes = [
  { label: 'A', name: 'Merah', color: '#ef4444', gradient: '#ef4444', shape: '▲' },
  { label: 'B', name: 'Biru', color: '#2563eb', gradient: '#2563eb', shape: '◆' },
  { label: 'C', name: 'Kuning', color: '#d97706', gradient: '#d97706', shape: '●' },
  { label: 'D', name: 'Hijau', color: '#059669', gradient: '#059669', shape: '■' }
]

const exitQuiz = async () => {
  stopPoller()
  if (isDosen.value) {
    router.push('/bank-soal')
  } else {
    router.push('/quizzes')
  }
}
</script>

<template>
  <div class="live-quiz-page light-theme">
    <!-- Top In-Page Header Toolbar -->
    <header class="live-header-toolbar glass-card-light">
      <div class="header-left">
        <button @click="exitQuiz" class="btn-exit-light" title="Kembali">
          <ArrowLeft class="icon-xs" />
          <span>Kembali</span>
        </button>
        <div class="live-badge-light">
          <Radio class="icon-pulse text-amber" />
          <span>LIVE QUIZ REAL-TIME</span>
        </div>
      </div>

      <div class="header-center" v-if="activeSession">
        <div class="title-pin-pill-light" @click="copyPin" title="Klik untuk salin PIN">
          <span class="pin-lbl">PIN:</span>
          <span class="pin-code-light">{{ activeSession.pin }}</span>
          <Copy class="icon-xs copy-icon" />
        </div>
        <span class="quiz-title-tag-light">{{ activeSession.title }}</span>
      </div>

      <div class="header-right" v-if="activeSession">
        <div class="participants-pill-light">
          <Users class="icon-xs" />
          <span>{{ participantCount }} Peserta</span>
        </div>
      </div>
    </header>

    <!-- ========================================================= -->
    <!-- MODE DOSEN SETUP: PENGATURAN & PELUNCUR SESI LIVE (DOSEN) -->
    <!-- ========================================================= -->
    <div v-if="mode === 'dosen-setup'" class="live-card-container animate-fade-in">
      <div class="setup-card glass-card-light">
        <div class="card-title-header">
          <div class="brand-logo-circle-gold">
            <Radio class="logo-icon text-white animate-pulse" />
          </div>
          <h2>Panel Peluncur Live Quiz Interaktif (Dosen)</h2>
          <p>Pak <strong>Rio Widyatmoko, S.Kom, M.M.S.I</strong> — Silakan tentukan mode input soal untuk ditampilkan pada layar proyektor perkuliahan.</p>
        </div>

        <!-- Dosen Setup Mode Tabs -->
        <div class="dosen-tab-selector">
          <button 
            type="button"
            @click="dosenSetupTab = 'existing'"
            :class="['tab-choice-btn', dosenSetupTab === 'existing' ? 'active' : '']"
          >
            <BookOpen class="icon-xs" />
            <span>Paket Kuis Online Ada</span>
          </button>
          <button 
            type="button"
            @click="dosenSetupTab = 'new_live'"
            :class="['tab-choice-btn', dosenSetupTab === 'new_live' ? 'active' : '']"
          >
            <Plus class="icon-xs" />
            <span>Input Sesi Live Baru</span>
          </button>
        </div>

        <form @submit.prevent="handleDosenCreateLive" class="setup-form">
          <!-- TAB A: Existing Quiz Package -->
          <div v-if="dosenSetupTab === 'existing'">
            <div class="form-group">
              <label><BookOpen class="icon-xs" /> Pilih Mata Kuliah</label>
              <select v-model="selectedCourseId" class="select-input-light" required>
                <option v-for="c in courses" :key="c.id" :value="c.id">
                  {{ c.name }} ({{ c.code }})
                </option>
              </select>
            </div>

            <div class="form-group mt-3">
              <label><HelpCircle class="icon-xs" /> Pilih Paket Soal Kuis Online</label>
              <select v-model="selectedQuizId" class="select-input-light" required>
                <option v-for="q in courseQuizzes" :key="q.id" :value="q.id">
                  {{ q.title }} ({{ q.questions?.length || 0 }} Soal)
                </option>
              </select>
            </div>
          </div>

          <!-- TAB B: Input New Live Quiz Session Directly -->
          <div v-else class="new-live-form-box">
            <div class="form-group">
              <label><BookOpen class="icon-xs" /> Pilih Mata Kuliah Target</label>
              <select v-model="selectedCourseId" class="select-input-light" required>
                <option v-for="c in courses" :key="c.id" :value="c.id">
                  {{ c.name }} ({{ c.code }})
                </option>
              </select>
            </div>

            <div class="form-group mt-3">
              <label><FileText class="icon-xs" /> Judul Sesi Live Quiz *</label>
              <input 
                v-model="newLiveTitle" 
                type="text" 
                placeholder="Contoh: Kuis Interaktif Pertemuan #5 — Agile Scrum" 
                class="text-input-light" 
                required 
              />
            </div>

            <!-- Instant Questions Input list -->
            <div class="instant-q-section">
              <label class="section-lbl">Daftar Pertanyaan Live Quiz (Pilihan Ganda):</label>

              <div 
                v-for="(qItem, qIdx) in newLiveQuestions" 
                :key="qIdx" 
                class="q-input-card-light"
              >
                <div class="q-card-head">
                  <span class="q-num-lbl">Pertanyaan #{{ qIdx + 1 }}</span>
                  <button 
                    v-if="newLiveQuestions.length > 1" 
                    type="button" 
                    @click="removeInstantQuestionRow(qIdx)" 
                    class="btn-remove-q"
                  >
                    Hapus
                  </button>
                </div>

                <input 
                  v-model="qItem.question" 
                  type="text" 
                  placeholder="Tuliskan teks pertanyaan di sini..." 
                  class="text-input-light mb-2" 
                  required 
                />

                <div class="opts-grid-2x2">
                  <input v-model="qItem.options[0]" type="text" placeholder="Opsi A" class="opt-input-light" required />
                  <input v-model="qItem.options[1]" type="text" placeholder="Opsi B" class="opt-input-light" required />
                  <input v-model="qItem.options[2]" type="text" placeholder="Opsi C" class="opt-input-light" required />
                  <input v-model="qItem.options[3]" type="text" placeholder="Opsi D" class="opt-input-light" required />
                </div>

                <div class="key-select-row">
                  <span>Kunci Jawaban Benar:</span>
                  <select v-model="qItem.correct_answer" class="select-key-light">
                    <option :value="0">Opsi A</option>
                    <option :value="1">Opsi B</option>
                    <option :value="2">Opsi C</option>
                    <option :value="3">Opsi D</option>
                  </select>
                </div>
              </div>

              <button type="button" @click="addInstantQuestionRow" class="btn-add-row-light">
                <Plus class="icon-xs" /> Tambah Pertanyaan Lain
              </button>
            </div>
          </div>

          <!-- Time Limit Per Question -->
          <div class="form-group mt-3">
            <label><Clock class="icon-xs" /> Alokasi Waktu per Soal (Detik)</label>
            <div class="time-options-grid-light">
              <button 
                v-for="t in [15, 20, 30, 45, 60]" 
                :key="t"
                type="button"
                @click="selectedTimeLimit = t"
                :class="['time-opt-btn-light', selectedTimeLimit === t ? 'active' : '']"
              >
                {{ t }}s
              </button>
            </div>
          </div>

          <button type="submit" class="btn btn-launch-live-gold">
            <Radio class="icon-sm animate-pulse" />
            <span>LUNCURKAN LIVE QUIZ (LAYAR PROYEKTOR) 🚀</span>
          </button>
        </form>

        <div class="join-existing-divider-light">
          <span>atau pantau sesi kuis yang sedang berjalan</span>
        </div>

        <div class="join-existing-row">
          <input 
            v-model="inputPin" 
            type="text" 
            maxlength="6" 
            placeholder="Masukkan PIN (6 Digit)" 
            class="pin-short-input-light" 
          />
          <button @click="handleJoinQuiz" class="btn btn-secondary-light">
            <span>Buka Presenter</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ========================================================= -->
    <!-- MODE MAHASISWA JOIN: GABUNG DENGAN PIN (MAHASISWA)       -->
    <!-- ========================================================= -->
    <div v-else-if="mode === 'join'" class="live-card-container animate-fade-in">
      <div class="join-card glass-card-light">
        <div class="join-header">
          <div class="brand-logo-circle-gold">
            <Zap class="logo-icon text-white" />
          </div>
          <h2>Gabung Live Quiz Interaktif (Mahasiswa)</h2>
          <p>Masukkan 6 digit PIN yang ditampilkan Dosen pada layar proyektor.</p>
        </div>

        <form @submit.prevent="handleJoinQuiz" class="join-form">
          <div class="form-group">
            <label>PIN Kuis Live (6 Digit)</label>
            <input 
              v-model="inputPin" 
              type="text" 
              maxlength="6" 
              placeholder="Contoh: 470077" 
              class="pin-input-light" 
              required 
            />
          </div>

          <div class="form-group">
            <label>Nama Lengkap</label>
            <input 
              v-model="inputStudentName" 
              type="text" 
              placeholder="Masukkan nama Anda" 
              class="text-input-light" 
              required 
            />
          </div>

          <div class="form-group">
            <label>NIM (Nomor Induk Mahasiswa)</label>
            <input 
              v-model="inputStudentNim" 
              type="text" 
              placeholder="Contoh: 20260801001" 
              class="text-input-light" 
            />
          </div>

          <button type="submit" class="btn btn-live-join-gold">
            <span>MASUK SESI KUIS 🚀</span>
          </button>
        </form>

        <div class="join-footer-light">
          <p>💡 <em>Tips: Menjawab lebih cepat memberikan bonus poin ekstra!</em></p>
        </div>
      </div>
    </div>

    <!-- ========================================================= -->
    <!-- ACTIVE SESSION DISPLAY (Host Presenter & Player Views)   -->
    <!-- ========================================================= -->
    <div v-else-if="activeSession" class="live-active-container">

      <!-- ================= LOBBY STATE ================= -->
      <div v-if="activeSession.status === 'LOBBY'" class="state-lobby animate-fade-in">
        <div class="lobby-banner glass-card-light">
          <div class="pin-display-box-light" @click="copyPin">
            <span class="pin-sub-gold">PIN GABUNG KUIS LIVE</span>
            <h1 class="pin-large-dark">{{ activeSession.pin }}</h1>
            <p class="pin-hint-light"><Copy class="icon-xs" /> Klik PIN untuk menyalin PIN kuis</p>
          </div>

          <div class="lobby-meta-light">
            <h3>{{ activeSession.title }}</h3>
            <p><Clock class="icon-xs" /> Alokasi Waktu: <strong>{{ activeSession.time_per_question }} Detik per Soal</strong></p>

            <div class="lobby-status-pill-gold">
              <span class="dot-pulse-amber"></span>
              <span>Menunggu Peserta Bergabung... ({{ participantCount }} Orang Terhubung)</span>
            </div>
          </div>
        </div>

        <!-- Participant Grid -->
        <div class="participants-section glass-card-light">
          <div class="section-title-dark">
            <Users class="icon-sm text-gold" />
            <h4>Peserta Terhubung ({{ participantCount }})</h4>
          </div>

          <div v-if="participantsList.length === 0" class="empty-participants-light">
            <p>Belum ada mahasiswa yang bergabung. Minta mahasiswa membuka menu <strong>Live Quiz Interaktif</strong> dan memasukkan PIN <strong>{{ activeSession.pin }}</strong></p>
          </div>

          <div v-else class="participants-grid">
            <div 
              v-for="p in participantsList" 
              :key="p.student_nim" 
              class="participant-chip-light animate-pop"
            >
              <div class="avatar-circle-blue">
                {{ p.student_name ? p.student_name.charAt(0).toUpperCase() : 'M' }}
              </div>
              <span class="p-name-dark">{{ p.student_name }}</span>
            </div>
          </div>
        </div>

        <!-- Host Controls (Dosen) -->
        <div v-if="mode === 'host' || isDosen" class="host-action-bar">
          <button @click="handleStartQuiz" class="btn btn-start-live-emerald">
            <Play class="icon-sm fill-white" />
            <span>MULAI KUIS SEKARANG {{ participantCount > 0 ? '(' + participantCount + ' Peserta Terhubung)' : '🚀' }}</span>
          </button>
        </div>

        <!-- Student Waiting View -->
        <div v-else class="student-waiting-box glass-card-light">
          <Sparkles class="sparkle-icon text-gold animate-bounce" />
          <h3 class="text-dark">Anda Sudah Masuk ke Lobby!</h3>
          <p class="text-muted">Bersiaplah! Sesi akan segera dimulai saat Dosen Pengampu menekan tombol "Mulai Kuis".</p>
        </div>
      </div>

      <!-- ================= QUESTION STATE ================= -->
      <div v-else-if="activeSession.status === 'QUESTION'" class="state-question animate-fade-in">
        <!-- Timer Bar Header -->
        <div class="question-header">
          <div class="progress-bar-container-light">
            <div 
              class="progress-bar-fill" 
              :style="{ width: (timeLeft / (activeSession.time_per_question || 20)) * 100 + '%' }"
              :class="{ 'warning-time': timeLeft <= 5 }"
            ></div>
          </div>

          <div class="timer-badge-gold" :class="{ 'urgent': timeLeft <= 5 }">
            <Clock class="icon-sm" />
            <span>{{ timeLeft }} Detik</span>
          </div>

          <div class="q-number-badge-light">
            Soal {{ activeSession.current_q_index + 1 }} dari {{ activeSession.quiz?.questions?.length || 1 }}
          </div>
        </div>

        <!-- Question Prompt -->
        <div class="question-card glass-card-light">
          <div class="q-badge-mini">PERTANYAAN #{{ activeSession.current_q_index + 1 }}</div>
          <h2 class="question-text-dark">{{ currentQuestion?.question }}</h2>
        </div>

        <!-- Host Presenter Screen (Dosen) -->
        <div v-if="mode === 'host' || isDosen" class="host-question-screen">
          <div class="votes-counter-card glass-card-light">
            <BarChart2 class="icon-md text-gold" />
            <div class="counter-text">
              <h3>{{ totalVotesForCurrentQ }} / {{ participantCount }} Mahasiswa</h3>
              <p>Telah mengirimkan jawaban</p>
            </div>
          </div>

          <div class="options-grid-presenter">
            <div 
              v-for="(opt, idx) in currentQuestion?.options" 
              :key="idx"
              class="option-card-presenter"
              :style="{ background: optionThemes[idx].gradient }"
            >
              <div class="shape-badge-player">
                <span>{{ optionThemes[idx].shape }}</span>
              </div>
              <span class="opt-text">{{ optionThemes[idx].label }}. {{ opt }}</span>
            </div>
          </div>

          <div class="host-controls-row">
            <button @click="handleShowReveal" class="btn btn-reveal-action">
              <CheckCircle2 class="icon-sm" />
              <span>TUTUP SOAL & TAMPILKAN JAWABAN</span>
            </button>
          </div>
        </div>

        <!-- Student Controller Screen (Mahasiswa) -->
        <div v-else class="student-question-screen">
          <div v-if="!hasVoted" class="options-grid-player">
            <button 
              v-for="(opt, idx) in currentQuestion?.options" 
              :key="idx"
              @click="handleVoteOption(idx)"
              class="btn-option-player"
              :style="{ background: optionThemes[idx].gradient }"
            >
              <div class="shape-badge-player">
                <span>{{ optionThemes[idx].shape }}</span>
              </div>
              <span class="opt-label-player">{{ optionThemes[idx].label }}. {{ opt }}</span>
            </button>
          </div>

          <div v-else class="voted-confirmation-box glass-card-light animate-pop">
            <div class="check-icon-circle-lg">
              <CheckCircle2 class="icon-lg text-white" />
            </div>
            <h3 class="text-dark">Jawaban Terkirim! ⚡</h3>
            <p class="text-muted">Anda telah memilih <strong>Option {{ optionThemes[selectedOpt]?.label }}</strong> dalam waktu <strong>{{ voteTimeSec }} Detik</strong>.</p>
            <div class="waiting-reveal-pill-light">
              <span class="dot-pulse-amber"></span>
              <span>Menunggu Dosen membuka kunci jawaban...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- ================= REVEAL STATE ================= -->
      <div v-else-if="activeSession.status === 'REVEAL'" class="state-reveal animate-fade-in">
        <div class="reveal-header">
          <h2 class="text-dark">Hasil Jawaban — Soal {{ activeSession.current_q_index + 1 }}</h2>
          <p class="question-sub-dark">{{ currentQuestion?.question }}</p>
        </div>

        <!-- Interactive Bar Graph -->
        <div class="chart-container glass-card-light">
          <div class="chart-bars">
            <div 
              v-for="(opt, idx) in currentQuestion?.options" 
              :key="idx"
              class="chart-bar-col"
            >
              <div class="bar-count-wrapper">
                <span v-if="idx === currentQuestion?.correct_answer" class="correct-badge-floating animate-bounce">✓</span>
                <span class="bar-count-dark">{{ voteDistribution[idx] }}</span>
              </div>

              <div class="bar-track-light">
                <div 
                  class="bar-fill" 
                  :style="{ 
                    height: totalVotesForCurrentQ > 0 ? Math.max((voteDistribution[idx] / totalVotesForCurrentQ * 100), 8) + '%' : '0%',
                    background: optionThemes[idx].gradient 
                  }"
                  :class="{ 'is-correct-bar': idx === currentQuestion?.correct_answer }"
                >
                </div>
              </div>

              <div class="bar-label-dark" :class="{ 'is-correct-text': idx === currentQuestion?.correct_answer }">
                <div class="option-shape-pill" :style="{ background: optionThemes[idx].gradient }">
                  <span class="shape-icon-sm">{{ optionThemes[idx].shape }}</span>
                  <span class="label-txt">{{ optionThemes[idx].label }}</span>
                </div>
                <div v-if="idx === currentQuestion?.correct_answer" class="correct-pill-tag">
                  <span>✓ BENAR</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Explanation Card -->
        <div class="explanation-card glass-card-light">
          <div class="explanation-header-box">
            <div class="check-icon-circle">
              <CheckCircle2 class="icon-sm text-white" />
            </div>
            <div class="explanation-title-text">
              <span class="exp-subtitle">KUNCI JAWABAN BENAR</span>
              <h4>Option {{ optionThemes[currentQuestion?.correct_answer]?.label }}: {{ currentQuestion?.options[currentQuestion?.correct_answer] }}</h4>
            </div>
          </div>
          <div v-if="currentQuestion?.explanation" class="explanation-body-box">
            <p>💡 <strong>Pembahasan:</strong> {{ currentQuestion.explanation }}</p>
          </div>
        </div>

        <!-- Personal Feedback (for Student) -->
        <div v-if="mode === 'player' && hasVoted" class="personal-feedback-card glass-card-light">
          <div v-if="selectedOpt === currentQuestion?.correct_answer" class="feedback-box correct">
            <div class="feedback-icon-wrap gold">
              <Sparkles class="icon-md text-white animate-bounce" />
            </div>
            <div class="feedback-text-content">
              <h3>Jawaban Anda Benar! 🎉</h3>
              <p>Waktu menjawab: <strong>{{ voteTimeSec }} Detik</strong> • Bonus poin bertambah!</p>
            </div>
          </div>
          <div v-else class="feedback-box wrong">
            <div class="feedback-icon-wrap red">
              <XCircle class="icon-md text-white" />
            </div>
            <div class="feedback-text-content">
              <h3>Jawaban Kurang Tepat!</h3>
              <p>Pilihan Anda: <strong>Option {{ optionThemes[selectedOpt]?.label }}</strong> • Kunci Jawaban: <strong>Option {{ optionThemes[currentQuestion?.correct_answer]?.label }}</strong></p>
            </div>
          </div>
        </div>

        <!-- Host Next Controls (Dosen) -->
        <div v-if="mode === 'host' || isDosen" class="host-action-bar">
          <button @click="handleShowLeaderboard" class="btn btn-leaderboard">
            <Trophy class="icon-sm text-gold" />
            <span>TAMPILKAN PAPAN PERINGKAT 🏆</span>
          </button>
        </div>
      </div>

      <!-- ================= LEADERBOARD STATE ================= -->
      <div v-else-if="activeSession.status === 'LEADERBOARD'" class="state-leaderboard animate-fade-in">
        <div class="leaderboard-header">
          <div class="trophy-circle-lg">
            <Trophy class="trophy-icon text-white animate-bounce" />
          </div>
          <h2 class="text-dark">Papan Peringkat Sementara</h2>
          <p class="text-muted">Top Skor Mahasiswa — Sesi Live Quiz</p>
        </div>

        <!-- Leaderboard List -->
        <div class="leaderboard-list glass-card-light">
          <div 
            v-for="(p, idx) in participantsList.slice(0, 10)" 
            :key="p.student_nim"
            class="leaderboard-row-light"
            :class="{ 
              'rank-1-light': idx === 0, 
              'rank-2-light': idx === 1, 
              'rank-3-light': idx === 2,
              'is-me': p.student_nim === inputStudentNim
            }"
          >
            <div 
              class="rank-badge-dark"
              :class="{ 'gold': idx === 0, 'silver': idx === 1, 'bronze': idx === 2 }"
            >
              <Crown v-if="idx === 0" class="crown-icon text-amber" />
              <span v-else>#{{ idx + 1 }}</span>
            </div>

            <div class="player-info">
              <div 
                class="avatar-small"
                :class="{ 'gold': idx === 0, 'silver': idx === 1, 'bronze': idx === 2 }"
              >
                {{ p.student_name ? p.student_name.charAt(0).toUpperCase() : 'M' }}
              </div>
              <div class="player-details">
                <span class="player-name-dark">{{ p.student_name }}</span>
                <span v-if="p.streak_count >= 2" class="streak-badge">
                  <Flame class="icon-xs text-orange" /> {{ p.streak_count }}x Streak!
                </span>
              </div>
            </div>

            <div class="player-score">
              <span class="score-val-dark">{{ p.total_score }}</span>
              <span class="pts-lbl-muted">Poin</span>
            </div>
          </div>
        </div>

        <!-- Personal Rank Card for Student -->
        <div v-if="mode === 'player' && myParticipantData" class="personal-rank-footer glass-card-light">
          <div class="rank-footer-left">
            <div class="my-rank-badge">#{{ myRank }}</div>
            <div class="my-rank-text">
              <h4>Posisi Anda Saat Ini</h4>
              <p>Dari total {{ participantCount }} Mahasiswa</p>
            </div>
          </div>
          <div class="rank-footer-right">
            <span class="my-score-num">{{ myParticipantData.total_score }}</span>
            <span class="my-score-lbl">Total Poin</span>
          </div>
        </div>

        <!-- Host Controls (Dosen) -->
        <div v-if="mode === 'host' || isDosen" class="host-action-bar">
          <button 
            v-if="activeSession.current_q_index + 1 < activeSession.quiz?.questions?.length" 
            @click="handleNextQuestion" 
            class="btn btn-next-q"
          >
            <span>SOAL BERIKUTNYA (#{{ activeSession.current_q_index + 2 }})</span>
            <ChevronRight class="icon-sm" />
          </button>
          <button 
            v-else 
            @click="updateSessionState('FINISHED', activeSession.current_q_index)" 
            class="btn btn-finish-quiz"
          >
            <Trophy class="icon-sm text-gold" />
            <span>SELESAIKAN KUIS & SELEBRASI PODIUM 🏁</span>
          </button>
        </div>
      </div>

      <!-- ================= FINISHED PODIUM STATE ================= -->
      <div v-else-if="activeSession.status === 'FINISHED'" class="state-finished animate-fade-in">
        <div class="finished-header">
          <div class="trophy-circle-lg gold">
            <Sparkles class="icon-lg text-white animate-bounce" />
          </div>
          <h2 class="text-dark">SELEBRASI JUARA KUIS LIVE</h2>
          <p class="text-muted">Selamat kepada para pemenang Sesi Live Quiz Real-Time!</p>
        </div>

        <!-- 3D Podium Display -->
        <div class="podium-container">
          <!-- 2nd Place -->
          <div v-if="participantsList[1]" class="podium-step step-2 animate-pop">
            <div class="avatar-podium silver">
              {{ participantsList[1].student_name.charAt(0) }}
            </div>
            <div class="podium-name-dark">{{ participantsList[1].student_name }}</div>
            <div class="podium-score-dark">{{ participantsList[1].total_score }} Poin</div>
            <div class="podium-block block-2-light">2</div>
          </div>

          <!-- 1st Place -->
          <div v-if="participantsList[0]" class="podium-step step-1 animate-pop">
            <Crown class="crown-winner text-gold" />
            <div class="avatar-podium gold">
              {{ participantsList[0].student_name.charAt(0) }}
            </div>
            <div class="podium-name-dark winner">{{ participantsList[0].student_name }}</div>
            <div class="podium-score-dark winner">{{ participantsList[0].total_score }} Poin</div>
            <div class="podium-block block-1-light">1</div>
          </div>

          <!-- 3rd Place -->
          <div v-if="participantsList[2]" class="podium-step step-3 animate-pop">
            <div class="avatar-podium bronze">
              {{ participantsList[2].student_name.charAt(0) }}
            </div>
            <div class="podium-name-dark">{{ participantsList[2].student_name }}</div>
            <div class="podium-score-dark">{{ participantsList[2].total_score }} Poin</div>
            <div class="podium-block block-3-light">3</div>
          </div>
        </div>

        <div class="finished-footer glass-card-light">
          <p class="text-muted mb-3">✅ <em>Nilai hasil kuis live ini telah otomatis tersimpan ke Rekap Nilai Kuis Mahasiswa.</em></p>
          <button @click="exitQuiz" class="btn btn-back-bank">
            <span>Kembali ke Halaman Utama</span>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* LIGHT THEME STYLES - SOLID CLEAN COLORS (NO GRADIENTS) */
.live-quiz-page.light-theme {
  min-height: calc(100vh - 4rem);
  background: #f8fafc;
  color: #0f172a;
  padding: 1.25rem 1rem 3rem;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  border-radius: 1.25rem;
  box-sizing: border-box;
}

/* Glassmorphic Light Cards */
.glass-card-light {
  background: #ffffff;
  border: 1px solid rgba(226, 232, 240, 0.9);
  border-radius: 1.25rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.05);
}

/* In-Page Header Toolbar Light */
.live-header-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1.25rem;
  margin-bottom: 1.5rem;
}

.header-left { display: flex; align-items: center; gap: 0.75rem; }

.btn-exit-light {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 0.4rem 0.85rem;
  border-radius: 999px;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-exit-light:hover {
  background: #fee2e2;
  color: #ef4444;
  border-color: #fca5a5;
}

.live-badge-light {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: #fef3c7;
  border: 1px solid #fde68a;
  color: #d97706;
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.header-center { display: flex; align-items: center; gap: 0.75rem; }

.title-pin-pill-light {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
  cursor: pointer;
  transition: transform 0.15s;
}

.title-pin-pill-light:hover {
  transform: scale(1.03);
}

.pin-lbl { font-size: 0.75rem; color: #b45309; font-weight: 800; }
.pin-code-light { font-size: 1.1rem; font-weight: 900; letter-spacing: 0.1em; color: #d97706; }
.quiz-title-tag-light { font-size: 0.85rem; color: #475569; font-weight: 700; }

.participants-pill-light {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #2563eb;
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 800;
}

/* Containers */
.live-card-container { max-width: 580px; margin: 1.5rem auto; }
.live-active-container { max-width: 860px; margin: 0 auto; }

/* Dosen Setup Tab Selector */
.dosen-tab-selector {
  display: flex;
  gap: 0.5rem;
  background: #f1f5f9;
  padding: 0.35rem;
  border-radius: 0.85rem;
  margin-bottom: 1.25rem;
}

.tab-choice-btn {
  flex: 1;
  border: none;
  background: transparent;
  color: #64748b;
  padding: 0.6rem;
  border-radius: 0.65rem;
  font-weight: 800;
  font-size: 0.85rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  transition: all 0.2s;
}

.tab-choice-btn.active {
  background: #ffffff;
  color: #d97706;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08);
}

.brand-logo-circle-gold {
  width: 3.5rem;
  height: 3.5rem;
  border-radius: 50%;
  background: #d97706;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25);
}

.logo-icon { width: 1.8rem; height: 1.8rem; color: #fff; }

.setup-card, .join-card { padding: 2rem 1.75rem; }
.card-title-header h2, .join-header h2 { font-size: 1.35rem; font-weight: 800; color: #0f172a; margin-bottom: 0.4rem; }
.card-title-header p, .join-header p { color: #64748b; font-size: 0.85rem; line-height: 1.4; }

.setup-form, .join-form { display: flex; flex-direction: column; gap: 1rem; }
.form-group label { font-size: 0.85rem; font-weight: 700; color: #334155; margin-bottom: 0.4rem; display: flex; align-items: center; gap: 0.4rem; }

.select-input-light, .text-input-light, .pin-input-light {
  width: 100%;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 0.75rem;
  padding: 0.75rem 1rem;
  color: #0f172a;
  font-size: 0.95rem;
  outline: none;
  box-sizing: border-box;
}

.select-input-light:focus, .text-input-light:focus {
  border-color: #f59e0b;
  box-shadow: 0 0 10px rgba(245, 158, 11, 0.2);
}

.pin-input-light {
  font-size: 1.8rem;
  font-weight: 900;
  letter-spacing: 0.2em;
  text-align: center;
  color: #d97706;
}

/* Instant Question Input Cards for Dosen */
.instant-q-section {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 0.85rem;
  padding: 1rem;
  margin-top: 0.5rem;
}

.section-lbl { font-size: 0.85rem; font-weight: 800; color: #475569; margin-bottom: 0.75rem; display: block; }

.q-input-card-light {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 0.75rem;
  padding: 0.85rem;
  margin-bottom: 0.85rem;
}

.q-card-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.q-num-lbl { font-size: 0.8rem; font-weight: 800; color: #d97706; }
.btn-remove-q { background: #fee2e2; color: #ef4444; border: none; padding: 0.2rem 0.5rem; border-radius: 0.4rem; font-size: 0.75rem; font-weight: 700; cursor: pointer; }

.opts-grid-2x2 { display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem; margin-bottom: 0.5rem; }
.opt-input-light { background: #f1f5f9; border: 1px solid #cbd5e1; border-radius: 0.5rem; padding: 0.5rem 0.75rem; font-size: 0.85rem; }

.key-select-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8rem; font-weight: 700; color: #475569; }
.select-key-light { background: #ffffff; border: 1px solid #cbd5e1; border-radius: 0.4rem; padding: 0.3rem 0.5rem; font-size: 0.8rem; }

.btn-add-row-light {
  width: 100%;
  background: #ffffff;
  border: 1px dashed #cbd5e1;
  color: #3b82f6;
  font-weight: 800;
  padding: 0.6rem;
  border-radius: 0.65rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
}

.time-options-grid-light { display: flex; gap: 0.5rem; }

.time-opt-btn-light {
  flex: 1;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #64748b;
  padding: 0.6rem;
  border-radius: 0.65rem;
  font-weight: 800;
  cursor: pointer;
}

.time-opt-btn-light.active {
  background: #d97706;
  color: #ffffff;
  border-color: #b45309;
}

/* Solid Colored Main Buttons */
.btn-launch-live-gold, .btn-live-join-gold {
  width: 100%;
  background: #d97706;
  color: #ffffff;
  font-size: 1rem;
  font-weight: 800;
  padding: 0.9rem;
  border-radius: 0.75rem;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25);
  margin-top: 0.5rem;
  transition: background-color 0.2s ease;
}

.btn-launch-live-gold:hover, .btn-live-join-gold:hover {
  background: #b45309;
}

.join-existing-divider-light { text-align: center; margin: 1.25rem 0 0.85rem; font-size: 0.75rem; color: #94a3b8; }
.join-existing-row { display: flex; gap: 0.5rem; }

.pin-short-input-light {
  flex: 1;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 0.65rem;
  padding: 0.6rem 0.85rem;
  color: #d97706;
  font-weight: 800;
  letter-spacing: 0.1em;
}

.btn-secondary-light {
  background: #f1f5f9;
  color: #334155;
  border: 1px solid #cbd5e1;
  padding: 0.6rem 1rem;
  border-radius: 0.65rem;
  font-weight: 700;
  cursor: pointer;
}

/* Lobby View Light */
.state-lobby { display: flex; flex-direction: column; gap: 1.25rem; }

.lobby-banner {
  padding: 2.25rem 1.75rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 1.25rem;
}

.pin-display-box-light {
  background: #fffbeb;
  border: 2px dashed #fde68a;
  padding: 1.25rem 3rem;
  border-radius: 1.5rem;
  cursor: pointer;
  transition: transform 0.2s;
}

.pin-display-box-light:hover {
  transform: scale(1.02);
}

.pin-sub-gold { font-size: 0.75rem; letter-spacing: 0.15em; color: #b45309; font-weight: 800; }
.pin-large-dark { font-size: 3.5rem; font-weight: 900; letter-spacing: 0.25em; color: #d97706; margin: 0.2rem 0; }
.pin-hint-light { font-size: 0.75rem; color: #94a3b8; display: flex; align-items: center; justify-content: center; gap: 0.3rem; }

.lobby-meta-light h3 { font-size: 1.3rem; font-weight: 800; color: #0f172a; margin-bottom: 0.3rem; }
.lobby-meta-light p { color: #64748b; font-size: 0.85rem; margin-bottom: 0.75rem; }

.lobby-status-pill-gold {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #fef3c7;
  color: #b45309;
  padding: 0.4rem 1rem;
  border-radius: 999px;
  font-weight: 700;
  font-size: 0.85rem;
  border: 1px solid #fde68a;
}

.dot-pulse-amber {
  width: 0.6rem;
  height: 0.6rem;
  border-radius: 50%;
  background: #f59e0b;
  animation: pulse 1s infinite alternate;
}

.participants-section { padding: 1.5rem 1.75rem; }
.section-title-dark { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1rem; color: #0f172a; font-size: 1rem; font-weight: 800; }

.participants-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.participant-chip-light {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  padding: 0.4rem 0.9rem 0.4rem 0.4rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.03);
}

.avatar-circle-blue {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: #2563eb;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
}

.p-name-dark { color: #1e293b; }

.btn-start-live-emerald {
  width: 100%;
  background: #059669;
  color: #ffffff;
  font-size: 1.15rem;
  font-weight: 900;
  padding: 1.1rem;
  border-radius: 1rem;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.6rem;
  box-shadow: 0 4px 12px rgba(5, 150, 105, 0.25);
  transition: background-color 0.2s ease;
}

.btn-start-live-emerald:hover {
  background: #047857;
}

.student-waiting-box { padding: 2.5rem; text-align: center; }

/* Question View Light */
.question-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.25rem; }

.progress-bar-container-light {
  flex: 1;
  height: 0.85rem;
  background: #e2e8f0;
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: #059669;
  transition: width 1s linear;
}

.progress-bar-fill.warning-time {
  background: #ef4444;
}

.timer-badge-gold {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
  padding: 0.45rem 1rem;
  border-radius: 999px;
  font-weight: 900;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.timer-badge-gold.urgent {
  background: #fee2e2;
  color: #ef4444;
  border-color: #fca5a5;
  animation: pulse 0.5s infinite alternate;
}

.q-number-badge-light {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 0.45rem 0.9rem;
  border-radius: 999px;
  font-weight: 700;
  font-size: 0.85rem;
}

.question-card {
  padding: 2rem 1.75rem;
  text-align: center;
  margin-bottom: 1.25rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 15px rgba(15, 23, 42, 0.04);
}

.q-badge-mini {
  display: inline-block;
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
  font-size: 0.75rem;
  font-weight: 800;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  margin-bottom: 0.75rem;
  letter-spacing: 0.05em;
}

.question-text-dark {
  font-size: 1.45rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.45;
  margin: 0;
}

.votes-counter-card { padding: 1rem 1.25rem; display: flex; align-items: center; gap: 0.85rem; margin-bottom: 1.25rem; }
.counter-text h3 { font-size: 1.15rem; font-weight: 800; color: #0f172a; margin: 0 0 0.1rem; }
.counter-text p { font-size: 0.8rem; color: #64748b; margin: 0; }

.options-grid-presenter, .options-grid-player { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-bottom: 1.5rem; }

.option-card-presenter {
  padding: 1.1rem 1.25rem;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  gap: 0.85rem;
  font-size: 1.05rem;
  font-weight: 800;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

.btn-option-player {
  padding: 1.15rem 1.25rem;
  border-radius: 1rem;
  border: none;
  color: #ffffff;
  font-size: 1.05rem;
  font-weight: 800;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.1);
  transition: transform 0.15s ease, opacity 0.15s ease;
  text-align: left;
}

.btn-option-player:hover {
  transform: translateY(-2px);
  opacity: 0.95;
}

.shape-badge-player {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  flex-shrink: 0;
}

.opt-label-player {
  line-height: 1.35;
  word-break: break-word;
}

.voted-confirmation-box {
  padding: 2.5rem 1.5rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.check-icon-circle-lg {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: #059669;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
  box-shadow: 0 4px 12px rgba(5, 150, 105, 0.3);
}

.waiting-reveal-pill-light {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  color: #b45309;
  padding: 0.4rem 1rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 700;
  margin-top: 1rem;
}

/* Reveal View Light */
.reveal-header { text-align: center; margin-bottom: 1.5rem; }
.reveal-header h2 { font-size: 1.4rem; font-weight: 800; margin-bottom: 0.3rem; }
.question-sub-dark { font-size: 0.9rem; color: #64748b; margin: 0; }

.chart-container { padding: 1.75rem 1.5rem 1.5rem; margin-bottom: 1.5rem; }

.chart-bars {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 190px;
  gap: 1rem;
}

.chart-bar-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  flex: 1;
  max-width: 100px;
}

.bar-count-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.2rem;
  margin-bottom: 0.4rem;
  min-height: 2.5rem;
  justify-content: flex-end;
}

.correct-badge-floating {
  background: #059669;
  color: #ffffff;
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
  font-size: 0.85rem;
  box-shadow: 0 2px 8px rgba(5, 150, 105, 0.4);
}

.bar-count-dark { font-weight: 900; font-size: 1rem; color: #0f172a; }

.bar-track-light {
  flex: 1;
  width: 100%;
  max-width: 65px;
  background: #f1f5f9;
  border-radius: 0.75rem 0.75rem 0 0;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.bar-fill {
  width: 100%;
  border-radius: 0.75rem 0.75rem 0 0;
  transition: height 0.8s ease;
}

.bar-fill.is-correct-bar {
  box-shadow: 0 -2px 8px rgba(5, 150, 105, 0.3);
}

.bar-label-dark {
  margin-top: 0.6rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.35rem;
}

.option-shape-pill {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  color: #ffffff;
  font-weight: 800;
  font-size: 0.85rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.08);
}

.correct-pill-tag {
  background: #d1fae5;
  color: #047857;
  border: 1px solid #6ee7b7;
  font-size: 0.7rem;
  font-weight: 900;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  white-space: nowrap;
}

/* Explanation Card */
.explanation-card {
  padding: 1.5rem;
  margin-bottom: 1.25rem;
  background: #ffffff;
  border: 1.5px solid #a7f3d0;
  border-left: 6px solid #059669;
}

.explanation-header-box {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  margin-bottom: 0.85rem;
}

.check-icon-circle {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background: #059669;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(5, 150, 105, 0.25);
  flex-shrink: 0;
}

.exp-subtitle {
  font-size: 0.7rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  color: #047857;
}

.explanation-title-text h4 {
  font-size: 1.05rem;
  font-weight: 800;
  color: #064e3b;
  margin: 0.1rem 0 0;
}

.explanation-body-box {
  background: #f8fafc;
  border: 1px solid #d1fae5;
  padding: 0.85rem 1rem;
  border-radius: 0.75rem;
  font-size: 0.9rem;
  color: #1e293b;
  line-height: 1.5;
}

/* Personal Feedback Card */
.personal-feedback-card {
  padding: 1rem 1.25rem;
  margin-bottom: 1.25rem;
}

.feedback-box {
  display: flex;
  align-items: center;
  gap: 1rem;
  border-radius: 0.85rem;
  padding: 0.85rem 1.15rem;
}

.feedback-box.correct {
  background: #fffbeb;
  border: 1.5px solid #fde68a;
}

.feedback-box.wrong {
  background: #fef2f2;
  border: 1.5px solid #fca5a5;
}

.feedback-icon-wrap {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feedback-icon-wrap.gold {
  background: #d97706;
  box-shadow: 0 2px 8px rgba(217, 119, 6, 0.25);
}

.feedback-icon-wrap.red {
  background: #dc2626;
  box-shadow: 0 2px 8px rgba(220, 38, 38, 0.25);
}

.feedback-text-content h3 {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 0.15rem;
}

.feedback-text-content p {
  font-size: 0.85rem;
  color: #475569;
  margin: 0;
}

/* Solid Action Buttons */
.btn-reveal-action, .btn-leaderboard, .btn-next-q, .btn-finish-quiz {
  width: 100%;
  padding: 1rem;
  font-size: 1.05rem;
  font-weight: 900;
  border-radius: 0.85rem;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: background-color 0.2s ease;
}

.btn-reveal-action { background: #d97706; color: #fff; box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25); }
.btn-reveal-action:hover { background: #b45309; }

.btn-leaderboard { background: #2563eb; color: #fff; box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25); }
.btn-leaderboard:hover { background: #1d4ed8; }

.btn-next-q { background: #059669; color: #fff; box-shadow: 0 4px 12px rgba(5, 150, 105, 0.25); }
.btn-next-q:hover { background: #047857; }

.btn-finish-quiz { background: #7c3aed; color: #fff; box-shadow: 0 4px 12px rgba(124, 58, 237, 0.25); }
.btn-finish-quiz:hover { background: #6d28d9; }

/* Leaderboard View Light */
.leaderboard-header { text-align: center; margin-bottom: 1.5rem; display: flex; flex-direction: column; align-items: center; }

.trophy-circle-lg {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: #d97706;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.75rem;
  box-shadow: 0 4px 14px rgba(217, 119, 6, 0.3);
}

.trophy-circle-lg.gold {
  background: #d97706;
}

.leaderboard-header h2 { font-size: 1.4rem; font-weight: 800; color: #0f172a; margin: 0 0 0.2rem; }
.leaderboard-header p { font-size: 0.85rem; color: #64748b; margin: 0; }

.leaderboard-list {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  margin-bottom: 1.25rem;
  max-height: 380px;
  overflow-y: auto;
}

.leaderboard-row-light {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 0.85rem;
  font-size: 0.95rem;
  transition: all 0.2s;
}

.leaderboard-row-light.rank-1-light {
  background: #fffbeb;
  border: 1.5px solid #fde68a;
}

.leaderboard-row-light.rank-2-light {
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
}

.leaderboard-row-light.rank-3-light {
  background: #fff7ed;
  border: 1.5px solid #fed7aa;
}

.leaderboard-row-light.is-me {
  border-color: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.rank-badge-dark {
  font-weight: 900;
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #475569;
  font-size: 0.85rem;
  margin-right: 0.5rem;
  flex-shrink: 0;
}

.rank-badge-dark.gold {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
}

.rank-badge-dark.silver {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #cbd5e1;
}

.rank-badge-dark.bronze {
  background: #ffedd5;
  color: #c2410c;
  border: 1px solid #fed7aa;
}

.crown-icon { width: 1.1rem; height: 1.1rem; }

.player-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.avatar-small {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background: #2563eb;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.9rem;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.2);
  flex-shrink: 0;
}

.avatar-small.gold { background: #d97706; }
.avatar-small.silver { background: #64748b; }
.avatar-small.bronze { background: #b45309; }

.player-details {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.player-name-dark { color: #0f172a; font-weight: 800; font-size: 0.95rem; }

.streak-badge {
  font-size: 0.7rem;
  color: #ea580c;
  font-weight: 800;
  display: flex;
  align-items: center;
  gap: 0.2rem;
}

.player-score {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
}

.score-val-dark {
  font-size: 1.15rem;
  font-weight: 900;
  color: #d97706;
}

.pts-lbl-muted {
  font-size: 0.75rem;
  font-weight: 700;
  color: #b45309;
}

/* Personal Rank Card for Student */
.personal-rank-footer {
  padding: 1.1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #ffffff;
  border: 1.5px solid #fde68a;
  box-shadow: 0 4px 15px rgba(245, 158, 11, 0.08);
  border-radius: 1rem;
  margin-bottom: 1.25rem;
}

.rank-footer-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.my-rank-badge {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background: #d97706;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: 900;
  box-shadow: 0 2px 8px rgba(217, 119, 6, 0.25);
}

.my-rank-text h4 {
  font-size: 0.95rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 0.15rem;
}

.my-rank-text p {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
}

.rank-footer-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.my-score-num {
  font-size: 1.6rem;
  font-weight: 900;
  color: #d97706;
  line-height: 1;
}

.my-score-lbl {
  font-size: 0.75rem;
  font-weight: 700;
  color: #b45309;
}

/* Podium Light */
.finished-header { text-align: center; margin-bottom: 1.5rem; display: flex; flex-direction: column; align-items: center; }

.podium-container {
  display: flex;
  justify-content: center;
  align-items: flex-end;
  gap: 1rem;
  margin-bottom: 1.5rem;
  height: 240px;
}

.podium-step { display: flex; flex-direction: column; align-items: center; width: 120px; }
.avatar-podium { width: 3.2rem; height: 3.2rem; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 1.4rem; font-weight: 900; margin-bottom: 0.3rem; box-shadow: 0 2px 8px rgba(0,0,0,0.12); }

.avatar-podium.gold { background: #d97706; color: #fff; }
.avatar-podium.silver { background: #64748b; color: #fff; }
.avatar-podium.bronze { background: #b45309; color: #fff; }

.podium-name-dark { font-weight: 800; font-size: 0.9rem; color: #0f172a; margin-bottom: 0.1rem; }
.podium-score-dark { font-size: 0.8rem; color: #d97706; font-weight: 900; margin-bottom: 0.5rem; }

.podium-block-1-light { height: 130px; width: 100%; border-radius: 0.75rem 0.75rem 0 0; background: #f59e0b; display: flex; align-items: center; justify-content: center; font-size: 2.2rem; font-weight: 900; color: #fff; box-shadow: 0 4px 12px rgba(245, 158, 11, 0.25); }
.podium-block-2-light { height: 90px; width: 100%; border-radius: 0.75rem 0.75rem 0 0; background: #94a3b8; display: flex; align-items: center; justify-content: center; font-size: 1.9rem; font-weight: 900; color: #fff; box-shadow: 0 4px 12px rgba(148, 163, 184, 0.25); }
.podium-block-3-light { height: 70px; width: 100%; border-radius: 0.75rem 0.75rem 0 0; background: #d97706; display: flex; align-items: center; justify-content: center; font-size: 1.7rem; font-weight: 900; color: #fff; box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25); }

.text-dark { color: #0f172a; }
.text-muted { color: #64748b; }
.text-gold { color: #d97706; }

.finished-footer { padding: 1.75rem; text-align: center; }
.btn-back-bank { display: inline-block; background: #2563eb; color: #fff; font-weight: 800; padding: 0.9rem 1.75rem; border-radius: 0.75rem; text-decoration: none; border: none; cursor: pointer; box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25); }

/* Animations & Responsive */
.animate-fade-in { animation: fadeIn 0.3s ease-out; }
.animate-pop { animation: pop 0.25s ease-out; }
.animate-bounce { animation: bounce 1.5s infinite; }
.icon-pulse { animation: pulse 1s infinite alternate; }

@keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }
@keyframes pop { from { opacity: 0; transform: scale(0.85); } to { opacity: 1; transform: scale(1); } }
@keyframes bounce { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-8px); } }
@keyframes pulse { from { opacity: 0.4; } to { opacity: 1; } }

@media (max-width: 768px) {
  .options-grid-player, .options-grid-presenter { grid-template-columns: 1fr; }
  .podium-container { height: 190px; }
  .podium-step { width: 85px; }
}
</style>
