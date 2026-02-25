<script setup lang="ts">
import Offcanvas, {type IOffcanvas} from '@/views/shared/components/Offcanvas.vue'
import {useRoute, useRouter} from 'vue-router'
import {computed, ref, toRef} from 'vue'
import {EPreloaderSize} from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import {useApp} from '@/app/app.ts'
import {ProductBrandModel} from '@/app/models/product_brand_model.ts'
import {ProductSubcategoryCreatedEvent} from '@/views/shared/types/events.ts'
import AddProductSubcategoryForm
  from "@/views/shared/components/Forms/AddProductSubcategoryForm.vue";

const app = useApp()
const route = useRoute()
const router = useRouter()

const categoryId = computed(() => route.params.categoryId as string)

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<typeof AddProductSubcategoryForm>()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

const handleHideOffcanvas = () => {
  offcanvasForm.value?.reset()
  router.back()
}

const handleOffcanvasFormSubmit = () => {
  offcanvas.value?.scrollTop()
}

const handleOffcanvasFormSubmitted = async (success: boolean, data?: ProductBrandModel) => {
  if (success && data) {
    app.events.emit(ProductSubcategoryCreatedEvent, data)
    await app.helpers.async.sleep(5000)
    offcanvas.value?.hide()
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
    title="Add Product Subcategory">
    <template #body>
      <add-product-subcategory-form
        ref="offcanvasForm"
        :categoryId="categoryId"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted"/>
    </template>
    <template v-if="offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true"/>
    </template>
  </offcanvas>
</template>
