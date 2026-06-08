<template>
  <div class="content-grid">
    <section class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="form.file_name" placeholder="文件名，例如 assets.csv" style="width: 240px" />
          <el-input v-model="form.file_url" placeholder="文件地址，例如 local://assets.csv" style="width: 280px" />
          <el-input v-model="form.operator_id" placeholder="操作人" style="width: 160px" />
          <el-button type="primary" :icon="UploadCloud" @click="create">创建导入任务</el-button>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="load">刷新</el-button>
          <el-button :icon="Download" @click="exportAssets">资产导出</el-button>
        </div>
      </div>

      <el-table :data="tasks" v-loading="loading" height="500">
        <el-table-column prop="task_no" label="任务号" min-width="260">
          <template #default="{ row }">
            <span class="mono">{{ row.task_no }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="file_name" label="文件名" width="180" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <StatusPill :value="row.status" />
          </template>
        </el-table-column>
        <el-table-column prop="total_count" label="总数" width="90" />
        <el-table-column prop="success_count" label="成功" width="90" />
        <el-table-column prop="failed_count" label="失败" width="90" />
        <el-table-column prop="operator_id" label="操作人" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="190" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="CircleAlert" @click="openErrors(row.id)">错误</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div style="display:flex; justify-content:flex-end; margin-top:14px">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 50]"
          :total="total"
          @change="load"
        />
      </div>
    </section>

    <el-drawer v-model="errorVisible" title="导入错误" size="640px">
      <el-table :data="errors">
        <el-table-column prop="row_number" label="行号" width="80" />
        <el-table-column prop="error_field" label="字段" width="120" />
        <el-table-column prop="error_message" label="错误信息" />
        <el-table-column prop="raw_data" label="原始数据" />
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { CircleAlert, Download, RefreshCw, UploadCloud } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { createImportTask, listImportErrors, listImportTasks } from '../services/api';
import type { ImportError, ImportTask } from '../types/api';

const form = reactive({ file_name: 'assets.csv', file_url: 'local://assets.csv', operator_id: 'admin' });
const tasks = ref<ImportTask[]>([]);
const errors = ref<ImportError[]>([]);
const loading = ref(false);
const errorVisible = ref(false);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

async function load() {
  loading.value = true;
  try {
    const res = await listImportTasks({ page: page.value, page_size: pageSize.value });
    tasks.value = res.items;
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function create() {
  if (!form.file_name) {
    ElMessage.warning('文件名不能为空');
    return;
  }
  await createImportTask(form);
  await load();
}

async function openErrors(id: number) {
  errors.value = await listImportErrors(id);
  errorVisible.value = true;
}

function exportAssets() {
  ElMessage.info('后端当前为导出预留接口，真实文件导出可在下一步接入。');
}

onMounted(load);
</script>
