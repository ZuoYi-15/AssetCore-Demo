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
    pending: '待处理',
    completed: '已完成',
    failed: '失败',
    passed: '通过',
    warning: '告警',
    strong: '强身份',
    medium: '中身份',
    weak: '弱身份'
  };
  return map[normalized.value] || normalized.value || '未知';
});
</script>
