<template>
  <div class="content-grid">
    <section class="panel">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-select v-model="status" clearable placeholder="筛选审批状态" style="width: 160px" @change="load">
            <el-option label="全部状态" value="" />
            <el-option label="待审批" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        </div>
      </div>

      <el-table :data="tasks" v-loading="loading" height="560">
        <el-table-column label="流程" width="110">
          <template #default="{ row }">{{ flowTypeLabel(row.instance?.flow_type) }}</template>
        </el-table-column>
        <el-table-column label="标题" min-width="220">
          <template #default="{ row }">{{ row.instance?.title }}</template>
        </el-table-column>
        <el-table-column prop="node_name" label="当前节点" width="150" />
        <el-table-column prop="approver_role" label="审批角色" width="130" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <StatusPill :value="row.status" />
          </template>
        </el-table-column>
        <el-table-column label="申请人" width="130">
          <template #default="{ row }">{{ row.instance?.applicant_name }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="190" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <span class="table-actions">
              <el-button size="small" :icon="Eye" @click="openTask(row)">详情</el-button>
              <el-button v-if="row.status === 'pending'" size="small" type="primary" :icon="Check" @click="openTask(row)">处理</el-button>
            </span>
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

    <el-dialog v-model="detailVisible" title="审批处理" width="620px">
      <div v-if="currentTask" class="workflow-detail">
        <div><strong>流程：</strong>{{ flowTypeLabel(currentTask.instance?.flow_type) }}</div>
        <div><strong>标题：</strong>{{ currentTask.instance?.title }}</div>
        <div><strong>节点：</strong>{{ currentTask.node_name }}</div>
        <div><strong>申请人：</strong>{{ currentTask.instance?.applicant_name }}</div>
        <div><strong>申请内容：</strong>{{ payloadText(currentTask.instance?.payload) }}</div>
      </div>
      <el-form v-if="currentTask?.status === 'pending'" label-position="top" style="margin-top: 16px">
        <el-form-item label="审批意见">
          <el-input v-model="comment" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button v-if="currentTask?.status === 'pending'" type="danger" :loading="submitting" @click="submit('reject')">驳回</el-button>
        <el-button v-if="currentTask?.status === 'pending'" type="primary" :loading="submitting" @click="submit('approve')">通过</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Check, Eye, RefreshCw } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { approveWorkflowTask, listWorkflowTasks } from '../services/api';
import type { WorkflowTask, WorkflowType } from '../types/api';

const tasks = ref<WorkflowTask[]>([]);
const loading = ref(false);
const submitting = ref(false);
const detailVisible = ref(false);
const currentTask = ref<WorkflowTask | null>(null);
const comment = ref('');
const status = ref('');
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

async function load() {
  loading.value = true;
  try {
    const res = await listWorkflowTasks({ page: page.value, page_size: pageSize.value, status: status.value });
    tasks.value = res.items;
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

function openTask(task: WorkflowTask) {
  currentTask.value = task;
  comment.value = '';
  detailVisible.value = true;
}

async function submit(action: 'approve' | 'reject') {
  if (!currentTask.value) return;
  submitting.value = true;
  try {
    await approveWorkflowTask(currentTask.value.id, action, comment.value);
    ElMessage.success(action === 'approve' ? '审批已通过' : '审批已驳回');
    detailVisible.value = false;
    await load();
  } finally {
    submitting.value = false;
  }
}

function flowTypeLabel(flowType?: WorkflowType) {
  const labels: Record<WorkflowType, string> = {
    purchase: '采购',
    transfer: '调拨',
    retire: '报废'
  };
  return flowType ? labels[flowType] : '';
}

function payloadText(payload?: string) {
  if (!payload) return '';
  try {
    const data = JSON.parse(payload) as Record<string, unknown>;
    return Object.entries(data)
      .map(([key, value]) => `${key}: ${String(value ?? '')}`)
      .join('；');
  } catch {
    return payload;
  }
}

onMounted(load);
</script>
