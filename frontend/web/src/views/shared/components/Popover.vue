<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { Popover } from 'bootstrap'

const popoverBtn = ref<HTMLElement | null>(null)
let popover: Popover | null = null

onMounted(() => {
  if (popoverBtn.value) {
    popover = new Popover(popoverBtn.value)
  }
})

onBeforeUnmount(() => {
  popover?.dispose()
})
</script>

<template>
  <button
    ref="popoverBtn"
    type="button"
    class="btn btn-secondary"
    data-bs-container="body"
    data-bs-toggle="popover"
    data-bs-placement="bottom"
    data-bs-content="Top popover">
    Popover on top
  </button>
</template>

<style scoped></style>
