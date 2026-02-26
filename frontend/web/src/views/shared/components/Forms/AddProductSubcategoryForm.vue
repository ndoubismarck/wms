<script setup lang="ts">
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import { computed, onMounted, ref, toRef } from 'vue'
import { useApp } from '@/app/app'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import { useProductCategoriesStore } from '@/stores/product_categories_store.ts'
import type { ProductSubcategoryModel } from '@/app/models/product_subcategory_model.ts'

const app = useApp()
const productCategoriesStore = useProductCategoriesStore()

const props = defineProps<{
  categoryId: string
}>()

const emits = defineEmits(['submit', 'submitted'])

const categoryId = toRef(props, 'categoryId')

const productCategory = computed(() => productCategoriesStore.getById(categoryId.value))

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)

const reset = () => {
  form.value?.reset()
}

const handleSubmit = async (data: VueFormData) => {
  if (!categoryId.value) {
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
    let payload: ProductSubcategoryModel | undefined
    try {
      data['category_id'] = categoryId.value
      const result = await app.services.products.addSubcategory(data)
      await app.helpers.async.sleep(2000)
      if (result.success) {
        success = true
        payload = result.subcategory
        form.value?.reset()
        alert.value = {
          type: EAlertMessageType.Success,
          body: {
            text: app.helpers.i18n.message('Product subcategory successfully added.'),
          },
        }
      } else {
        alert.value = {
          type: EAlertMessageType.Error,
          body: {
            text: app.helpers.i18n.message('Product subcategory could not be added.'),
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
  if (!categoryId.value) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.message('Parent product category is missing.'),
      },
    }
  } else {
    try {
      const result = await app.services.products.getCategory(categoryId.value)
      await app.helpers.async.sleep(2000)
      if (result.success) {
        productCategoriesStore.set(result.category)
      } else {
        alert.value = {
          type: EAlertMessageType.Error,
          body: {
            text: app.helpers.i18n.message('Parent product category could not be fetched.'),
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
            :disabled="isSubmitting || !productCategory"
            name="name"
            type="text"
            :label="`${productCategory?.name} Subcategory Name`"
            :placeholder="`Enter ${productCategory?.name} Subcategory Name`" />
        </div>
        <button :disabled="isSubmitting || !productCategory" class="btn btn-primary">Submit</button>
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
