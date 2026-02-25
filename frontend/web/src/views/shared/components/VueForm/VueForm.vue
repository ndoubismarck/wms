<script setup lang="ts">
import { Form as VeeForm, type FormActions as VeeFormActions } from 'vee-validate'
import { onMounted, provide, ref, toRef } from 'vue'

export type VueFormData = { [key: string]: any }

export type VueFormActions = VeeFormActions<any>

const emit = defineEmits<{
  (e: 'submit', data: VueFormData, actions: VueFormActions): void
}>()

const props = defineProps<{
  keepValues?: boolean
  validationSchema?: { [key: string]: any }
}>()

const keepValues = toRef(props, 'keepValues', true)
const validationSchema = toRef(props, 'validationSchema')

const form = ref<typeof VeeForm>()
const resetCount = ref<number>(0)

const reset = () => {
  resetCount.value++
  form.value?.resetForm()
}

const resetField = (field: string) => {
  form.value?.resetField(field)
}

const handleSubmit = (data: VueFormData, actions: VueFormActions) => {
  emit('submit', data, actions)
}

provide('formReset', resetCount)

defineExpose({
  reset,
  resetField,
})

onMounted(() => {})
</script>

<template>
  <vee-form ref="form" @submit="handleSubmit" :keepValues="keepValues" :validationSchema="validationSchema">
    <slot></slot>
  </vee-form>
</template>

<style scoped></style>
