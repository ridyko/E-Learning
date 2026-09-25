<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Clock, CheckCircle2, AlertCircle, Award, RotateCcw, Send, UserCheck } from 'lucide-vue-next'

import { showSuccess, showWarning } from '../utils/swal.js'

const props = defineProps({
  quiz: Object,
  courseTitle: String
})

const emit = defineEmits(['finishQuiz'])

const currentUser = ref(null)
const studentNim = ref('')
const studentName = ref('')
const isStarted = ref(false)
const isSubmitted = ref(false)
const selectedAnswers = ref({})
const essayAnswers = ref({})
const timeLeft = ref(0)
let timerInterval = null

const score = ref(0)
const totalQuestions = computed(() => props.quiz?.questions?.length || 0)

const loadUserInfo = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const u = JSON.parse(userStr)
      currentUser.value = u
      if (u.username) studentNim.value = u.username
      if (u.name) studentName.value = u.name
    } catch (e) {
      console.warn('Load user error:', e)
    }
  }
}

onMounted(() => {
  loadUserInfo()
})

const startQuiz = () => {
  if (currentUser.value) {
    if (!studentNim.value) studentNim.value = currentUser.value.username || ''
    if (!studentName.value) studentName.value = currentUser.value.name || ''
  }

  if (!studentNim.value || !studentName.value) {
    showWarning('Perhatian', 'Harap isi NIM dan Nama terlebih dahulu!')
    return
  }
  isStarted.value = true
  timeLeft.value = (props.quiz?.time_limit || 15) * 60

  timerInterval = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      submitQuiz()
    }
  }, 1000)
}

const formattedTime = computed(() => {
  const mins = Math.floor(timeLeft.value / 60)
  const secs = timeLeft.value % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
})

const selectOption = (qId, optionIdx) => {
  if (isSubmitted.value) return
  selectedAnswers.value[qId] = optionIdx
}

const submitQuiz = async () => {
  if (timerInterval) clearInterval(timerInterval)

  let correct = 0
  let mcCount = 0
  props.quiz.questions.forEach((q) => {
    if (q.type !== 'essay') {
      mcCount++
      if (selectedAnswers.value[q.id] === q.correct_answer) {
        correct++
      }
    }
  })

  if (mcCount > 0) {
    score.value = (correct / mcCount) * 100
  } else {
    score.value = 100
  }
  isSubmitted.value = true

  showSuccess('Kuis Selesai Diterima! 🎯', `Nilai Uji Pilihan Ganda: ${score.value.toFixed(0)} / 100. Jawaban berhasil dikirim.`)

  try {
    await fetch('/api/quizzes/submit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        quiz_id: props.quiz.id,
        student_nim: studentNim.value,
        student_name: studentName.value,
        answers: selectedAnswers.value,
        essay_answers: essayAnswers.value,
        score: score.value,
        total_score: 100
      })
    })
  } catch (err) {
    console.warn('Backend submit fallback:', err)
  }

  emit('finishQuiz', {
    quizId: props.quiz.id,
    score: score.value,
    studentName: studentName.value
  })
}

const resetQuiz = () => {
  isStarted.value = false
  isSubmitted.value = false
  selectedAnswers.value = {}
  if (timerInterval) clearInterval(timerInterval)
}

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})
</script>

<template>
  <div class="glass-card quiz-card" v-if="quiz">
    <!-- Header -->
    <div class="quiz-header">
      <div>
        <span class="badge badge-gold">{{ courseTitle }}</span>
        <h3>{{ quiz.title }}</h3>
        <p class="quiz-desc">{{ quiz.description }}</p>
      </div>

      <div v-if="isStarted && !isSubmitted" class="timer-box">
        <Clock class="timer-icon" />
        <span>Sisa Waktu: <strong>{{ formattedTime }}</strong></span>
      </div>
    </div>

    <!-- Step 1: Form Identitas Mahasiswa -->
    <div v-if="!isStarted" class="identity-step">
      <!-- Auto verified badge if user is logged in -->
      <div v-if="currentUser" class="user-identity-card">
        <div class="user-identity-main">
          <UserCheck class="user-badge-icon text-gold" />
          <div class="user-identity-info">
            <span class="user-identity-name">👤 Peserta Kuis: {{ studentName }}</span>
            <span class="user-identity-meta">NIM: <strong>{{ studentNim }}</strong> | {{ currentUser.prodi || 'Teknik Informatika' }}</span>
          </div>
        </div>
        <span class="badge badge-emerald">AKUN TERVERIFIKASI</span>
      </div>

      <!-- Fallback manual inputs if guest -->
      <div v-else class="form-grid">
        <div>
          <label class="input-label">NIM Mahasiswa</label>
          <input v-model="studentNim" class="glass-input" placeholder="Contoh: 20260801001" />
        </div>
        <div>
          <label class="input-label">Nama Lengkap</label>
          <input v-model="studentName" class="glass-input" placeholder="Masukkan nama anda" />
        </div>
      </div>

      <div class="rules-box">
        <h4>📋 Ketentuan Kuis Online:</h4>
        <ul>
          <li>Jumlah Soal: {{ totalQuestions }} Soal Pilihan Ganda</li>
          <li>Durasi Waktu: {{ quiz.time_limit }} Menit</li>
          <li>Nilai kelulusan minimal: 70 / 100</li>
        </ul>
      </div>

      <button @click="startQuiz" class="btn btn-gold btn-lg">
        Mulai Kuis Sekarang
      </button>
    </div>

    <!-- Step 2: Halaman Soal Kuis -->
    <div v-else-if="isStarted && !isSubmitted" class="questions-step">
      <div 
        v-for="(q, idx) in quiz.questions" 
        :key="q.id" 
        class="question-block"
      >
        <div class="question-title">
          <span class="q-num">Soal {{ idx + 1 }}:</span>
          <span>{{ q.question }}</span>
          <span v-if="q.type === 'essay'" class="badge badge-purple" style="margin-left: 0.5rem;">ESSAY</span>
        </div>

        <!-- Render Essay Input Textarea if question type is essay -->
        <div v-if="q.type === 'essay'" class="essay-input-box" style="margin-top: 1rem;">
          <textarea 
            v-model="essayAnswers[q.id]" 
            class="glass-input textarea" 
            placeholder="Tuliskan jawaban essay / uraian Anda secara lengkap dan terstruktur di sini..." 
            rows="4"
            style="width: 100%; border-radius: 0.85rem; padding: 1rem; font-family: inherit;"
          ></textarea>
        </div>

        <!-- Render Multiple Choice Option Buttons if question type is MC -->
        <div v-else class="options-grid">
          <button 
            v-for="(opt, optIdx) in q.options" 
            :key="optIdx"
            @click="selectOption(q.id, optIdx)"
            :class="['option-btn', selectedAnswers[q.id] === optIdx ? 'selected' : '']"
          >
            <span class="opt-prefix">{{ String.fromCharCode(65 + optIdx) }}</span>
            <span>{{ opt }}</span>
          </button>
        </div>
      </div>

      <div class="submit-footer">
        <button @click="submitQuiz" class="btn btn-primary btn-lg">
          <Send class="btn-icon-sm" />
          Kirim Jawaban Kuis
        </button>
      </div>
    </div>

    <!-- Step 3: Hasil & Review Kuis -->
    <div v-else-if="isSubmitted" class="result-step">
      <div class="score-banner">
        <Award class="award-icon" />
        <div>
          <span class="score-label">Hasil Kuis {{ studentName }} ({{ studentNim }})</span>
          <div class="score-number">{{ Math.round(score) }} <span class="max-score">/ 100</span></div>
          <p v-if="score >= 70" class="text-emerald">🎉 Selamat! Anda Lulus Kuis Matakuliah ini!</p>
          <p v-else class="text-rose">⚠️ Tetap semangat, silakan pelajari kembali modul pertemuan.</p>
        </div>
      </div>

      <!-- Detail Pembahasan Jawaban -->
      <div class="review-section">
        <h4>Pembahasan Jawaban:</h4>
        <div 
          v-for="(q, idx) in quiz.questions" 
          :key="q.id" 
          class="review-block"
        >
          <div class="review-q">
            <strong>Soal {{ idx + 1 }}:</strong> {{ q.question }}
          </div>
          <div class="review-status">
            <span v-if="selectedAnswers[q.id] === q.correct_answer" class="status-pass">
              <CheckCircle2 class="icon-sm" /> Jawaban Anda Benar ({{ q.options[q.correct_answer] }})
            </span>
            <span v-else class="status-fail">
              <AlertCircle class="icon-sm" /> Jawaban Anda: {{ q.options[selectedAnswers[q.id]] || 'Tidak Dijawab' }} | Jawaban Benar: {{ q.options[q.correct_answer] }}
            </span>
          </div>
          <p class="explanation-box">💡 <strong>Pembahasan Pak Rio:</strong> {{ q.explanation }}</p>
        </div>
      </div>

      <button @click="emit('finishQuiz')" class="btn btn-gold btn-lg w-full margin-top">
        <CheckCircle2 class="btn-icon-sm" />
        Kembali ke Daftar Kuis
      </button>
    </div>
  </div>
</template>

<style scoped>
.quiz-card {
  padding: 2rem;
  margin-bottom: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.quiz-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 1rem;
}

.quiz-header h3 {
  font-size: 1.4rem;
  font-weight: 800;
  color: #0f172a;
  margin-top: 0.3rem;
}

.quiz-desc {
  color: #475569;
  font-size: 0.9rem;
}

.timer-box {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #fffbeb;
  color: #b45309;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  border: 1px solid #fde68a;
  font-family: var(--font-heading);
}

.timer-icon {
  width: 18px;
  height: 18px;
}

.identity-step {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-grid {
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

.rules-box {
  background: #f8fafc;
  padding: 1.25rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.rules-box h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 0.5rem;
}

.rules-box ul {
  padding-left: 1.25rem;
  font-size: 0.88rem;
  color: #475569;
}

.questions-step {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.question-block {
  background: #ffffff;
  padding: 1.25rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #e2e8f0;
}

.question-title {
  font-weight: 700;
  font-size: 1rem;
  color: #0f172a;
  margin-bottom: 1rem;
}

.q-num {
  color: #d97706;
  margin-right: 0.4rem;
}

.options-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.option-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: var(--radius-sm);
  color: #0f172a;
  font-size: 0.9rem;
  cursor: pointer;
  text-align: left;
  transition: var(--transition);
}

.option-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.option-btn.selected {
  background: #eff6ff;
  border-color: #2563eb;
  color: #1d4ed8;
  font-weight: 700;
}

.opt-prefix {
  width: 24px;
  height: 24px;
  background: #e2e8f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 0.8rem;
  color: #0f172a;
}

.submit-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}

.btn-lg {
  padding: 0.85rem 1.75rem;
  font-size: 1rem;
}

.result-step {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.score-banner {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: #fffbeb;
  padding: 1.5rem;
  border-radius: var(--radius-md);
  border: 1px solid #fde68a;
}

.award-icon {
  width: 50px;
  height: 50px;
  color: #d97706;
}

.score-label {
  font-size: 0.9rem;
  color: #64748b;
  font-weight: 600;
}

.score-number {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 800;
  color: #d97706;
}

.max-score {
  font-size: 1.2rem;
  color: #94a3b8;
}

.review-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.review-block {
  background: #f8fafc;
  padding: 1rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.review-q {
  color: #0f172a;
}

.status-pass {
  color: #047857;
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.88rem;
  font-weight: 600;
}

.status-fail {
  color: #be123c;
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.88rem;
  font-weight: 600;
}

.explanation-box {
  margin-top: 0.5rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 0.6rem 0.85rem;
  border-radius: var(--radius-sm);
  font-size: 0.88rem;
  color: #1e3a8a;
}

.icon-sm {
  width: 16px;
  height: 16px;
}

.text-emerald {
  color: #059669;
  font-weight: 700;
}

.text-rose {
  color: #e11d48;
  font-weight: 700;
}

.user-identity-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.1rem 1.4rem;
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-left: 5px solid #2563eb;
  border-radius: var(--radius-md);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.user-identity-main {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-badge-icon {
  width: 32px;
  height: 32px;
}

.user-identity-info {
  display: flex;
  flex-direction: column;
}

.user-identity-name {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
}

.user-identity-meta {
  font-size: 0.85rem;
  color: #2563eb;
  font-weight: 600;
}

.badge-emerald {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}
</style>
