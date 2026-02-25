<script setup lang="ts">
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormButton from '@/views/shared/components/VueForm/VueFormButton.vue'
import { ref, toRef } from 'vue'
import { useApp } from '@/app/app'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import type { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model.ts'

const app = useApp()

const props = defineProps<{
  shelfId: string
}>()
const emits = defineEmits(['close', 'submit', 'submitted'])

const shelfId = toRef(props, 'shelfId')

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)
const validationSchema = {}

const reset = () => {
  form.value?.reset()
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')
  isSubmitting.value = true
  let success = false
  let payload: LocationShelfLevelModel | undefined
  try {
    data = { shelf_id: shelfId.value, ...data }
    const result = await app.services.locations.addShelfLevel(data)
    await app.helpers.async.sleep(2000)
    if (result.success) {
      success = true
      payload = result.shelfLevel
      form.value?.reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Shelf Level successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Shelf Level could not be added.'),
        },
      }
    }
  } catch (ex) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.messageFromError(ex),
      },
    }
    app.services.logger.error(ex)
  }
  isSubmitting.value = false
  emits('submitted', success, payload)
  await app.helpers.async.sleep(5000)
  if (success) {
    emits('close')
  }
}

defineExpose({
  reset,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit" :validationSchema="validationSchema">
        <alert-message v-if="alert" :params="alert" />
        <div class="mb-6">
          <vue-form-input
            :disabled="isSubmitting"
            inputMode="numeric"
            name="code"
            type="text"
            label="Shelf Level Code"
            placeholder="Shelf Level Code" />
        </div>
        <div class="mb-6">
          <vue-form-input
            :disabled="isSubmitting"
            name="name"
            type="text"
            label="Shelf Level Name"
            placeholder="Shelf Level Name" />
        </div>
        <div class="mb-6">
          <vue-form-textarea-input
            :disabled="isSubmitting"
            name="description"
            label="Shelf Level Description"
            placeholder="Shelf Level Description" />
        </div>
        <vue-form-button :disabled="isSubmitting" text="Submit" />
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
