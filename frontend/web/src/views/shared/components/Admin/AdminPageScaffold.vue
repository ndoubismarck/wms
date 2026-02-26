<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title mb-2">{{ title }}</h5>
      <p v-if="description" class="text-body-secondary mb-0">{{ description }}</p>
    </div>
  </div>
</template>

<style scoped></style>
