<template>
  <section class="panel register-panel">
    <div class="panel-header">
      <div>
        <div class="panel-title">注册账号</div>
        <div class="empty-hint">超级管理员可创建管理员或普通用户账号。</div>
      </div>
    </div>

    <el-form class="register-form" label-position="top" @submit.prevent="submit">
      <div class="form-grid">
        <el-form-item label="用户名">
          <el-input v-model="form.username" autocomplete="new-username" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" autocomplete="new-password" show-password />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_code" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </div>
      <div class="register-actions">
        <el-button type="primary" :icon="UserPlus" :loading="loading" native-type="submit" @click="submit">创建账号</el-button>
      </div>
    </el-form>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { UserPlus } from 'lucide-vue-next';
import { registerUser } from '../services/auth';
import type { RegisterPayload } from '../types/api';

const loading = ref(false);
const form = reactive<RegisterPayload>({
  username: '',
  password: '',
  display_name: '',
  role_code: 'user'
});

async function submit() {
  if (!form.username || !form.password || !form.role_code) {
    ElMessage.warning('请完整填写账号信息');
    return;
  }
  loading.value = true;
  try {
    await registerUser(form);
    ElMessage.success('账号创建成功');
    form.username = '';
    form.password = '';
    form.display_name = '';
    form.role_code = 'user';
  } finally {
    loading.value = false;
  }
}
</script>
