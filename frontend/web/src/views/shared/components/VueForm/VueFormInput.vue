<script setup lang="ts">
import { useApp } from '@/app/app'
import { onMounted, ref, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField, type RuleExpression } from 'vee-validate'

type TInputMode = 'text' | 'email' | 'search' | 'tel' | 'url' | 'none' | 'numeric' | 'decimal' | undefined

type TValueFormat = 'lowercase' | 'uppercase'

const props = defineProps<{
  name: string
  type: string
  max?: number
  label: string
  value?: string
  disabled?: boolean
  inputMode?: TInputMode
  validation?: RuleExpression<any>
  placeholder?: string
  valueFormat?: TValueFormat
}>()

const app = useApp()

const model = ref<any | null>(null)

const type = toRef(props, 'type')
const name = toRef(props, 'name')
const label = toRef(props, 'label')
const value = toRef(props, 'value')
const disabled = toRef(props, 'disabled', false)
const inputMode = toRef(props, 'inputMode')
const valueFormat = toRef(props, 'valueFormat')
const validation = toRef(props, 'validation')
const placeholder = toRef(props, 'placeholder')

const toggleInputType = ref<string>(type.value)

watch(model, () => {
  if (valueFormat.value && type.value == 'text') {
    if (valueFormat.value == 'lowercase') {
      model.value = `${model.value}`.toLowerCase()
    }
    if (valueFormat.value == 'uppercase') {
      model.value = `${model.value}`.toUpperCase()
    }
  }
})

onMounted(async () => {
  if (value.value) {
    model.value = value.value
  }
})

const handleTogglePasswordInput = (evt: Event) => {
  if (evt.isTrusted) {
    toggleInputType.value = toggleInputType.value == 'password' ? 'text' : 'password'
  }
}
</script>

<template>
  <div class="w-100">
    <div :class="`input-${type}`">
      <vee-field v-slot="{ field, meta }" :name="name" :rules="validation" v-model="model">
        <div class="form-floating">
          <input
            v-bind="field"
            :type="toggleInputType"
            :disabled="disabled"
            :placeholder="placeholder"
            :inputmode="inputMode"
            :max="max"
            class="form-control"
            :class="{ 'is-invalid': meta.touched && !meta.valid }" />
          <label>{{ label }}</label>
        </div>
      </vee-field>
      <template v-if="type == 'password'">
        <span @click="handleTogglePasswordInput" class="toggler">
          <i
            class="icon-base bx"
            :class="{
              'bx-hide': toggleInputType == 'text',
              'bx-show': toggleInputType == 'password',
            }" />
        </span>
      </template>
    </div>
    <vee-error-message v-slot="{ message }" :name="name">
      <span class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </span>
    </vee-error-message>
  </div>
</template>

<style scoped>
.input-password {
  position: relative;
}

.input-password input {
  padding-right: 35px;
}

.input-password .toggler {
  top: calc(50% - 10px);
  right: 10px;
  position: absolute;
  cursor: pointer;
}

.input-password .toggler:hover {
  color: var(--bs-primary);
}
</style>
