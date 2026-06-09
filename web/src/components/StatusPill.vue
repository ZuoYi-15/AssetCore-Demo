<template>
  <span class="status-pill" :class="`status-${normalized}`">
    <span class="status-dot" />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ value?: string }>();

const normalized = computed(() => props.value || 'unknown');
const label = computed(() => {
  const map: Record<string, string> = {
    discovered: '已发现',
    registered: '已登记',
    verified: '已验证',
    abnormal: '异常',
    retired: '已退役',
    active: '启用',
    inactive: '禁用',
    pending: '待审批',
    approved: '已通过',
    rejected: '已驳回',
    completed: '已完成',
    failed: '失败',
    passed: '通过',
    warning: '警告',
    strong: '强身份',
    medium: '中身份',
    weak: '弱身份',
    running: '运行中'
  };
  return map[normalized.value] || normalized.value || '未知';
});
</script>
