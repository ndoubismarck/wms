<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField } from 'vee-validate'
import { useApp } from '@/app/app'

interface IInputOption {
  key: string
  value: any
}

const props = defineProps<{
  name: string
  hint?: string
  rules?: string
  label?: string
  value?: number[] | string[]
  options: { [key: number | string]: any }
  columnClass?: string
}>()

const app = useApp()

const model = ref<number[] | string[]>([])
const fields = ref<(typeof VeeField)[] | null>(null)

const name = toRef(props, 'name')
const hint = toRef(props, 'hint')
const rules = toRef(props, 'rules')
const label = toRef(props, 'label')
const value = toRef(props, 'value')
const options = toRef(props, 'options')
const columnClass = toRef(props, 'columnClass', 'col-sm-6 col-md-4')

const optionsToKeyValue = (options: { [key: number | string]: any }): IInputOption[] => {
  return Object.entries(options).map(([key, value]) => ({
    key: key,
    value: value,
  }))
}

const removeValueDuplicates = (value: undefined | number[] | string[]): number[] | string[] => {
  if (!value) {
    return []
  }
  return Array.from(new Set(value as any[]))
}

watch(value, () => {
  model.value = Array.from(new Set(removeValueDuplicates(value.value) as any[]))
})

onMounted(async () => {})

onBeforeUnmount(() => {
  fields.value = []
})
</script>

<template>
  <label v-if="label" class="form-label">{{ label }}</label>
  <div v-if="hint" class="form-hint mb-2">{{ hint }}</div>
  <div class="row">
    <div v-for="item in optionsToKeyValue(options)" :class="columnClass">
      <div class="form-check">
        <div class="d-flex align-items-center">
          <div class="me-2">
            <vee-field
              ref="fields"
              v-slot="{ field, meta }"
              v-model="model"
              :name="name"
              :value="item.key"
              :rules="rules"
              type="checkbox">
              <input
                v-bind="field"
                type="checkbox"
                class="form-check-input"
                :class="{ 'is-invalid': meta.touched && !meta.valid }"
                :id="`dropdown-checkbox-${name}-${item.value.key}`" />
            </vee-field>
          </div>
          <label class="form-check-label" :for="`dropdown-checkbox-${name}-${item.value.key}`">
            {{ `${item.value.charAt(0).toUpperCase()}${item.value.slice(1)}` }}
          </label>
        </div>
      </div>
    </div>
  </div>
  <vee-error-message v-slot="{ message }" :name="name">
    <p class="invalid-feedback d-block p-0 m-0">
      {{ app.helpers.veeValidate.errorMessage(message) }}
    </p>
  </vee-error-message>
</template>
