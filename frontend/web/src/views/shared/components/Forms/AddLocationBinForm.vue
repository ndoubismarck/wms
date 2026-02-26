<script setup lang="ts">
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import VueFormButton from '@/views/shared/components/VueForm/VueFormButton.vue'
import { ref, toRef } from 'vue'
import { useApp } from '@/app/app.ts'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import type { LocationBinModel } from '@/app/models/location_bin_model.ts'

const app = useApp()

const props = defineProps<{
  shelfLevelId?: string
}>()

const emits = defineEmits(['close', 'submit', 'submitted'])

const shelfLevelId = toRef(props, 'shelfLevelId')
const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)
const isLoading = ref<boolean>(false)

const reset = () => {
  form.value?.reset()
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')
  if (!shelfLevelId.value) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.message('Shelf Level is required to create a bin.'),
      },
    }
    return
  }

  isSubmitting.value = true
  let success = false
  let payload: LocationBinModel | undefined
  try {
    const result = await app.services.locations.addBin({
      ...data,
      shelf_level_id: shelfLevelId.value,
    })
    await app.helpers.async.sleep(2000)
    if (result.success) {
      success = true
      payload = result.bin
      form.value?.reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Bin successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Bin could not be added.'),
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
  isLoading,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert" />

        <div class="mb-6">
          <vue-form-input :disabled="isSubmitting" name="name" type="text" label="Bin Name" placeholder="Bin Name" />
        </div>

        <div class="mb-6">
          <vue-form-textarea-input
            :disabled="isSubmitting"
            name="description"
            label="Bin Description"
            placeholder="Bin Description" />
        </div>

        <vue-form-button :disabled="isSubmitting || !shelfLevelId" text="Submit" />
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
