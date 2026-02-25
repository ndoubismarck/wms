<script setup lang="ts">
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import { ref } from 'vue'
import { useApp } from '@/app/app'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import { ProductBrandModel } from '@/app/models/product_brand_model.ts'

const app = useApp()

const emits = defineEmits(['submit', 'submitted'])

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)

const reset = () => {
  form.value?.reset()
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')
  isSubmitting.value = true
  let success = false
  let payload: ProductBrandModel | undefined
  try {
    const result = await app.services.products.addBrand(data)
    await app.helpers.async.sleep(2000)
    if (result.success) {
      success = true
      payload = result.brand
      form.value?.reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Product brand successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Product brand could not be added.'),
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
}

defineExpose({
  reset,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert" />
        <div class="mb-6">
          <vue-form-input
            :disabled="isSubmitting"
            name="name"
            type="text"
            label="Brand Name"
            placeholder="Enter Brand Name" />
        </div>
        <button :disabled="isSubmitting" class="btn btn-primary">Submit</button>
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
