<script setup lang="ts">
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import { computed, onMounted, ref, toRef } from 'vue'
import { useApp } from '@/app/app'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import { ProductCategoryModel } from '@/app/models/product_category_model.ts'
import { useProductBrandsStore } from '@/stores/product_brands_store.ts'

const app = useApp()
const productBrandsStore = useProductBrandsStore()

const props = defineProps<{
  brandId: string
}>()

const emits = defineEmits(['submit', 'submitted'])

const brandId = toRef(props, 'brandId')

const productBrand = computed(() => productBrandsStore.getById(brandId.value))

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)

const reset = () => {
  form.value?.reset()
}

const handleSubmit = async (data: VueFormData) => {
  if (!brandId.value) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.message('Parent product brand is missing.'),
      },
    }
  } else {
    emits('submit')
    isSubmitting.value = true
    let success = false
    let payload: ProductCategoryModel | undefined
    try {
      data['brand_id'] = brandId.value
      const result = await app.services.products.addCategory(data)
      await app.helpers.async.sleep(2000)
      if (result.success) {
        success = true
        payload = result.category
        form.value?.reset()
        alert.value = {
          type: EAlertMessageType.Success,
          body: {
            text: app.helpers.i18n.message('Product category successfully added.'),
          },
        }
      } else {
        alert.value = {
          type: EAlertMessageType.Error,
          body: {
            text: app.helpers.i18n.message('Product category could not be added.'),
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
}

const getProductBrand = async () => {
  if (!brandId.value) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.message('Parent product brand is missing.'),
      },
    }
  } else {
    try {
      const result = await app.services.products.getBrand(brandId.value)
      await app.helpers.async.sleep(2000)
      if (result.success) {
        productBrandsStore.set(result.brand)
      } else {
        alert.value = {
          type: EAlertMessageType.Error,
          body: {
            text: app.helpers.i18n.message('Parent product brand could not be fetched.'),
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
  }
}

defineExpose({
  reset,
  isSubmitting,
})

onMounted(async () => {
  await getProductBrand()
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert" />
        <div class="mb-6">
          <vue-form-input
            :disabled="isSubmitting || !productBrand"
            name="name"
            type="text"
            :label="`${productBrand?.name} Brand Category Name`"
            :placeholder="`Enter ${productBrand?.name} Brand Category Name`" />
        </div>
        <button :disabled="isSubmitting || !productBrand" class="btn btn-primary">Submit</button>
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
