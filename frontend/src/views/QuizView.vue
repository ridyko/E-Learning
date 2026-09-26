<script setup>
import { ref, computed, onMounted } from 'vue'
import QuizSimulator from '../components/QuizSimulator.vue'
import { CheckSquare, BookOpen, Code, Clock, Award, ShieldCheck, Plus, CheckCircle2, Lock, Radio, Zap } from 'lucide-vue-next'
import { showWarning } from '../utils/swal.js'

const quizzes = ref([])
const quizSubmissions = ref([])
const activeQuiz = ref(null)
const currentUser = ref(null)

const checkUser = () => {
  const u = localStorage.getItem('user')
  if (u) {
    try { currentUser.value = JSON.parse(u) } catch { currentUser.value = null }
  }
}

const isDosen = computed(() => currentUser.value?.role === 'dosen')

const fetchQuizzesData = async () => {
  try {
    const [resQ, resS] = await Promise.all([
      fetch('/api/quizzes'),
      fetch('/api/quizzes/submissions')
    ])
    quizzes.value = await resQ.json()
    quizSubmissions.value = await resS.json()
  } catch (err) {
    console.warn('Fetch quizzes error:', err)
  }
}

const getStudentSubmission = (quizId) => {
  if (!currentUser.value || currentUser.value.role === 'dosen') return null
  return quizSubmissions.value.find(s => s.quiz_id === quizId && s.student_nim === currentUser.value.username)
}

onMounted(() => {
  checkUser()
  fetchQuizzesData()
})

const selectQuiz = (quiz) => {
  const existing = getStudentSubmission(quiz.id)
  if (existing) {
    showWarning(
      'Kuis Sudah Selesai Dikerjakan! 🔒', 
      `Anda sudah menyelesaikan kuis '${quiz.title}' dengan nilai ${Math.round(existing.score)}/100. Setiap mahasiswa hanya memiliki 1 kali kesempatan pengerjaan.`
    )
    return
  }
  activeQuiz.value = quiz
}

const handleQuizFinish = () => {
  activeQuiz.value = null
  fetchQuizzesData()
}
</script>

<template>
  <div class="quiz-container animate-fade-in">
    <!-- Dosen Quick Notice Banner -->
    <div v-if="isDosen" class="dosen-quiz-banner">
      <div class="banner-info">
        <ShieldCheck class="banner-icon text-gold" />
        <div>
          <h4>Mode Dosen Pengampu</h4>
          <p>Untuk menginput soal kuis baru atau melihat rekap nilai mahasiswa, silakan masuk ke menu Bank Soal & Nilai Kuis.</p>
        </div>
      </div>
      <router-link to="/bank-soal" class="btn btn-gold btn-banner-action">
        <Plus class="icon-xs" />
        <span>Buka Bank Soal & Nilai</span>
      </router-link>
    </div>

    <!-- Live Quiz Quick Join Banner for Students & Guest -->
    <div class="live-quiz-join-banner glassmorphism" style="background: #fffbeb; border: 1px solid #fde68a; border-radius: 1.25rem; padding: 1.25rem 1.5rem; margin-bottom: 2rem; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 1rem; box-shadow: 0 4px 15px rgba(0, 0, 0, 0.03);">
      <div style="display: flex; align-items: center; gap: 1rem;">
        <div style="width: 3.2rem; height: 3.2rem; border-radius: 50%; background: #d97706; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25);">
          <Radio style="width: 1.6rem; height: 1.6rem; color: white;" class="animate-pulse" />
        </div>
        <div>
          <h4 style="margin: 0 0 0.2rem 0; font-weight: 800; font-size: 1.15rem; color: #d97706;">Live Quiz Interaktif Real-Time 🔥</h4>
          <p style="margin: 0; font-size: 0.9rem; color: #475569;">Punya PIN 6-digit kuis dari Dosen Pengampu? Gabung sesi kuis live secara langsung!</p>
        </div>
      </div>
      <router-link to="/live-quiz" class="btn" style="background: #d97706; color: white; border: none; padding: 0.75rem 1.4rem; border-radius: 0.85rem; font-weight: 800; text-decoration: none; display: flex; align-items: center; gap: 0.5rem; box-shadow: 0 4px 12px rgba(217, 119, 6, 0.25);">
        <Zap style="width: 1rem; height: 1rem;" />
        <span>GABUNG DENGAN PIN 🚀</span>
      </router-link>
    </div>

    <!-- Header -->
    <div class="page-header">
      <div class="title-row">
        <CheckSquare class="header-icon text-gold" />
        <h2>Evaluasi & Kuis Online Mandiri</h2>
      </div>
      <p class="subtitle">Silakan pilih kuis sesuai mata kuliah yang sedang Anda tempuh di ITB Swadharma.</p>
    </div>

    <!-- Active Quiz Engine -->
    <QuizSimulator 
      v-if="activeQuiz" 
      :quiz="activeQuiz" 
      courseTitle="Kuis Perkuliahan"
      @finishQuiz="handleQuizFinish"
    />

    <!-- Quiz Selection Cards Grid -->
    <div v-else class="quiz-grid">
      <div 
        v-for="q in quizzes" 
        :key="q.id" 
        class="glass-card quiz-select-card"
        :class="{ 'quiz-completed-card': getStudentSubmission(q.id) }"
      >
        <div class="card-header">
          <span :class="['badge', q.course_id === 'rpl-2026' ? 'badge-blue' : 'badge-gold']">
            {{ q.course_id === 'rpl-2026' ? 'Rekayasa Perangkat Lunak' : 'Pemrograman Web' }}
          </span>
          <div class="time-tag">
            <Clock class="icon-xs" />
            <span>{{ q.time_limit }} Menit</span>
          </div>
        </div>

        <h3>{{ q.title }}</h3>
        <p class="desc">{{ q.description }}</p>

        <!-- Status Box if Already Submitted -->
        <div v-if="getStudentSubmission(q.id)" class="completed-info-box">
          <div class="completed-header">
            <CheckCircle2 class="icon-sm text-emerald" />
            <span>Status: <strong>Sudah Dikerjakan</strong></span>
          </div>
          <div class="completed-score">
            <span>Nilai Anda: <strong>{{ Math.round(getStudentSubmission(q.id).score) }} / 100</strong></span>
            <span :class="['badge', getStudentSubmission(q.id).score >= 70 ? 'badge-emerald' : 'badge-rose']">
              {{ getStudentSubmission(q.id).score >= 70 ? 'LULUS' : 'REMIDIAL' }}
            </span>
          </div>
        </div>

        <div v-else class="info-row">
          <span>📝 {{ q.questions?.length }} Soal Pilihan Ganda</span>
          <span>🎯 KKM: 70</span>
        </div>

        <!-- Button disabled if already completed -->
        <button 
          v-if="getStudentSubmission(q.id)" 
          @click="selectQuiz(q)"
          class="btn btn-secondary btn-completed w-full"
        >
          <Lock class="icon-xs" />
          <span>Kuis Sudah Selesai (1x Kesempatan)</span>
        </button>

        <button 
          v-else 
          @click="selectQuiz(q)" 
          class="btn btn-gold w-full"
        >
          Mulai Kuis Ini
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.quiz-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  margin-bottom: 2rem;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-icon {
  width: 32px;
  height: 32px;
}

.title-row h2 {
  font-size: 1.8rem;
  font-weight: 800;
  color: #0f172a;
}

.subtitle {
  color: #475569;
  font-size: 0.95rem;
  margin-top: 0.3rem;
}

.quiz-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 1.5rem;
}

.quiz-select-card {
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.time-tag {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.8rem;
  color: #b45309;
  font-weight: 700;
}

.icon-xs {
  width: 14px;
  height: 14px;
}

.quiz-select-card h3 {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.5rem;
}

.desc {
  color: #475569;
  font-size: 0.9rem;
  margin-bottom: 1.25rem;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.82rem;
  color: #334155;
  font-weight: 600;
  margin-bottom: 1.5rem;
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.w-full {
  width: 100%;
}

.text-gold {
  color: #d97706;
}

.dosen-quiz-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: #fefce8;
  border: 1.5px solid #fef08a;
  border-left: 5px solid #d97706;
  border-radius: var(--radius-md);
  margin-bottom: 2rem;
  box-shadow: 0 2px 4px rgba(217, 119, 6, 0.05);
}

.banner-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.banner-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.banner-info h4 {
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.banner-info p {
  font-size: 0.88rem;
  color: #713f12;
  margin: 0.2rem 0 0 0;
}

.btn-banner-action {
  white-space: nowrap;
  font-size: 0.88rem;
  font-weight: 700;
}

.completed-info-box {
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-left: 4px solid #059669;
  border-radius: var(--radius-sm);
  padding: 0.85rem;
  margin-bottom: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.completed-header {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.82rem;
  color: #334155;
}

.completed-score {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.88rem;
  color: #0f172a;
}

.btn-completed {
  background: #e2e8f0;
  color: #64748b;
  border: 1px solid #cbd5e1;
  cursor: not-allowed;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
}

.btn-completed:hover {
  background: #cbd5e1;
}

.badge-emerald {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.badge-rose {
  background: #fff1f2;
  color: #e11d48;
  border: 1px solid #fecdd3;
}

.text-emerald {
  color: #059669;
}

@media (max-width: 768px) {
  .dosen-quiz-banner {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
