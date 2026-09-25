import { createRouter, createWebHistory } from 'vue-router'
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
    { path: '/', name: 'home', component: HomeView },
    { path: '/course/:id', name: 'course', component: CourseView },
    { path: '/quizzes', name: 'quizzes', component: QuizView },
    { path: '/live-quiz', name: 'live-quiz', component: LiveQuizView },
    { path: '/live-quiz/:pin', name: 'live-quiz-pin', component: LiveQuizView },
    { path: '/playground', name: 'playground', component: PlaygroundView },
    { path: '/presensi', name: 'presensi', component: PresensiView },
    { path: '/mahasiswa', name: 'mahasiswa', component: MahasiswaView },
    { path: '/bank-soal', name: 'bank-soal', component: BankSoalView },
    { path: '/mata-kuliah', name: 'mata-kuliah', component: ManageCoursesView },
    { path: '/pengumuman', name: 'pengumuman', component: ManageAnnouncementsView },
    { path: '/admin', name: 'admin', component: AdminView },
    { path: '/profile', name: 'profile', component: ProfileView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
  ]
})

export default router
