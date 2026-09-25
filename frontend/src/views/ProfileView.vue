<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  User, 
  Mail, 
  Phone, 
  MapPin, 
  Clock, 
  Building2, 
  Award, 
  KeyRound, 
  Save, 
  CheckCircle2, 
  ShieldCheck,
  Camera,
  UploadCloud,
  Image,
  BookOpen,
  GraduationCap
} from 'lucide-vue-next'

import { showSuccess, showError, showWarning } from '../utils/swal.js'

const currentUser = ref(null)

const isDosen = computed(() => currentUser.value?.role === 'dosen')
const isMahasiswa = computed(() => currentUser.value?.role === 'mahasiswa')

// Dosen Profile State
const profile = ref({
  name: '',
  degree: '',
  title: '',
  nidn: '',
  institution: '',
  faculty: '',
  department: '',
  email: '',
  phone: '',
  office_hours: '',
  room: '',
  avatar: '',
  bio: ''
})

// Student Profile State
const studentForm = ref({
  name: '',
  username: '',
  email: '',
  prodi: 'Teknik Informatika',
  phone: '',
  avatar: ''
})

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const saveSuccess = ref(false)
const passwordSuccess = ref(false)
const errorMessage = ref('')
const fileInputRef = ref(null)

const dosenAvatarPresets = [
  { name: 'Formal Male Lecturer 1', url: 'https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&q=80&w=400' },
  { name: 'Executive Male Lecturer 2', url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=400' },
  { name: 'Academic Male Lecturer 3', url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&q=80&w=400' },
  { name: 'Modern Lecturer 4', url: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=crop&q=80&w=400' }
]

const studentAvatarPresets = [
  { name: 'Student 1', url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=400' },
  { name: 'Student 2', url: 'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?auto=format&fit=crop&q=80&w=400' },
  { name: 'Student 3', url: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?auto=format&fit=crop&q=80&w=400' },
  { name: 'Student 4', url: 'https://images.unsplash.com/photo-1524504388940-b1c1722653e1?auto=format&fit=crop&q=80&w=400' }
]

const triggerFileInput = () => {
  if (fileInputRef.value) {
    fileInputRef.value.click()
  }
}

const onFileSelected = (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    showWarning('Perhatian', 'Harap pilih file gambar (JPG, PNG, WEBP).')
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    showWarning('Perhatian', 'Ukuran foto maksimal 5 MB.')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    if (isMahasiswa.value) {
      studentForm.value.avatar = e.target.result
    } else {
      profile.value.avatar = e.target.result
    }
    showSuccess('Foto Terpilih! 📷', 'Foto profil baru berhasil diunggah. Klik "Simpan Perubahan Profil".')
  }
  reader.readAsDataURL(file)
}

const selectAvatarPreset = (url) => {
  if (isMahasiswa.value) {
    studentForm.value.avatar = url
  } else {
    profile.value.avatar = url
  }
  showSuccess('Foto Profil Terpilih!', 'Foto profil baru telah diterapkan.')
}

const loadUserData = async () => {
  const stored = localStorage.getItem('user')
  if (stored) {
    try {
      currentUser.value = JSON.parse(stored)
    } catch { currentUser.value = null }
  }

  if (isMahasiswa.value && currentUser.value) {
    studentForm.value = {
      name: currentUser.value.name || 'Ahmad Fauzi',
      username: currentUser.value.username || '20260801001',
      email: currentUser.value.email || 'ahmad.fauzi@student.swadharma.ac.id',
      prodi: currentUser.value.prodi || 'Teknik Informatika',
      phone: currentUser.value.phone || '+62 812-3456-7890',
      avatar: currentUser.value.avatar || studentAvatarPresets[0].url
    }
  } else {
    fetchDosenProfile()
  }
}

const fetchDosenProfile = async () => {
  try {
    const res = await fetch('/api/profile')
    if (res.ok) {
      profile.value = await res.json()
      if (!profile.value.avatar) {
        profile.value.avatar = dosenAvatarPresets[0].url
      }
    }
  } catch (err) {
    console.warn('Fetch profile error:', err)
  }
}

onMounted(() => {
  loadUserData()
})

const updateStudentProfile = () => {
  if (!studentForm.value.name || !studentForm.value.email) {
    showWarning('Perhatian', 'Nama Lengkap dan Email wajib diisi!')
    return
  }

  const stored = localStorage.getItem('user')
  if (stored) {
    try {
      const u = JSON.parse(stored)
      u.name = studentForm.value.name
      u.email = studentForm.value.email
      u.prodi = studentForm.value.prodi
      u.phone = studentForm.value.phone
      u.avatar = studentForm.value.avatar
      localStorage.setItem('user', JSON.stringify(u))
      currentUser.value = u
      window.dispatchEvent(new Event('auth-changed'))

      saveSuccess.value = true
      showSuccess('Profil Mahasiswa Diperbarui! 👤', 'Data profil dan informasi akun Anda berhasil diperbarui secara realtime.')
      setTimeout(() => saveSuccess.value = false, 3500)
    } catch (e) {
      showError('Gagal!', 'Gagal menyimpan data profil.')
    }
  }
}

const updateProfileData = async () => {
  if (isMahasiswa.value) {
    updateStudentProfile()
    return
  }

  try {
    const res = await fetch('/api/profile/update', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(profile.value)
    })

    if (res.ok) {
      const stored = localStorage.getItem('user')
      if (stored) {
        try {
          const u = JSON.parse(stored)
          u.name = profile.value.name + (profile.value.degree ? ', ' + profile.value.degree : '')
          u.email = profile.value.email
          u.avatar = profile.value.avatar
          localStorage.setItem('user', JSON.stringify(u))
          window.dispatchEvent(new Event('auth-changed'))
        } catch (e) {}
      }

      saveSuccess.value = true
      showSuccess('Profil Diperbarui! 🎉', 'Foto profil & data diri akademik Anda telah berhasil disimpan.')
      setTimeout(() => saveSuccess.value = false, 3500)
    }
  } catch (err) {
    showError('Gagal Menyimpan!', 'Terjadi kesalahan sistem saat memperbarui profil.')
  }
}

const updatePassword = () => {
  errorMessage.value = ''
  if (!currentPassword.value) {
    errorMessage.value = 'Harap masukkan Password saat ini!'
    showWarning('Perhatian', errorMessage.value)
    return
  }
  if (!newPassword.value || newPassword.value.length < 6) {
    errorMessage.value = 'Password baru minimal 6 karakter!'
    showWarning('Perhatian', errorMessage.value)
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Konfirmasi password baru tidak cocok!'
    showWarning('Perhatian', errorMessage.value)
    return
  }

  const stored = localStorage.getItem('user')
  if (stored) {
    try {
      const u = JSON.parse(stored)
      u.password = newPassword.value
      localStorage.setItem('user', JSON.stringify(u))
    } catch (e) {}
  }

  passwordSuccess.value = true
  showSuccess('Password Berhasil Diubah! 🔑', 'Silakan gunakan password baru ini untuk login berikutnya.')
  currentPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  setTimeout(() => passwordSuccess.value = false, 3500)
}
</script>

<template>
  <div class="profile-container animate-fade-in">
    <!-- Header for Mahasiswa -->
    <div v-if="isMahasiswa" class="page-header">
      <div class="header-content">
        <div class="avatar-large-preview">
          <img :src="studentForm.avatar || studentAvatarPresets[0].url" alt="Foto Mahasiswa" class="avatar-img-preview" />
        </div>
        <div>
          <h2>Edit Profil Mahasiswa — {{ studentForm.name }}</h2>
          <p class="subtitle">NIM: <strong>{{ studentForm.username }}</strong> | Prodi: <strong>{{ studentForm.prodi }}</strong> — Kelola data diri, kontak, foto avatar, dan kata sandi akun Mahasiswa ITB Swadharma.</p>
        </div>
      </div>
    </div>

    <!-- Header for Dosen -->
    <div v-else class="page-header">
      <div class="header-content">
        <div class="avatar-large-preview">
          <img :src="profile.avatar || dosenAvatarPresets[0].url" alt="Foto Dosen" class="avatar-img-preview" />
        </div>
        <div>
          <h2>Edit Profil Akademik — {{ profile.name }}</h2>
          <p class="subtitle">NIP: <strong>{{ profile.nidn || '21099001' }}</strong> — Kelola foto profil, data diri, kontak, jam konsultasi, dan kata sandi akun ITB Swadharma.</p>
        </div>
      </div>
    </div>

    <!-- Main Form Block -->
    <div class="profile-grid">
      <!-- Left Card: Edit Main Profile -->
      <div class="glass-card profile-card">
        <div class="card-title-row">
          <User class="title-icon text-gold" />
          <div>
            <h3>{{ isMahasiswa ? 'Informasi Diri & Akun Mahasiswa' : 'Informasi Diri & Akademik Dosen' }}</h3>
            <p class="card-sub">{{ isMahasiswa ? 'Perbarui foto profil, nama lengkap, email, dan kontak WhatsApp Anda.' : 'Perbarui foto profil, nama, gelar, departemen, dan kontak publik.' }}</p>
          </div>
        </div>

        <div v-if="saveSuccess" class="alert-success">
          <CheckCircle2 class="icon-sm" />
          <span>Profil berhasil diperbarui secara realtime!</span>
        </div>

        <!-- FORM UNTUK MAHASISWA -->
        <form v-if="isMahasiswa" @submit.prevent="updateStudentProfile" class="profile-form">
          <!-- Upload Foto Profil Section -->
          <div class="avatar-upload-box">
            <div class="avatar-preview-wrapper">
              <img :src="studentForm.avatar || studentAvatarPresets[0].url" alt="Foto Avatar Mahasiswa" class="avatar-preview-img" />
              <button type="button" @click="triggerFileInput" class="btn-avatar-cam" title="Upload Foto dari Perangkat">
                <Camera class="cam-icon" />
              </button>
            </div>
            
            <div class="avatar-upload-info">
              <h4>📷 Foto Profil Mahasiswa (Avatar)</h4>
              <p>Upload foto diri Anda dari komputer (Format: JPG, PNG, WEBP, Maks. 5MB).</p>
              
              <input 
                type="file" 
                ref="fileInputRef" 
                accept="image/*" 
                @change="onFileSelected" 
                style="display: none" 
              />
              
              <div class="avatar-actions">
                <button type="button" @click="triggerFileInput" class="btn btn-secondary btn-sm">
                  <UploadCloud class="btn-icon-xs" />
                  <span>Pilih & Upload Foto...</span>
                </button>
              </div>

              <!-- Avatar Presets Selection -->
              <div class="preset-avatars-group">
                <span class="preset-label">Atau pilih sampel foto avatar mahasiswa:</span>
                <div class="preset-list">
                  <img 
                    v-for="(av, idx) in studentAvatarPresets" 
                    :key="idx" 
                    :src="av.url" 
                    :alt="av.name"
                    class="preset-thumb"
                    :class="{ 'selected': studentForm.avatar === av.url }"
                    @click="selectAvatarPreset(av.url)"
                    :title="av.name"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Nama Lengkap *</label>
              <input v-model="studentForm.name" class="glass-input" required placeholder="Nama Mahasiswa..." />
            </div>
            <div>
              <label class="input-label">NIM (Nomor Induk Mahasiswa)</label>
              <input v-model="studentForm.username" class="glass-input disabled-input" readonly title="NIM tidak dapat diubah" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Program Studi *</label>
              <select v-model="studentForm.prodi" class="glass-input">
                <option value="Teknik Informatika">Teknik Informatika</option>
                <option value="Sistem Informasi">Sistem Informasi</option>
                <option value="Desain Komunikasi Visual">Desain Komunikasi Visual</option>
              </select>
            </div>
            <div>
              <label class="input-label">Email Resmi Mahasiswa *</label>
              <input v-model="studentForm.email" type="email" class="glass-input" required placeholder="student@swadharma.ac.id" />
            </div>
          </div>

          <div>
            <label class="input-label">Nomor WhatsApp Direct</label>
            <input v-model="studentForm.phone" class="glass-input" placeholder="+62 812-3456-7890" />
          </div>

          <button type="submit" class="btn btn-gold btn-lg">
            <Save class="btn-icon-sm" />
            <span>Simpan Perubahan Profil Mahasiswa</span>
          </button>
        </form>

        <!-- FORM UNTUK DOSEN (RIO WIDYATMOKO) -->
        <form v-else @submit.prevent="updateProfileData" class="profile-form">
          <!-- Upload Foto Profil Section -->
          <div class="avatar-upload-box">
            <div class="avatar-preview-wrapper">
              <img :src="profile.avatar || dosenAvatarPresets[0].url" alt="Foto Profil Dosen" class="avatar-preview-img" />
              <button type="button" @click="triggerFileInput" class="btn-avatar-cam" title="Upload Foto dari Perangkat">
                <Camera class="cam-icon" />
              </button>
            </div>
            
            <div class="avatar-upload-info">
              <h4>📷 Foto Profil Dosen (Avatar)</h4>
              <p>Upload foto resmi Pak Rio dari komputer (Format: JPG, PNG, WEBP, Maks. 5MB).</p>
              
              <input 
                type="file" 
                ref="fileInputRef" 
                accept="image/*" 
                @change="onFileSelected" 
                style="display: none" 
              />
              
              <div class="avatar-actions">
                <button type="button" @click="triggerFileInput" class="btn btn-secondary btn-sm">
                  <UploadCloud class="btn-icon-xs" />
                  <span>Pilih & Upload Foto...</span>
                </button>
              </div>

              <!-- Avatar Presets Selection -->
              <div class="preset-avatars-group">
                <span class="preset-label">Atau pilih sampel foto dosen profesional:</span>
                <div class="preset-list">
                  <img 
                    v-for="(p, idx) in dosenAvatarPresets" 
                    :key="idx" 
                    :src="p.url" 
                    :alt="p.name"
                    class="preset-thumb"
                    :class="{ 'selected': profile.avatar === p.url }"
                    @click="selectAvatarPreset(p.url)"
                    :title="p.name"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Nama Lengkap *</label>
              <input v-model="profile.name" class="glass-input" required placeholder="Contoh: Rio Widyatmoko" />
            </div>
            <div>
              <label class="input-label">Gelar Akademik</label>
              <input v-model="profile.degree" class="glass-input" placeholder="Contoh: S.Kom, M.M.S.I" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Jabatan Akademik / Title</label>
              <input v-model="profile.title" class="glass-input" placeholder="Dosen Pengampu ITB Swadharma" />
            </div>
            <div>
              <label class="input-label">NIP Dosen *</label>
              <input v-model="profile.nidn" class="glass-input" required placeholder="21099001" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Email Akademik</label>
              <input v-model="profile.email" type="email" class="glass-input" required placeholder="rio.widyatmoko@swadharma.ac.id" />
            </div>
            <div>
              <label class="input-label">Nomor HP / WhatsApp Direct</label>
              <input v-model="profile.phone" class="glass-input" placeholder="+62 812-9876-5432" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Institusi Kampus</label>
              <input v-model="profile.institution" class="glass-input" placeholder="Institut Teknologi dan Bisnis Swadharma" />
            </div>
            <div>
              <label class="input-label">Fakultas / Program Studi</label>
              <input v-model="profile.department" class="glass-input" placeholder="Teknik Informatika / Sistem Informasi" />
            </div>
          </div>

          <div class="form-row">
            <div>
              <label class="input-label">Jam Konsultasi Dosen</label>
              <input v-model="profile.office_hours" class="glass-input" placeholder="Senin & Rabu: 13.00 - 16.00 WIB" />
            </div>
            <div>
              <label class="input-label">Lokasi Ruang Dosen</label>
              <input v-model="profile.room" class="glass-input" placeholder="Ruang Dosen Lantai 3 - Kampus Swadharma" />
            </div>
          </div>

          <div>
            <label class="input-label">Bio & Profil Singkat Dosen</label>
            <textarea v-model="profile.bio" class="glass-input textarea" placeholder="Deskripsi profil dosen..."></textarea>
          </div>

          <button type="submit" class="btn btn-gold btn-lg">
            <Save class="btn-icon-sm" />
            <span>Simpan Perubahan Profil Dosen</span>
          </button>
        </form>
      </div>

      <!-- Right Card: Keamanan Akun & Kata Sandi -->
      <div class="glass-card profile-card">
        <div class="card-title-row">
          <KeyRound class="title-icon text-gold" />
          <div>
            <h3>Keamanan Akun & Kata Sandi</h3>
            <p class="card-sub">Ubah password akun login Anda secara berkala.</p>
          </div>
        </div>

        <div v-if="passwordSuccess" class="alert-success">
          <CheckCircle2 class="icon-sm" />
          <span>Password berhasil diperbarui!</span>
        </div>

        <form @submit.prevent="updatePassword" class="password-form">
          <div>
            <label class="input-label">Password Saat Ini *</label>
            <input 
              v-model="currentPassword" 
              type="password" 
              class="glass-input" 
              required 
              placeholder="Masukkan password saat ini (cth: 123456)" 
            />
          </div>

          <div>
            <label class="input-label">Password Baru *</label>
            <input 
              v-model="newPassword" 
              type="password" 
              class="glass-input" 
              required 
              placeholder="Minimal 6 karakter" 
            />
          </div>

          <div>
            <label class="input-label">Konfirmasi Password Baru *</label>
            <input 
              v-model="confirmPassword" 
              type="password" 
              class="glass-input" 
              required 
              placeholder="Ulangi password baru" 
            />
          </div>

          <button type="submit" class="btn btn-primary w-full">
            <KeyRound class="btn-icon-sm" />
            <span>Ubah Password Sekarang</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

.page-header {
  margin-bottom: 2rem;
  background: #ffffff;
  padding: 1.5rem;
  border-radius: var(--radius-md);
  border: 1px solid #e2e8f0;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.avatar-large-preview {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid #d97706;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.avatar-img-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.page-header h2 {
  font-size: 1.6rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.subtitle {
  color: #475569;
  font-size: 0.9rem;
  margin-top: 0.3rem;
}

.profile-grid {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 1.5rem;
}

.profile-card {
  padding: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid #f1f5f9;
}

.title-icon {
  width: 24px;
  height: 24px;
}

.card-title-row h3 {
  font-size: 1.2rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.card-sub {
  color: #475569;
  font-size: 0.85rem;
  margin: 0.2rem 0 0 0;
}

.alert-success {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #047857;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
  font-size: 0.88rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.profile-form, .password-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.input-label {
  display: block;
  font-size: 0.85rem;
  color: #334155;
  font-weight: 700;
  margin-bottom: 0.35rem;
}

.textarea {
  min-height: 90px;
  resize: vertical;
}

.disabled-input {
  background: #f1f5f9;
  color: #64748b;
  cursor: not-allowed;
}

.avatar-upload-box {
  display: flex;
  align-items: flex-start;
  gap: 1.5rem;
  padding: 1.25rem;
  background: #f8fafc;
  border-radius: var(--radius-sm);
  border: 1.5px solid #e2e8f0;
  margin-bottom: 0.5rem;
}

.avatar-preview-wrapper {
  position: relative;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  border: 3px solid #2563eb;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  flex-shrink: 0;
}

.avatar-preview-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.btn-avatar-cam {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 32px;
  height: 32px;
  background: #d97706;
  border: 2px solid #ffffff;
  border-radius: 50%;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  transition: all 0.2s ease;
}

.btn-avatar-cam:hover {
  background: #b45309;
  transform: scale(1.08);
}

.cam-icon {
  width: 16px;
  height: 16px;
}

.avatar-upload-info h4 {
  font-size: 0.95rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.avatar-upload-info p {
  font-size: 0.82rem;
  color: #64748b;
  margin: 0.25rem 0 0.75rem 0;
}

.avatar-actions {
  margin-bottom: 1rem;
}

.preset-avatars-group {
  margin-top: 0.75rem;
}

.preset-label {
  display: block;
  font-size: 0.8rem;
  color: #475569;
  font-weight: 700;
  margin-bottom: 0.4rem;
}

.preset-list {
  display: flex;
  gap: 0.6rem;
}

.preset-thumb {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #cbd5e1;
  cursor: pointer;
  transition: all 0.2s ease;
}

.preset-thumb:hover {
  transform: scale(1.1);
  border-color: #2563eb;
}

.preset-thumb.selected {
  border-color: #d97706;
  box-shadow: 0 0 0 3px rgba(217, 119, 6, 0.3);
}

.w-full {
  width: 100%;
}

.text-gold {
  color: #d97706;
}

.icon-sm {
  width: 16px;
  height: 16px;
}

@media (max-width: 900px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }
}
</style>
