<script setup lang="ts">
import { useApp } from '@/app/app'
import { onMounted, ref, toRef } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField, type RuleExpression } from 'vee-validate'

const props = defineProps<{
  name: string
  label: string
  value?: string
  disabled?: boolean
  placeholder?: string
  mode?: 'date' | 'datetime-local'
  validation?: RuleExpression<any>
}>()

const app = useApp()

const model = ref<string>('')

const name = toRef(props, 'name')
const label = toRef(props, 'label')
const value = toRef(props, 'value')
const disabled = toRef(props, 'disabled', false)
const placeholder = toRef(props, 'placeholder')
const mode = toRef(props, 'mode', 'datetime-local')
const validation = toRef(props, 'validation')

onMounted(() => {
  if (value.value) {
    model.value = value.value
  }
})
</script>

<template>
  <div class="w-100">
    <vee-field v-slot="{ field, meta }" :name="name" :rules="validation" v-model="model">
      <div class="form-control-fi">
        <input
          v-bind="field"
          :type="mode"
          :disabled="disabled"
          :placeholder="placeholder"
          class="form-control"
          :class="{ 'is-invalid': meta.touched && meta.validated && !meta.valid }" />
        <label>{{ label }}</label>
      </div>
    </vee-field>
    <vee-error-message v-slot="{ message }" :name="name">
      <span class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </span>
    </vee-error-message>
  </div>
</template>
