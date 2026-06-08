<template>
  <div class="content-grid">
    <div class="metric-grid">
      <MetricCard :icon="Database" label="资产总数" :value="summary.total" hint="来自资产台账的当前记录" />
      <MetricCard :icon="ShieldCheck" label="已验证资产" :value="summary.verified" hint="状态为 verified 的设备" />
      <MetricCard :icon="AlertTriangle" label="异常资产" :value="summary.abnormal" hint="需人工复核或隔离" />
      <MetricCard :icon="Fingerprint" label="可信身份覆盖" :value="identityCoverage" hint="拥有身份 ID 的资产比例" />
    </div>

    <div class="split-grid">
      <section class="panel">
        <div class="panel-header">
          <div class="panel-title">最近资产</div>
          <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        </div>
        <el-table :data="assets" height="360" v-loading="loading">
          <el-table-column prop="asset_name" label="资产名称" min-width="170" />
          <el-table-column prop="asset_type" label="类型" width="120" />
          <el-table-column prop="ip_address" label="IP" width="140" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <StatusPill :value="row.status" />
            </template>
          </el-table-column>
          <el-table-column prop="updated_at" label="更新时间" min-width="190" />
        </el-table>
      </section>

      <section class="panel">
        <div class="panel-header">
          <div class="panel-title">平台态势</div>
        </div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="API 地址">http://127.0.0.1:8080</el-descriptions-item>
          <el-descriptions-item label="身份策略">序列号、厂商、型号、MAC 指纹</el-descriptions-item>
          <el-descriptions-item label="验证策略">字段完整性和身份绑定检查</el-descriptions-item>
          <el-descriptions-item label="Kafka">默认关闭，事件发布已预留</el-descriptions-item>
        </el-descriptions>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { AlertTriangle, Database, Fingerprint, RefreshCw, ShieldCheck } from 'lucide-vue-next';
import MetricCard from '../components/MetricCard.vue';
import StatusPill from '../components/StatusPill.vue';
import { listAssets } from '../services/api';
import type { Asset } from '../types/api';

const assets = ref<Asset[]>([]);
const total = ref(0);
const loading = ref(false);

const summary = computed(() => ({
  total: total.value,
  verified: assets.value.filter((item) => item.status === 'verified').length,
  abnormal: assets.value.filter((item) => item.status === 'abnormal').length
}));

const identityCoverage = computed(() => {
  if (!assets.value.length) return '0%';
  const count = assets.value.filter((item) => item.identity_id).length;
  return `${Math.round((count / assets.value.length) * 100)}%`;
});

async function load() {
  loading.value = true;
  try {
    const res = await listAssets({ page: 1, page_size: 8 });
    assets.value = res.items;
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>
