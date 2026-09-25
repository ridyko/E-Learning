<script setup>
import { 
  UserCheck, 
  Mail, 
  Phone, 
  Clock, 
  MapPin, 
  Award, 
  BookOpen, 
  MessageSquare,
  ExternalLink
} from 'lucide-vue-next'

defineProps({
  profile: Object
})
</script>

<template>
  <div class="glass-card profile-banner" v-if="profile">
    <div class="profile-layout">
      <!-- Avatar & Badge -->
      <div class="avatar-section">
        <div class="avatar-wrapper">
          <img :src="profile.avatar" :alt="profile.name" class="avatar-img" />
          <div class="verified-badge" title="Dosen Terverifikasi ITB Swadharma">
            <UserCheck class="badge-icon" />
          </div>
        </div>
        <div class="institution-pill">
          <Award class="pill-icon" />
          <span>NIP: {{ profile.nidn || '21099001' }}</span>
        </div>
      </div>

      <!-- Detail Dosen -->
      <div class="info-section">
        <div class="header-name">
          <span class="badge badge-gold">DOSEN PENGAMPU ITB SWADHARMA</span>
          <h2>{{ profile.name }}, {{ profile.degree }}</h2>
          <p class="title-sub">{{ profile.title }} — {{ profile.department }}</p>
        </div>

        <p class="bio-text">{{ profile.bio }}</p>

        <!-- Expertise Tags -->
        <div class="expertise-list">
          <span v-for="(tag, idx) in profile.expertise_areas" :key="idx" class="tag-item">
            <BookOpen class="tag-icon" />
            {{ tag }}
          </span>
        </div>

        <!-- Contact & Office Info Grid -->
        <div class="meta-grid">
          <div class="meta-item">
            <div class="icon-circle">
              <Mail class="meta-icon" />
            </div>
            <div>
              <span class="meta-label">Email Resmi</span>
              <a :href="'mailto:' + profile.email" class="meta-value">{{ profile.email }}</a>
            </div>
          </div>

          <div class="meta-item">
            <div class="icon-circle">
              <Clock class="meta-icon" />
            </div>
            <div>
              <span class="meta-label">Jam Konsultasi</span>
              <span class="meta-value">{{ profile.office_hours }}</span>
            </div>
          </div>

          <div class="meta-item">
            <div class="icon-circle">
              <MapPin class="meta-icon" />
            </div>
            <div>
              <span class="meta-label">Lokasi Ruang</span>
              <span class="meta-value">{{ profile.room }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Action CTA -->
      <div class="action-section">
        <a 
          :href="'https://wa.me/' + profile.phone?.replace(/[^0-9]/g, '') + '?text=Halo%20Pak%20Rio%20Widyatmoko,%20saya%20mahasiswa%20ITB%20Swadharma%20...'" 
          target="_blank" 
          class="btn btn-gold wa-btn"
        >
          <MessageSquare class="btn-icon-sm" />
          <span>Chat WhatsApp Pak Rio</span>
          <ExternalLink class="btn-icon-xs" />
        </a>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-banner {
  padding: 2.25rem;
  margin-bottom: 2rem;
  background: #ffffff;
  border: 1px solid #e2e8f0;
}

.profile-layout {
  display: grid;
  grid-template-columns: 200px 1fr 220px;
  gap: 2rem;
  align-items: center;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.85rem;
}

.avatar-wrapper {
  position: relative;
  width: 140px;
  height: 140px;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #2563eb;
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.25);
}

.verified-badge {
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 34px;
  height: 34px;
  background: #d97706;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.badge-icon {
  width: 18px;
  height: 18px;
}

.institution-pill {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: #f8fafc;
  padding: 0.35rem 0.85rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 700;
  color: #334155;
  border: 1px solid #cbd5e1;
}

.pill-icon {
  width: 14px;
  height: 14px;
  color: #d97706;
}

.header-name h2 {
  font-size: 1.8rem;
  font-weight: 800;
  margin-top: 0.4rem;
  color: #0f172a;
}

.title-sub {
  color: #2563eb;
  font-weight: 700;
  font-size: 0.95rem;
}

.bio-text {
  color: #334155;
  font-size: 0.95rem;
  margin: 0.8rem 0;
  line-height: 1.6;
}

.expertise-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.25rem;
}

.tag-item {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  background: #eff6ff;
  color: #1d4ed8;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-sm);
  border: 1px solid #bfdbfe;
}

.tag-icon {
  width: 13px;
  height: 13px;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
  gap: 1.5rem;
  padding-top: 1.25rem;
  border-top: 1px solid #e2e8f0;
}

.meta-item {
  display: flex;
  align-items: flex-start;
  gap: 0.85rem;
  min-width: 0;
}

.meta-item > div:last-child {
  min-width: 0;
  flex: 1;
}

.icon-circle {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: #fffbeb;
  border: 1.5px solid #fde68a;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.meta-icon {
  width: 19px;
  height: 19px;
  color: #d97706;
}

.meta-label {
  display: block;
  font-size: 0.75rem;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.05em;
  margin-bottom: 0.2rem;
}

.meta-value {
  font-size: 0.88rem;
  color: #0f172a;
  font-weight: 700;
  text-decoration: none;
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.45;
  display: block;
}

.meta-value:hover {
  color: #1d4ed8;
}

.action-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.wa-btn {
  width: 100%;
  padding: 0.85rem 1rem;
}

.btn-icon-sm {
  width: 18px;
  height: 18px;
}

.btn-icon-xs {
  width: 14px;
  height: 14px;
}

@media (max-width: 1024px) {
  .profile-layout {
    grid-template-columns: 1fr;
    text-align: center;
  }
  .avatar-section {
    margin-bottom: 1rem;
  }
  .expertise-list {
    justify-content: center;
  }
  .meta-grid {
    text-align: left;
  }
}
</style>
