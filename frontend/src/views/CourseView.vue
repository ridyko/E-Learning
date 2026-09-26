<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  BookOpen, 
  Calendar, 
  CheckCircle, 
  FileText, 
  Download, 
  Upload, 
  Code,
  Lock,
  Unlock,
  Clock,
  Settings,
  AlertCircle,
  Check,
  Play,
  Edit3,
  PlusCircle,
  FileUp,
  Link,
  Video,
  Layers
} from 'lucide-vue-next'
import Swal from 'sweetalert2'
import { showSuccess, showError, showWarning } from '../utils/swal.js'

const route = useRoute()
const router = useRouter()
const course = ref(null)
const allCourses = ref([])
const loading = ref(true)
const activeMeeting = ref(1)

const fetchAllCourses = async () => {
  try {
    const res = await fetch('/api/courses')
    if (res.ok) {
      allCourses.value = await res.json()
    }
  } catch (err) {
    console.warn('Fetch all courses error:', err)
  }
}

const switchCourse = (courseId) => {
  if (course.value && course.value.id === courseId) return
  activeMeeting.value = 1
  router.push(`/course/${courseId}`)
}

const currentUser = ref(null)
const isDosen = computed(() => currentUser.value?.role === 'dosen')
const isMahasiswa = computed(() => currentUser.value?.role === 'mahasiswa')

const isCourseInactiveForStudent = computed(() => {
  if (isDosen.value) return false
  return course.value?.status === 'Non Aktif'
})

const visibleAllCourses = computed(() => {
  if (isDosen.value) return allCourses.value
  return allCourses.value.filter(c => c.status !== 'Non Aktif')
})

const profile = ref(null)
const fetchProfile = async () => {
  try {
    const res = await fetch('/api/profile')
    if (res.ok) {
      profile.value = await res.json()
    }
  } catch {}
}

const dosenDisplayName = computed(() => {
  if (!profile.value) return 'Rio Widyatmoko'
  const deg = (profile.value.degree && profile.value.degree.trim() !== '' && profile.value.degree.trim() !== '-')
    ? ', ' + profile.value.degree.trim()
    : ''
  return (profile.value.name || 'Rio Widyatmoko').trim() + deg
})

// Assignment state
const assignments = ref([])
const showAssignmentModal = ref(false)
const selectedModule = ref(null)
const studentNim = ref('')
const studentName = ref('')
const repoLink = ref('')
const studentNotes = ref('')
const submissionSuccess = ref(false)

// Edit / Upload Module Modal State (For Dosen)
const showEditModuleModal = ref(false)
const isSavingModule = ref(false)
const editModuleForm = ref({
  meeting_number: 1,
  title: '',
  description: '',
  topicsStr: '',
  slide_url: '',
  pdf_url: '',
  video_url: '',
  code_sample: '',
  task_due_date: '',
  scheduled_at: '',
  status: 'terbuka',
  is_active: true
})

const loadUserData = () => {
  const uStr = localStorage.getItem('user')
  if (uStr) {
    try {
      currentUser.value = JSON.parse(uStr)
      if (currentUser.value?.username && currentUser.value?.role === 'mahasiswa') {
        studentNim.value = currentUser.value.username
        studentName.value = currentUser.value.name
      }
    } catch {
      currentUser.value = null
    }
  }
}

const fetchAssignments = async () => {
  try {
    const res = await fetch('/api/assignments')
    if (res.ok) {
      assignments.value = await res.json()
    }
  } catch (err) {
    console.warn('Fetch assignments error:', err)
  }
}

const fetchCourse = async () => {
  loading.value = true
  try {
    const courseId = route.params.id || 'rpl-2026'
    const res = await fetch(`/api/courses/${courseId}`)
    course.value = await res.json()
  } catch (err) {
    console.warn('Fetch course error:', err)
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, () => {
  fetchCourse()
  fetchAssignments()
})

const studentCount = ref(0)
const fetchStudentsCount = async () => {
  try {
    const res = await fetch('/api/students')
    if (res.ok) {
      const list = await res.json()
      if (list && list.length > 0) {
        studentCount.value = list.length
      }
    }
  } catch {}
}

onMounted(() => {
  loadUserData()
  fetchCourse()
  fetchAssignments()
  fetchAllCourses()
  fetchStudentsCount()
  fetchProfile()
  window.addEventListener('course-changed', fetchAllCourses)
})

// Check if module access is open for students
const isModuleOpenForStudents = (mod) => {
  if (!mod) return false
  if (typeof mod.is_active === 'boolean') {
    return mod.is_active
  }
  if (mod.status === 'terkunci') return false
  if (mod.status === 'terbuka') return true
  return mod.meeting_number === 1
}

// Check if module is unlocked for current user in UI
const isModuleUnlocked = (mod) => {
  if (isDosen.value) return true // Dosen can view and manage modules in UI
  return isModuleOpenForStudents(mod)
}

// Get all student submissions for a specific meeting (for Dosen view)
const getMeetingSubmissions = (meetingNo) => {
  if (!assignments.value || !course.value) return []
  return assignments.value.filter(a => a.course_id === course.value.id && a.meeting_no === meetingNo)
}

// Get student's assignment for specific meeting
const getStudentAssignment = (meetingNo) => {
  if (!currentUser.value?.username || !assignments.value.length || !course.value) return null
  return assignments.value.find(a => 
    a.course_id === course.value.id && 
    a.meeting_no === meetingNo && 
    (a.student_nim === currentUser.value.username || currentUser.value.role === 'dosen')
  )
}

// Check if task deadline has passed
const isTaskClosed = (mod) => {
  if (!mod || !mod.task_due_date) return false
  if (mod.task_due_date.toLowerCase().includes('ditutup') || mod.task_due_date.toLowerCase().includes('berakhir')) {
    return true
  }
  try {
    const dateParsed = new Date(mod.task_due_date)
    if (!isNaN(dateParsed.getTime())) {
      return new Date() > dateParsed
    }
  } catch {
    return false
  }
  return false
}

// Dosen Controls: Toggle Module Unlock/Lock
const toggleModuleStatus = async (mod) => {
  const currentOpen = isModuleOpenForStudents(mod)
  const newStatus = !currentOpen
  const statusText = newStatus ? 'terbuka' : 'terkunci'

  try {
    const res = await fetch('/api/courses/module/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: course.value.id,
        meeting_number: mod.meeting_number,
        is_active: newStatus,
        status: statusText
      })
    })

    if (res.ok) {
      mod.is_active = newStatus
      mod.status = statusText
      if (newStatus) {
        showSuccess('Akses Modul Dibuka! 🟢', `Pertemuan ke-${mod.meeting_number} sekarang TERBUKA untuk mahasiswa.`)
      } else {
        showSuccess('Akses Modul Ditutup 🔒', `Pertemuan ke-${mod.meeting_number} berhasil DITUTUP untuk mahasiswa.`)
      }
    } else {
      showError('Gagal!', 'Gagal memperbarui status modul.')
    }
  } catch (err) {
    showError('Error!', 'Terjadi kesalahan pada server.')
  }
}

// Schedule Modal State
const showScheduleModal = ref(false)
const scheduleForm = ref({ meeting_number: 1, scheduled_at: '' })
const isSavingSchedule = ref(false)

const openScheduleModal = (mod) => {
  scheduleForm.value = {
    meeting_number: mod.meeting_number,
    scheduled_at: mod.scheduled_at || 'Senin, 08.00 WIB'
  }
  showScheduleModal.value = true
}

const saveSchedule = async () => {
  if (!scheduleForm.value.scheduled_at) {
    showWarning('Perhatian', 'Jadwal Buka Wajib Diisi!')
    return
  }
  isSavingSchedule.value = true
  try {
    const res = await fetch('/api/courses/module/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: course.value.id,
        meeting_number: scheduleForm.value.meeting_number,
        is_active: course.value.modules.find(m => m.meeting_number === scheduleForm.value.meeting_number)?.is_active || false,
        status: 'dijadwalkan',
        scheduled_at: scheduleForm.value.scheduled_at
      })
    })

    if (res.ok) {
      const mod = course.value.modules.find(m => m.meeting_number === scheduleForm.value.meeting_number)
      if (mod) {
        mod.scheduled_at = scheduleForm.value.scheduled_at
        mod.status = 'dijadwalkan'
      }
      showSuccess('Jadwal Disimpan! ⏰', `Modul akan otomatis dibuka pada: ${scheduleForm.value.scheduled_at}`)
      showScheduleModal.value = false
    } else {
      showError('Gagal!', 'Gagal menyimpan jadwal.')
    }
  } catch (err) {
    showError('Error', 'Gagal menyimpan jadwal.')
  } finally {
    isSavingSchedule.value = false
  }
}

// Deadline Modal State
const showDeadlineModal = ref(false)
const deadlineForm = ref({ meeting_number: 1, task_due_date: '' })
const isSavingDeadline = ref(false)

const openDeadlineModal = (mod) => {
  deadlineForm.value = {
    meeting_number: mod.meeting_number,
    task_due_date: mod.task_due_date || 'Senin, 28 Sep 2026 - 23:59 WIB'
  }
  showDeadlineModal.value = true
}

const saveDeadline = async () => {
  if (!deadlineForm.value.task_due_date) {
    showWarning('Perhatian', 'Batas Waktu Wajib Diisi!')
    return
  }
  isSavingDeadline.value = true
  try {
    const res = await fetch('/api/courses/module/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: course.value.id,
        meeting_number: deadlineForm.value.meeting_number,
        is_active: course.value.modules.find(m => m.meeting_number === deadlineForm.value.meeting_number)?.is_active || false,
        task_due_date: deadlineForm.value.task_due_date
      })
    })

    if (res.ok) {
      const mod = course.value.modules.find(m => m.meeting_number === deadlineForm.value.meeting_number)
      if (mod) {
        mod.task_due_date = deadlineForm.value.task_due_date
      }
      showSuccess('Deadline Disimpan! 📅', `Batas pengumpulan tugas diperbarui: ${deadlineForm.value.task_due_date}`)
      showDeadlineModal.value = false
    } else {
      showError('Gagal!', 'Gagal menyimpan deadline.')
    }
  } catch (err) {
    showError('Error', 'Gagal memperbarui deadline.')
  } finally {
    isSavingDeadline.value = false
  }
}

// Open Edit & Upload Module Modal for Dosen
const openEditModuleModal = (mod) => {
  editModuleForm.value = {
    meeting_number: mod.meeting_number,
    title: mod.title || `Pertemuan ${mod.meeting_number}`,
    description: mod.description || '',
    topicsStr: mod.topics ? mod.topics.join(', ') : '',
    slide_url: mod.slide_url || '',
    pdf_url: mod.pdf_url || '',
    video_url: mod.video_url || '',
    code_sample: mod.code_sample || '',
    task_due_date: mod.task_due_date || 'Senin, 28 Sep 2026 - 23:59 WIB',
    scheduled_at: mod.scheduled_at || '',
    status: mod.status || (mod.is_active ? 'terbuka' : 'terkunci'),
    is_active: mod.is_active
  }
  showEditModuleModal.value = true
}

// Real File Uploaders for Slide & PDF
const isUploadingSlide = ref(false)
const isUploadingPdf = ref(false)

const handleSlideFile = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  isUploadingSlide.value = true
  showSuccess('Mengunggah Slide... ⏳', `Sedang mengunggah file "${file.name}" ke server...`)
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res = await fetch('/api/upload', {
      method: 'POST',
      body: formData
    })
    const data = await res.json()
    if (res.ok && data.url) {
      editModuleForm.value.slide_url = data.url
      showSuccess('Slide Berhasil Diunggah! 📄', `File slide tersimpan di server.`)
    } else {
      showError('Gagal Upload', data.error || 'Gagal mengunggah slide')
    }
  } catch (err) {
    showError('Gagal Upload', 'Terjadi kesalahan koneksi saat upload.')
  } finally {
    isUploadingSlide.value = false
  }
}

const handlePdfFile = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  isUploadingPdf.value = true
  showSuccess('Mengunggah Modul PDF... ⏳', `Sedang mengunggah file "${file.name}" ke server...`)
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res = await fetch('/api/upload', {
      method: 'POST',
      body: formData
    })
    const data = await res.json()
    if (res.ok && data.url) {
      editModuleForm.value.pdf_url = data.url
      showSuccess('Modul PDF Berhasil Diunggah! 📄', `File modul tersimpan di server.`)
    } else {
      showError('Gagal Upload', data.error || 'Gagal mengunggah modul PDF')
    }
  } catch (err) {
    showError('Gagal Upload', 'Terjadi kesalahan koneksi saat upload.')
  } finally {
    isUploadingPdf.value = false
  }
}

// Save Full Module Content (Title, Description, Topics, Slide, PDF, Code Sample, Deadline)
const saveModuleContent = async () => {
  if (!editModuleForm.value.title || !editModuleForm.value.description) {
    showWarning('Perhatian', 'Judul dan Deskripsi Modul Wajib Diisi!')
    return
  }

  isSavingModule.value = true
  const topicsArray = editModuleForm.value.topicsStr
    .split(',')
    .map(t => t.trim())
    .filter(t => t.length > 0)

  const isActiveBool = editModuleForm.value.status === 'terbuka'

  try {
    const res = await fetch('/api/courses/module/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: course.value.id,
        meeting_number: editModuleForm.value.meeting_number,
        title: editModuleForm.value.title,
        description: editModuleForm.value.description,
        topics: topicsArray,
        slide_url: editModuleForm.value.slide_url,
        pdf_url: editModuleForm.value.pdf_url,
        video_url: editModuleForm.value.video_url,
        code_sample: editModuleForm.value.code_sample,
        task_due_date: editModuleForm.value.task_due_date,
        scheduled_at: editModuleForm.value.scheduled_at,
        status: editModuleForm.value.status,
        is_active: isActiveBool
      })
    })

    if (res.ok) {
      showSuccess('Materi & Deskripsi Disimpan! 🚀', `Modul Pertemuan ke-${editModuleForm.value.meeting_number} berhasil diunggah & diperbarui di portal.`)
      showEditModuleModal.value = false
      fetchCourse()
    } else {
      showError('Gagal!', 'Gagal menyimpan perubahan modul.')
    }
  } catch (err) {
    showError('Error', 'Terjadi kesalahan saat menyimpan modul.')
  } finally {
    isSavingModule.value = false
  }
}

// Open Submit Assignment Modal
const openSubmitModal = (mod) => {
  if (!isModuleUnlocked(mod)) {
    showWarning('Terkunci! 🔒', 'Modul ini belum dimulai oleh Dosen.')
    return
  }
  if (isTaskClosed(mod)) {
    showWarning('Batas Waktu Telah Berakhir! 🔴', 'Pengumpulan tugas untuk pertemuan ini sudah ditutup.')
    return
  }
  selectedModule.value = mod
  showAssignmentModal.value = true
  submissionSuccess.value = false
}

// Submit Assignment Handler
const submitAssignment = async () => {
  if (!studentNim.value || !studentName.value || !repoLink.value) {
    showWarning('Perhatian', 'Harap lengkapi NIM, Nama, dan Link Repository / File!')
    return
  }

  try {
    const res = await fetch('/api/assignments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: course.value.id,
        meeting_no: selectedModule.value.meeting_number,
        title: selectedModule.value.title,
        student_nim: studentNim.value,
        student_name: studentName.value,
        repo_link: repoLink.value,
        notes: studentNotes.value
      })
    })

    if (res.ok) {
      submissionSuccess.value = true
      showSuccess('Tugas Terkirim! 🚀', 'Tugas Anda telah berhasil diunggah ke portal Dosen.')
      fetchAssignments()
      setTimeout(() => {
        showAssignmentModal.value = false
      }, 1800)
    } else {
      showError('Gagal!', 'Gagal mengirim tugas!')
    }
  } catch (err) {
    showError('Gagal!', 'Gagal mengirim tugas!')
  }
}
</script>

<template>
  <div class="course-container animate-fade-in" v-if="course">
    <!-- Course Switcher Bar -->
    <div v-if="visibleAllCourses.length > 0" class="course-switcher-bar glass-card">
      <div class="switcher-header-row">
        <div class="switcher-label">
          <Layers class="switcher-icon text-gold" />
          <span>Pilih Mata Kuliah:</span>
        </div>
        <div class="switcher-tabs">
          <button 
            v-for="c in visibleAllCourses" 
            :key="c.id"
            @click="switchCourse(c.id)"
            :class="['switcher-tab-btn', (course && course.id === c.id) ? 'active' : '']"
          >
            <component :is="(c.name.toLowerCase().includes('web') || c.id.includes('web')) ? Code : BookOpen" class="tab-icon" />
            <span>{{ c.name }}</span>
            <span class="code-tag">{{ c.code }}</span>
            <span v-if="c.status === 'Non Aktif'" class="code-tag" style="background: #fee2e2; color: #dc2626;">Non-Aktif</span>
            <span v-else-if="course && course.id === c.id" class="active-indicator">Aktif</span>
          </button>
        </div>
        <router-link v-if="isDosen" to="/mata-kuliah" class="btn-manage-matkul" title="Kelola & Tambah Mata Kuliah Baru">
          <PlusCircle class="manage-icon" />
          <span>Kelola / Tambah Matkul</span>
        </router-link>
      </div>
    </div>

    <!-- Course Banner -->
    <div class="glass-card course-header">
      <div class="header-badges">
        <span class="badge badge-gold">{{ course.code }}</span>
        <span class="badge badge-blue">{{ course.sks }} SKS</span>
        <span :class="course.status === 'Non Aktif' ? 'badge badge-rose' : 'badge badge-emerald'">
          Status: {{ course.status || 'Aktif' }}
        </span>
        <span class="badge badge-emerald">Dosen: {{ dosenDisplayName }}</span>
      </div>
      <h2>{{ course.name }}</h2>
      <p class="header-desc">{{ course.description }}</p>

      <div class="meta-row">
        <span>📅 {{ course.class_time }}</span>
        <span>📍 {{ course.room }}</span>
        <span>👥 {{ studentCount || course.total_students }} Mahasiswa</span>
      </div>
    </div>

    <!-- INACTIVE COURSE WARNING BANNER FOR STUDENTS -->
    <div v-if="isCourseInactiveForStudent" class="glass-card" style="padding: 2.5rem; text-align: center; margin-bottom: 2rem; background: #fff5f5; border: 1.5px solid #fecaca; border-radius: 1rem;">
      <div style="width: 4rem; height: 4rem; background: #fee2e2; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem;">
        <Lock style="width: 2.2rem; height: 2.2rem; color: #dc2626;" />
      </div>
      <h3 style="font-size: 1.5rem; font-weight: 800; color: #991b1b; margin-bottom: 0.5rem;">Mata Kuliah Non-Aktif 🔒</h3>
      <p style="font-size: 1rem; color: #7f1d1d; max-width: 600px; margin: 0 auto 1.5rem; line-height: 1.6;">
        Mata kuliah <strong>{{ course.name }}</strong> saat ini berstatus <strong>Non-Aktif</strong> oleh Dosen Pengampu. Modul pertemuan, materi, dan tugas belum dapat diakses oleh Mahasiswa.
      </p>
      <router-link to="/" class="btn btn-gold">
        <span>Kembali ke Beranda Utama</span>
      </router-link>
    </div>

    <!-- Main Content Layout -->
    <div v-else class="course-layout">
      <!-- Left Sidebar: Meetings Navigation (1-14) -->
      <div class="meetings-sidebar glass-card">
        <div class="sidebar-header">
          <h3>Modul Perkuliahan (1-14)</h3>
          <span class="sub-hint">Klik pertemuan untuk memilih</span>
        </div>
        <div class="meeting-nav">
          <button 
            v-for="mod in course.modules" 
            :key="mod.meeting_number"
            @click="activeMeeting = mod.meeting_number"
            :class="[
              'meeting-nav-btn', 
              activeMeeting === mod.meeting_number ? 'active' : '',
              !isModuleOpenForStudents(mod) ? 'nav-locked' : 'nav-open'
            ]"
          >
            <div class="m-nav-info">
              <span class="m-num">P{{ mod.meeting_number }}</span>
              <span class="m-title">{{ mod.title }}</span>
            </div>
            
            <!-- Lock / Status Pill Badge -->
            <div class="m-status-pill">
              <span v-if="isModuleOpenForStudents(mod)" class="pill pill-active" title="Modul Terbuka Mahasiswa">
                <Unlock class="pill-icon" /> Terbuka
              </span>
              <span v-else-if="mod.scheduled_at" class="pill pill-scheduled" title="Dijadwalkan">
                <Clock class="pill-icon" /> Scheduled
              </span>
              <span v-else class="pill pill-locked" title="Modul Dikunci Mahasiswa">
                <Lock class="pill-icon" /> Terkunci
              </span>
            </div>
          </button>
        </div>
      </div>

      <!-- Right: Meeting Detail Content -->
      <div class="meeting-content" v-if="course.modules">
        <div 
          v-for="mod in course.modules" 
          :key="mod.meeting_number"
          v-show="activeMeeting === mod.meeting_number"
          class="glass-card detail-card"
        >
          <!-- DOSEN CONTROL PANEL BANNER (Only for Dosen) -->
          <div v-if="isDosen" class="dosen-panel-banner">
            <div class="dosen-panel-header">
              <div class="dosen-badge">
                <Settings class="d-icon" />
                <span>PANEL KONTROL DOSEN - PAK RIO</span>
              </div>
              <span class="module-state-tag" :class="isModuleOpenForStudents(mod) ? 'tag-open' : 'tag-locked'">
                Status Access: {{ isModuleOpenForStudents(mod) ? '🟢 TERBUKA UNTUK MAHASISWA' : '🔴 DITUTUP (TERKUNCI)' }}
              </span>
            </div>

            <div class="dosen-panel-actions">
              <!-- Upload & Edit Materi Modul Main Action Button -->
              <button @click="openEditModuleModal(mod)" class="btn btn-sm btn-gold btn-highlight">
                <Edit3 class="btn-icon-xs" />
                ✏️ Upload & Edit Materi Modul P{{ mod.meeting_number }}
              </button>

              <button 
                @click="toggleModuleStatus(mod)" 
                :class="['btn', 'btn-sm', isModuleOpenForStudents(mod) ? 'btn-danger-outline' : 'btn-emerald']"
              >
                <component :is="isModuleOpenForStudents(mod) ? Lock : Play" class="btn-icon-xs" />
                {{ isModuleOpenForStudents(mod) ? '🔴 Kunci Modul Kembali' : '🟢 Buka Akses Modul' }}
              </button>

              <button @click="openScheduleModal(mod)" class="btn btn-sm btn-blue-outline">
                <Clock class="btn-icon-xs" />
                Atur Jadwal Buka
              </button>

              <button @click="openDeadlineModal(mod)" class="btn btn-sm btn-amber-outline">
                <Calendar class="btn-icon-xs" />
                Atur Deadline Tugas
              </button>
            </div>
          </div>

          <!-- MEETING HEADER -->
          <div class="detail-header">
            <div class="detail-badges-row">
              <span class="badge badge-gold">Pertemuan ke-{{ mod.meeting_number }}</span>
              <span v-if="isModuleOpenForStudents(mod)" class="badge badge-emerald">
                <Unlock class="badge-icon" /> Status Akses: Terbuka
              </span>
              <span v-else class="badge badge-rose">
                <Lock class="badge-icon" /> Status Akses: Ditutup Dosen
              </span>
            </div>
            <h3>{{ mod.title }}</h3>
            <p class="mod-desc">{{ mod.description }}</p>
          </div>

          <!-- IF MODULE IS LOCKED (For Mahasiswa) -->
          <div v-if="!isDosen && !isModuleOpenForStudents(mod)" class="locked-module-notice">
            <div class="locked-icon-wrapper">
              <Lock class="locked-hero-icon" />
            </div>
            <h4>Modul Pertemuan Ke-{{ mod.meeting_number }} Ditutup / Belum Dimulai Dosen</h4>
            <p>
              Materi slide presentation, modul PDF, dan pengumpulan tugas untuk pertemuan ini 
              baru akan dibuka setelah Dosen membuka sesi perkuliahan atau sesuai jadwal otomatis.
            </p>

            <div v-if="mod.scheduled_at" class="scheduled-banner-info">
              <Clock class="sch-icon" />
              <span>Jadwal Pembukaan Modul: <strong>{{ mod.scheduled_at }}</strong></span>
            </div>
          </div>

          <!-- TOPICS & CPMK LIST (Visible when Unlocked for Mahasiswa or always for Dosen) -->
          <div v-if="isDosen || isModuleOpenForStudents(mod)" class="topics-box">
            <h4>📌 Pokok Bahasan & CPMK Pertemuan Ini:</h4>
            <div class="topics-grid">
              <span v-for="(tp, idx) in mod.topics" :key="idx" class="topic-chip">
                <CheckCircle class="chip-icon text-emerald" />
                {{ tp }}
              </span>
            </div>
          </div>

          <!-- CONTENT MATERIALS (Visible when Unlocked for Mahasiswa or always for Dosen) -->
          <template v-if="isDosen || isModuleOpenForStudents(mod)">
            <!-- Code Sample Preview -->
            <div v-if="mod.code_sample" class="code-box">
              <div class="code-header">
                <Code class="code-icon" />
                <span>Contoh Kode Praktikum</span>
              </div>
              <pre class="code-block"><code>{{ mod.code_sample }}</code></pre>
            </div>

            <!-- Material Download Resources -->
            <div class="resources-row">
              <template v-if="mod.slide_url && mod.slide_url !== '#'">
                <a :href="mod.slide_url" target="_blank" class="btn btn-secondary" title="Download / Buka Slide Presentation">
                  <FileText class="btn-icon-xs text-gold" />
                  <span>Download Slide Presentation</span>
                </a>
              </template>
              <template v-else>
                <button v-if="isDosen" @click="openEditModuleModal(mod)" class="btn btn-secondary btn-dashed" title="Klik untuk mengunggah materi slide">
                  <PlusCircle class="btn-icon-xs text-gold" />
                  <span>+ Pasang Slide Presentasi</span>
                </button>
                <span v-else class="resource-pending-pill">
                  <FileText class="pill-icon text-muted" />
                  <span>Slide belum diunggah Dosen</span>
                </span>
              </template>

              <template v-if="mod.pdf_url && mod.pdf_url !== '#'">
                <a :href="mod.pdf_url" target="_blank" class="btn btn-secondary" title="Download / Buka Modul PDF">
                  <Download class="btn-icon-xs text-blue" />
                  <span>Download Modul PDF</span>
                </a>
              </template>
              <template v-else>
                <button v-if="isDosen" @click="openEditModuleModal(mod)" class="btn btn-secondary btn-dashed" title="Klik untuk mengunggah modul praktikum">
                  <PlusCircle class="btn-icon-xs text-blue" />
                  <span>+ Pasang Modul PDF</span>
                </button>
                <span v-else class="resource-pending-pill">
                  <Download class="pill-icon text-muted" />
                  <span>Modul PDF belum diunggah Dosen</span>
                </span>
              </template>

              <a v-if="mod.video_url && mod.video_url.trim() !== ''" :href="mod.video_url" target="_blank" class="btn btn-secondary">
                <Video class="btn-icon-xs text-rose" />
                <span>Tonton Video Tutorial</span>
              </a>
            </div>

            <!-- ASSIGNMENT / TUGAS SECTION -->
            <div class="assignment-section glass-card">
              <div class="asg-header">
                <div class="asg-title">
                  <Upload class="asg-icon" />
                  <h4>Pengumpulan Tugas Pertemuan {{ mod.meeting_number }}</h4>
                </div>

                <div class="asg-deadline-tag">
                  <Clock class="deadline-icon" />
                  <span>Batas Waktu: <strong>{{ mod.task_due_date || 'Senin, 28 Sep 2026 - 23:59 WIB' }}</strong></span>
                </div>
              </div>

              <!-- DOSEN VIEW: REKAP PENGUMPULAN TUGAS MAHASISWA -->
              <div v-if="isDosen" class="dosen-assignment-view">
                <div class="dosen-asg-summary">
                  <span class="sub-count-badge">
                    👥 {{ getMeetingSubmissions(mod.meeting_number).length }} dari {{ studentCount || course.total_students }} Mahasiswa Sudah Mengumpulkan
                  </span>
                </div>

                <div v-if="getMeetingSubmissions(mod.meeting_number).length > 0" class="table-responsive">
                  <table class="dosen-submissions-table">
                    <thead>
                      <tr>
                        <th>No</th>
                        <th>NIM</th>
                        <th>Nama Mahasiswa</th>
                        <th>Waktu Pengumpulan</th>
                        <th>Link Repository / File</th>
                        <th>Catatan</th>
                        <th>Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="(sub, index) in getMeetingSubmissions(mod.meeting_number)" :key="sub.id || index">
                        <td>{{ index + 1 }}</td>
                        <td><code class="nim-code">{{ sub.student_nim }}</code></td>
                        <td class="student-name-td"><strong>{{ sub.student_name }}</strong></td>
                        <td class="time-td">
                          {{ sub.submitted_at ? new Date(sub.submitted_at).toLocaleString('id-ID') : '-' }}
                        </td>
                        <td>
                          <a :href="sub.repo_link" target="_blank" class="repo-link-btn" title="Buka Link Tugas">
                            <Link class="link-icon-sm" />
                            Buka Link Tugas
                          </a>
                        </td>
                        <td class="notes-td">{{ sub.notes || '-' }}</td>
                        <td>
                          <span class="badge badge-emerald">Terkumpul</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div v-else class="empty-submissions-box">
                  <AlertCircle class="empty-sub-icon" />
                  <div>
                    <h5>Belum Ada Mahasiswa Mengumpulkan</h5>
                    <p>Belum ada hasil pengumpulan tugas dari mahasiswa untuk Pertemuan ke-{{ mod.meeting_number }}.</p>
                  </div>
                </div>
              </div>

              <!-- MAHASISWA VIEW: STATUS & FORM KUMPUL TUGAS -->
              <template v-else>
                <!-- State 1: Student already submitted -->
                <div v-if="getStudentAssignment(mod.meeting_number)" class="asg-status-box status-submitted">
                  <div class="status-left">
                    <CheckCircle class="st-icon text-emerald" />
                    <div>
                      <h5>Tugas Terkumpul (Tepat Waktu) ✅</h5>
                      <p class="st-desc">
                        Link Repository: <a :href="getStudentAssignment(mod.meeting_number).repo_link" target="_blank" class="repo-link-text">{{ getStudentAssignment(mod.meeting_number).repo_link }}</a>
                      </p>
                      <span class="st-meta">Dikirim pada: {{ new Date(getStudentAssignment(mod.meeting_number).submitted_at).toLocaleString('id-ID') }}</span>
                    </div>
                  </div>
                  <button @click="openSubmitModal(mod)" class="btn btn-sm btn-secondary">
                    ✏️ Update Link Tugas
                  </button>
                </div>

                <!-- State 2: Deadline passed and not submitted -->
                <div v-else-if="isTaskClosed(mod)" class="asg-status-box status-closed">
                  <div class="status-left">
                    <AlertCircle class="st-icon text-rose" />
                    <div>
                      <h5>Pengumpulan Tugas Ditutup 🔴</h5>
                      <p class="st-desc">Batas waktu pengumpulan telah berakhir pada {{ mod.task_due_date }}. Hubungi Dosen jika memerlukan perpanjangan waktu.</p>
                    </div>
                  </div>
                  <button disabled class="btn btn-secondary btn-disabled">
                    🔒 Pengumpulan Ditutup (Lewat Deadline)
                  </button>
                </div>

                <!-- State 3: Open & Ready for Submission -->
                <div v-else class="asg-status-box status-open">
                  <div class="status-left">
                    <Calendar class="st-icon text-blue" />
                    <div>
                      <h5>Status Pengumpulan: TERBUKA 🟢</h5>
                      <p class="st-desc">Silakan unggah link repository GitHub / Google Drive hasil tugas pengerjaan Anda sebelum deadline.</p>
                    </div>
                  </div>

                  <button @click="openSubmitModal(mod)" class="btn btn-gold">
                    <Upload class="btn-icon-xs" />
                    Kumpul Tugas Pertemuan Ini
                  </button>
                </div>
              </template>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- MODAL 1: UPLOAD & EDIT MATERI MODUL (DOSEN PAK RIO) -->
    <Teleport to="body">
      <div v-if="showEditModuleModal" class="modal-overlay" @click.self="showEditModuleModal = false">
        <div class="glass-card modal-card modal-lg">
          <div class="modal-header-banner">
            <div class="modal-title-box">
              <Edit3 class="modal-title-icon text-gold" />
              <div>
                <h3>Upload & Edit Materi Pertemuan ke-{{ editModuleForm.meeting_number }}</h3>
                <p class="modal-sub">Kelola Judul, Deskripsi, Slide, PDF, Video Tutorial, & Tugas Modul</p>
              </div>
            </div>
            <button @click="showEditModuleModal = false" class="btn-close-modal">✕</button>
          </div>

          <form @submit.prevent="saveModuleContent" class="modal-form-grid">
            <!-- Left Column -->
            <div class="form-col">
              <div class="form-group">
                <label class="input-label">Judul Modul Perkuliahan *</label>
                <input v-model="editModuleForm.title" class="glass-input" required placeholder="Contoh: Pengantar Rekayasa Perangkat Lunak" />
              </div>

              <div class="form-group">
                <label class="input-label">Deskripsi Lengkap Modul *</label>
                <textarea v-model="editModuleForm.description" class="glass-input textarea-large" required placeholder="Jelaskan cakupan materi dan tujuan pembelajaran..."></textarea>
              </div>

              <div class="form-group">
                <label class="input-label">Pokok Bahasan & CPMK (Pisahkan dengan Koma)</label>
                <input v-model="editModuleForm.topicsStr" class="glass-input" placeholder="Definisi RPL, SDLC Waterfall, Agile Scrum, SRS" />
              </div>

              <div class="form-group">
                <label class="input-label">Status Akses Mahasiswa</label>
                <select v-model="editModuleForm.status" class="glass-input">
                  <option value="terbuka">🟢 Terbuka (Mahasiswa Bisa Mengakses)</option>
                  <option value="terkunci">🔒 Terkunci (Belum Dimulai)</option>
                  <option value="dijadwalkan">⏳ Dijadwalkan (Buka Otomatis)</option>
                </select>
              </div>
            </div>

            <!-- Right Column -->
            <div class="form-col">
              <div class="form-group">
                <label class="input-label">Upload / Link Slide Presentation (PDF/PPTX)</label>
                <div class="upload-input-group">
                  <input v-model="editModuleForm.slide_url" class="glass-input" placeholder="https://drive.google.com/... atau pilih file untuk upload" />
                  <label class="btn-upload-file" :class="{ 'btn-uploading': isUploadingSlide }">
                    <FileUp class="up-icon" /> {{ isUploadingSlide ? 'Mengunggah...' : 'Pilih & Upload File' }}
                    <input type="file" @change="handleSlideFile" :disabled="isUploadingSlide" hidden accept=".pdf,.pptx,.ppt,.zip" />
                  </label>
                </div>
              </div>

              <div class="form-group">
                <label class="input-label">Upload / Link Modul PDF Praktikum</label>
                <div class="upload-input-group">
                  <input v-model="editModuleForm.pdf_url" class="glass-input" placeholder="https://drive.google.com/... atau pilih file untuk upload" />
                  <label class="btn-upload-file" :class="{ 'btn-uploading': isUploadingPdf }">
                    <FileUp class="up-icon" /> {{ isUploadingPdf ? 'Mengunggah...' : 'Pilih & Upload File' }}
                    <input type="file" @change="handlePdfFile" :disabled="isUploadingPdf" hidden accept=".pdf,.doc,.docx" />
                  </label>
                </div>
              </div>

              <div class="form-group">
                <label class="input-label">Link Video Tutorial (Youtube / Drive)</label>
                <input v-model="editModuleForm.video_url" class="glass-input" placeholder="https://www.youtube.com/watch?v=..." />
              </div>

              <div class="form-group">
                <label class="input-label">Batas Waktu Pengumpulan Tugas (Deadline)</label>
                <input v-model="editModuleForm.task_due_date" class="glass-input" placeholder="Senin, 28 Sep 2026 - 23:59 WIB" />
              </div>

              <div class="form-group">
                <label class="input-label">Contoh Kode Praktikum (Opsional)</label>
                <textarea v-model="editModuleForm.code_sample" class="glass-input code-textarea" placeholder="Paste kode contoh praktikum di sini..."></textarea>
              </div>
            </div>

            <!-- Full Width Actions -->
            <div class="modal-actions-full">
              <button type="button" @click="showEditModuleModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" :disabled="isSavingModule" class="btn btn-gold btn-lg">
                <Check class="btn-icon-xs" />
                {{ isSavingModule ? 'Menyimpan...' : 'Simpan & Publikasikan Modul' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- MODAL 2: ASSIGNMENT SUBMISSION MODAL (STUDENT) -->
    <Teleport to="body">
      <div v-if="showAssignmentModal" class="modal-overlay" @click.self="showAssignmentModal = false">
        <div class="glass-card modal-card">
          <div class="modal-header">
            <h3>Upload Tugas Pertemuan {{ selectedModule?.meeting_number }}</h3>
            <p class="modal-sub">{{ selectedModule?.title }}</p>
          </div>

          <div v-if="submissionSuccess" class="alert-success">
            ✅ Tugas berhasil dikirim ke Dosen!
          </div>

          <form v-else @submit.prevent="submitAssignment" class="modal-form">
            <div class="form-group">
              <label class="input-label">NIM Mahasiswa</label>
              <input v-model="studentNim" class="glass-input" required placeholder="Contoh: 20260801001" />
            </div>

            <div class="form-group">
              <label class="input-label">Nama Lengkap Mahasiswa</label>
              <input v-model="studentName" class="glass-input" required placeholder="Nama anda" />
            </div>

            <div class="form-group">
              <label class="input-label">Link Repository GitHub / Drive Tugas</label>
              <input v-model="repoLink" class="glass-input" required placeholder="https://github.com/username/project" />
            </div>

            <div class="form-group">
              <label class="input-label">Catatan Tambahan untuk Dosen</label>
              <textarea v-model="studentNotes" class="glass-input textarea" placeholder="Catatan atau kendala pengerjaan..."></textarea>
            </div>

            <div class="modal-actions">
              <button type="button" @click="showAssignmentModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" class="btn btn-gold">Kirimkan Tugas</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- MODAL 3: ATUR JADWAL BUKA MODUL -->
    <Teleport to="body">
      <div v-if="showScheduleModal" class="modal-overlay" @click.self="showScheduleModal = false">
        <div class="glass-card modal-card">
          <div class="modal-header-banner">
            <div class="modal-title-box">
              <Clock class="modal-title-icon text-blue" />
              <div>
                <h3>Atur Jadwal Buka Modul</h3>
                <p class="modal-sub">Pertemuan ke-{{ scheduleForm.meeting_number }}</p>
              </div>
            </div>
            <button @click="showScheduleModal = false" class="btn-close-modal">✕</button>
          </div>

          <form @submit.prevent="saveSchedule" class="modal-form">
            <div class="form-group">
              <label class="input-label">Jadwal Pembukaan Otomatis *</label>
              <input 
                v-model="scheduleForm.scheduled_at" 
                class="glass-input" 
                required 
                placeholder="Contoh: Senin, 28 Sep 2026 - 08:00 WIB" 
              />
              <span class="help-text">Mahasiswa dapat mengakses modul secara otomatis pada jadwal ini.</span>
            </div>

            <!-- Quick Preset Buttons -->
            <div class="preset-group">
              <span class="preset-label">Rekomendasi Waktu:</span>
              <div class="preset-buttons">
                <button type="button" @click="scheduleForm.scheduled_at = 'Senin, 08.00 WIB'" class="btn-preset">
                  Senin, 08.00 WIB
                </button>
                <button type="button" @click="scheduleForm.scheduled_at = 'Rabu, 08.00 WIB'" class="btn-preset">
                  Rabu, 08.00 WIB
                </button>
                <button type="button" @click="scheduleForm.scheduled_at = 'Setiap Jam Kuliah'" class="btn-preset">
                  Jam Kuliah
                </button>
              </div>
            </div>

            <div class="modal-actions">
              <button type="button" @click="showScheduleModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" :disabled="isSavingSchedule" class="btn btn-blue">
                <Check class="btn-icon-xs" />
                {{ isSavingSchedule ? 'Menyimpan...' : 'Simpan Jadwal' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- MODAL 4: ATUR DEADLINE TUGAS -->
    <Teleport to="body">
      <div v-if="showDeadlineModal" class="modal-overlay" @click.self="showDeadlineModal = false">
        <div class="glass-card modal-card">
          <div class="modal-header-banner">
            <div class="modal-title-box">
              <Calendar class="modal-title-icon text-gold" />
              <div>
                <h3>Atur Batas Waktu Tugas (Deadline)</h3>
                <p class="modal-sub">Pertemuan ke-{{ deadlineForm.meeting_number }}</p>
              </div>
            </div>
            <button @click="showDeadlineModal = false" class="btn-close-modal">✕</button>
          </div>

          <form @submit.prevent="saveDeadline" class="modal-form">
            <div class="form-group">
              <label class="input-label">Batas Waktu Pengumpulan (Deadline) *</label>
              <input 
                v-model="deadlineForm.task_due_date" 
                class="glass-input" 
                required 
                placeholder="Contoh: Senin, 28 Sep 2026 - 23:59 WIB" 
              />
              <span class="help-text">Batas waktu pengumpulan menentukan kapan tugas dikunci bagi mahasiswa.</span>
            </div>

            <!-- Quick Preset Buttons -->
            <div class="preset-group">
              <span class="preset-label">Pilihan Cepat Deadline:</span>
              <div class="preset-buttons">
                <button type="button" @click="deadlineForm.task_due_date = 'Senin, 28 Sep 2026 - 23:59 WIB'" class="btn-preset">
                  Senin (23:59 WIB)
                </button>
                <button type="button" @click="deadlineForm.task_due_date = '1 Minggu Setelah Perkuliahan'" class="btn-preset">
                  1 Minggu
                </button>
                <button type="button" @click="deadlineForm.task_due_date = 'Pengumpulan Ditutup'" class="btn-preset btn-preset-danger">
                  Tutup Tugas
                </button>
              </div>
            </div>

            <div class="modal-actions">
              <button type="button" @click="showDeadlineModal = false" class="btn btn-secondary">Batal</button>
              <button type="submit" :disabled="isSavingDeadline" class="btn btn-gold">
                <Check class="btn-icon-xs" />
                {{ isSavingDeadline ? 'Menyimpan...' : 'Simpan Batas Waktu' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.course-container {
  max-width: 1350px;
  margin: 0 auto;
  padding: 1.5rem;
}

/* Course Switcher Bar at Top */
.course-switcher-bar {
  padding: 0.85rem 1.25rem;
  margin-bottom: 1.25rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-sm);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.02);
}

.switcher-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.switcher-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 0.9rem;
  color: #0f172a;
}

.switcher-icon {
  width: 18px;
  height: 18px;
}

.switcher-tabs {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  flex-wrap: wrap;
  flex: 1;
}

.switcher-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.45rem 0.9rem;
  border-radius: 8px;
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  color: #334155;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.switcher-tab-btn:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #1d4ed8;
}

.switcher-tab-btn.active {
  background: #1d4ed8;
  border-color: #1e40af;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(29, 78, 216, 0.25);
}

.switcher-tab-btn.active .code-tag {
  background: rgba(255, 255, 255, 0.25);
  color: #ffffff;
}

.tab-icon {
  width: 16px;
  height: 16px;
}

.code-tag {
  font-size: 0.72rem;
  font-weight: 800;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  background: #e2e8f0;
  color: #475569;
}

.active-indicator {
  font-size: 0.68rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: #10b981;
  color: #ffffff;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
}

.btn-manage-matkul {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.85rem;
  background: #fffbeb;
  border: 1.5px solid #fde68a;
  border-radius: 8px;
  color: #b45309;
  font-size: 0.82rem;
  font-weight: 700;
  text-decoration: none;
  transition: all 0.2s ease;
}

.btn-manage-matkul:hover {
  background: #fef3c7;
  border-color: #f59e0b;
  color: #92400e;
}

.manage-icon {
  width: 15px;
  height: 15px;
}

.course-header {
  padding: 2rem;
  margin-bottom: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.header-badges {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  flex-wrap: wrap;
}

.course-header h2 {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  color: #0f172a;
}

.header-desc {
  color: #475569;
  font-size: 0.95rem;
  margin-bottom: 1.25rem;
}

.meta-row {
  display: flex;
  gap: 1.5rem;
  font-size: 0.88rem;
  color: #64748b;
  font-weight: 500;
}

.course-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 1.5rem;
}

.meetings-sidebar {
  padding: 1.25rem;
  height: max-content;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.sidebar-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.sub-hint {
  font-size: 0.75rem;
  color: #94a3b8;
  display: block;
  margin-bottom: 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.meeting-nav {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 620px;
  overflow-y: auto;
  padding-right: 0.25rem;
}

.meeting-nav-btn {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-sm);
  color: #475569;
  cursor: pointer;
  text-align: left;
  transition: var(--transition);
}

.meeting-nav-btn:hover {
  background: #f8fafc;
  color: #1d4ed8;
}

.meeting-nav-btn.active {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #1d4ed8;
  font-weight: 700;
}

.m-nav-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  overflow: hidden;
}

.m-num {
  font-family: var(--font-heading);
  font-weight: 800;
  font-size: 0.75rem;
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  flex-shrink: 0;
}

.m-title {
  font-size: 0.82rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Status Pills in Sidebar */
.pill {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  font-size: 0.68rem;
  padding: 0.15rem 0.4rem;
  border-radius: 99px;
  font-weight: 700;
  text-transform: uppercase;
}

.pill-icon {
  width: 10px;
  height: 10px;
}

.pill-active {
  background: #dcfce7;
  color: #15803d;
  border: 1px solid #86efac;
}

.pill-scheduled {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
}

.pill-locked {
  background: #f1f5f9;
  color: #64748b;
  border: 1px solid #cbd5e1;
}

.detail-card {
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

/* Dosen Panel Control Banner */
.dosen-panel-banner {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-left: 4px solid #1d4ed8;
  border-radius: var(--radius-sm);
  padding: 1rem 1.25rem;
  margin-bottom: 1.5rem;
}

.dosen-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.dosen-badge {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.82rem;
  font-weight: 800;
  color: #1e3a8a;
  letter-spacing: 0.5px;
}

.d-icon {
  width: 16px;
  height: 16px;
}

.module-state-tag {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.tag-open {
  background: #dcfce7;
  color: #15803d;
}

.tag-locked {
  background: #ffe4e6;
  color: #be123c;
}

.dosen-panel-actions {
  display: flex;
  gap: 0.6rem;
  flex-wrap: wrap;
}

.btn-sm {
  padding: 0.4rem 0.75rem;
  font-size: 0.8rem;
}

.btn-highlight {
  box-shadow: 0 4px 12px rgba(180, 83, 9, 0.15);
  font-weight: 700;
}

.btn-emerald {
  background: #059669;
  color: white;
  border: none;
}
.btn-emerald:hover {
  background: #047857;
}

.btn-danger-outline {
  background: #fff1f2;
  border: 1px solid #fecdd3;
  color: #be123c;
}
.btn-danger-outline:hover {
  background: #ffe4e6;
}

.btn-blue-outline {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1d4ed8;
}
.btn-blue-outline:hover {
  background: #dbeafe;
}

.btn-amber-outline {
  background: #fffbeb;
  border: 1px solid #fde68a;
  color: #b45309;
}
.btn-amber-outline:hover {
  background: #fef3c7;
}

.detail-badges-row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.badge-rose {
  background: #ffe4e6;
  color: #be123c;
  border: 1px solid #fecdd3;
}

.badge-icon {
  width: 12px;
  height: 12px;
  display: inline-block;
  margin-right: 0.2rem;
}

.detail-header h3 {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.5rem;
}

.mod-desc {
  color: #475569;
  font-size: 0.95rem;
  margin-bottom: 1.5rem;
}

/* Locked Notice Box for Student */
.locked-module-notice {
  background: #fef2f2;
  border: 1px solid #fecdd3;
  border-radius: var(--radius-sm);
  padding: 2rem;
  text-align: center;
  margin-bottom: 1.5rem;
}

.locked-icon-wrapper {
  width: 56px;
  height: 56px;
  background: #ffe4e6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
}

.locked-hero-icon {
  width: 28px;
  height: 28px;
  color: #e11d48;
}

.locked-module-notice h4 {
  font-size: 1.15rem;
  font-weight: 800;
  color: #9f1239;
  margin-bottom: 0.5rem;
}

.locked-module-notice p {
  font-size: 0.9rem;
  color: #be123c;
  max-width: 600px;
  margin: 0 auto 1rem;
}

.scheduled-banner-info {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #ffffff;
  border: 1px solid #fca5a5;
  padding: 0.4rem 1rem;
  border-radius: 99px;
  font-size: 0.85rem;
  color: #9f1239;
}

.sch-icon {
  width: 16px;
  height: 16px;
}

.topics-box {
  background: #f8fafc;
  padding: 1.25rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
  margin-bottom: 1.5rem;
}

.topics-box h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 0.75rem;
}

.topics-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.topic-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 0.4rem 0.75rem;
  border-radius: 99px;
  font-size: 0.85rem;
  color: #334155;
  font-weight: 500;
}

.chip-icon {
  width: 14px;
  height: 14px;
}

.code-box {
  background: #0f172a;
  border-radius: var(--radius-sm);
  border: 1px solid #334155;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.code-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1rem;
  background: #1e293b;
  font-size: 0.85rem;
  font-weight: 700;
  color: #fbbf24;
}

.code-icon {
  width: 16px;
  height: 16px;
}

.code-block {
  padding: 1rem;
  font-family: var(--font-code);
  font-size: 0.88rem;
  color: #38bdf8;
  overflow-x: auto;
}

.resources-row {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 2rem;
}

/* Assignment Section Container */
.assignment-section {
  padding: 1.5rem;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-top: 4px solid #b45309;
  border-radius: var(--radius-sm);
}

.asg-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.asg-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.asg-icon {
  width: 20px;
  height: 20px;
  color: #b45309;
}

.asg-title h4 {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.asg-deadline-tag {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  padding: 0.35rem 0.8rem;
  border-radius: 99px;
  font-size: 0.85rem;
  color: #92400e;
}

.deadline-icon {
  width: 14px;
  height: 14px;
}

/* Assignment Status Cards */
.asg-status-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border-radius: var(--radius-sm);
  gap: 1rem;
  flex-wrap: wrap;
}

.status-left {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.st-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
  margin-top: 0.1rem;
}

.asg-status-box h5 {
  font-size: 0.98rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 0.25rem;
}

.st-desc {
  font-size: 0.85rem;
  color: #475569;
  margin-bottom: 0.25rem;
}

.st-meta {
  font-size: 0.75rem;
  color: #64748b;
}

.repo-link-text {
  color: #1d4ed8;
  font-weight: 600;
  text-decoration: underline;
}

.status-submitted {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
}

.status-closed {
  background: #fff1f2;
  border: 1px solid #fecdd3;
}

.status-open {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
}

.btn-disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Modal Styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(6px);
  z-index: 200;
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
  border: 1px solid #e2e8f0;
}

.modal-lg {
  max-width: 820px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header-banner {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.modal-title-box {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.modal-title-icon {
  width: 24px;
  height: 24px;
  margin-top: 0.2rem;
}

.btn-close-modal {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #64748b;
  cursor: pointer;
  padding: 0.2rem 0.5rem;
}

.btn-close-modal:hover {
  color: #0f172a;
}

.modal-form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

.form-col {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.textarea-large {
  min-height: 110px;
}

.code-textarea {
  min-height: 90px;
  font-family: var(--font-code);
  font-size: 0.82rem;
}

.upload-input-group {
  display: flex;
  gap: 0.5rem;
}

.btn-upload-file {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  white-space: nowrap;
}

.btn-upload-file:hover {
  background: #e2e8f0;
}

.up-icon {
  width: 14px;
  height: 14px;
}

.modal-actions-full {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}

.btn-lg {
  padding: 0.65rem 1.5rem;
  font-size: 0.92rem;
}

.modal-header h3 {
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
}

.modal-sub {
  color: #475569;
  font-size: 0.9rem;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  margin-bottom: 0.35rem;
  color: #334155;
  font-weight: 600;
}

.textarea {
  min-height: 80px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}

.alert-success {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #047857;
  padding: 1rem;
  border-radius: var(--radius-sm);
  text-align: center;
  font-weight: 700;
}

.text-emerald {
  color: #059669;
}

.text-rose {
  color: #e11d48;
}

.text-blue {
  color: #2563eb;
}

.text-gold {
  color: #b45309;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}

@media (max-width: 900px) {
  .course-layout {
    grid-template-columns: 1fr;
  }
  .modal-form-grid {
    grid-template-columns: 1fr;
  }
}

/* Dosen Assignment Submissions View */
.dosen-assignment-view {
  margin-top: 0.5rem;
}

.dosen-asg-summary {
  margin-bottom: 1rem;
}

.sub-count-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.45rem 0.85rem;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1d4ed8;
  font-weight: 700;
  font-size: 0.85rem;
  border-radius: 8px;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
  border-radius: var(--radius-sm);
  border: 1px solid #e2e8f0;
}

.dosen-submissions-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.88rem;
  text-align: left;
}

.dosen-submissions-table th {
  background: #f8fafc;
  color: #475569;
  font-weight: 700;
  padding: 0.75rem 1rem;
  border-bottom: 2px solid #cbd5e1;
  white-space: nowrap;
}

.dosen-submissions-table td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  vertical-align: middle;
}

.dosen-submissions-table tbody tr:hover {
  background: #f8fafc;
}

.nim-code {
  font-family: var(--font-code);
  font-size: 0.8rem;
  background: #f1f5f9;
  color: #0f172a;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  border: 1px solid #cbd5e1;
}

.student-name-td {
  color: #0f172a;
}

.time-td {
  font-size: 0.8rem;
  color: #64748b;
  white-space: nowrap;
}

.repo-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.35rem 0.7rem;
  background: #1d4ed8;
  color: #ffffff;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
  text-decoration: none;
  transition: background 0.2s ease;
}

.repo-link-btn:hover {
  background: #1e40af;
}

.link-icon-sm {
  width: 14px;
  height: 14px;
}

.notes-td {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.82rem;
  color: #64748b;
}

.empty-submissions-box {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: var(--radius-sm);
  color: #92400e;
}

.empty-sub-icon {
  width: 28px;
  height: 28px;
  flex-shrink: 0;
  color: #d97706;
}

.empty-submissions-box h5 {
  font-size: 0.95rem;
  font-weight: 800;
  margin-bottom: 0.25rem;
  color: #92400e;
}

.empty-submissions-box p {
  font-size: 0.85rem;
  color: #b45309;
  margin: 0;
}

.help-text {
  display: block;
  font-size: 0.78rem;
  color: #64748b;
  margin-top: 0.35rem;
}

.preset-group {
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.preset-label {
  display: block;
  font-size: 0.78rem;
  font-weight: 700;
  color: #475569;
  margin-bottom: 0.5rem;
}

.preset-buttons {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
}

.btn-preset {
  padding: 0.35rem 0.65rem;
  background: #ffffff;
  border: 1.5px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-preset:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #1d4ed8;
}

.btn-preset-danger {
  color: #be123c;
  border-color: #fecdd3;
}

.btn-preset-danger:hover {
  background: #fff1f2;
  border-color: #fda4af;
  color: #be123c;
}

.btn-blue {
  background: #1d4ed8;
  color: #ffffff;
  border: none;
}
.btn-blue:hover {
  background: #1e40af;
}

.btn-dashed {
  border: 1.5px dashed #cbd5e1;
  background: #f8fafc;
  color: #334155;
  font-size: 0.85rem;
}
.btn-dashed:hover {
  background: #eff6ff;
  border-color: #3b82f6;
  color: #1d4ed8;
}

.btn-uploading {
  opacity: 0.7;
  cursor: wait;
}

.resource-pending-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.85rem;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
  color: #64748b;
  font-weight: 500;
}
</style>
