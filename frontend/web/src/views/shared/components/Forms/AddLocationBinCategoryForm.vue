<script setup lang="ts">
import { ref } from 'vue'

const emits = defineEmits(['close', 'submit', 'submitted'])

const isSubmitting = ref<boolean>(false)

const reset = () => {
  isSubmitting.value = false
}

const handleClose = () => {
  emits('close')
}

defineExpose({
  reset,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <p class="mb-4 text-body-secondary">
        Bin categories are not configured yet in this build.
      </p>
      <button type="button" class="btn btn-outline-primary" @click="handleClose">Close</button>
    </div>
  </div>
</template>

<style scoped></style>
