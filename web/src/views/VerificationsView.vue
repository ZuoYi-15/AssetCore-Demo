<template>
  <div class="content-grid">
    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">选择资产发起交叉认证</div>
          <div class="empty-hint">按资产名称、序列号、MAC、IP 或身份 ID 搜索资产后发起认证。</div>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="loadAssets">刷新资产</el-button>
        </div>
      </div>

      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="assetQuery.keyword" clearable placeholder="搜索资产名称、序列号、MAC、IP、身份 ID" style="width: 320px" @keyup.enter="loadAssets" />
          <el-select v-model="assetQuery.status" clearable placeholder="资产状态" style="width: 140px" @change="loadAssets">
            <el-option label="已登记" value="registered" />
            <el-option label="已验证" value="verified" />
            <el-option label="异常" value="abnormal" />
            <el-option label="已退役" value="retired" />
          </el-select>
          <el-button :icon="Search" @click="loadAssets">查询</el-button>
        </div>
      </div>

      <el-table :data="assets" v-loading="assetLoading" height="360" highlight-current-row>
        <el-table-column prop="asset_name" label="资产名称" min-width="170" />
        <el-table-column prop="serial_number" label="序列号" min-width="150" />
        <el-table-column label="身份 ID" min-width="170">
          <template #default="{ row }">
            <span v-if="row.identity_id" class="mono">{{ shortText(row.identity_id) }}</span>
            <el-tag v-else type="info" effect="plain">未绑定</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="可信等级" width="110">
          <template #default="{ row }"><StatusPill :value="row.trust_level" /></template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" :icon="Workflow" :loading="verifyingAssetID === row.id" @click="create(row)">开始认证</el-button>
          </template>
        </el-table-column>
      </el-table>
    </section>

    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">认证结果</div>
      </div>
      <template v-if="result">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务号" :span="2"><span class="mono">{{ result.task.task_no }}</span></el-descriptions-item>
          <el-descriptions-item label="任务状态"><StatusPill :value="result.task.status" /></el-descriptions-item>
          <el-descriptions-item label="结果"><StatusPill :value="result.task.result" /></el-descriptions-item>
          <el-descriptions-item label="分数">{{ result.task.score }}</el-descriptions-item>
          <el-descriptions-item label="冲突数">{{ result.conflicts.length }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />
        <div class="panel-header">
          <div class="panel-title">冲突字段</div>
        </div>
        <el-table :data="result.conflicts" height="260">
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
      <p v-else class="empty-hint">选择资产并开始认证后，这里会显示评分和冲突字段。</p>
    </section>

    <section class="panel full-row">
      <div class="panel-header">
        <div>
          <div class="panel-title">认证记录</div>
          <div class="empty-hint">按认证结果筛选历史记录，点击详情查看完整冲突信息。</div>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="loadRecords">刷新记录</el-button>
        </div>
      </div>

      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="recordQuery.keyword" clearable placeholder="搜索资产名称、序列号、身份 ID、任务号" style="width: 320px" @keyup.enter="loadRecords" />
          <el-select v-model="recordQuery.result" clearable placeholder="筛选认证结果" style="width: 150px" @change="loadRecords">
            <el-option label="全部结果" value="" />
            <el-option label="通过" value="passed" />
            <el-option label="预警" value="warning" />
            <el-option label="失败" value="failed" />
          </el-select>
          <el-button :icon="Search" @click="loadRecords">查询</el-button>
        </div>
      </div>

      <el-table :data="records" v-loading="recordLoading" height="320">
        <el-table-column prop="asset_name" label="资产名称" min-width="170" />
        <el-table-column prop="serial_number" label="序列号" min-width="150" />
        <el-table-column prop="owner_department" label="部门" width="130" />
        <el-table-column label="身份 ID" min-width="180">
          <template #default="{ row }">
            <span v-if="row.identity_id" class="mono">{{ shortText(row.identity_id) }}</span>
            <el-tag v-else type="info" effect="plain">未绑定</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="结果" width="110">
          <template #default="{ row }"><StatusPill :value="row.result" /></template>
        </el-table-column>
        <el-table-column prop="score" label="分数" width="90" />
        <el-table-column prop="created_at" label="认证时间" width="190" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Eye" @click="openRecord(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div style="display:flex; justify-content:flex-end; margin-top:14px">
        <el-pagination
          v-model:current-page="recordPage"
          v-model:page-size="recordPageSize"
          layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 50]"
          :total="recordTotal"
          @change="loadRecords"
        />
      </div>
    </section>

    <el-dialog v-model="recordDetailVisible" title="认证记录详情" width="760px">
      <template v-if="recordResult">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务号" :span="2"><span class="mono">{{ recordResult.task.task_no }}</span></el-descriptions-item>
          <el-descriptions-item label="任务状态"><StatusPill :value="recordResult.task.status" /></el-descriptions-item>
          <el-descriptions-item label="结果"><StatusPill :value="recordResult.task.result" /></el-descriptions-item>
          <el-descriptions-item label="分数">{{ recordResult.task.score }}</el-descriptions-item>
          <el-descriptions-item label="冲突数">{{ recordResult.conflicts.length }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />
        <div class="panel-header">
          <div class="panel-title">冲突字段</div>
        </div>
        <el-table :data="recordResult.conflicts" height="320">
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
      <template #footer>
        <el-button @click="recordDetailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Eye, RefreshCw, Search, Workflow } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { createVerification, getVerification, listAssets, listVerifications } from '../services/api';
import type { Asset, VerificationRecord, VerificationResult } from '../types/api';

const assets = ref<Asset[]>([]);
const records = ref<VerificationRecord[]>([]);
const result = ref<VerificationResult | null>(null);
const recordResult = ref<VerificationResult | null>(null);
const assetLoading = ref(false);
const recordLoading = ref(false);
const recordDetailVisible = ref(false);
const verifyingAssetID = ref<number | null>(null);
const assetQuery = reactive({ keyword: '', status: '' });
const recordQuery = reactive({ keyword: '', result: '' });
const recordPage = ref(1);
const recordPageSize = ref(10);
const recordTotal = ref(0);

async function loadAssets() {
  assetLoading.value = true;
  try {
    const res = await listAssets({ page: 1, page_size: 20, keyword: assetQuery.keyword, status: assetQuery.status, asset_type: '' });
    assets.value = res.items;
  } finally {
    assetLoading.value = false;
  }
}

async function loadRecords() {
  recordLoading.value = true;
  try {
    const res = await listVerifications({
      page: recordPage.value,
      page_size: recordPageSize.value,
      keyword: recordQuery.keyword,
      result: recordQuery.result,
      status: ''
    });
    records.value = res.items;
    recordTotal.value = res.total;
  } finally {
    recordLoading.value = false;
  }
}

async function create(asset: Asset) {
  verifyingAssetID.value = asset.id;
  try {
    result.value = await createVerification(asset.id);
    ElMessage.success(`已完成 ${asset.asset_name} 的交叉认证`);
    await loadRecords();
  } finally {
    verifyingAssetID.value = null;
  }
}

async function openRecord(record: VerificationRecord) {
  recordResult.value = await getVerification(record.id);
  recordDetailVisible.value = true;
}

function shortText(value: string) {
  if (value.length <= 24) {
    return value;
  }
  return `${value.slice(0, 15)}...${value.slice(-6)}`;
}

onMounted(async () => {
  await Promise.all([loadAssets(), loadRecords()]);
});
</script>
