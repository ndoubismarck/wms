<script setup lang="ts">
import { onMounted, ref, toRef } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField } from 'vee-validate'
import { useApp } from '@/app/app'

const props = defineProps<{
  name: string
  hint?: string
  label: string
}>()

const app = useApp()

const model = ref<string>('')

const name = toRef(props, 'name')
const hint = toRef(props, 'hint')
const label = toRef(props, 'label')

onMounted(async () => {})
</script>

<template>
  <vee-field v-slot="{ field }" :name="name" v-model="model">
    <label class="form-label d-block mb-0">{{ label }}</label>
    <small v-if="hint" class="text-muted d-block mb-1">
      {{ hint }}
    </small>
    <div class="form-check form-check-inline me-4">
      <input v-bind="field" class="form-check-input" type="radio" value="yes" />
      <label class="form-check-label">Yes</label>
    </div>
    <div class="form-check form-check-inline">
      <input v-bind="field" class="form-check-input" type="radio" value="no" />
      <label class="form-check-label">No</label>
    </div>
  </vee-field>
  <vee-error-message v-slot="{ message }" :name="name">
    <p class="invalid-feedback d-block p-0 m-0">
      {{ app.helpers.veeValidate.errorMessage(message) }}
    </p>
  </vee-error-message>
</template>

<style scoped></style>
