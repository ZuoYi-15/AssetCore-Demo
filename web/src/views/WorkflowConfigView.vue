<template>
  <div class="content-grid">
    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">审批流程配置</div>
          <div class="empty-hint">超级管理员可配置采购、调拨、报废流程的审批节点和审批角色。</div>
        </div>
        <el-button :icon="RefreshCw" @click="load">刷新</el-button>
      </div>

      <el-tabs v-model="activeFlow">
        <el-tab-pane v-for="flow in flowTypes" :key="flow.value" :label="flow.label" :name="flow.value">
          <el-form label-position="top" class="workflow-form">
            <div class="form-grid">
              <el-form-item label="流程名称">
                <el-input v-model="draft.name" />
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="draft.status" style="width: 100%">
                  <el-option label="启用" value="active" />
                  <el-option label="停用" value="inactive" />
                </el-select>
              </el-form-item>
            </div>
          </el-form>

          <div class="panel-header compact-header">
            <div class="panel-title">审批节点</div>
            <el-button :icon="Plus" @click="addNode">新增节点</el-button>
          </div>

          <el-table :data="draft.nodes" border>
            <el-table-column label="顺序" width="90">
              <template #default="{ row, $index }">
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

          <div class="form-actions">
            <el-button type="primary" :icon="Save" :loading="saving" @click="save">保存配置</el-button>
          </div>
        </el-tab-pane>
      </el-tabs>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { Plus, RefreshCw, Save, Trash2 } from 'lucide-vue-next';
import { listWorkflowDefinitions, saveWorkflowDefinition } from '../services/api';
import type { WorkflowDefinition, WorkflowNode, WorkflowType } from '../types/api';

const flowTypes: Array<{ value: WorkflowType; label: string }> = [
  { value: 'purchase', label: '采购审批' },
  { value: 'transfer', label: '调拨审批' },
  { value: 'retire', label: '报废审批' }
];

const definitions = ref<WorkflowDefinition[]>([]);
const activeFlow = ref<WorkflowType>('purchase');
const saving = ref(false);
const draft = reactive<{
  name: string;
  status: 'active' | 'inactive';
  nodes: Array<Pick<WorkflowNode, 'node_name' | 'approver_role' | 'sort_order'>>;
}>({
  name: '',
  status: 'active',
  nodes: []
});

const currentDefinition = computed(() => definitions.value.find((item) => item.flow_type === activeFlow.value));

async function load() {
  definitions.value = await listWorkflowDefinitions();
  syncDraft();
}

function syncDraft() {
  const item = currentDefinition.value;
  const fallback = flowTypes.find((flow) => flow.value === activeFlow.value);
  draft.name = item?.name || fallback?.label || '';
  draft.status = item?.status || 'active';
  draft.nodes = (item?.nodes || [{ node_name: '管理员审批', approver_role: 'admin', sort_order: 1 }]).map((node, index) => ({
    node_name: node.node_name,
    approver_role: node.approver_role,
    sort_order: node.sort_order || index + 1
  }));
  sortNodes();
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
    await saveWorkflowDefinition({
      flow_type: activeFlow.value,
      name: draft.name,
      status: draft.status,
      nodes: draft.nodes
    });
    ElMessage.success('流程配置已保存');
    await load();
  } finally {
    saving.value = false;
  }
}

watch(activeFlow, syncDraft);
onMounted(load);
</script>
