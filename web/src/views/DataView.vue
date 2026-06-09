<template>
  <div class="content-grid">
    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">导入任务</div>
          <div class="empty-hint">查看资产批量导入任务和错误明细。</div>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        </div>
      </div>

      <el-table :data="tasks" v-loading="loading" height="500">
        <el-table-column prop="task_no" label="任务号" min-width="260">
          <template #default="{ row }">
            <span class="mono">{{ row.task_no }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="file_name" label="文件名" width="190" />
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
            <el-button v-if="row.failed_count > 0" size="small" :icon="CircleAlert" @click="openErrors(row.id)">错误</el-button>
            <span v-else class="empty-hint">无错误</span>
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
import { onMounted, ref } from 'vue';
import { CircleAlert, RefreshCw } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { listImportErrors, listImportTasks } from '../services/api';
import type { ImportError, ImportTask } from '../types/api';

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

async function openErrors(id: number) {
  errors.value = await listImportErrors(id);
  errorVisible.value = true;
}

onMounted(load);
</script>
