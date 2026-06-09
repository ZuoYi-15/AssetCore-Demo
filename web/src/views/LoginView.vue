<template>
  <main class="auth-page">
    <section class="auth-panel">
      <div class="auth-brand">
        <div class="brand-mark">
          <ShieldCheck :size="22" />
        </div>
        <div>
          <div class="brand-title">Asset-Core</div>
          <div class="brand-subtitle">可信设备资产管理</div>
        </div>
      </div>

      <el-form class="auth-form" label-position="top" @submit.prevent="submit">
        <h1>登录</h1>
        <el-form-item label="用户名">
          <el-input v-model="form.username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <el-button type="primary" :loading="loading" native-type="submit" @click="submit">登录</el-button>
        <div class="auth-hint">默认超级管理员：superadmin / Admin@123456</div>
      </el-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { ShieldCheck } from 'lucide-vue-next';
import { login } from '../services/auth';

const router = useRouter();
const loading = ref(false);
const form = reactive({ username: 'superadmin', password: 'Admin@123456' });

async function submit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码');
    return;
  }
  loading.value = true;
  try {
    await login(form.username, form.password);
    await router.push('/dashboard');
  } finally {
    loading.value = false;
  }
}
</script>
