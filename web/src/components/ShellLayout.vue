<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">
          <ShieldCheck :size="22" />
        </div>
        <div>
          <div class="brand-title">Asset-Core</div>
          <div class="brand-subtitle">可信设备资产管理</div>
        </div>
      </div>

      <nav class="nav">
        <RouterLink v-for="item in navItems" :key="item.path" class="nav-link" :to="item.path">
          <component :is="item.icon" :size="18" />
          <span>{{ item.label }}</span>
        </RouterLink>
        <RouterLink v-if="canRegisterUser" class="nav-link" to="/register">
          <UserPlus :size="18" />
          <span>注册账号</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="main">
      <header class="topbar">
        <div>
          <div class="page-kicker">Trusted Device Management</div>
          <h1>{{ currentTitle }}</h1>
        </div>
        <div class="topbar-actions">
          <el-tag :type="healthState === 'ok' ? 'success' : 'danger'" effect="dark">
            {{ healthState === 'ok' ? '服务在线' : '服务异常' }}
          </el-tag>
          <el-tag type="info">{{ currentUserLabel }}</el-tag>
          <el-button :icon="RefreshCw" @click="checkHealth">刷新</el-button>
          <el-button :icon="LogOut" @click="logout">退出</el-button>
        </div>
      </header>

      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router';
import { Database, FileStack, Fingerprint, Gauge, LogOut, RefreshCw, ShieldCheck, UserPlus, Workflow } from 'lucide-vue-next';
import { health } from '../services/api';
import { authState, clearSession, hasPermission } from '../services/auth';

const route = useRoute();
const healthState = ref<'ok' | 'down'>('down');

const navItems = [
  { path: '/dashboard', label: '运行总览', icon: Gauge },
  { path: '/assets', label: '资产台账', icon: Database },
  { path: '/identities', label: '身份 ID', icon: Fingerprint },
  { path: '/verifications', label: '交叉验证', icon: Workflow },
  { path: '/data', label: '数据管理', icon: FileStack }
];

const titleMap: Record<string, string> = {
  dashboard: '运行总览',
  assets: '资产台账',
  identities: '身份 ID',
  verifications: '交叉验证',
  data: '数据管理',
  register: '注册账号'
};

const currentTitle = computed(() => titleMap[String(route.name)] || 'Asset-Core');
const canRegisterUser = computed(() => hasPermission('user:create'));
const currentUserLabel = computed(() => authState.user?.display_name || authState.user?.username || '未登录');

async function checkHealth() {
  try {
    await health();
    healthState.value = 'ok';
  } catch {
    healthState.value = 'down';
  }
}

function logout() {
  clearSession();
  location.href = '/login';
}

onMounted(checkHealth);
</script>
