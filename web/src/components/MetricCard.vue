<template>
  <section class="metric-card">
    <div class="metric-icon">
      <component :is="icon" :size="20" />
    </div>
    <div class="metric-content">
      <div class="metric-label">{{ label }}</div>
      <div class="metric-value">{{ value }}</div>
      <div class="metric-hint">{{ hint }}</div>
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  icon: object;
  label: string;
  value: string | number;
  hint: string;
}>();
</script>
