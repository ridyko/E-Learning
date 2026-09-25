<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { LogIn, ShieldCheck, UserCheck, Key, User, UserPlus, AlertCircle } from 'lucide-vue-next'
import { showToast, showError } from '../utils/swal.js'

const router = useRouter()
const route = useRoute()

const activeTab = ref('dosen')
const username = ref('21099001')
const password = ref('123456')
const errorMsg = ref('')
const loading = ref(false)

const setRoleTab = (role) => {
  activeTab.value = role
  errorMsg.value = ''
  if (role === 'dosen') {
    username.value = '21099001'
    password.value = '123456'
  } else {
    username.value = '20260801001'
    password.value = '123456'
  }
}

const handleLogin = async () => {
  if (!username.value || !password.value) {
    errorMsg.value = 'Harap isi Username (NIP/NIM) dan Password!'
    showError('Login Gagal', errorMsg.value)
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value,
        password: password.value
      })
    })

    const data = await res.json()

    if (!res.ok) {
      errorMsg.value = data.error || 'Login gagal, periksa NIP/NIM dan password Anda.'
      showError('Login Gagal', errorMsg.value)
      return
    }

    localStorage.setItem('user', JSON.stringify(data.user))
    localStorage.setItem('token', data.token)
    window.dispatchEvent(new Event('auth-changed'))

    showToast(`Selamat datang, ${data.user.name}! 👋`)

    if (route.query.redirect) {
      router.push(route.query.redirect)
    } else if (data.user.role === 'dosen') {
      router.push('/admin')
    } else {
      router.push('/')
    }
  } catch (err) {
    errorMsg.value = 'Terjadi kesalahan koneksi ke server backend.'
    showError('Koneksi Gagal', errorMsg.value)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrapper animate-fade-in">
    <div class="glass-card login-card">
      <div class="login-header">
        <div class="badge badge-gold">ITB SWADHARMA</div>
        <h2>Portal Masuk E-Learning</h2>
        <p>Silakan masuk menggunakan kredensial akun Anda</p>
      </div>

      <div class="tab-switcher">
        <button 
          @click="setRoleTab('dosen')" 
          :class="['tab-btn', activeTab === 'dosen' ? 'active-dosen' : '']"
        >
          <ShieldCheck class="tab-icon text-gold" />
          <span>Login Dosen (Pak Rio)</span>
        </button>
        <button 
          @click="setRoleTab('mahasiswa')" 
          :class="['tab-btn', activeTab === 'mahasiswa' ? 'active-student' : '']"
        >
          <UserCheck class="tab-icon text-blue" />
          <span>Login Mahasiswa</span>
        </button>
      </div>

      <div class="info-alert">
        <AlertCircle class="alert-icon text-gold" />
        <div v-if="activeTab === 'dosen'">
          <strong>Login Dosen:</strong> Username: <code>21099001</code> | Password: <code>123456</code>
        </div>
        <div v-else>
          <strong>Login Mahasiswa:</strong> NIM: <code>20260801001</code> | Password Default: <code>123456</code>
        </div>
      </div>

      <div v-if="errorMsg" class="error-alert">
        ❌ {{ errorMsg }}
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div>
          <label class="input-label">
            {{ activeTab === 'dosen' ? 'NIP Dosen' : 'NIM Mahasiswa' }}
          </label>
          <div class="input-wrapper">
            <User class="field-icon" />
            <input 
              v-model="username" 
              class="glass-input with-icon" 
              required 
              :placeholder="activeTab === 'dosen' ? 'Masukkan NIP (21099001)' : 'Masukkan NIM Anda'"
            />
          </div>
        </div>

        <div>
          <label class="input-label">Password</label>
          <div class="input-wrapper">
            <Key class="field-icon" />
            <input 
              v-model="password" 
              type="password" 
              class="glass-input with-icon" 
              required 
              placeholder="Masukkan password"
            />
          </div>
        </div>

        <button type="submit" class="btn btn-gold btn-lg w-full" :disabled="loading">
          <LogIn class="btn-icon-sm" />
          <span>{{ loading ? 'Memproses Login...' : 'Masuk ke Portal E-Learning' }}</span>
        </button>
      </form>

      <div class="register-footer">
        <p>Belum memiliki akun mahasiswa?</p>
        <router-link to="/register" class="btn btn-secondary btn-sm">
          <UserPlus class="btn-icon-xs" />
          Pendaftaran Mahasiswa Mandiri
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-wrapper {
  max-width: 500px;
  margin: 3rem auto;
  padding: 1rem;
}

.login-card {
  padding: 2.5rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.login-header {
  text-align: center;
  margin-bottom: 1.75rem;
}

.login-header h2 {
  font-size: 1.75rem;
  font-weight: 800;
  color: #0f172a;
  margin-top: 0.5rem;
}

.login-header p {
  color: #475569;
  font-size: 0.9rem;
}

.tab-switcher {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  background: #f8fafc;
  padding: 0.35rem;
  border-radius: var(--radius-sm);
  margin-bottom: 1.25rem;
  border: 1px solid #e2e8f0;
}

.tab-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.6rem;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  color: #475569;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
  transition: var(--transition);
}

.tab-btn:hover {
  color: #0f172a;
}

.active-dosen {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}

.active-student {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.tab-icon {
  width: 16px;
  height: 16px;
}

.info-alert {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  background: #fffbeb;
  border: 1px solid #fde68a;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-sm);
  font-size: 0.83rem;
  color: #0f172a;
  margin-bottom: 1.25rem;
}

.alert-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

code {
  background: #e2e8f0;
  color: #0f172a;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-family: var(--font-code);
  font-weight: bold;
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

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
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

.register-footer {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e2e8f0;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

.register-footer p {
  font-size: 0.88rem;
  color: #475569;
}

.text-gold {
  color: #d97706;
}

.text-blue {
  color: #1d4ed8;
}

.btn-icon-sm {
  width: 18px;
  height: 18px;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}
</style>
