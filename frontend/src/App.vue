<script setup>
import { ref, onMounted } from 'vue'
import Navbar from './components/Navbar.vue'

const isLoggedIn = ref(false)
const isSidebarCollapsed = ref(false)

const checkAuth = () => {
  isLoggedIn.value = !!localStorage.getItem('user')
}

const syncSidebarState = () => {
  isSidebarCollapsed.value = document.body.classList.contains('sidebar-collapsed')
}

onMounted(() => {
  checkAuth()
  syncSidebarState()
  window.addEventListener('auth-changed', checkAuth)
  window.addEventListener('sidebar-toggled', syncSidebarState)
})
</script>

<template>
  <div class="app-root" :class="{ 'with-sidebar': isLoggedIn, 'sidebar-collapsed': isSidebarCollapsed }">
    <!-- Left Sidebar & Top Header Navigation -->
    <Navbar />

    <!-- Main Content Area -->
    <main class="main-content">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="footer-wrapper">
      <div class="footer-container">
        <div class="footer-brand">
          <h4>Institut Teknologi dan Bisnis Swadharma</h4>
          <p>Portal E-Learning Resmi Pengampu Dosen: <strong>Rio Widyatmoko, S.Kom, M.M.S.I</strong></p>
          <p class="courses-tag">Mata Kuliah: <span>Rekayasa Perangkat Lunak</span> & <span>Pemrograman Web</span></p>
        </div>
        <div class="footer-copy">
          <p>&copy; 2026 Rio Widyatmoko — Powered by Golang REST API & Vue.js 3</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.app-root {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding-bottom: 3rem;
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.footer-wrapper {
  background: #ffffff;
  border-top: 1px solid rgba(15, 23, 42, 0.08);
  padding: 2.5rem 1.5rem 1.5rem;
  color: #475569;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.02);
  transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Desktop Sidebar Offset & Collapse state */
@media (min-width: 992px) {
  .app-root.with-sidebar:not(.sidebar-collapsed) .main-content,
  .app-root.with-sidebar:not(.sidebar-collapsed) .footer-wrapper {
    margin-left: 270px;
  }
  .app-root.with-sidebar.sidebar-collapsed .main-content,
  .app-root.with-sidebar.sidebar-collapsed .footer-wrapper {
    margin-left: 0;
  }
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  align-items: center;
  text-align: center;
}

.footer-brand h4 {
  color: #d97706;
  font-size: 1.1rem;
  margin-bottom: 0.3rem;
  font-weight: 800;
}

.courses-tag {
  font-size: 0.85rem;
  margin-top: 0.4rem;
}

.courses-tag span {
  color: #1d4ed8;
  font-weight: 700;
}

.footer-copy {
  font-size: 0.82rem;
  color: #94a3b8;
}
</style>
