<template>
  <div class="split-grid">
    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">创建交叉验证</div>
      </div>
      <el-form label-position="top">
        <el-form-item label="资产 ID">
          <el-input-number v-model="assetID" :min="1" />
        </el-form-item>
        <el-button type="primary" :icon="Workflow" @click="create">开始验证</el-button>
      </el-form>

      <el-divider />

      <div class="panel-header">
        <div class="panel-title">查询验证任务</div>
      </div>
      <el-form label-position="top">
        <el-form-item label="验证任务 ID">
          <el-input-number v-model="verificationID" :min="1" />
        </el-form-item>
        <el-button :icon="Search" @click="lookup">查询任务</el-button>
      </el-form>
    </section>

    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">验证结果</div>
      </div>
      <template v-if="result">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务号" :span="2"><span class="mono">{{ result.task.task_no }}</span></el-descriptions-item>
          <el-descriptions-item label="资产 ID">{{ result.task.asset_id }}</el-descriptions-item>
          <el-descriptions-item label="任务状态"><StatusPill :value="result.task.status" /></el-descriptions-item>
          <el-descriptions-item label="结果"><StatusPill :value="result.task.result" /></el-descriptions-item>
          <el-descriptions-item label="分数">{{ result.task.score }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />
        <div class="panel-header">
          <div class="panel-title">冲突字段</div>
        </div>
        <el-table :data="result.conflicts" height="310">
          <el-table-column prop="field" label="字段" width="140" />
          <el-table-column prop="expected" label="期望" />
          <el-table-column prop="actual" label="实际" />
          <el-table-column label="等级" width="110">
            <template #default="{ row }">
              <StatusPill :value="row.severity === 'high' ? 'failed' : 'warning'" />
            </template>
          </el-table-column>
        </el-table>
      </template>
      <p v-else class="empty-hint">选择资产并开始验证后，这里会显示评分和冲突字段。</p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Search, Workflow } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { createVerification, getVerification } from '../services/api';
import type { VerificationResult } from '../types/api';

const assetID = ref(1);
const verificationID = ref(1);
const result = ref<VerificationResult | null>(null);

async function create() {
  result.value = await createVerification(assetID.value);
  verificationID.value = result.value.task.id;
}

async function lookup() {
  result.value = await getVerification(verificationID.value);
}
</script>
