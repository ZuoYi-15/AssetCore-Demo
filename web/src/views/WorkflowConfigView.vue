<template>
  <section class="panel">
    <div class="panel-header">
      <div>
        <div class="panel-title">审批流程列表</div>
        <div class="empty-hint">已配置的采购、调拨、报废审批流会在这里展示；禁用后可删除。</div>
      </div>
      <div class="toolbar-right">
        <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        <el-button type="primary" :icon="Plus" @click="openCreate">新建流程</el-button>
      </div>
    </div>

    <el-table :data="definitions" v-loading="loading" border>
      <el-table-column label="流程类型" width="130">
        <template #default="{ row }">{{ flowTypeLabel(row.flow_type) }}</template>
      </el-table-column>
      <el-table-column prop="name" label="流程名称" min-width="180" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <StatusPill :value="row.status" />
        </template>
      </el-table-column>
      <el-table-column label="节点" min-width="260">
        <template #default="{ row }">
          <span class="workflow-node-list">
            <span v-for="node in row.nodes" :key="node.id || node.sort_order">
              {{ node.sort_order }}. {{ node.node_name }} / {{ roleLabel(node.approver_role) }}
            </span>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <span class="table-actions">
            <el-button size="small" :icon="Pencil" @click="openEdit(row)">修改</el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Trash2"
              :disabled="row.status !== 'inactive'"
              @click="removeDefinition(row)"
            >
              删除
            </el-button>
          </span>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingID ? '修改审批流' : '新建审批流'" width="860px">
      <el-form label-position="top" class="workflow-form">
        <div class="form-grid">
          <el-form-item label="流程类型">
            <el-select v-model="draft.flow_type" style="width: 100%">
              <el-option v-for="flow in flowTypes" :key="flow.value" :label="flow.label" :value="flow.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="流程名称">
            <el-input v-model="draft.name" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="draft.status" style="width: 100%">
              <el-option label="启用" value="active" />
              <el-option label="禁用" value="inactive" />
            </el-select>
          </el-form-item>
        </div>
      </el-form>

      <div class="panel-header compact-header">
        <div class="panel-title">审批节点</div>
        <el-button :icon="Plus" @click="addNode">新增节点</el-button>
      </div>

      <el-table :data="draft.nodes" border>
        <el-table-column label="顺序" width="110">
          <template #default="{ row }">
            <el-input-number v-model="row.sort_order" :min="1" :max="99" size="small" @change="sortNodes" />
          </template>
        </el-table-column>
        <el-table-column label="节点名称">
          <template #default="{ row }">
            <el-input v-model="row.node_name" />
          </template>
        </el-table-column>
        <el-table-column label="审批角色" width="220">
          <template #default="{ row }">
            <el-select v-model="row.approver_role" style="width: 100%">
              <el-option label="超级管理员" value="super_admin" />
              <el-option label="管理员" value="admin" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ $index }">
            <el-button type="danger" size="small" :icon="Trash2" @click="removeNode($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :icon="Save" :loading="saving" @click="save">保存配置</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Pencil, Plus, RefreshCw, Save, Trash2 } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { deleteWorkflowDefinition, listWorkflowDefinitions, saveWorkflowDefinition } from '../services/api';
import type { WorkflowDefinition, WorkflowNode, WorkflowType } from '../types/api';

const flowTypes: Array<{ value: WorkflowType; label: string }> = [
  { value: 'purchase', label: '采购审批' },
  { value: 'transfer', label: '调拨审批' },
  { value: 'retire', label: '报废审批' }
];

const definitions = ref<WorkflowDefinition[]>([]);
const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingID = ref<number | null>(null);
const draft = reactive<{
  flow_type: WorkflowType;
  name: string;
  status: 'active' | 'inactive';
  nodes: Array<Pick<WorkflowNode, 'node_name' | 'approver_role' | 'sort_order'>>;
}>({
  flow_type: 'purchase',
  name: '采购审批',
  status: 'active',
  nodes: [{ node_name: '管理员审批', approver_role: 'admin', sort_order: 1 }]
});

async function load() {
  loading.value = true;
  try {
    definitions.value = await listWorkflowDefinitions();
  } finally {
    loading.value = false;
  }
}

function openEdit(item: WorkflowDefinition) {
  editingID.value = item.id;
  draft.flow_type = item.flow_type;
  draft.name = item.name;
  draft.status = item.status;
  draft.nodes = (item.nodes || []).map((node, index) => ({
    node_name: node.node_name,
    approver_role: node.approver_role,
    sort_order: node.sort_order || index + 1
  }));
  if (draft.nodes.length === 0) {
    draft.nodes = [{ node_name: '管理员审批', approver_role: 'admin', sort_order: 1 }];
  }
  sortNodes();
  dialogVisible.value = true;
}

function openCreate() {
  editingID.value = null;
  const firstMissing = flowTypes.find((flow) => !definitions.value.some((item) => item.flow_type === flow.value)) || flowTypes[0];
  draft.flow_type = firstMissing.value;
  draft.name = firstMissing.label;
  draft.status = 'active';
  draft.nodes = [{ node_name: '管理员审批', approver_role: 'admin', sort_order: 1 }];
  dialogVisible.value = true;
}

function addNode() {
  draft.nodes.push({ node_name: `审批节点 ${draft.nodes.length + 1}`, approver_role: 'admin', sort_order: draft.nodes.length + 1 });
}

function removeNode(index: number) {
  draft.nodes.splice(index, 1);
  sortNodes();
}

function sortNodes() {
  draft.nodes.sort((a, b) => a.sort_order - b.sort_order);
  draft.nodes.forEach((node, index) => {
    node.sort_order = index + 1;
  });
}

async function save() {
  if (!draft.name || draft.nodes.length === 0) {
    ElMessage.warning('请填写流程名称并至少配置一个节点');
    return;
  }
  saving.value = true;
  try {
    const saved = await saveWorkflowDefinition({
      flow_type: draft.flow_type,
      name: draft.name,
      status: draft.status,
      nodes: draft.nodes
    });
    ElMessage.success('流程配置已保存');
    await load();
    editingID.value = saved.id;
    dialogVisible.value = false;
  } finally {
    saving.value = false;
  }
}

async function removeDefinition(item: WorkflowDefinition) {
  if (item.status !== 'inactive') {
    ElMessage.warning('只有已禁用的审批流可以删除');
    return;
  }
  await ElMessageBox.confirm(`确认删除审批流 ${item.name}？`, '删除确认', { type: 'warning' });
  await deleteWorkflowDefinition(item.id);
  ElMessage.success('审批流已删除');
  editingID.value = null;
  await load();
}

function flowTypeLabel(flowType: WorkflowType) {
  return flowTypes.find((item) => item.value === flowType)?.label || flowType;
}

function roleLabel(role: string) {
  const labels: Record<string, string> = {
    super_admin: '超级管理员',
    admin: '管理员'
  };
  return labels[role] || role;
}

onMounted(load);
</script>
