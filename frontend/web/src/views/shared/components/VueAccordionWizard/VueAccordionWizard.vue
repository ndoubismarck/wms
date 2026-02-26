<script setup lang="ts">
import {onBeforeUnmount, provide, ref} from 'vue'

const items = ref<number[]>([])
const currentIndex = ref(0)
const transitionDirection = ref<'forward' | 'backward'>('forward')

function register(uid: number) {
  if (!items.value.includes(uid)) {
    items.value.push(uid)
  }
  return items.value.indexOf(uid)
}

function unregister(uid: number) {
  const index = items.value.indexOf(uid)
  if (index < 0) {
    return
  }

  items.value.splice(index, 1)

  if (currentIndex.value > index) {
    currentIndex.value -= 1
  }

  if (currentIndex.value >= items.value.length) {
    currentIndex.value = Math.max(0, items.value.length - 1)
  }
}

function openByIndex(index: number) {
  if (index >= 0 && index < items.value.length) {
    transitionDirection.value = index >= currentIndex.value ? 'forward' : 'backward'
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
  unregister,
  currentIndex,
  transitionDirection,
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
  transitionDirection.value = 'forward'
})
</script>

<template>
  <div class="accordion">
    <slot :next="next" :prev="prev" :hasNext="hasNext" :hasPrev="hasPrev"/>
  </div>
</template>
