<template>
  <div class="content-grid">
    <div class="metric-grid">
      <MetricCard :icon="Database" label="资产总数" :value="summary.total" hint="资产台账当前记录总量" />
      <MetricCard :icon="ShieldCheck" label="已验证资产" :value="summary.verified" hint="状态为 verified 的资产" />
      <MetricCard :icon="AlertTriangle" label="异常资产" :value="summary.abnormal" hint="需要人工复核或处置" />
      <MetricCard :icon="Fingerprint" label="身份覆盖率" :value="identityCoverage" hint="已绑定身份 ID 的资产占比" />
    </div>

    <div class="dashboard-grid">
      <section class="panel chart-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">资产状态分布</div>
            <div class="empty-hint">按登记、验证、异常、退役状态统计</div>
          </div>
          <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        </div>
        <div class="donut-layout">
          <div class="donut-chart" :style="{ background: statusDonut }">
            <div class="donut-center">
              <strong>{{ sampleTotal }}</strong>
              <span>统计样本</span>
            </div>
          </div>
          <div class="chart-legend">
            <div v-for="item in statusStats" :key="item.key" class="legend-row">
              <span class="legend-dot" :style="{ background: item.color }" />
              <span>{{ item.label }}</span>
              <strong>{{ item.value }}</strong>
            </div>
          </div>
        </div>
      </section>

      <section class="panel chart-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">资产类型排行</div>
            <div class="empty-hint">按资产类型聚合，取前 6 类</div>
          </div>
        </div>
        <div class="bar-list">
          <div v-for="item in typeStats" :key="item.label" class="bar-row">
            <div class="bar-meta">
              <span>{{ item.label }}</span>
              <strong>{{ item.value }}</strong>
            </div>
            <div class="bar-track">
              <div class="bar-fill" :style="{ width: `${item.percent}%` }" />
            </div>
          </div>
          <div v-if="!typeStats.length" class="empty-hint">暂无资产类型数据</div>
        </div>
      </section>

      <section class="panel chart-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">可信等级分布</div>
            <div class="empty-hint">展示身份可信等级覆盖情况</div>
          </div>
        </div>
        <div class="trust-grid">
          <div v-for="item in trustStats" :key="item.key" class="trust-item">
            <div class="trust-value">{{ item.value }}</div>
            <StatusPill :value="item.key" />
            <div class="trust-percent">{{ item.percent }}%</div>
          </div>
        </div>
      </section>

      <section class="panel chart-panel">
        <div class="panel-header">
          <div>
            <div class="panel-title">健康度雷达</div>
            <div class="empty-hint">综合身份覆盖、验证率、异常控制、活跃资产</div>
          </div>
        </div>
        <div class="radar-wrap">
          <svg viewBox="0 0 220 180" class="radar-chart" role="img" aria-label="资产健康度雷达图">
            <polygon points="110,18 190,68 160,154 60,154 30,68" class="radar-grid" />
            <polygon points="110,42 166,78 145,132 75,132 54,78" class="radar-grid inner" />
            <polygon :points="radarPoints" class="radar-area" />
            <g class="radar-labels">
              <text x="110" y="12" text-anchor="middle">身份</text>
              <text x="198" y="70">验证</text>
              <text x="165" y="172">异常</text>
              <text x="42" y="172">活跃</text>
              <text x="2" y="70">完整</text>
            </g>
          </svg>
        </div>
      </section>
    </div>

    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">最近资产</div>
          <div class="empty-hint">最近更新的资产记录</div>
        </div>
      </div>
      <el-table :data="recentAssets" height="360" v-loading="loading">
        <el-table-column prop="asset_name" label="资产名称" min-width="170" />
        <el-table-column prop="asset_type" label="类型" width="120" />
        <el-table-column prop="ip_address" label="IP" width="140" />
        <el-table-column prop="owner_department" label="部门" width="140" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <StatusPill :value="row.status" />
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" min-width="190" />
      </el-table>
    </section>
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

const statusConfig = [
  { key: 'registered', label: '已登记', color: '#175cd3' },
  { key: 'verified', label: '已验证', color: '#087443' },
  { key: 'abnormal', label: '异常', color: '#c01048' },
  { key: 'retired', label: '已退役', color: '#667085' },
  { key: 'discovered', label: '已发现', color: '#b54708' }
];

const summary = computed(() => ({
  total: total.value,
  verified: countBy('status', 'verified'),
  abnormal: countBy('status', 'abnormal')
}));

const sampleTotal = computed(() => assets.value.length);
const recentAssets = computed(() => assets.value.slice(0, 10));

const identityCoverage = computed(() => `${ratio(assets.value.filter((item) => item.identity_id).length, sampleTotal.value)}%`);

const statusStats = computed(() =>
  statusConfig.map((item) => ({
    ...item,
    value: countBy('status', item.key)
  }))
);

const statusDonut = computed(() => {
  if (!sampleTotal.value) return 'conic-gradient(#e5e7eb 0deg 360deg)';
  let cursor = 0;
  const parts = statusStats.value
    .filter((item) => item.value > 0)
    .map((item) => {
      const degrees = (item.value / sampleTotal.value) * 360;
      const part = `${item.color} ${cursor}deg ${cursor + degrees}deg`;
      cursor += degrees;
      return part;
    });
  return `conic-gradient(${parts.join(', ')})`;
});

const typeStats = computed(() => {
  const counts = groupCount(assets.value.map((item) => item.asset_type || '未分类'));
  const max = Math.max(1, ...counts.map((item) => item.value));
  return counts.slice(0, 6).map((item) => ({
    ...item,
    percent: Math.round((item.value / max) * 100)
  }));
});

const trustStats = computed(() => {
  const keys = ['strong', 'medium', 'weak', ''];
  return keys.map((key) => {
    const value = key ? countBy('trust_level', key) : assets.value.filter((item) => !item.trust_level).length;
    return {
      key: key || 'unknown',
      value,
      percent: ratio(value, sampleTotal.value)
    };
  });
});

const radarPoints = computed(() => {
  const scores = [
    ratio(assets.value.filter((item) => item.identity_id).length, sampleTotal.value),
    ratio(countBy('status', 'verified'), sampleTotal.value),
    100 - ratio(countBy('status', 'abnormal'), sampleTotal.value),
    100 - ratio(countBy('status', 'retired'), sampleTotal.value),
    ratio(assets.value.filter((item) => item.asset_name && item.serial_number).length, sampleTotal.value)
  ];
  const center = { x: 110, y: 100 };
  const radius = 76;
  return scores
    .map((score, index) => {
      const angle = (-90 + index * 72) * (Math.PI / 180);
      const length = radius * (score / 100);
      return `${center.x + Math.cos(angle) * length},${center.y + Math.sin(angle) * length}`;
    })
    .join(' ');
});

function countBy(field: keyof Asset, value: string) {
  return assets.value.filter((item) => item[field] === value).length;
}

function ratio(value: number, base: number) {
  if (!base) return 0;
  return Math.round((value / base) * 100);
}

function groupCount(values: string[]) {
  const map = new Map<string, number>();
  values.forEach((value) => map.set(value, (map.get(value) || 0) + 1));
  return Array.from(map.entries())
    .map(([label, value]) => ({ label, value }))
    .sort((a, b) => b.value - a.value);
}

async function load() {
  loading.value = true;
  try {
    const res = await listAssets({ page: 1, page_size: 200 });
    assets.value = res.items;
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>
