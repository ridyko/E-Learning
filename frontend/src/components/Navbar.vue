<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  GraduationCap, 
  BookOpen, 
  Code, 
  CheckSquare, 
  ClipboardList, 
  ShieldCheck, 
  Sparkles,
  LogIn,
  LogOut,
  User,
  Users,
  UserPlus,
  Menu,
  X,
  Edit3,
  HelpCircle,
  Megaphone,
  Layers,
  Radio,
  ChevronDown,
  ChevronRight,
  Folder
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const currentUser = ref(null)
const courses = ref([])
const isMobileMenuOpen = ref(false)
const isUserDropdownOpen = ref(false)
const isSidebarCollapsed = ref(false)

// Expandable Menu Sections State
const expandedSections = ref({
  panelDosen: true,
  materi: true,
  akademik: true,
  navUtama: true,
  matkul: true,
  akademikMahasiswa: true
})

const toggleSection = (sectionKey) => {
  expandedSections.value[sectionKey] = !expandedSections.value[sectionKey]
}

const toggleSidebarToggle = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
  isMobileMenuOpen.value = !isMobileMenuOpen.value
  if (isSidebarCollapsed.value) {
    document.body.classList.add('sidebar-collapsed')
  } else {
    document.body.classList.remove('sidebar-collapsed')
  }
  window.dispatchEvent(new Event('sidebar-toggled'))
}

const toggleMobileMenu = () => {
  toggleSidebarToggle()
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

const toggleUserDropdown = () => {
  isUserDropdownOpen.value = !isUserDropdownOpen.value
}

const closeUserDropdown = () => {
  isUserDropdownOpen.value = false
}

const autoExpandActiveSection = (currentPath) => {
  if (['/', '/playground', '/presensi', '/admin', '/mata-kuliah', '/pengumuman'].includes(currentPath)) {
    expandedSections.value.panelDosen = true
    expandedSections.value.navUtama = true
  }
  if (currentPath.startsWith('/course/')) {
    expandedSections.value.materi = true
    expandedSections.value.matkul = true
  }
  if (['/live-quiz', '/bank-soal', '/quizzes', '/mahasiswa'].includes(currentPath)) {
    expandedSections.value.akademik = true
    expandedSections.value.akademikMahasiswa = true
  }
}

// Automatically close menus & expand active menu on route change
watch(() => route.path, (newPath) => {
  closeMobileMenu()
  closeUserDropdown()
  autoExpandActiveSection(newPath)
}, { immediate: true })

const fetchUserData = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      currentUser.value = JSON.parse(userStr)
    } catch {
      currentUser.value = null
    }
  } else {
    currentUser.value = null
  }
}

const fetchCourses = async () => {
  try {
    const res = await fetch('/api/courses')
    if (res.ok) {
      courses.value = await res.json()
    }
  } catch (err) {
    console.warn('Navbar fetch courses error:', err)
  }
}

onMounted(() => {
  fetchUserData()
  fetchCourses()
  window.addEventListener('auth-changed', fetchUserData)
  window.addEventListener('course-changed', fetchCourses)
})

const isDosen = computed(() => currentUser.value?.role === 'dosen')

const userInitial = computed(() => {
  if (!currentUser.value?.name) return 'R'
  return currentUser.value.name.trim().charAt(0).toUpperCase()
})

const handleLogout = () => {
  localStorage.removeItem('user')
  localStorage.removeItem('token')
  currentUser.value = null
  closeMobileMenu()
  closeUserDropdown()
  window.dispatchEvent(new Event('auth-changed'))
  router.push('/login')
}

// Student Menu Item Lists
const studentNavMain = computed(() => [
  { name: 'Beranda', path: '/', icon: GraduationCap },
  { name: 'Live Code Editor', path: '/playground', icon: Sparkles },
  { name: 'Presensi', path: '/presensi', icon: ClipboardList }
])

const studentNavCourses = computed(() => {
  return courses.value.map(c => ({
    name: c.name,
    code: c.code,
    path: `/course/${c.id}`,
    icon: (c.name.toLowerCase().includes('web') || c.id.includes('web')) ? Code : BookOpen
  }))
})

const studentNavAcad = computed(() => [
  { name: 'Live Quiz Interaktif', path: '/live-quiz', icon: Radio },
  { name: 'Kuis Online', path: '/quizzes', icon: CheckSquare }
])

// Dosen Menu Item Lists
const dosenNavLinks = computed(() => [
  { name: 'Beranda', path: '/', icon: GraduationCap },
  { name: 'Live Code Editor', path: '/playground', icon: Sparkles },
  { name: 'Presensi', path: '/presensi', icon: ClipboardList },
  { name: 'Panel Kontrol Dashboard', path: '/admin', icon: ShieldCheck, highlight: true },
  { name: 'Kelola Mata Kuliah', path: '/mata-kuliah', icon: Layers },
  { name: 'Terbitkan Pengumuman', path: '/pengumuman', icon: Megaphone }
])

const dosenCourseLinks = computed(() => {
  return courses.value.map(c => ({
    name: c.name,
    code: c.code,
    path: `/course/${c.id}`,
    icon: (c.name.toLowerCase().includes('web') || c.id.includes('web')) ? Code : BookOpen
  }))
})

const dosenAcadLinks = computed(() => [
  { name: 'Live Quiz Interaktif', path: '/live-quiz', icon: Radio },
  { name: 'Bank Soal & Nilai Kuis', path: '/bank-soal', icon: HelpCircle },
  { name: 'Data Mahasiswa', path: '/mahasiswa', icon: Users }
])
</script>

<template>
  <!-- TOP HEADER FOR GUESTS (PRE-LOGIN) -->
  <header v-if="!currentUser" class="guest-navbar">
    <div class="guest-container">
      <router-link to="/" class="brand-logo">
        <div class="logo-icon"><span>ITB</span></div>
        <div class="brand-text">
          <span class="institution">ITB SWADHARMA</span>
          <span class="portal-name">E-Learning Pak Rio</span>
        </div>
      </router-link>
      <div class="auth-buttons">
        <router-link to="/login" class="btn btn-gold btn-sm">
          <LogIn class="btn-icon-xs" />
          <span>Login Portal</span>
        </router-link>
        <router-link to="/register" class="btn btn-secondary btn-sm">
          <UserPlus class="btn-icon-xs" />
          <span>Daftar NIM</span>
        </router-link>
      </div>
    </div>
  </header>

  <!-- TOP HEADER FOR LOGGED IN USERS -->
  <header v-else class="top-header">
    <div class="top-header-container">
      <!-- Left: Toggle Sidebar Button -->
      <button @click="toggleMobileMenu" class="btn-menu-toggle" aria-label="Toggle Navigation">
        <X v-if="isMobileMenuOpen" class="toggle-icon" />
        <Menu v-else class="toggle-icon" />
      </button>

      <!-- Center: ITB Brand Logo -->
      <router-link to="/" class="brand-logo-center">
        <div class="logo-icon-sm"><span>ITB</span></div>
        <span class="portal-name-center">E-Learning Pak Rio</span>
      </router-link>

      <!-- Right: User Avatar Circle 'R' with Profile & Logout Dropdown -->
      <div class="user-avatar-wrapper">
        <button 
          @click="toggleUserDropdown" 
          class="user-initial-btn" 
          :class="{ 'active': isUserDropdownOpen }"
          title="Menu Profil & Keluar"
        >
          <span>{{ userInitial }}</span>
        </button>

        <!-- Transparent click outside backdrop for dropdown -->
        <div 
          v-if="isUserDropdownOpen" 
          class="dropdown-backdrop" 
          @click="closeUserDropdown"
        ></div>

        <!-- Dropdown Menu -->
        <div v-if="isUserDropdownOpen" class="user-dropdown-menu animate-fade-in">
          <div class="dropdown-header">
            <div class="dropdown-avatar-circle">
              <span>{{ userInitial }}</span>
            </div>
            <div class="dropdown-user-info">
              <span class="dropdown-name">{{ currentUser.name }}</span>
              <span class="dropdown-role">
                {{ isDosen ? 'Dosen (21099001)' : 'Mahasiswa (' + currentUser.username + ')' }}
              </span>
            </div>
          </div>

          <div class="dropdown-divider"></div>

          <router-link to="/profile" class="dropdown-item" @click="closeUserDropdown">
            <Edit3 class="dropdown-icon text-blue" />
            <span>Edit Profile</span>
          </router-link>

          <router-link v-if="isDosen" to="/admin" class="dropdown-item" @click="closeUserDropdown">
            <ShieldCheck class="dropdown-icon text-gold" />
            <span>Panel Kontrol Dashboard</span>
          </router-link>

          <router-link v-if="isDosen" to="/bank-soal" class="dropdown-item" @click="closeUserDropdown">
            <HelpCircle class="dropdown-icon text-gold" />
            <span>Bank Soal & Nilai Kuis</span>
          </router-link>

          <router-link v-if="isDosen" to="/mahasiswa" class="dropdown-item" @click="closeUserDropdown">
            <Users class="dropdown-icon text-blue" />
            <span>Data Mahasiswa</span>
          </router-link>

          <div class="dropdown-divider"></div>

          <button @click="handleLogout" class="dropdown-item dropdown-logout">
            <LogOut class="dropdown-icon text-red" />
            <span>Keluar / Logout</span>
          </button>
        </div>
      </div>
    </div>
  </header>

  <!-- BACKDROP OVERLAY FOR MOBILE MENU ONLY -->
  <div 
    v-if="currentUser && isMobileMenuOpen" 
    class="sidebar-backdrop" 
    @click="closeMobileMenu"
  ></div>

  <!-- FIXED LEFT SIDEBAR FOR NAVIGATION -->
  <aside 
    v-if="currentUser" 
    class="sidebar-wrapper"
    :class="{ 'mobile-open': isMobileMenuOpen, 'desktop-collapsed': isSidebarCollapsed }"
  >
    <div class="sidebar-container">
      <!-- Sidebar Brand Header -->
      <div class="sidebar-brand">
        <router-link to="/" class="brand-logo" @click="closeMobileMenu">
          <div class="logo-icon"><span>ITB</span></div>
          <div class="brand-text">
            <span class="institution">ITB SWADHARMA</span>
            <span class="portal-name">E-Learning Pak Rio</span>
          </div>
        </router-link>
      </div>

      <!-- Sidebar Navigation Menu with Accordion Menu & Submenu Structure -->
      <nav class="sidebar-nav">
        <!-- DOSEN MENU & SUBMENU -->
        <template v-if="isDosen">
          <!-- Menu 1: Panel Dosen -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('panelDosen')"
            >
              <div class="group-title">
                <ShieldCheck class="group-icon text-gold" />
                <span>PANEL DOSEN</span>
              </div>
              <ChevronDown v-if="expandedSections.panelDosen" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.panelDosen" class="submenu-container">
                <router-link 
                  v-for="link in dosenNavLinks" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  :class="{ 'admin-sidebar-link': link.highlight }"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon" :class="{ 'text-gold': link.highlight }" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- Menu 2: Upload & Kelola Materi -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('materi')"
            >
              <div class="group-title">
                <Folder class="group-icon text-blue" />
                <span>UPLOAD & KELOLA MATERI</span>
              </div>
              <ChevronDown v-if="expandedSections.materi" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.materi" class="submenu-container">
                <router-link 
                  v-for="link in dosenCourseLinks" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon text-gold" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- Menu 3: Fitur Akademik -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('akademik')"
            >
              <div class="group-title">
                <Radio class="group-icon text-gold" />
                <span>FITUR AKADEMIK</span>
              </div>
              <ChevronDown v-if="expandedSections.akademik" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.akademik" class="submenu-container">
                <router-link 
                  v-for="link in dosenAcadLinks" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>
        </template>

        <!-- MAHASISWA MENU & SUBMENU -->
        <template v-else>
          <!-- Menu 1: Navigasi Utama -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('navUtama')"
            >
              <div class="group-title">
                <GraduationCap class="group-icon text-blue" />
                <span>NAVIGASI UTAMA</span>
              </div>
              <ChevronDown v-if="expandedSections.navUtama" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.navUtama" class="submenu-container">
                <router-link 
                  v-for="link in studentNavMain" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- Menu 2: Mata Kuliah Saya -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('matkul')"
            >
              <div class="group-title">
                <BookOpen class="group-icon text-gold" />
                <span>MATA KULIAH SAYA</span>
              </div>
              <ChevronDown v-if="expandedSections.matkul" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.matkul" class="submenu-container">
                <router-link 
                  v-for="link in studentNavCourses" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon text-gold" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- Menu 3: Fitur Akademik -->
          <div class="nav-menu-group">
            <button 
              type="button" 
              class="nav-group-header" 
              @click="toggleSection('akademikMahasiswa')"
            >
              <div class="group-title">
                <Radio class="group-icon text-gold" />
                <span>FITUR AKADEMIK</span>
              </div>
              <ChevronDown v-if="expandedSections.akademikMahasiswa" class="chevron-icon" />
              <ChevronRight v-else class="chevron-icon" />
            </button>

            <transition name="submenu-slide">
              <div v-show="expandedSections.akademikMahasiswa" class="submenu-container">
                <router-link 
                  v-for="link in studentNavAcad" 
                  :key="link.path" 
                  :to="link.path"
                  class="sidebar-nav-item submenu-item"
                  active-class="active"
                  @click="closeMobileMenu"
                >
                  <component :is="link.icon" class="nav-icon" />
                  <span>{{ link.name }}</span>
                </router-link>
              </div>
            </transition>
          </div>
        </template>
      </nav>

      <!-- Sidebar Quick Profile Link Footer -->
      <div class="sidebar-footer">
        <router-link to="/profile" class="sidebar-profile-link" @click="closeMobileMenu">
          <Edit3 class="nav-icon text-blue" />
          <span>Pengaturan Profil & Sandi</span>
        </router-link>
      </div>
    </div>
  </aside>
</template>

<style scoped>
/* Guest Top Navbar */
.guest-navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid #e2e8f0;
  padding: 0.75rem 1.5rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
}

.guest-container {
  max-width: 1350px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* TOP HEADER FOR LOGGED IN USERS */
.top-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  height: 64px;
  background: #ffffff;
  border-bottom: 1.5px solid #e2e8f0;
  padding: 0 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
}

.top-header-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-menu-toggle {
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: var(--radius-sm);
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0f172a;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-menu-toggle:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.toggle-icon {
  width: 20px;
  height: 20px;
}

.brand-logo-center {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  text-decoration: none;
}

.logo-icon-sm {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #1d4ed8;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 800;
  font-size: 0.8rem;
  box-shadow: 0 2px 8px rgba(29, 78, 216, 0.25);
}

.portal-name-center {
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 1.05rem;
  color: #0f172a;
}

/* User Avatar Initial 'R' Button & Dropdown */
.user-avatar-wrapper {
  position: relative;
}

.user-initial-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #eff6ff;
  border: 2px solid #bfdbfe;
  color: #1d4ed8;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 1.05rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 6px rgba(29, 78, 216, 0.12);
}

.user-initial-btn:hover, .user-initial-btn.active {
  background: #1d4ed8;
  color: #ffffff;
  border-color: #1e40af;
  transform: scale(1.05);
}

.user-dropdown-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 260px;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.12);
  padding: 0.75rem;
  z-index: 20;
}

.dropdown-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
}

.dropdown-avatar-circle {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1d4ed8;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 1rem;
  flex-shrink: 0;
}

.dropdown-user-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.dropdown-name {
  font-size: 0.88rem;
  font-weight: 800;
  color: #0f172a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dropdown-role {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 600;
}

.dropdown-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 0.5rem 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.65rem 0.75rem;
  border-radius: var(--radius-sm);
  color: #334155;
  text-decoration: none;
  font-size: 0.88rem;
  font-weight: 700;
  background: transparent;
  border: none;
  width: 100%;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dropdown-item:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.dropdown-logout:hover {
  background: #fef2f2;
  color: #dc2626;
}

.dropdown-icon {
  width: 18px;
  height: 18px;
}

.text-red {
  color: #dc2626;
}

/* Mobile Sidebar Backdrop (No blur) */
.sidebar-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1040;
  background: rgba(15, 23, 42, 0.3);
  animation: fadeIn 0.2s ease-out;
}

@media (min-width: 992px) {
  .sidebar-backdrop {
    display: none !important;
  }
}

.dropdown-backdrop {
  position: fixed;
  inset: 0;
  z-index: 10;
  background: transparent;
}

/* Left Sidebar Wrapper (Desktop & Mobile Drawer) */
.sidebar-wrapper {
  position: fixed;
  top: 64px;
  left: 0;
  bottom: 0;
  width: 270px;
  z-index: 1050;
  background: #ffffff;
  border-right: 1.5px solid #e2e8f0;
  box-shadow: 4px 0 25px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-wrapper.desktop-collapsed {
  transform: translateX(-270px) !important;
}

.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 1.25rem 0.85rem;
  overflow-y: auto;
}

.sidebar-brand {
  padding-bottom: 1.25rem;
  border-bottom: 1.5px solid #f1f5f9;
  margin-bottom: 1rem;
  display: none;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
}

.logo-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #1d4ed8;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 0.9rem;
  box-shadow: 0 4px 12px rgba(29, 78, 216, 0.3);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.institution {
  font-size: 0.72rem;
  font-weight: 800;
  color: #d97706;
  letter-spacing: 0.08em;
}

.portal-name {
  font-family: var(--font-heading);
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
}

/* Navigation Links */
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

/* Menu & Submenu Groups */
.nav-menu-group {
  margin-bottom: 0.25rem;
}

.nav-group-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: transparent;
  border: 1px solid transparent;
  padding: 0.55rem 0.65rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: #475569;
  font-weight: 800;
  font-size: 0.74rem;
  letter-spacing: 0.06em;
  transition: all 0.2s ease;
  user-select: none;
}

.nav-group-header:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.group-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.group-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.chevron-icon {
  width: 16px;
  height: 16px;
  color: #94a3b8;
  transition: transform 0.2s ease;
}

.submenu-container {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  padding-left: 0.5rem;
  margin-top: 0.25rem;
  margin-left: 0.8rem;
  border-left: 2px solid #e2e8f0;
}

.submenu-item {
  display: flex;
  align-items: center;
  gap: 0.7rem;
  padding: 0.55rem 0.75rem;
  border-radius: 0.6rem;
  color: #475569;
  text-decoration: none;
  font-size: 0.87rem;
  font-weight: 600;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.submenu-item:hover {
  color: #1d4ed8;
  background: #f1f5f9;
}

.submenu-item.active {
  color: #1d4ed8;
  background: #eff6ff;
  border-color: #bfdbfe;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(29, 78, 216, 0.06);
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.admin-sidebar-link {
  color: #b45309;
  background: #fffbeb;
  border-color: #fde68a;
}

.admin-sidebar-link:hover {
  background: #fef3c7;
  color: #92400e;
}

.admin-sidebar-link.active {
  color: #78350f;
  background: #fef3c7;
  border-color: #f59e0b;
}

/* Sidebar Footer Profile Link */
.sidebar-footer {
  padding-top: 1rem;
  border-top: 1.5px solid #f1f5f9;
  margin-top: auto;
}

.sidebar-profile-link {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  padding: 0.65rem 0.85rem;
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-sm);
  color: #1d4ed8;
  text-decoration: none;
  font-size: 0.88rem;
  font-weight: 700;
  transition: all 0.2s ease;
}

.sidebar-profile-link:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.text-gold {
  color: #d97706;
}

.text-blue {
  color: #1d4ed8;
}

.auth-buttons {
  display: flex;
  gap: 0.5rem;
}

.btn-icon-xs {
  width: 16px;
  height: 16px;
}

/* Vue Submenu Slide Animation */
.submenu-slide-enter-active,
.submenu-slide-leave-active {
  transition: all 0.25s ease-out;
  max-height: 500px;
  overflow: hidden;
}

.submenu-slide-enter-from,
.submenu-slide-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(-4px);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* RESPONSIVE MEDIA QUERIES */
@media (max-width: 991px) {
  .sidebar-wrapper {
    top: 64px;
    transform: translateX(-100%);
  }

  .sidebar-wrapper.mobile-open {
    transform: translateX(0);
  }
}
</style>
