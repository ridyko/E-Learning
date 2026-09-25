import { createRouter, createWebHistory } from 'vue-router'
import { showWarning } from '../utils/swal.js'

import HomeView from '../views/HomeView.vue'
import CourseView from '../views/CourseView.vue'
import QuizView from '../views/QuizView.vue'
import PlaygroundView from '../views/PlaygroundView.vue'
import PresensiView from '../views/PresensiView.vue'
import AdminView from '../views/AdminView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import ProfileView from '../views/ProfileView.vue'
import MahasiswaView from '../views/MahasiswaView.vue'
import BankSoalView from '../views/BankSoalView.vue'
import ManageCoursesView from '../views/ManageCoursesView.vue'
import ManageAnnouncementsView from '../views/ManageAnnouncementsView.vue'
import LiveQuizView from '../views/LiveQuizView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Public Access (Landing page for guests / Dashboard for logged in users)
    { path: '/', name: 'home', component: HomeView },

    // Guest Only Routes (Login & Register)
    { path: '/login', name: 'login', component: LoginView, meta: { guestOnly: true } },
    { path: '/register', name: 'register', component: RegisterView, meta: { guestOnly: true } },

    // Protected Routes (Requires Login - Accessible by Mahasiswa & Dosen)
    { path: '/course/:id', name: 'course', component: CourseView, meta: { requiresAuth: true } },
    { path: '/quizzes', name: 'quizzes', component: QuizView, meta: { requiresAuth: true } },
    { path: '/live-quiz', name: 'live-quiz', component: LiveQuizView, meta: { requiresAuth: true } },
    { path: '/live-quiz/:pin', name: 'live-quiz-pin', component: LiveQuizView, meta: { requiresAuth: true } },
    { path: '/playground', name: 'playground', component: PlaygroundView, meta: { requiresAuth: true } },
    { path: '/presensi', name: 'presensi', component: PresensiView, meta: { requiresAuth: true } },
    { path: '/profile', name: 'profile', component: ProfileView, meta: { requiresAuth: true } },

    // Protected Dosen-Only Routes (Requires Login & Dosen Role)
    { path: '/admin', name: 'admin', component: AdminView, meta: { requiresAuth: true, requiresRole: 'dosen' } },
    { path: '/mata-kuliah', name: 'mata-kuliah', component: ManageCoursesView, meta: { requiresAuth: true, requiresRole: 'dosen' } },
    { path: '/pengumuman', name: 'pengumuman', component: ManageAnnouncementsView, meta: { requiresAuth: true, requiresRole: 'dosen' } },
    { path: '/bank-soal', name: 'bank-soal', component: BankSoalView, meta: { requiresAuth: true, requiresRole: 'dosen' } },
    { path: '/mahasiswa', name: 'mahasiswa', component: MahasiswaView, meta: { requiresAuth: true, requiresRole: 'dosen' } },

    // Catch-all wildcard redirect to Home
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

// Navigation Guard to enforce login & role authorization
router.beforeEach((to, from, next) => {
  const userStr = localStorage.getItem('user')
  let user = null
  if (userStr) {
    try {
      user = JSON.parse(userStr)
    } catch (e) {
      user = null
    }
  }

  // 1. If route requires authentication & user is NOT logged in -> Redirect to /login
  if (to.meta.requiresAuth && !user) {
    showWarning('Akses Ditolak 🔒', 'Anda harus login terlebih dahulu untuk mengakses menu ini.')
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // 2. If route requires specific role ('dosen') & user is NOT Dosen -> Redirect to /
  if (to.meta.requiresRole === 'dosen' && user?.role !== 'dosen') {
    showWarning('Akses Terbatas 🔒', 'Menu ini khusus untuk Dosen Pengampu.')
    return next({ path: '/' })
  }

  // 3. If route is guest-only (login/register) & user IS logged in -> Redirect to dashboard
  if (to.meta.guestOnly && user) {
    if (user.role === 'dosen') {
      return next({ path: '/admin' })
    }
    return next({ path: '/' })
  }

  next()
})

export default router

