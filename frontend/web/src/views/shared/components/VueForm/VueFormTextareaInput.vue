<script setup lang="ts">
import {onMounted, ref, toRef} from 'vue'
import {ErrorMessage as VeeErrorMessage, Field as VeeField} from 'vee-validate'
import {useApp} from '@/app/app'

const props = defineProps<{
  name: string
  rules?: string
  label: string
  value?: string
  disabled?: boolean
  placeholder?: string
}>()

const app = useApp()

const model = ref<string>('')

const name = toRef(props, 'name')
const label = toRef(props, 'label')
const rules = toRef(props, 'rules')
const value = toRef(props, 'value')
const placeholder = toRef(props, 'placeholder')

onMounted(async () => {
  if (value.value) {
    model.value = value.value
  }
})
</script>

<template>
  <vee-field v-slot="{ field, meta }" :name="name" :rules="rules" v-model="model">
    <div class="form-control-fi">
      <textarea
        v-bind="field"
        :disabled="disabled"
        :placeholder="placeholder"
        class="form-control"
        style="min-height: 100px"
        :class="{ 'is-invalid': meta.validated && !meta.valid }"/>
      <label>{{ label }}</label>
    </div>
  </vee-field>
  <vee-error-message v-slot="{ message }" :name="name">
    <p class="invalid-feedback d-block p-0 m-0">{{
        app.helpers.veeValidate.errorMessage(message)
      }}</p>
  </vee-error-message>
</template>
