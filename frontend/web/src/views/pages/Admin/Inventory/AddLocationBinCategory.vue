<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRouter } from 'vue-router'
import { ref, toRef } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useApp } from '@/app/app.ts'
import AddLocationBinCategoryForm from '@/views/shared/components/Forms/AddLocationBinCategoryForm.vue'
import type { ProductCategoryModel } from '@/app/models/product_category_model.ts'

const app = useApp()
const router = useRouter()

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<any | null>()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

const handleHideOffcanvas = () => {
  offcanvasForm.value?.reset()
  router.back()
}

const handleOffcanvasFormClose = () => {
  router.back()
}

const handleOffcanvasFormSubmit = () => {
  offcanvas.value?.scrollTop()
}

const handleOffcanvasFormSubmitted = async (success: boolean, data?: ProductCategoryModel) => {
  if (success && data) {
    app.events.emit('forms.location.add.bin.category.submitted', data)
  }
}
</script>

<template>
  <offcanvas
    ref="offcanvas"
    :show="true"
    @hide="handleHideOffcanvas"
    :parent="parent"
    :overflow="offcanvasForm?.isSubmitting"
    title="Add Bin Category">
    <template #body>
      <add-location-bin-category-form
        ref="offcanvasForm"
        @close="handleOffcanvasFormClose"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted" />
    </template>
    <template v-if="offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
  <router-view />
</template>
