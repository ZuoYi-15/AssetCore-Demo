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
          <el-button :icon="RefreshCw" @click="checkHealth">刷新</el-button>
        </div>
      </header>

      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router';
import { Database, FileStack, Fingerprint, Gauge, RefreshCw, ShieldCheck, Workflow } from 'lucide-vue-next';
import { health } from '../services/api';

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
  data: '数据管理'
};

const currentTitle = computed(() => titleMap[String(route.name)] || 'Asset-Core');

async function checkHealth() {
  try {
    await health();
    healthState.value = 'ok';
  } catch {
    healthState.value = 'down';
  }
}

onMounted(checkHealth);
</script>
