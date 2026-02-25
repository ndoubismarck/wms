<script setup lang="ts">
import { onBeforeUnmount, provide, ref } from 'vue'

const items = ref<number[]>([])
const currentIndex = ref(0)

function register(uid: number) {
  if (!items.value.includes(uid)) {
    items.value.push(uid)
  }
  return items.value.indexOf(uid)
}

function openByIndex(index: number) {
  if (index >= 0 && index < items.value.length) {
    currentIndex.value = index
  }
}

function next() {
  openByIndex(currentIndex.value + 1)
}

function prev() {
  openByIndex(currentIndex.value - 1)
}

const hasNext = (): boolean => {
  const index = currentIndex.value + 1
  return index < items.value.length
}

const hasPrev = (): boolean => {
  const index = currentIndex.value - 1
  return index > 0
}

provide('accordion', {
  next,
  prev,
  hasNext,
  hasPrev,
  register,
  currentIndex,
})

defineExpose({
  next,
  prev,
  hasNext,
  hasPrev,
})

onBeforeUnmount(() => {
  items.value = []
  currentIndex.value = 0
})
</script>

<template>
  <div class="accordion">
    <slot :next="next" :prev="prev" :hasNext="hasNext" :hasPrev="hasPrev" />
  </div>
</template>

<style scoped></style>
