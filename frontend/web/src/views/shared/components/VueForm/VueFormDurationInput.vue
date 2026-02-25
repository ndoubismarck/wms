<script setup lang="ts">
import { useApp } from '@/app/app'
import { onMounted, ref, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField } from 'vee-validate'

const props = defineProps<{
  name: string
  type: string
  hint?: string
  rules?: string
  label: string
  value?: string
  disabled: boolean
  placeholder?: string
}>()

const app = useApp()

const model = ref<string>('')
const duration = ref<number | null>(null)
const durationUnit = ref<string | null>('H')

const value = toRef(props, 'value')

watch(model, () => {
  if (model.value.length == 0) {
    duration.value = null
    durationUnit.value = 'H'
  }
})

watch(duration, () => {
  if (duration.value && durationUnit.value) {
    if (duration.value > 0 && durationUnit.value.length > 0) {
      model.value = `${duration.value} ${durationUnit.value}`
    }
  }
})

watch(durationUnit, () => {
  if (duration.value && durationUnit.value) {
    if (duration.value > 0 && durationUnit.value.length > 0) {
      model.value = `${duration.value} ${durationUnit.value}`
    }
  }
})

const name = toRef(props, 'name')
const hint = toRef(props, 'hint')
const label = toRef(props, 'label')
const rules = toRef(props, 'rules')
const placeholder = toRef(props, 'placeholder')

onMounted(async () => {
  if (value.value) {
    const split = value.value.split(' ')
    if (split[0] && split[0].trim().length > 0) {
      duration.value = parseInt(split[0].trim())
    }
    if (split[1] && split[1].trim().length > 0) {
      durationUnit.value = split[1].trim().slice(0, 1).toUpperCase()
    }
  }
})
</script>

<template>
  <small v-if="hint" class="text-muted d-block">{{ hint }}</small>
  <vee-field v-slot="{ field, meta }" :name="name" :rules="rules" v-model="model">
    <div class="input-group">
      <div class="form-floating w-75">
        <input
          v-model="duration"
          type="number"
          :disabled="disabled"
          :placeholder="placeholder"
          class="form-control no-icon"
          :class="{ 'is-invalid': meta.validated && !meta.valid }" />
        <label>{{ label }}</label>
      </div>
      <select
        v-model="durationUnit"
        class="form-control w-25"
        :disabled="disabled"
        :class="{ 'is-invalid': meta.validated && !meta.valid }">
        <option value="H">Hours</option>
        <option value="D">Days</option>
        <option value="W">Weeks</option>
        <option value="M">Months</option>
        <option value="Y">Years</option>
      </select>
      <input v-bind="field" type="hidden" class="d-none" />
    </div>
  </vee-field>
  <vee-error-message v-slot="{ message }" :name="name">
    <p class="invalid-feedback d-block p-0 m-0">
      {{ app.helpers.veeValidate.errorMessage(message) }}
    </p>
  </vee-error-message>
</template>
