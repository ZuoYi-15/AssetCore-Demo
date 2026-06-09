<template>
  <section class="panel">
    <div class="panel-header">
      <div>
        <div class="panel-title">账号管理</div>
        <div class="empty-hint">超级管理员可创建账号、查看用户权限，并修改用户信息与密码。</div>
      </div>
      <div class="toolbar-right">
        <el-button :icon="RefreshCw" @click="loadUsers">刷新</el-button>
        <el-button type="primary" :icon="UserPlus" @click="openCreate">新增账号</el-button>
      </div>
    </div>

    <div class="toolbar user-filter">
      <div class="toolbar-left">
        <el-input v-model="query.keyword" clearable placeholder="搜索用户名、显示名称、邮箱" style="width: 260px" @keyup.enter="loadUsers" />
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 140px">
          <el-option label="启用" value="active" />
          <el-option label="禁用" value="disabled" />
        </el-select>
        <el-select v-model="query.role" clearable placeholder="角色" style="width: 160px">
          <el-option label="超级管理员" value="super_admin" />
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
        </el-select>
        <el-button :icon="Search" @click="loadUsers">查询</el-button>
      </div>
    </div>

    <el-table :data="users" v-loading="tableLoading" height="560">
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="display_name" label="显示名称" width="160" />
      <el-table-column prop="email" label="邮箱" min-width="200">
        <template #default="{ row }">{{ row.email || '未设置' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <StatusPill :value="row.status" />
        </template>
      </el-table-column>
      <el-table-column label="角色" width="170">
        <template #default="{ row }">
          <span class="tag-list">
            <el-tag v-for="role in row.roles" :key="role" size="small">{{ roleLabel(role) }}</el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="权限" min-width="360">
        <template #default="{ row }">
          <span class="permission-list">
            <el-tag v-for="permission in row.permissions" :key="permission" size="small" type="info">
              {{ permissionLabel(permission) }}
            </el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :icon="Pencil" @click="openEdit(row)">修改</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingUser ? '修改账号' : '新增账号'" width="720px">
      <el-form label-position="top" @submit.prevent="submit">
        <div class="form-grid">
          <el-form-item label="用户名">
            <el-input v-model="form.username" autocomplete="new-username" />
          </el-form-item>
          <el-form-item label="显示名称">
            <el-input v-model="form.display_name" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="form.email" autocomplete="email" />
          </el-form-item>
          <el-form-item :label="editingUser ? '新密码' : '密码'">
            <el-input v-model="form.password" type="password" autocomplete="new-password" show-password />
          </el-form-item>
          <el-form-item label="角色">
            <el-select v-model="form.role_code" style="width: 100%">
              <el-option label="超级管理员" value="super_admin" />
              <el-option label="管理员" value="admin" />
              <el-option label="普通用户" value="user" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="form.status" style="width: 100%">
              <el-option label="启用" value="active" />
              <el-option label="禁用" value="disabled" />
            </el-select>
          </el-form-item>
          <el-form-item label="用户权限" class="full-row">
            <el-select v-model="form.permission_codes" multiple filterable collapse-tags collapse-tags-tooltip style="width: 100%">
              <el-option
                v-for="permission in permissions"
                :key="permission.code"
                :label="permissionLabel(permission.code)"
                :value="permission.code"
              />
            </el-select>
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :icon="editingUser ? Save : UserPlus" :loading="loading" @click="submit">
          {{ editingUser ? '保存修改' : '创建账号' }}
        </el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Pencil, RefreshCw, Save, Search, UserPlus } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { listPermissions, listUsers, registerUser, updateUser } from '../services/auth';
import type { AuthUser, Permission, RegisterPayload, UpdateUserPayload } from '../types/api';

type AccountForm = RegisterPayload & Pick<UpdateUserPayload, 'status' | 'permission_codes'>;

const loading = ref(false);
const tableLoading = ref(false);
const dialogVisible = ref(false);
const users = ref<AuthUser[]>([]);
const permissions = ref<Permission[]>([]);
const editingUser = ref<AuthUser | null>(null);
const query = reactive({ keyword: '', status: '', role: '' });
const form = reactive<AccountForm>({
  username: '',
  password: '',
  display_name: '',
  email: '',
  role_code: 'user',
  status: 'active',
  permission_codes: []
});

async function loadUsers() {
  tableLoading.value = true;
  try {
    users.value = await listUsers({
      keyword: query.keyword,
      status: query.status,
      role: query.role
    });
  } finally {
    tableLoading.value = false;
  }
}

async function loadPermissions() {
  permissions.value = await listPermissions();
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(user: AuthUser) {
  editingUser.value = user;
  form.username = user.username;
  form.password = '';
  form.display_name = user.display_name;
  form.email = user.email;
  form.role_code = (user.roles[0] || 'user') as AccountForm['role_code'];
  form.status = (user.status || 'active') as AccountForm['status'];
  form.permission_codes = [...user.permissions];
  dialogVisible.value = true;
}

function resetForm() {
  editingUser.value = null;
  form.username = '';
  form.password = '';
  form.display_name = '';
  form.email = '';
  form.role_code = 'user';
  form.status = 'active';
  form.permission_codes = [];
}

async function submit() {
  if (!form.username || !form.role_code) {
    ElMessage.warning('请完整填写账号信息');
    return;
  }
  if (!editingUser.value && !form.password) {
    ElMessage.warning('创建账号时密码不能为空');
    return;
  }
  loading.value = true;
  try {
    if (editingUser.value) {
      await updateUser(editingUser.value.id, {
        username: form.username,
        password: form.password || undefined,
        display_name: form.display_name,
        email: form.email,
        role_code: form.role_code,
        status: form.status,
        permission_codes: form.permission_codes
      });
      ElMessage.success('账号已更新');
    } else {
      await registerUser(form);
      ElMessage.success('账号创建成功');
    }
    dialogVisible.value = false;
    resetForm();
    await loadUsers();
  } finally {
    loading.value = false;
  }
}

function roleLabel(role: string) {
  const labels: Record<string, string> = {
    super_admin: '超级管理员',
    admin: '管理员',
    user: '普通用户'
  };
  return labels[role] || role;
}

function permissionLabel(permission: string) {
  const labels: Record<string, string> = {
    'asset:read': '查看资产',
    'asset:create': '新增资产',
    'asset:update': '编辑资产',
    'asset:delete': '删除资产',
    'user:create': '账号管理',
    'workflow:config': '配置审批',
    'workflow:start': '发起审批',
    'workflow:approve': '处理审批'
  };
  return labels[permission] || permission;
}

onMounted(async () => {
  await Promise.all([loadUsers(), loadPermissions()]);
});
</script>
