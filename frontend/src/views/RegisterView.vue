<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { UserPlus, CheckCircle2, User, Mail, Lock, BookOpen, ArrowLeft, Phone, Calendar } from 'lucide-vue-next'

const router = useRouter()

const nim = ref('')
const name = ref('')
const email = ref('')
const prodi = ref('Teknik Informatika (S1)')
const phone = ref('')
const password = ref('123456')

const successMsg = ref('')
const errorMsg = ref('')
const loading = ref(false)

import { showSuccess, showError, showWarning } from '../utils/swal.js'

const handleRegister = async () => {
  if (!nim.value || !name.value || !email.value) {
    errorMsg.value = 'Harap lengkapi NIM, Nama Lengkap, dan Email!'
    showWarning('Perhatian', errorMsg.value)
    return
  }

  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    const res = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: nim.value,
        name: name.value,
        email: email.value,
        prodi: prodi.value,
        phone: phone.value,
        password: password.value || '123456'
      })
    })

    const data = await res.json()

    if (!res.ok) {
      errorMsg.value = data.error || 'Pendaftaran gagal!'
      showError('Pendaftaran Gagal', errorMsg.value)
      return
    }

    successMsg.value = `✅ Pendaftaran Akun Mahasiswa Berhasil! NIM: ${nim.value}. Mengalihkan ke halaman login...`
    showSuccess('Pendaftaran Berhasil! 🎉', `Akun NIM: ${nim.value} terdaftar. Silakan login dengan password default 123456.`)

    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err) {
    errorMsg.value = 'Terjadi kesalahan koneksi ke server.'
    showError('Koneksi Gagal', errorMsg.value)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-wrapper animate-fade-in">
    <div class="glass-card register-card">
      <router-link to="/login" class="back-link">
        <ArrowLeft class="icon-xs" />
        Kembali ke Login
      </router-link>

      <div class="register-header">
        <div class="badge badge-gold">ITB SWADHARMA</div>
        <h2>Pendaftaran Mahasiswa Mandiri</h2>
        <p>Isi data Anda untuk membuat akun E-Learning ITB Swadharma</p>
      </div>

      <div v-if="successMsg" class="success-alert">
        {{ successMsg }}
      </div>

      <div v-if="errorMsg" class="error-alert">
        ❌ {{ errorMsg }}
      </div>

      <form v-else @submit.prevent="handleRegister" class="register-form">
        <div>
          <label class="input-label">NIM Mahasiswa</label>
          <div class="input-wrapper">
            <User class="field-icon" />
            <input 
              v-model="nim" 
              class="glass-input with-icon" 
              required 
              placeholder="Contoh: 20260801003"
            />
          </div>
        </div>

        <div>
          <label class="input-label">Nama Lengkap Mahasiswa</label>
          <div class="input-wrapper">
            <User class="field-icon" />
            <input 
              v-model="name" 
              class="glass-input with-icon" 
              required 
              placeholder="Masukkan nama lengkap anda"
            />
          </div>
        </div>

        <div class="form-row">
          <div>
            <label class="input-label">Email Mahasiswa</label>
            <div class="input-wrapper">
              <Mail class="field-icon" />
              <input 
                v-model="email" 
                type="email"
                class="glass-input with-icon" 
                required 
                placeholder="Contoh: nama@gmail.com / student@swadharma.ac.id"
              />
            </div>
          </div>

          <div>
            <label class="input-label">No. WhatsApp</label>
            <div class="input-wrapper">
              <Phone class="field-icon" />
              <input 
                v-model="phone" 
                class="glass-input with-icon" 
                placeholder="Contoh: 62812345678"
              />
            </div>
          </div>
        </div>

        <div class="form-row">
          <div>
            <label class="input-label">Program Studi</label>
            <select v-model="prodi" class="glass-input">
              <option value="Teknik Informatika (S1)">Teknik Informatika (S1)</option>
              <option value="Sistem Informasi (S1)">Sistem Informasi (S1)</option>
              <option value="Teknologi Informasi (D3)">Teknologi Informasi (D3)</option>
            </select>
          </div>

          <div>
            <label class="input-label">Password (Default: 123456)</label>
            <div class="input-wrapper">
              <Lock class="field-icon" />
              <input 
                v-model="password" 
                type="password"
                class="glass-input with-icon" 
                placeholder="123456"
              />
            </div>
          </div>
        </div>

        <button type="submit" class="btn btn-gold btn-lg w-full" :disabled="loading">
          <UserPlus class="btn-icon-sm" />
          <span>{{ loading ? 'Mendaftarkan...' : 'Daftar Akun Mahasiswa Now' }}</span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.register-wrapper {
  max-width: 550px;
  margin: 2.5rem auto;
  padding: 1rem;
}

.register-card {
  padding: 2.5rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  color: #475569;
  font-weight: 600;
  text-decoration: none;
  margin-bottom: 1.25rem;
  transition: var(--transition);
}

.back-link:hover {
  color: #d97706;
}

.icon-xs {
  width: 14px;
  height: 14px;
}

.register-header {
  text-align: center;
  margin-bottom: 1.75rem;
}

.register-header h2 {
  font-size: 1.75rem;
  font-weight: 800;
  color: #0f172a;
  margin-top: 0.5rem;
}

.register-header p {
  color: #475569;
  font-size: 0.9rem;
}

.success-alert {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #047857;
  padding: 1.25rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  text-align: center;
}

.error-alert {
  background: #fff1f2;
  border: 1px solid #fecdd3;
  color: #be123c;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
  font-size: 0.88rem;
  margin-bottom: 1.25rem;
}

.register-form {
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
  font-weight: 600;
  margin-bottom: 0.35rem;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.field-icon {
  position: absolute;
  left: 1rem;
  width: 18px;
  height: 18px;
  color: #64748b;
}

.glass-input.with-icon {
  padding-left: 2.75rem;
}

.w-full {
  width: 100%;
}

.btn-icon-sm {
  width: 18px;
  height: 18px;
}
</style>
