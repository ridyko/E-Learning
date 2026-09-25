<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  Users, 
  UserPlus, 
  Search, 
  Edit3, 
  Trash2, 
  GraduationCap, 
  Mail, 
  Phone, 
  CheckCircle2, 
  XCircle, 
  BookOpen,
  Filter,
  MessageSquare
} from 'lucide-vue-next'

import Swal, { showSuccess, showError, showConfirm } from '../utils/swal.js'

// Student Data List
const students = ref([
  {
    id: 'std-1',
    nim: '221112019',
    name: 'LINTANG ANGEL STEFANI',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'lintang.angel@swadharma.ac.id',
    phone: '6281234567890',
    ipk: '3.75'
  },
  {
    id: 'std-2',
    nim: '221112020',
    name: 'IGNATION SENSEKO MANGGUR',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'ignation.senseko@swadharma.ac.id',
    phone: '6281298765432',
    ipk: '3.60'
  },
  {
    id: 'std-3',
    nim: '231112028',
    name: 'FAIZ IJLAL ARAYYAN',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'faiz.ijlal@swadharma.ac.id',
    phone: '6281311223344',
    ipk: '3.85'
  },
  {
    id: 'std-4',
    nim: '241112001',
    name: 'ALDINUS NDRURU',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'aldinus.ndruru@swadharma.ac.id',
    phone: '6281344556677',
    ipk: '3.50'
  },
  {
    id: 'std-5',
    nim: '241112002',
    name: 'OVAROLDUS SUPRATMAN',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'ovaroldus.s@swadharma.ac.id',
    phone: '6281377889900',
    ipk: '3.65'
  },
  {
    id: 'std-6',
    nim: '241112003',
    name: 'RADEN DHAFA ADHITYA SOSIAWAN',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'raden.dhafa@swadharma.ac.id',
    phone: '6281388990011',
    ipk: '3.90'
  },
  {
    id: 'std-7',
    nim: '20260801001',
    name: 'AHMAD FAUZI',
    prodi: 'Teknik Informatika (S1)',
    status: 'Aktif',
    email: 'ahmad.fauzi@swadharma.ac.id',
    phone: '6281512345678',
    ipk: '3.70'
  },
  {
    id: 'std-8',
    nim: '20260801002',
    name: 'SITI NURHALIZA',
    prodi: 'Sistem Informasi (S1)',
    status: 'Aktif',
    email: 'siti.nurhaliza@swadharma.ac.id',
    phone: '6281698765432',
    ipk: '3.80'
  }
])

const searchQuery = ref('')
const selectedProdi = ref('all')
const selectedStatus = ref('all')

// Filtered Students List
const filteredStudents = computed(() => {
  return students.value.filter(s => {
    const matchesSearch = !searchQuery.value || 
      s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
      s.nim.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.email.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesProdi = selectedProdi.value === 'all' || s.prodi.includes(selectedProdi.value)
    const matchesStatus = selectedStatus.value === 'all' || s.status === selectedStatus.value

    return matchesSearch && matchesProdi && matchesStatus
  })
})

// Summary Stats
const stats = computed(() => {
  const total = students.value.length
  const aktif = students.value.filter(s => s.status === 'Aktif').length
  const tif = students.value.filter(s => s.prodi.includes('Informatika')).length
  const si = students.value.filter(s => s.prodi.includes('Sistem')).length
  return { total, aktif, tif, si }
})

// Add New Student Form Popup (SweetAlert2)
const openAddStudentModal = async () => {
  const { value: formValues } = await Swal.fire({
    title: '🎓 Tambah Data Mahasiswa Baru',
    html: `
      <div style="text-align: left; display: flex; flex-direction: column; gap: 0.65rem; padding: 0.25rem 0.5rem;">
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">NIM Mahasiswa</label>
          <input id="swal-nim" class="swal2-input" placeholder="Contoh: 241112099" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.88rem;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Nama Lengkap Mahasiswa</label>
          <input id="swal-name" class="swal2-input" placeholder="Masukkan Nama Lengkap" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.88rem;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Program Studi</label>
          <select id="swal-prodi" class="swal2-input" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem; padding: 0.5rem;">
            <option value="Teknik Informatika (S1)">Teknik Informatika (S1)</option>
            <option value="Sistem Informasi (S1)">Sistem Informasi (S1)</option>
            <option value="Teknologi Informasi (D3)">Teknologi Informasi (D3)</option>
          </select>
        </div>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem;">
          <div>
            <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Email Kampus</label>
            <input id="swal-email" class="swal2-input" placeholder="email@swadharma.ac.id" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem;">
          </div>
          <div>
            <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">No. WhatsApp</label>
            <input id="swal-phone" class="swal2-input" placeholder="6281234567890" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem;">
          </div>
        </div>
      </div>
    `,
    focusConfirm: false,
    showCancelButton: true,
    confirmButtonText: 'Simpan Mahasiswa',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#2563eb',
    cancelButtonColor: '#64748b',
    preConfirm: () => {
      const nim = document.getElementById('swal-nim').value.trim()
      const name = document.getElementById('swal-name').value.trim()
      const prodi = document.getElementById('swal-prodi').value
      const email = document.getElementById('swal-email').value.trim()
      const phone = document.getElementById('swal-phone').value.trim()

      if (!nim || !name) {
        Swal.showValidationMessage('NIM dan Nama Mahasiswa wajib diisi!')
        return false
      }
      return { nim, name: name.toUpperCase(), prodi, email, phone }
    }
  })

  if (formValues) {
    const newStudent = {
      id: 'std-' + Date.now(),
      ...formValues,
      status: 'Aktif',
      ipk: '3.70'
    }
    students.value.push(newStudent)
    showSuccess('Mahasiswa Ditambahkan! 🎓', `${newStudent.name} (${newStudent.nim}) berhasil didaftarkan ke sistem.`)
  }
}

// Edit Student Modal Popup (SweetAlert2)
const openEditStudentModal = async (student) => {
  const { value: formValues } = await Swal.fire({
    title: `✏️ Edit Biodata Mahasiswa`,
    html: `
      <div style="text-align: left; display: flex; flex-direction: column; gap: 0.65rem; padding: 0.25rem 0.5rem;">
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">NIM Mahasiswa</label>
          <input id="swal-edit-nim" class="swal2-input" value="${student.nim}" placeholder="Contoh: 221112019" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.88rem;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Nama Lengkap Mahasiswa</label>
          <input id="swal-edit-name" class="swal2-input" value="${student.name}" placeholder="Masukkan Nama Lengkap" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.88rem;">
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Program Studi</label>
          <select id="swal-edit-prodi" class="swal2-input" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem; padding: 0.5rem;">
            <option value="Teknik Informatika (S1)" ${student.prodi.includes('Informatika') ? 'selected' : ''}>Teknik Informatika (S1)</option>
            <option value="Sistem Informasi (S1)" ${student.prodi.includes('Sistem') ? 'selected' : ''}>Sistem Informasi (S1)</option>
          </select>
        </div>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.65rem;">
          <div>
            <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Email Kampus</label>
            <input id="swal-edit-email" class="swal2-input" value="${student.email}" placeholder="Email Kampus" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem;">
          </div>
          <div>
            <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">No. WhatsApp</label>
            <input id="swal-edit-phone" class="swal2-input" value="${student.phone}" placeholder="Contoh: 62812345678" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem;">
          </div>
        </div>
        <div>
          <label style="display: block; font-size: 0.8rem; font-weight: 700; color: #475569; margin-bottom: 0.25rem;">Status Akademik</label>
          <select id="swal-edit-status" class="swal2-input" style="margin: 0; width: 100%; box-sizing: border-box; font-size: 0.85rem; padding: 0.5rem;">
            <option value="Aktif" ${student.status === 'Aktif' ? 'selected' : ''}>Aktif</option>
            <option value="Cuti" ${student.status === 'Cuti' ? 'selected' : ''}>Cuti</option>
          </select>
        </div>
      </div>
    `,
    focusConfirm: false,
    showCancelButton: true,
    confirmButtonText: 'Simpan Perubahan',
    cancelButtonText: 'Batal',
    confirmButtonColor: '#2563eb',
    cancelButtonColor: '#64748b',
    preConfirm: () => {
      const nim = document.getElementById('swal-edit-nim').value.trim()
      const name = document.getElementById('swal-edit-name').value.trim()
      const prodi = document.getElementById('swal-edit-prodi').value
      const email = document.getElementById('swal-edit-email').value.trim()
      const phone = document.getElementById('swal-edit-phone').value.trim()
      const status = document.getElementById('swal-edit-status').value

      if (!nim || !name) {
        Swal.showValidationMessage('NIM dan Nama Mahasiswa wajib diisi!')
        return false
      }
      return { nim, name: name.toUpperCase(), prodi, email, phone, status }
    }
  })

  if (formValues) {
    Object.assign(student, formValues)
    showSuccess('Perubahan Disimpan! 💾', `Data mahasiswa ${student.name} berhasil diperbarui.`)
  }
}

// Delete Student (SweetAlert2 Confirm)
const deleteStudent = async (student) => {
  const confirmed = await showConfirm(
    'Hapus Data Mahasiswa?',
    `Apakah Anda yakin ingin menghapus mahasiswa ${student.name} (${student.nim}) dari database?`,
    'Ya, Hapus Mahasiswa'
  )
  if (confirmed) {
    students.value = students.value.filter(s => s.id !== student.id)
    showSuccess('Mahasiswa Dihapus 🗑️', `Data ${student.name} berhasil dihapus dari sistem.`)
  }
}
</script>

<template>
  <div class="mahasiswa-container animate-fade-in">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-main">
        <div class="header-title-group">
          <GraduationCap class="header-icon text-blue" />
          <div>
            <h2>Data & Roster Mahasiswa ITB Swadharma</h2>
            <p class="subtitle">Manajemen data profil dan biodata akademik mahasiswa Dosen Rio Widyatmoko.</p>
          </div>
        </div>

        <button @click="openAddStudentModal" class="btn btn-primary-add">
          <UserPlus class="btn-icon-xs" />
          <span>Tambah Mahasiswa Baru</span>
        </button>
      </div>
    </div>

    <!-- Summary Stats Chips -->
    <div class="stats-row">
      <div class="stat-card stat-total">
        <div class="stat-icon-wrapper text-blue">
          <Users class="stat-icon" />
        </div>
        <div>
          <span class="stat-label">Total Mahasiswa</span>
          <h3 class="stat-value">{{ stats.total }}</h3>
        </div>
      </div>

      <div class="stat-card stat-aktif">
        <div class="stat-icon-wrapper text-emerald">
          <CheckCircle2 class="stat-icon" />
        </div>
        <div>
          <span class="stat-label">Mahasiswa Aktif</span>
          <h3 class="stat-value">{{ stats.aktif }}</h3>
        </div>
      </div>

      <div class="stat-card stat-2024">
        <div class="stat-icon-wrapper text-gold">
          <BookOpen class="stat-icon" />
        </div>
        <div>
          <span class="stat-label">Teknik Informatika</span>
          <h3 class="stat-value">{{ stats.tif }}</h3>
        </div>
      </div>

      <div class="stat-card stat-2023">
        <div class="stat-icon-wrapper text-sky">
          <GraduationCap class="stat-icon" />
        </div>
        <div>
          <span class="stat-label">Sistem Informasi</span>
          <h3 class="stat-value">{{ stats.si }}</h3>
        </div>
      </div>
    </div>

    <!-- Filter & Toolbar Card -->
    <div class="filter-card glass-card">
      <div class="search-box">
        <Search class="search-icon" />
        <input 
          v-model="searchQuery" 
          class="search-input" 
          placeholder="Cari Nama Mahasiswa, NIM, atau Email..." 
        />
      </div>

      <div class="filter-controls">
        <div class="filter-item">
          <label class="filter-lbl">Program Studi:</label>
          <select v-model="selectedProdi" class="filter-select">
            <option value="all">Semua Prodi</option>
            <option value="Informatika">Teknik Informatika</option>
            <option value="Sistem">Sistem Informasi</option>
          </select>
        </div>

        <div class="filter-item">
          <label class="filter-lbl">Status:</label>
          <select v-model="selectedStatus" class="filter-select">
            <option value="all">Semua Status</option>
            <option value="Aktif">Aktif</option>
            <option value="Cuti">Cuti</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Students Table Card -->
    <div class="table-card glass-card">
      <div class="table-responsive">
        <table class="students-table">
          <thead>
            <tr>
              <th class="col-no">No</th>
              <th class="col-nim">NIM</th>
              <th class="col-name">Nama Mahasiswa</th>
              <th class="col-prodi">Program Studi</th>
              <th class="col-status">Status</th>
              <th class="col-contact">Kontak / Email</th>
              <th class="col-actions">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(student, idx) in filteredStudents" :key="student.id">
              <td class="col-no">{{ idx + 1 }}</td>
              <td class="col-nim">
                <span class="nim-badge">{{ student.nim }}</span>
              </td>
              <td class="col-name">
                <div class="student-name-cell">
                  <h4 class="student-name-text">{{ student.name }}</h4>
                </div>
              </td>
              <td class="col-prodi">
                <span class="prodi-text">{{ student.prodi }}</span>
              </td>
              <td class="col-status">
                <span :class="['status-pill', student.status === 'Aktif' ? 'pill-aktif' : 'pill-cuti']">
                  {{ student.status }}
                </span>
              </td>
              <td class="col-contact">
                <div class="contact-info">
                  <div class="email-line">
                    <Mail class="contact-icon" />
                    <span>{{ student.email }}</span>
                  </div>
                  <a 
                    :href="'https://wa.me/' + student.phone" 
                    target="_blank" 
                    class="wa-link-btn"
                    title="Kirim pesan WhatsApp ke Mahasiswa"
                  >
                    <MessageSquare class="wa-icon" />
                    <span>WhatsApp</span>
                  </a>
                </div>
              </td>
              <td class="col-actions">
                <div class="action-buttons">
                  <button 
                    @click="openEditStudentModal(student)" 
                    class="btn-action btn-edit"
                    title="Edit Biodata Mahasiswa"
                  >
                    <Edit3 class="action-icon" />
                    <span>Edit</span>
                  </button>
                  <button 
                    @click="deleteStudent(student)" 
                    class="btn-action btn-delete"
                    title="Hapus Mahasiswa"
                  >
                    <Trash2 class="action-icon" />
                    <span>Hapus</span>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredStudents.length === 0">
              <td colspan="8" class="empty-cell">
                <div class="empty-state">
                  <Users class="empty-icon" />
                  <p>Tidak ada data mahasiswa yang cocok dengan kriteria pencarian.</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mahasiswa-container {
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

.btn-primary-add {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.6rem 1.25rem;
  border-radius: var(--radius-sm);
  border: none;
  background: #2563eb;
  color: #ffffff;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);
  transition: all 0.2s ease;
}

.btn-primary-add:hover {
  background: #1d4ed8;
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.35);
  transform: translateY(-1px);
}

/* Stats Row */
.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.25rem;
  margin-bottom: 1.5rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.02);
}

.stat-icon-wrapper {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon {
  width: 22px;
  height: 22px;
}

.stat-label {
  display: block;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
  margin-bottom: 0.15rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

/* Filter Card */
.filter-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.25rem;
  background: #ffffff;
  border: 1.5px solid #e2e8f0;
  border-radius: var(--radius-md);
  margin-bottom: 1.25rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 280px;
}

.search-icon {
  position: absolute;
  left: 0.85rem;
  width: 18px;
  height: 18px;
  color: #94a3b8;
}

.search-input {
  width: 100%;
  padding: 0.6rem 0.85rem 0.6rem 2.4rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #cbd5e1;
  font-size: 0.9rem;
  outline: none;
  transition: all 0.2s ease;
}

.search-input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.filter-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.filter-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.filter-lbl {
  font-size: 0.85rem;
  font-weight: 700;
  color: #475569;
}

.filter-select {
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-sm);
  border: 1.5px solid #cbd5e1;
  font-size: 0.85rem;
  font-weight: 600;
  color: #0f172a;
  outline: none;
  background: #ffffff;
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

.students-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.students-table th {
  background: #ffffff;
  padding: 1rem 1.25rem;
  font-size: 0.85rem;
  font-weight: 800;
  color: #334155;
  border-bottom: 2px solid #e2e8f0;
}

.students-table td {
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.students-table tr:hover {
  background: #f8fafc;
}

.col-no {
  width: 45px;
  font-weight: 700;
  color: #64748b;
  text-align: center;
}

.col-nim {
  min-width: 130px;
}

.nim-badge {
  font-weight: 800;
  color: #2563eb;
  font-size: 0.88rem;
  letter-spacing: 0.02em;
}

.col-name {
  min-width: 240px;
}

.student-name-text {
  font-size: 0.92rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
  letter-spacing: 0.01em;
}

.col-prodi {
  min-width: 180px;
}

.prodi-text {
  font-size: 0.85rem;
  color: #475569;
  font-weight: 600;
}

.col-angkatan {
  min-width: 100px;
}

.angkatan-badge {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 700;
  color: #334155;
}

.col-status {
  min-width: 90px;
}

.status-pill {
  display: inline-block;
  padding: 0.25rem 0.65rem;
  border-radius: 99px;
  font-size: 0.8rem;
  font-weight: 700;
}

.pill-aktif {
  background: #dcfce7;
  color: #15803d;
}

.pill-cuti {
  background: #fef3c7;
  color: #b45309;
}

.col-contact {
  min-width: 260px;
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.email-line {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.82rem;
  color: #475569;
  font-weight: 600;
}

.contact-icon {
  width: 14px;
  height: 14px;
  color: #64748b;
}

.wa-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.78rem;
  font-weight: 700;
  color: #16a34a;
  text-decoration: none;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  padding: 0.2rem 0.55rem;
  border-radius: 6px;
  width: fit-content;
  transition: all 0.2s ease;
}

.wa-link-btn:hover {
  background: #16a34a;
  color: #ffffff;
}

.wa-icon {
  width: 13px;
  height: 13px;
}

.col-actions {
  min-width: 140px;
}

.action-buttons {
  display: flex;
  gap: 0.4rem;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 0.35rem 0.65rem;
  border-radius: var(--radius-xs);
  font-size: 0.8rem;
  font-weight: 700;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-edit {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.btn-edit:hover {
  background: #2563eb;
  color: #ffffff;
}

.btn-delete {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.btn-delete:hover {
  background: #dc2626;
  color: #ffffff;
}

.action-icon {
  width: 14px;
  height: 14px;
}

.empty-cell {
  text-align: center;
  padding: 3rem 1rem !important;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  color: #94a3b8;
}

.empty-icon {
  width: 48px;
  height: 48px;
}
</style>
