<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { 
  ClipboardList, 
  CheckCircle2, 
  XCircle, 
  FileText, 
  Activity, 
  QrCode, 
  Search, 
  Save, 
  Users, 
  Info, 
  Sparkles, 
  UserCheck, 
  Calendar,
  X,
  Trash2,
  UserPlus,
  Smartphone
} from 'lucide-vue-next'

import Swal, { showSuccess, showError, showWarning, showConfirm } from '../utils/swal.js'

const currentUser = ref(null)
const activeTab = ref('dosen') // 'dosen' or 'mahasiswa'
const selectedCourse = ref('rpl-2026')
const selectedMeeting = ref(1)
const searchQuery = ref('')
const showQrModal = ref(false)
const courses = ref([])

const fetchCourses = async () => {
  try {
    const res = await fetch('/api/courses')
    if (res.ok) {
      courses.value = await res.json()
    }
  } catch (err) {
    console.warn('Presensi fetch courses error:', err)
  }
}

const visibleCourses = computed(() => {
  if (isDosen.value) return courses.value
  return courses.value.filter(c => c.status !== 'Non Aktif')
})

const checkUser = () => {
  const u = localStorage.getItem('user')
  if (u) {
    try { currentUser.value = JSON.parse(u) } catch { currentUser.value = null }
  }
}

const isDosen = computed(() => currentUser.value?.role === 'dosen')
const isMahasiswa = computed(() => currentUser.value?.role === 'mahasiswa')

const qrCodeUrl = computed(() => {
  const c = selectedCourse.value || 'rpl-2026'
  const m = selectedMeeting.value || 1
  const baseUrl = window.location.origin
  return `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(`${baseUrl}/presensi?checkin=${c}-p${m}`)}`
})

// Class Roster / Students List for Dosen View (Default empty status per meeting unless saved)
const classRoster = ref([
  { id: 'std-1', name: 'LINTANG ANGEL STEFANI', nim: '221112019', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-2', name: 'IGNATION SENSEKO MANGGUR', nim: '221112020', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-3', name: 'FAIZ IJLAL ARAYYAN', nim: '231112028', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-4', name: 'ALDINUS NDRURU', nim: '241112001', prodi: 'Sistem Informasi (S1)', status: '' },
  { id: 'std-5', name: 'OVAROLDUS SUPRATMAN', nim: '241112002', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-6', name: 'RADEN DHAFA ADHITYA SOSIAWAN', nim: '241112003', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-7', name: 'AHMAD FAUZI', nim: '20260801001', prodi: 'Teknik Informatika (S1)', status: '' },
  { id: 'std-8', name: 'SITI NURHALIZA', nim: '20260801002', prodi: 'Sistem Informasi (S1)', status: '' },
])

const selectedProdi = ref('ALL') // 'ALL', 'Teknik Informatika (S1)', 'Sistem Informasi (S1)'

// Student Check-In Form State
const studentNim = ref('')
const studentName = ref('')
const studentStatus = ref('Hadir')

// Registered Attendances Log from Backend
const attendances = ref([])

const loadMeetingAttendance = () => {
  const cId = selectedCourse.value
  const mNo = parseInt(selectedMeeting.value)

  const savedForMeeting = attendances.value.filter(
    a => a.course_id === cId && parseInt(a.meeting_no) === mNo
  )

  classRoster.value.forEach(student => {
    const rec = savedForMeeting.find(a => a.student_nim === student.nim)
    if (rec) {
      student.status = rec.status
    } else {
      student.status = ''
    }
  })
}

watch([selectedCourse, selectedMeeting], () => {
  loadMeetingAttendance()
})

const fetchStudents = async () => {
  const local = localStorage.getItem('elearning_students')
  if (local) {
    try {
      const parsed = JSON.parse(local)
      parsed.forEach(st => {
        const targetNim = st.nim || st.username
        const existing = classRoster.value.find(r => r.nim === targetNim)
        if (existing) {
          if (st.name) existing.name = st.name.toUpperCase()
          if (st.prodi) existing.prodi = st.prodi
        } else {
          classRoster.value.push({
            id: st.id || 'std-' + targetNim,
            name: (st.name || targetNim).toUpperCase(),
            nim: targetNim,
            prodi: st.prodi || 'Teknik Informatika (S1)',
            status: ''
          })
        }
      })
    } catch {}
  }

  try {
    const res = await fetch('/api/students')
    if (res.ok) {
      const data = await res.json()
      data.forEach(st => {
        const existing = classRoster.value.find(r => r.nim === st.username)
        if (existing) {
          if (st.name) existing.name = st.name.toUpperCase()
          if (st.prodi) existing.prodi = st.prodi
        } else {
          classRoster.value.push({
            id: st.id || 'std-' + st.username,
            name: (st.name || st.username).toUpperCase(),
            nim: st.username,
            prodi: st.prodi || 'Teknik Informatika (S1)',
            status: ''
          })
        }
      })
    }
  } catch (err) {
    console.warn('Fetch students error:', err)
  }
}

const fetchAttendance = async () => {
  try {
    const res = await fetch('/api/attendance')
    if (res.ok) {
      attendances.value = await res.json()
      localStorage.setItem('elearning_attendances', JSON.stringify(attendances.value))
    }
  } catch (err) {
    console.warn('Fetch attendance error:', err)
    const local = localStorage.getItem('elearning_attendances')
    if (local) {
      try { attendances.value = JSON.parse(local) } catch {}
    }
  }
  loadMeetingAttendance()
}

onMounted(() => {
  checkUser()
  fetchCourses()
  fetchStudents()
  fetchAttendance()
})

// Filtered Students List by Search Query and Program Studi
const filteredRoster = computed(() => {
  let list = classRoster.value

  if (selectedProdi.value !== 'ALL') {
    list = list.filter(s => s.prodi === selectedProdi.value)
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(s => 
      s.name.toLowerCase().includes(q) || 
      s.nim.toLowerCase().includes(q) ||
      (s.prodi && s.prodi.toLowerCase().includes(q))
    )
  }
  return list
})

// Summary Stats
const stats = computed(() => {
  const currentList = filteredRoster.value
  const total = currentList.length
  const hadir = currentList.filter(s => s.status === 'Hadir').length
  const absen = currentList.filter(s => s.status === 'Absen').length
  const izin = currentList.filter(s => s.status === 'Izin').length
  const sakit = currentList.filter(s => s.status === 'Sakit').length
  const belum = currentList.filter(s => !s.status).length
  const percentage = total > 0 ? ((hadir / total) * 100).toFixed(1) : 0

  const totalAll = classRoster.value.length
  const countTif = classRoster.value.filter(s => s.prodi === 'Teknik Informatika (S1)').length
  const countSi = classRoster.value.filter(s => s.prodi === 'Sistem Informasi (S1)').length

  return { total, hadir, absen, izin, sakit, belum, percentage, totalAll, countTif, countSi }
})

// Quick Batch Action: Set All Students Status
const setAllStatus = (targetStatus) => {
  filteredRoster.value.forEach(student => {
    student.status = targetStatus
  })
  showSuccess('Set Kehadiran Massal!', `Semua ${filteredRoster.value.length} mahasiswa berhasil diubah statusnya menjadi '${targetStatus}'.`)
}

// Single Student Status Change
const setStudentStatus = (student, newStatus) => {
  student.status = newStatus
}

// Reset / Cancel Changes with Confirmation
const cancelOrResetChanges = async () => {
  const confirmed = await showConfirm(
    'Batalkan Perubahan?',
    'Apakah Anda yakin ingin membatalkan perubahan rekap presensi pertemuan ini?',
    'Ya, Batalkan'
  )
  if (confirmed) {
    loadMeetingAttendance()
    showSuccess('Dibatalkan! 🔄', 'Perubahan presensi berhasil dibatalkan dan dikembalikan.')
  }
}

// Delete / Remove Student from Class Roster
const deleteStudent = async (student) => {
  const confirmed = await showConfirm(
    'Hapus Mahasiswa?',
    `Apakah Anda yakin ingin menghapus ${student.name} (${student.nim}) dari daftar presensi kelas?`,
    'Ya, Hapus Mahasiswa'
  )
  if (confirmed) {
    classRoster.value = classRoster.value.filter(s => s.id !== student.id)
    showSuccess('Mahasiswa Dihapus 🗑️', `${student.name} telah dihapus dari lembar presensi.`)
  }
}

// Add New Student Prompt Modal
const addNewStudentPrompt = async () => {
  const { value: formValues } = await Swal.fire({
    title: 'Tambah Mahasiswa Baru',
    html:
      '<div style="text-align: left; font-size: 0.88rem; color: #475569; margin-bottom: 0.5rem; font-weight: 700;">Masukkan Data Mahasiswa:</div>' +
      '<input id="swal-input-name" class="swal2-input" placeholder="Nama Lengkap Mahasiswa" style="margin: 0.5rem 0; width: 100%; box-sizing: border-box;">' +
      '<input id="swal-input-nim" class="swal2-input" placeholder="NIM Mahasiswa (cth: 221112099)" style="margin: 0.5rem 0; width: 100%; box-sizing: border-box;">' +
      '<select id="swal-input-prodi" class="swal2-select" style="margin: 0.5rem 0; width: 100%;">' +
        '<option value="Teknik Informatika (S1)">Teknik Informatika (S1)</option>' +
        '<option value="Sistem Informasi (S1)">Sistem Informasi (S1)</option>' +
      '</select>',
    focusConfirm: false,
    showCancelButton: true,
    confirmButtonText: 'Tambah ke Daftar',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#2563eb',
    cancelButtonColor: '#64748b',
    preConfirm: () => {
      const name = document.getElementById('swal-input-name').value
      const nim = document.getElementById('swal-input-nim').value
      const prodi = document.getElementById('swal-input-prodi').value
      if (!name || !nim) {
        Swal.showValidationMessage('Nama dan NIM wajib diisi!')
        return false
      }
      return { name, nim, prodi }
    }
  })

  if (formValues) {
    const newStudent = {
      id: 'std-' + Date.now(),
      name: formValues.name.toUpperCase(),
      nim: formValues.nim,
      prodi: formValues.prodi,
      status: 'Hadir'
    }
    classRoster.value.push(newStudent)
    showSuccess('Mahasiswa Ditambahkan 🎓', `${newStudent.name} (${newStudent.nim}) - ${newStudent.prodi} berhasil ditambahkan.`)
  }
}

// Simulate QR Code Scan Action
const simulateQrScan = () => {
  setAllStatus('Hadir')
  showSuccess('Scan QR Berhasil! 📱', `Presensi Pertemuan ${selectedMeeting.value} berhasil dicatat secara otomatis.`)
  showQrModal.value = false
}

// Batch Save Attendance to Backend
const saveBatchAttendance = async () => {
  try {
    const cId = selectedCourse.value
    const mNo = parseInt(selectedMeeting.value)

    const promises = classRoster.value.map(s => {
      if (!s.status) return Promise.resolve()
      return fetch('/api/attendance', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          course_id: cId,
          meeting_no: mNo,
          student_nim: s.nim,
          student_name: s.name,
          status: s.status
        })
      })
    })
    await Promise.all(promises)

    classRoster.value.forEach(s => {
      if (!s.status) return
      const existingIdx = attendances.value.findIndex(a => a.course_id === cId && parseInt(a.meeting_no) === mNo && a.student_nim === s.nim)
      if (existingIdx !== -1) {
        attendances.value[existingIdx].status = s.status
      } else {
        attendances.value.push({
          id: 'att-' + Date.now() + '-' + s.nim,
          course_id: cId,
          meeting_no: mNo,
          student_nim: s.nim,
          student_name: s.name,
          status: s.status,
          check_in_time: new Date().toISOString()
        })
      }
    })
    localStorage.setItem('elearning_attendances', JSON.stringify(attendances.value))

    const filledCount = classRoster.value.filter(s => s.status).length
    showSuccess('Rekap Kehadiran Disimpan! 💾', `Presensi kelas Pertemuan ${selectedMeeting.value} untuk ${filledCount} mahasiswa berhasil disimpan ke server.`)
    fetchAttendance()
  } catch (err) {
    showError('Gagal Menyimpan!', 'Terjadi kesalahan saat menyimpan rekap presensi.')
  }
}

// Student Self Check-In
const submitStudentCheckIn = async () => {
  if (!studentNim.value || !studentName.value) {
    showWarning('Perhatian', 'Harap isi NIM dan Nama Mahasiswa!')
    return
  }

  try {
    const res = await fetch('/api/attendance', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        course_id: selectedCourse.value,
        meeting_no: parseInt(selectedMeeting.value),
        student_nim: studentNim.value,
        student_name: studentName.value,
        status: studentStatus.value
      })
    })

    if (res.ok) {
      showSuccess('Check-In Berhasil! ✅', `Presensi Pertemuan ${selectedMeeting.value} atas nama ${studentName.value} (${studentNim.value}) berhasil dicatat.`)
      fetchAttendance()
      studentNim.value = ''
      studentName.value = ''
    }
  } catch (err) {
    showError('Gagal!', 'Gagal mengirim presensi!')
  }
}

// Mahasiswa Personal Presensi Helpers
const isMeetingActive = (mNo) => {
  const status = getStudentMeetingStatus(mNo)
  return ['Hadir', 'Izin', 'Sakit', 'Absen'].includes(status)
}

// Get the latest active meeting number from course modules
const currentActiveMeetingNo = computed(() => {
  const thisCourse = courses.value.find(c => c.id === selectedCourse.value)
  if (!thisCourse || !thisCourse.modules) return 1
  const activeModules = thisCourse.modules.filter(m => m.is_active || m.status === 'terbuka')
  if (activeModules.length === 0) return 1
  const maxActive = Math.max(...activeModules.map(m => m.meeting_number))
  return maxActive
})

// Generate meeting date for a given meeting number
const getMeetingDateLabel = (mNo) => {
  const startDate = new Date('2026-09-15')
  const mDate = new Date(startDate.getTime() + (mNo - 1) * 7 * 24 * 60 * 60 * 1000)
  return mDate.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const getMeetingTopic = (courseId, mNo) => {
  if (!isMeetingActive(mNo)) {
    return '—'
  }

  if (courseId === 'rpl-2026') {
    const topics = [
      'Pengantar Rekayasa Perangkat Lunak', 'Model Proses SDLC (Waterfall & Agile)',
      'Manajemen Proyek Perangkat Lunak', 'Analisis Kebutuhan & Dokumen SRS',
      'Pemodelan UML & Use Case Diagram', 'Diagram Activity & Sequence',
      'Diagram Class & Structural Modeling', 'Ujian Tengah Semester (UTS)',
      'Perancangan Arsitektur Microservices', 'UI/UX Design Patterns',
      'Prinsip OOP & SOLID', 'Software Testing & QA',
      'DevOps & CI/CD Pipeline', 'Presentasi Proyek Akhir (UAS)'
    ]
    return topics[mNo - 1] || `Materi Pertemuan ${mNo}`
  } else {
    const topics = [
      'Pengantar Web & Protokol HTTP', 'HTML5 Semantik & Structure',
      'CSS Flexbox & Responsive Layout', 'JavaScript ES6+ Fundamentals',
      'DOM Manipulation & Local Storage', 'Async JS, Promises & Fetch API',
      'Pengantar Vue.js 3 Reactive', 'Ujian Tengah Semester (UTS)',
      'Vue Components & State Management', 'Pengantar Backend Golang REST API',
      'Routing & Controller Backend', 'Database MySQL & ORM Integration',
      'Web Security & Auth JWT', 'Presentasi Web App Final (UAS)'
    ]
    return topics[mNo - 1] || `Materi Pertemuan ${mNo}`
  }
}

const getMeetingDate = (mNo) => {
  const startDate = new Date('2026-09-15')
  const mDate = new Date(startDate.getTime() + (mNo - 1) * 7 * 24 * 60 * 60 * 1000)
  return mDate.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const getStudentMeetingRecord = (mNo) => {
  if (!currentUser.value) return null
  return attendances.value.find(a => 
    a.course_id === selectedCourse.value && 
    a.meeting_no === mNo && 
    a.student_nim === currentUser.value.username
  )
}

const getStudentMeetingStatus = (mNo) => {
  const rec = getStudentMeetingRecord(mNo)
  if (rec && rec.status) return rec.status
  return 'Belum Mulai'
}

const getStudentMeetingTime = (mNo) => {
  const rec = getStudentMeetingRecord(mNo)
  if (rec && rec.check_in_time) {
    return new Date(rec.check_in_time).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WIB'
  }
  return '—'
}

const myHadirCount = computed(() => {
  let count = 0
  for (let m = 1; m <= 14; m++) {
    if (getStudentMeetingStatus(m) === 'Hadir') count++
  }
  return count
})

const myIzinSakitCount = computed(() => {
  let count = 0
  for (let m = 1; m <= 14; m++) {
    const st = getStudentMeetingStatus(m)
    if (st === 'Izin' || st === 'Sakit') count++
  }
  return count
})

const myAbsenCount = computed(() => {
  let count = 0
  for (let m = 1; m <= 14; m++) {
    if (getStudentMeetingStatus(m) === 'Absen') count++
  }
  return count
})

const myAttendancePercentage = computed(() => {
  const completedMeetings = myHadirCount.value + myIzinSakitCount.value + myAbsenCount.value
  if (completedMeetings === 0) return 100
  return Math.round((myHadirCount.value / completedMeetings) * 100)
})

const openSelfScanModal = async () => {
  const activeMeetingNo = currentActiveMeetingNo.value
  const { value: tokenCode } = await Swal.fire({
    title: '📱 Scan QR / Input Kode Token Presensi',
    html: `
      <div style="text-align: left; font-size: 0.88rem; color: #475569; margin-bottom: 0.4rem; font-weight: 600;">
        Pertemuan Aktif: <strong>Pertemuan ${activeMeetingNo}</strong> (${getMeetingDateLabel(activeMeetingNo)})
      </div>
      <div style="text-align: left; font-size: 0.88rem; color: #475569; margin-bottom: 0.75rem;">
        Masukkan kode token presensi yang ditampilkan Dosen di layar proyektor:
      </div>
      <input id="swal-token" class="swal2-input" placeholder="Contoh: 2026 / 8890" style="margin: 0.5rem 0; width: 100%; box-sizing: border-box; text-align: center; font-size: 1.5rem; letter-spacing: 4px; font-weight: 800;" maxlength="6">
    `,
    showCancelButton: true,
    confirmButtonText: '✅ Kirim Presensi Hadir',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#059669',
    preConfirm: () => {
      const code = document.getElementById('swal-token').value
      if (!code) {
        Swal.showValidationMessage('Kode Token Presensi wajib diisi!')
        return false
      }
      return code
    }
  })

  if (tokenCode) {
    try {
      const res = await fetch('/api/attendance', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          course_id: selectedCourse.value,
          meeting_no: activeMeetingNo,
          student_nim: currentUser.value.username,
          student_name: currentUser.value.name,
          status: 'Hadir'
        })
      })

      if (res.ok) {
        showSuccess('Presensi Berhasil Terverifikasi! 🎉', `Status kehadiran Anda atas nama ${currentUser.value.name} (${currentUser.value.username}) berhasil dicatat sebagai HADIR untuk Pertemuan ${activeMeetingNo}.`)
        fetchAttendance()
      }
    } catch (err) {
      showError('Gagal!', 'Gagal memproses presensi token.')
    }
  }
}

const openIzinModal = async () => {
  // Generate all 14 meeting options dynamically
  const meetingOptions = Array.from({ length: 14 }, (_, i) => i + 1)
    .map(m => `<option value="${m}">Pertemuan ${m} (${getMeetingDateLabel(m)})</option>`)
    .join('')

  const { value: formValues } = await Swal.fire({
    title: '✉️ Form Pengajuan Surat Izin / Sakit',
    html: `
      <div style="text-align: left; font-size: 0.88rem; color: #475569; margin-bottom: 0.5rem; font-weight: 600;">
        Pilih Pertemuan & Jenis Keterangan:
      </div>
      <select id="swal-meeting" class="swal2-select" style="margin: 0.5rem 0; width: 100%;">
        ${meetingOptions}
      </select>
      <select id="swal-status" class="swal2-select" style="margin: 0.5rem 0; width: 100%;">
        <option value="Izin">Izin (Keperluan / Acara)</option>
        <option value="Sakit">Sakit (Surat Dokter)</option>
      </select>
      <textarea id="swal-notes" class="swal2-textarea" placeholder="Alasan ketidakhadiran..." style="margin: 0.5rem 0; width: 100%; box-sizing: border-box;"></textarea>
    `,
    showCancelButton: true,
    confirmButtonText: 'Kirim Pengajuan Izin',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#d97706',
    preConfirm: () => {
      const mNo = document.getElementById('swal-meeting').value
      const st = document.getElementById('swal-status').value
      const notes = document.getElementById('swal-notes').value
      if (!notes) {
        Swal.showValidationMessage('Alasan ketidakhadiran wajib diisi!')
        return false
      }
      return { mNo, st, notes }
    }
  })

  if (formValues) {
    try {
      const res = await fetch('/api/attendance', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          course_id: selectedCourse.value,
          meeting_no: parseInt(formValues.mNo),
          student_nim: currentUser.value.username,
          student_name: currentUser.value.name,
          status: formValues.st
        })
      })

      if (res.ok) {
        showSuccess('Pengajuan Dikirim! 📨', `Surat Keterangan ${formValues.st} Pertemuan ${formValues.mNo} telah dikirim ke Dosen Pengampu.`)
        fetchAttendance()
      }
    } catch (err) {
      showError('Gagal!', 'Gagal memproses pengajuan izin.')
    }
  }
}
</script>

<template>
  <div class="presensi-container animate-fade-in">
    <!-- Header -->
    <div class="page-header">
      <div class="header-main">
        <div class="header-title-group">
          <ClipboardList class="header-icon text-gold" />
          <div>
            <h2>{{ isMahasiswa ? '📊 Presensi & Syarat Ujian Mahasiswa' : 'Presensi & Kehadiran Perkuliahan' }}</h2>
            <p class="subtitle">{{ isMahasiswa ? 'Pantau persentase kehadiran Anda per mata kuliah, lakukan scan QR presensi, dan ajukan surat izin.' : 'Kelola dan rekap kehadiran mahasiswa ITB Swadharma Dosen Rio Widyatmoko.' }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- MAHASISWA PERSONAL ATTENDANCE DASHBOARD -->
    <div v-if="isMahasiswa" class="mahasiswa-attendance-wrapper">
      <!-- Student Banner -->
      <div class="glass-card std-welcome-card">
        <div class="std-profile-info">
          <div class="std-avatar-box">
            <UserCheck class="std-avatar-icon text-blue" />
          </div>
          <div>
            <h3>Peserta: {{ currentUser.name }}</h3>
            <p class="std-meta">NIM: <strong>{{ currentUser.username }}</strong> | Prodi: <strong>{{ currentUser.prodi || 'Teknik Informatika' }}</strong></p>
          </div>
        </div>

        <div class="std-quick-actions">
          <button @click="openSelfScanModal" class="btn btn-gold btn-lg">
            <QrCode class="btn-icon-sm" />
            <span>📱 Scan QR / Input Token Presensi</span>
          </button>
          <button @click="openIzinModal" class="btn btn-secondary btn-lg">
            <FileText class="btn-icon-sm" />
            <span>✉️ Ajukan Surat Izin / Sakit</span>
          </button>
        </div>
      </div>

      <!-- Course Select & Stats Overview -->
      <div class="course-selector-bar glass-card">
        <div class="select-course-box">
          <label class="filter-label">Pilih Mata Kuliah Perkuliahan:</label>
          <select v-model="selectedCourse" class="glass-input select-lg">
            <template v-if="visibleCourses.length > 0">
              <option v-for="c in visibleCourses" :key="c.id" :value="c.id">
                {{ c.name }} ({{ c.code }})
              </option>
            </template>
            <template v-else>
              <option value="rpl-2026">Rekayasa Perangkat Lunak (TIF-301)</option>
              <option value="webdev-2026">Pemrograman Web (TIF-302)</option>
            </template>
          </select>
        </div>

        <!-- Personal Attendance Stats Cards -->
        <div class="std-stats-row">
          <div class="std-stat-card border-blue">
            <span class="stat-title">Persentase Kehadiran</span>
            <span class="stat-value text-blue">{{ myAttendancePercentage }}%</span>
            <span :class="['badge', myAttendancePercentage >= 75 ? 'badge-emerald' : 'badge-rose']">
              {{ myAttendancePercentage >= 75 ? '✅ LULUS SYARAT UAS' : '⚠️ KURANG DARI 75%' }}
            </span>
          </div>

          <div class="std-stat-card border-emerald">
            <span class="stat-title">Hadir</span>
            <span class="stat-value text-emerald">{{ myHadirCount }} Pertemuan</span>
            <span class="stat-sub">Sesuai Jadwal</span>
          </div>

          <div class="std-stat-card border-gold">
            <span class="stat-title">Izin / Sakit</span>
            <span class="stat-value text-gold">{{ myIzinSakitCount }} Pertemuan</span>
            <span class="stat-sub">Surat Terverifikasi</span>
          </div>

          <div class="std-stat-card border-rose">
            <span class="stat-title">Tanpa Keterangan</span>
            <span class="stat-value text-rose">{{ myAbsenCount }} Pertemuan</span>
            <span class="stat-sub">Batas Maksimal: 3x</span>
          </div>
        </div>
      </div>

      <!-- Table 14 Pertemuan -->
      <div class="glass-card meeting-history-card">
        <div class="history-card-header">
          <div>
            <h3>📅 Riwayat Presensi Pertemuan (1 - 14)</h3>
            <p class="card-sub">Detail status kehadiran Anda untuk mata kuliah {{ selectedCourse === 'rpl-2026' ? 'Rekayasa Perangkat Lunak' : 'Pemrograman Web' }}.</p>
          </div>
        </div>

        <div class="meeting-table-container">
          <table class="meeting-table">
            <thead>
              <tr>
                <th>Pertemuan</th>
                <th>Tanggal & Topik Materi</th>
                <th>Status Kehadiran</th>
                <th>Metode Presensi</th>
                <th>Waktu Check-In</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in 14" :key="m" class="meeting-tr">
                <td>
                  <span class="meeting-pill">Pertemuan {{ m }}</span>
                </td>
                <td>
                  <div class="topic-info">
                    <span class="topic-title">
                      <template v-if="isMeetingActive(m)">
                        Pertemuan {{ m }}: {{ getMeetingTopic(selectedCourse, m) }}
                      </template>
                      <template v-else>
                        Pertemuan {{ m }}: <span style="color: #94a3b8; font-weight: 500; font-style: italic;">— (Belum Dimulai)</span>
                      </template>
                    </span>
                    <span class="topic-date">📅 {{ getMeetingDate(m) }}</span>
                  </div>
                </td>
                <td>
                  <span :class="['badge', getStudentMeetingStatus(m) === 'Hadir' ? 'badge-emerald' : getStudentMeetingStatus(m) === 'Izin' || getStudentMeetingStatus(m) === 'Sakit' ? 'badge-gold' : getStudentMeetingStatus(m) === 'Absen' ? 'badge-rose' : 'badge-blue']">
                    {{ getStudentMeetingStatus(m) }}
                  </span>
                </td>
                <td>
                  <span class="method-text">
                    {{ getStudentMeetingStatus(m) === 'Hadir' ? '📱 Scan QR / Check-In' : getStudentMeetingStatus(m) === 'Izin' || getStudentMeetingStatus(m) === 'Sakit' ? '✉️ Surat Keterangan' : '—' }}
                  </span>
                </td>
                <td>
                  <span class="checkin-time">
                    {{ getStudentMeetingTime(m) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- DOSEN CLASS ATTENDANCE SHEET -->
    <div v-else class="dosen-attendance-sheet">
      
      <!-- Top Notice Banner Card (Matching Reference Image Layout with Modernized UX) -->
      <div class="notice-card-header">
        <div class="notice-card-content">
          <div class="info-icon-badge">
            <Info class="info-icon" />
          </div>
          <span class="notice-text">
            Pengajar melakukan presensi kehadiran peserta secara manual pada pertemuan ini
          </span>
        </div>
        <div class="notice-card-actions">
          <button @click="showQrModal = true" class="btn-qr-outline">
            <QrCode class="btn-icon-sm text-blue" />
            <span>Generate QR Code</span>
          </button>
          <span class="action-divider">|</span>
          <button @click="cancelOrResetChanges" class="btn-cancel">
            Batal
          </button>
          <button @click="saveBatchAttendance" class="btn-simpan">
            Simpan kehadiran
          </button>
        </div>
      </div>

      <!-- Course Filter & Stats Bar -->
      <div class="filter-card glass-card">
        <div class="filter-group">
          <div>
            <label class="filter-label">Mata Kuliah</label>
            <select v-model="selectedCourse" class="glass-input select-lg">
              <template v-if="courses.length > 0">
                <option v-for="c in courses" :key="c.id" :value="c.id">
                  {{ c.name }} ({{ c.code }}) {{ c.status === 'Non Aktif' ? '[NON-AKTIF]' : '' }}
                </option>
              </template>
              <template v-else>
                <option value="rpl-2026">Rekayasa Perangkat Lunak (TIF-301)</option>
                <option value="webdev-2026">Pemrograman Web (TIF-302)</option>
              </template>
            </select>
          </div>

          <div>
            <label class="filter-label">Pertemuan Ke-</label>
            <select v-model="selectedMeeting" class="glass-input select-lg">
              <option v-for="n in 14" :key="n" :value="n">Pertemuan {{ n }}</option>
            </select>
          </div>

          <div>
            <label class="filter-label">Program Studi</label>
            <select v-model="selectedProdi" class="glass-input select-lg">
              <option value="ALL">Semua Program Studi ({{ stats.totalAll }})</option>
              <option value="Teknik Informatika (S1)">💻 Teknik Informatika (S1) ({{ stats.countTif }})</option>
              <option value="Sistem Informasi (S1)">📊 Sistem Informasi (S1) ({{ stats.countSi }})</option>
            </select>
          </div>
        </div>

        <!-- Search & Batch Tools -->
        <div class="filter-tools">
          <div class="search-input-box">
            <Search class="search-icon" />
            <input v-model="searchQuery" class="search-input" placeholder="Cari Nama Mahasiswa / NIM / Prodi..." />
          </div>
          <button @click="addNewStudentPrompt" class="btn-add-student">
            <UserPlus class="btn-icon-xs" />
            <span>Tambah Mahasiswa</span>
          </button>
        </div>

        <!-- Summary Stat Badges -->
        <div class="summary-stats">
          <div class="stat-pill stat-total">
            <Users class="pill-icon" />
            <span>Total: <strong>{{ stats.total }}</strong></span>
          </div>
          <div class="stat-pill stat-hadir">
            <CheckCircle2 class="pill-icon text-emerald" />
            <span>Hadir: <strong>{{ stats.hadir }}</strong></span>
          </div>
          <div class="stat-pill stat-absen">
            <XCircle class="pill-icon text-rose" />
            <span>Absen: <strong>{{ stats.absen }}</strong></span>
          </div>
          <div class="stat-pill stat-izin">
            <FileText class="pill-icon text-sky" />
            <span>Izin: <strong>{{ stats.izin }}</strong></span>
          </div>
          <div class="stat-pill stat-sakit">
            <Activity class="pill-icon text-amber" />
            <span>Sakit: <strong>{{ stats.sakit }}</strong></span>
          </div>
          <div v-if="stats.belum > 0" class="stat-pill stat-belum">
            <Info class="pill-icon text-slate" />
            <span>Belum Diisi: <strong>{{ stats.belum }}</strong></span>
          </div>
        </div>
      </div>

      <!-- Attendance Table (Designed based on reference image + enhanced UX) -->
      <div class="table-card glass-card">
        <div class="table-responsive">
          <table class="attendance-table">
            <thead>
              <tr>
                <th class="col-no">No</th>
                <th class="col-student">Mahasiswa & Program Studi</th>
                <th class="col-status-badge">Kehadiran</th>
                <th class="col-actions">
                  <div class="status-header-wrapper">
                    <span class="header-title">Status</span>
                    <div class="batch-radio-selectors">
                      <button @click="setAllStatus('Hadir')" class="batch-pill h-pill" title="Set semua Mahasiswa Hadir">
                        <span class="radio-ring"><span class="radio-dot"></span></span> H
                      </button>
                      <button @click="setAllStatus('Absen')" class="batch-pill a-pill" title="Set semua Mahasiswa Absen">
                        <span class="radio-ring"><span class="radio-dot"></span></span> A
                      </button>
                      <button @click="setAllStatus('Izin')" class="batch-pill i-pill" title="Set semua Mahasiswa Izin">
                        <span class="radio-ring"><span class="radio-dot"></span></span> I
                      </button>
                      <button @click="setAllStatus('Sakit')" class="batch-pill s-pill" title="Set semua Mahasiswa Sakit">
                        <span class="radio-ring"><span class="radio-dot"></span></span> S
                      </button>
                    </div>
                  </div>
                </th>
                <th class="col-delete-header">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(student, idx) in filteredRoster" 
                :key="student.id"
                :class="{ 'row-absen': student.status === 'Absen', 'row-hadir': student.status === 'Hadir' }"
              >
                <td class="col-no">{{ idx + 1 }}</td>
                <td class="col-student">
                  <div class="student-name-group">
                    <div class="name-badge-row">
                      <h4 class="student-name-text">{{ student.name }}</h4>
                      <span 
                        class="prodi-tag-pill"
                        :class="student.prodi === 'Sistem Informasi (S1)' ? 'tag-si' : 'tag-tif'"
                      >
                        {{ student.prodi === 'Sistem Informasi (S1)' ? '📊 Sistem Informasi' : '💻 Teknik Informatika' }}
                      </span>
                    </div>
                    <span class="student-nim-text">NIM: {{ student.nim }}</span>
                  </div>
                </td>
                <td class="col-status-badge">
                  <span 
                    class="status-indicator-badge"
                    :class="{
                      'badge-hadir': student.status === 'Hadir',
                      'badge-absen': student.status === 'Absen',
                      'badge-izin': student.status === 'Izin',
                      'badge-sakit': student.status === 'Sakit',
                      'badge-empty': !student.status
                    }"
                  >
                    <span>{{ student.status || 'Belum Diisi' }}</span>
                  </span>
                </td>
                <td class="col-actions">
                  <!-- Enhanced Radio Button Pills (H / A / I / S) -->
                  <div class="row-status-selectors">
                    <button 
                      type="button"
                      @click="setStudentStatus(student, 'Hadir')"
                      class="status-radio-pill pill-h"
                      :class="{ 'selected': student.status === 'Hadir' }"
                      title="Set Hadir"
                    >
                      <span class="radio-ring"><span class="radio-dot"></span></span>
                      <span class="lbl">H</span>
                    </button>
                    <button 
                      type="button"
                      @click="setStudentStatus(student, 'Absen')"
                      class="status-radio-pill pill-a"
                      :class="{ 'selected': student.status === 'Absen' }"
                      title="Set Absen"
                    >
                      <span class="radio-ring"><span class="radio-dot"></span></span>
                      <span class="lbl">A</span>
                    </button>
                    <button 
                      type="button"
                      @click="setStudentStatus(student, 'Izin')"
                      class="status-radio-pill pill-i"
                      :class="{ 'selected': student.status === 'Izin' }"
                      title="Set Izin"
                    >
                      <span class="radio-ring"><span class="radio-dot"></span></span>
                      <span class="lbl">I</span>
                    </button>
                    <button 
                      type="button"
                      @click="setStudentStatus(student, 'Sakit')"
                      class="status-radio-pill pill-s"
                      :class="{ 'selected': student.status === 'Sakit' }"
                      title="Set Sakit"
                    >
                      <span class="radio-ring"><span class="radio-dot"></span></span>
                      <span class="lbl">S</span>
                    </button>
                  </div>
                </td>
                <td class="col-delete-cell">
                  <button 
                    type="button"
                    @click="deleteStudent(student)"
                    class="btn-delete-student"
                    title="Hapus Mahasiswa ini dari daftar presensi"
                  >
                    <Trash2 class="icon-trash" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>



    <!-- QR CODE POPUP MODAL -->
    <Teleport to="body">
      <div v-if="showQrModal" class="modal-backdrop" @click.self="showQrModal = false">
        <div class="modal-box glass-card">
          <div class="modal-header">
            <h3>📲 QR Code Presensi Pertemuan {{ selectedMeeting }}</h3>
            <button @click="showQrModal = false" class="btn-close-modal">
              <X class="icon-sm" />
            </button>
          </div>
          <div class="modal-body">
            <div class="qr-display-box">
              <!-- High-Quality Scannable 2D QR Code -->
              <div class="qr-image-frame">
                <img 
                  :src="qrCodeUrl" 
                  alt="QR Code Presensi Pertemuan" 
                  class="real-qr-code"
                />
              </div>
            </div>
            <p class="qr-instruction">
              Scan QR Code ini menggunakan Kamera Smartphone Mahasiswa (atau WhatsApp/Google Lens) untuk Check-In Presensi otomatis Pertemuan {{ selectedMeeting }}.
            </p>
            <div class="qr-timer-badge">
              <Sparkles class="icon-xs text-gold" />
              <span>Berlaku 15 Menit — Status: Aktif & Siap Scan</span>
            </div>
          </div>
          <div class="modal-footer">
            <div class="modal-footer-grid">
              <button @click="simulateQrScan" class="btn-qr-simulate">
                <Smartphone class="btn-icon-xs" />
                <span>Simulasi Scan</span>
              </button>
              <button @click="showQrModal = false" class="btn-qr-close">
                <span>Tutup QR Code</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.presensi-container {
  max-width: 1250px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  margin-bottom: 1.5rem;
}

.header-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.header-title-group {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.header-icon {
  width: 36px;
  height: 36px;
}

.header-title-group h2 {
  font-family: var(--font-heading);
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 0.15rem 0;
}

.subtitle {
  color: #64748b;
  font-size: 0.88rem;
  margin: 0;
}

.mode-tabs {
  display: flex;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  padding: 0.25rem;
  border-radius: var(--radius-sm);
  gap: 0.25rem;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-xs);
  border: none;
  background: transparent;
  color: #475569;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn.active {
  background: #ffffff;
  color: #1d4ed8;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
}

.tab-icon {
  width: 16px;
  height: 16px;
}

/* Filter Card */
.filter-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  margin-bottom: 1rem;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  gap: 1.25rem;
}

.filter-label {
  display: block;
  font-size: 0.75rem;
  font-weight: 800;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.3rem;
}

.select-lg {
  min-width: 220px;
  font-weight: 700;
}

.summary-stats {
  display: flex;
  gap: 0.65rem;
  flex-wrap: wrap;
}

.stat-pill {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.4rem 0.75rem;
  border-radius: 99px;
  font-size: 0.82rem;
  border: 1px solid #cbd5e1;
  background: #f8fafc;
}

.stat-total { background: #f1f5f9; color: #0f172a; border-color: #cbd5e1; }
.stat-hadir { background: #ecfdf5; color: #065f46; border-color: #a7f3d0; }
.stat-absen { background: #fef2f2; color: #991b1b; border-color: #fecaca; }
.stat-izin { background: #f0f9ff; color: #075985; border-color: #bae6fd; }
.stat-sakit { background: #fffbeb; color: #92400e; border-color: #fde68a; }

.pill-icon {
  width: 15px;
  height: 15px;
}

/* Top Notice Banner Card (Matching User Image Layout) */
.notice-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-left: 5px solid #2563eb;
  padding: 1rem 1.35rem;
  border-radius: var(--radius-md);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.02);
  margin-bottom: 1.25rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.notice-card-content {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  flex: 1;
  min-width: 280px;
}

.info-icon-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #2563eb;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.info-icon {
  width: 18px;
  height: 18px;
}

.notice-text {
  font-size: 0.92rem;
  font-weight: 600;
  color: #2563eb;
  line-height: 1.4;
}

.notice-card-actions {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  flex-shrink: 0;
}

.btn-qr-outline {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.5rem 0.95rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #2563eb;
  background: #ffffff;
  color: #2563eb;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-qr-outline:hover {
  background: #eff6ff;
}

.action-divider {
  color: #cbd5e1;
  font-weight: 300;
  font-size: 1.25rem;
}

.btn-cancel {
  padding: 0.5rem 1.1rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.btn-simpan {
  padding: 0.55rem 1.35rem;
  border-radius: var(--radius-sm);
  border: none;
  background: #2563eb;
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);
  transition: all 0.2s ease;
}

.btn-simpan:hover {
  background: #1d4ed8;
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.35);
  transform: translateY(-1px);
}

.filter-tools {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.btn-add-student {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 0.85rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #2563eb;
  background: #eff6ff;
  color: #2563eb;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-student:hover {
  background: #2563eb;
  color: #ffffff;
}

.search-input-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  width: 16px;
  height: 16px;
  color: #94a3b8;
}

.search-input {
  padding: 0.5rem 0.85rem 0.5rem 2.25rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #cbd5e1;
  font-size: 0.88rem;
  outline: none;
  min-width: 220px;
  transition: all 0.2s ease;
}

.search-input:focus {
  border-color: #1d4ed8;
  box-shadow: 0 0 0 3px rgba(29, 78, 216, 0.1);
}

.col-delete-header {
  width: 60px;
  text-align: center;
}

.col-delete-cell {
  text-align: center;
}

.btn-delete-student {
  background: transparent;
  border: none;
  color: #94a3b8;
  padding: 0.4rem;
  border-radius: 6px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.btn-delete-student:hover {
  background: #fef2f2;
  color: #dc2626;
  transform: scale(1.1);
}

.icon-trash {
  width: 17px;
  height: 17px;
}

/* Table Design */
.table-card {
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.03);
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.attendance-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.attendance-table th {
  background: #ffffff;
  padding: 1rem 1.25rem;
  font-size: 0.85rem;
  font-weight: 800;
  color: #334155;
  border-bottom: 2px solid #e2e8f0;
}

.attendance-table td {
  padding: 0.85rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.attendance-table tr:hover {
  background: #f8fafc;
}

.attendance-table tr.row-absen {
  background: #ffffff;
}

.attendance-table tr.row-hadir {
  background: #ffffff;
}

.col-no {
  width: 50px;
  font-weight: 700;
  color: #64748b;
  text-align: center;
}

.col-student {
  min-width: 260px;
}

.student-name-group {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.name-badge-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.2rem;
}

.prodi-tag-pill {
  display: inline-flex;
  align-items: center;
  padding: 0.15rem 0.55rem;
  border-radius: 99px;
  font-size: 0.72rem;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.tag-tif {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.tag-si {
  background: #f5f3ff;
  color: #6d28d9;
  border: 1px solid #ddd6fe;
}

.student-name-text {
  font-size: 0.92rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  letter-spacing: 0.02em;
}

.student-nim-text {
  font-size: 0.82rem;
  color: #64748b;
  font-weight: 600;
}

.col-status-badge {
  width: 130px;
}

.status-indicator-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  border-radius: 99px;
  font-size: 0.8rem;
  font-weight: 700;
}

.stat-sakit { background: #fffbeb; color: #92400e; border-color: #fde68a; }
.stat-belum { background: #f1f5f9; color: #475569; border-color: #cbd5e1; }

.badge-hadir { background: #dcfce7; color: #15803d; }
.badge-absen { background: #fee2e2; color: #dc2626; }
.badge-izin { background: #e0f2fe; color: #0369a1; }
.badge-sakit { background: #fef3c7; color: #b45309; }
.badge-empty { background: #f1f5f9; color: #64748b; border: 1px dashed #cbd5e1; }

.col-actions {
  min-width: 260px;
}

/* Header Status & Batch Controls */
.status-header-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title {
  font-weight: 800;
}

.batch-radio-selectors {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.batch-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 800;
  color: #475569;
  transition: all 0.2s ease;
}

.batch-pill:hover {
  color: #2563eb;
  transform: scale(1.08);
}

/* Row Status Radio Pills (H / A / I / S) */
.row-status-selectors {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.status-radio-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.3rem 0.5rem;
  border-radius: 99px;
  border: 1.5px solid transparent;
  background: transparent;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 700;
  color: #475569;
  transition: all 0.2s ease;
}

.radio-ring {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid #cbd5e1;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.radio-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: transparent;
  transition: all 0.2s ease;
}

.status-radio-pill:hover .radio-ring {
  border-color: #2563eb;
}

/* Selected Pill States */
.status-radio-pill.selected.pill-h {
  color: #15803d;
}
.status-radio-pill.selected.pill-h .radio-ring {
  border-color: #2563eb;
  background: #ffffff;
}
.status-radio-pill.selected.pill-h .radio-dot {
  background: #2563eb;
}

.status-radio-pill.selected.pill-a {
  color: #dc2626;
}
.status-radio-pill.selected.pill-a .radio-ring {
  border-color: #2563eb;
  background: #ffffff;
}
.status-radio-pill.selected.pill-a .radio-dot {
  background: #2563eb;
}

.status-radio-pill.selected.pill-i {
  color: #0369a1;
}
.status-radio-pill.selected.pill-i .radio-ring {
  border-color: #2563eb;
  background: #ffffff;
}
.status-radio-pill.selected.pill-i .radio-dot {
  background: #2563eb;
}

.status-radio-pill.selected.pill-s {
  color: #b45309;
}
.status-radio-pill.selected.pill-s .radio-ring {
  border-color: #2563eb;
  background: #ffffff;
}
.status-radio-pill.selected.pill-s .radio-dot {
  background: #2563eb;
}

/* Student Mode Grid */
.student-checkin-view .presensi-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-card, .log-card {
  padding: 2rem;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
}

.form-card h3, .log-card h3 {
  font-size: 1.3rem;
  font-weight: 800;
  color: #0f172a;
}

.card-sub {
  color: #64748b;
  font-size: 0.88rem;
  margin-bottom: 1.5rem;
}

.presensi-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  color: #334155;
  font-weight: 700;
  margin-bottom: 0.35rem;
}

.status-options {
  display: flex;
  gap: 1.5rem;
}

.status-radio {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.9rem;
  color: #0f172a;
  font-weight: 600;
  cursor: pointer;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.log-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.85rem 1rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-sm);
}

.log-info {
  display: flex;
  flex-direction: column;
}

.log-name {
  font-size: 0.9rem;
  font-weight: 800;
  color: #0f172a;
}

.log-detail {
  font-size: 0.8rem;
  color: #64748b;
}

/* Modal QR Code */
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 1100;
  background: rgba(15, 23, 42, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.modal-box {
  width: 100%;
  max-width: 420px;
  background: #ffffff;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
}

.modal-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.btn-close-modal {
  background: #f1f5f9;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qr-display-box {
  display: flex;
  justify-content: center;
  margin-bottom: 1.25rem;
}

.qr-image-frame {
  padding: 0.85rem;
  background: #ffffff;
  border: 2px solid #cbd5e1;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.real-qr-code {
  width: 200px;
  height: 200px;
  display: block;
}

.modal-footer-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
  width: 100%;
}

.btn-qr-simulate {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  height: 44px;
  padding: 0 0.85rem;
  border-radius: var(--radius-sm);
  border: none;
  background: #16a34a;
  color: #ffffff;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(22, 163, 74, 0.2);
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.btn-qr-simulate:hover {
  background: #15803d;
  box-shadow: 0 6px 16px rgba(22, 163, 74, 0.3);
  transform: translateY(-1px);
}

.btn-qr-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 44px;
  padding: 0 0.85rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
  box-sizing: border-box;
}

.btn-qr-close:hover {
  background: #f8fafc;
  border-color: #94a3b8;
  color: #0f172a;
}

.modal-footer {
  margin-top: 1.25rem;
}

.qr-instruction {
  font-size: 0.88rem;
  color: #475569;
  text-align: center;
  margin: 0 0 1.25rem 0;
  line-height: 1.5;
}

.qr-timer-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  color: #b45309;
  padding: 0.65rem 1rem;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
  font-weight: 700;
  margin-bottom: 1.25rem;
}

.mahasiswa-attendance-wrapper {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.std-welcome-card {
  padding: 1.5rem 1.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-left: 5px solid #2563eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1.5rem;
}

.std-profile-info {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}

.std-avatar-box {
  width: 52px;
  height: 52px;
  background: #eff6ff;
  border: 1.5px solid #bfdbfe;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.std-avatar-icon {
  width: 26px;
  height: 26px;
}

.std-profile-info h3 {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.std-meta {
  font-size: 0.88rem;
  color: #475569;
  margin: 0.2rem 0 0 0;
}

.std-quick-actions {
  display: flex;
  gap: 1rem;
}

.course-selector-bar {
  padding: 1.5rem 1.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.select-course-box {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f1f5f9;
}

.std-stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.25rem;
}

.std-stat-card {
  padding: 1.25rem;
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.border-blue { border-left: 4px solid #2563eb; }
.border-emerald { border-left: 4px solid #059669; }
.border-gold { border-left: 4px solid #d97706; }
.border-rose { border-left: 4px solid #e11d48; }

.stat-title {
  font-size: 0.82rem;
  color: #64748b;
  font-weight: 700;
}

.stat-value {
  font-size: 1.6rem;
  font-weight: 800;
}

.stat-sub {
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.meeting-history-card {
  padding: 1.75rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.history-card-header h3 {
  font-size: 1.2rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.meeting-table-container {
  overflow-x: auto;
  margin-top: 1.25rem;
}

.meeting-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.meeting-table th {
  padding: 0.85rem 1rem;
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
  font-size: 0.82rem;
  font-weight: 800;
  color: #334155;
  text-transform: uppercase;
}

.meeting-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.88rem;
  vertical-align: middle;
}

.meeting-tr:hover {
  background: #f8fafc;
}

.meeting-pill {
  font-weight: 800;
  font-size: 0.85rem;
  color: #2563eb;
  background: #eff6ff;
  padding: 0.3rem 0.65rem;
  border-radius: var(--radius-sm);
  border: 1px solid #bfdbfe;
}

.topic-info {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.topic-title {
  font-weight: 700;
  color: #0f172a;
}

.topic-date {
  font-size: 0.78rem;
  color: #64748b;
}

.method-text {
  font-size: 0.82rem;
  color: #475569;
  font-weight: 600;
}

.checkin-time {
  font-size: 0.82rem;
  color: #0f172a;
  font-weight: 700;
}

@media (max-width: 900px) {
  .std-welcome-card {
    flex-direction: column;
    align-items: stretch;
  }
  .std-quick-actions {
    flex-direction: column;
  }
  .filter-card, .toolbar-card, .header-main {
    flex-direction: column;
    align-items: stretch;
  }
  .summary-stats {
    justify-content: space-between;
  }
  .student-checkin-view .presensi-grid {
    grid-template-columns: 1fr;
  }
}
</style>
