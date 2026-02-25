<script setup lang="ts">
import { ref, toRef, watch } from 'vue'

const emits = defineEmits(['click'])

const props = defineProps<{
  text: string
  class?: string
  disabled?: boolean
  isLoading?: boolean
}>()

const disabled = toRef(props, 'disabled', false)
const isLoading = toRef(props, 'isLoading', false)
const btnClass = toRef(props, 'class', 'btn-outline-primary')

const element = ref<HTMLButtonElement | null>(null)
const elementParams = ref<{ width?: number; height?: number }>({})

watch(element, () => {
  if (element.value) {
    elementParams.value.width = element.value.offsetWidth
    elementParams.value.height = element.value.offsetHeight
  }
})

watch(isLoading, () => {
  if (element.value) {
    if (elementParams.value.width) {
      element.value.style.width = `${elementParams.value.width}px`
    }
    if (elementParams.value.height) {
      element.value.style.height = `${elementParams.value.height}px`
    }
  }
})

const handleClick = (evt: Event) => {
  emits('click', evt)
}
</script>

<template>
  <button ref="element" @click="handleClick" :disabled="disabled || isLoading" type="submit" :class="`btn ${btnClass}`">
    <span v-if="isLoading" class="spinner-border spinner-border-sm text-primary" />
    <span v-else>{{ text }}</span>
  </button>
</template>

<style scoped></style>
