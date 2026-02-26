<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRoute, useRouter } from 'vue-router'
import AddProductForm from '@/views/shared/components/Forms/AddProductForm.vue'
import { onBeforeMount, ref } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useApp } from '@/app/app.ts'
import type { ProductModel } from '@/app/models/product_model.ts'
import { useProductsStore } from '@/stores/products_store.ts'

const app = useApp()
const route = useRoute()
const router = useRouter()

const productsStore = useProductsStore()

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<any>()
const showChildView = ref<boolean>(false)

const handleAddProductBrand = () => {
  showChildView.value = true
  router.push({
    name: 'admin:inventory:products:add:brand',
  })
}

const handleAddProductCategory = async (brandId: string) => {
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:add:category',
    params: { brandId: brandId },
  })
}

const handleAddProductSubcategory = async (categoryId: string) => {
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:add:subcategory',
    params: { categoryId: categoryId },
  })
}

const handleAddLocationAisle = async () => {
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:aisle:add',
  })
}

const handleAddLocationBay = async (aisleId?: string) => {
  if (!aisleId) {
    return
  }
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:bay:add',
    params: { aisleId },
  })
}

const handleAddLocationShelf = async (bayId?: string) => {
  if (!bayId) {
    return
  }
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:shelf:add',
    params: { bayId },
  })
}

const handleAddLocationShelfLevel = async (shelfId?: string) => {
  if (!shelfId) {
    return
  }
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:shelf:level:add',
    params: { shelfId },
  })
}

const handleAddLocationBin = async (shelfLevelId?: string) => {
  if (!shelfLevelId) {
    return
  }
  showChildView.value = true
  await app.helpers.async.sleep(100)
  await router.push({
    name: 'admin:inventory:products:bin:add',
    params: { shelfLevelId },
  })
}

const handleHideOffcanvas = () => {
  if (offcanvasForm.value) {
    offcanvasForm.value.reset()
  }
  router.back()
}

const handleOffcanvasFormSubmit = () => {
  offcanvas.value?.scrollTop()
}

const handleOffcanvasFormSubmitted = async (success: boolean, data?: ProductModel) => {
  if (success && data) {
    productsStore.set(data)
    await app.helpers.async.sleep(5000)
    offcanvas.value?.hide()
  }
}

onBeforeMount(() => {
  const name = route.name?.toString()
  if (name != 'admin:inventory:products:add') {
    showChildView.value = false
    router.replace({
      name: 'admin:inventory:products:add',
    })
  }
})
</script>

<template>
  <offcanvas
    ref="offcanvas"
    :show="true"
    @hide="handleHideOffcanvas"
    :overflow="!offcanvasForm?.isLoading || !offcanvasForm?.isSubmitting"
    title="Add Product">
    <template #body>
      <add-product-form
        ref="offcanvasForm"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted"
        @addProductBrand="handleAddProductBrand"
        @addProductCategory="handleAddProductCategory"
        @addProductSubcategory="handleAddProductSubcategory"
        @addLocationAisle="handleAddLocationAisle"
        @addLocationBay="handleAddLocationBay"
        @addLocationShelf="handleAddLocationShelf"
        @addLocationShelfLevel="handleAddLocationShelfLevel"
        @addLocationBin="handleAddLocationBin" />
    </template>
    <template v-if="offcanvasForm?.isLoading || offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
  <router-view v-if="offcanvas && showChildView" :parent="offcanvas" />
</template>
