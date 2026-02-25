<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRouter } from 'vue-router'
import { ref, toRef } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import AddLocationBinForm from '@/views/shared/components/Forms/AddLocationBinForm.vue'
import type { LocationBinModel } from '@/app/models/location_bin_model.ts'
import { useApp } from '@/app/app.ts'

const app = useApp()
const router = useRouter()

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<any | null>()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

const handleAddAisle = () => {
  router.push({
    name: 'admin:inventory:products:aisle:add',
  })
}

const handleAddBay = (aisleId: string) => {
  router.push({
    name: 'admin:inventory:products:bay:add',
    params: { aisleId: aisleId },
  })
}

const handleAddShelf = (bayId: string) => {
  router.push({
    name: 'admin:inventory:products:shelf:add',
    params: { bayId: bayId },
  })
}

const handleAddShelfLevel = (shelfId: string) => {
  router.push({
    name: 'admin:inventory:products:shelf:level:add',
    params: { shelfId: shelfId },
  })
}

const handleAddBinCategory = () => {
  router.push({
    name: 'admin:inventory:products:bin:category:add',
  })
}

const handleHideOffcanvas = () => {
  if (offcanvasForm.value) {
    offcanvasForm.value.reset()
  }
  router.back()
}

const handleOffcanvasFormClose = () => {
  router.back()
}

const handleOffcanvasFormSubmit = () => {
  offcanvas.value?.scrollTop()
}

const handleOffcanvasFormSubmitted = async (success: boolean, data?: LocationBinModel) => {
  if (success && data) {
    app.events.emit('forms.location.add.bin.submitted', data)
  }
}
</script>

<template>
  <offcanvas
    ref="offcanvas"
    :show="true"
    @hide="handleHideOffcanvas"
    :parent="parent"
    :overflow="!offcanvasForm?.isSubmitting || !offcanvasForm?.isSubmitting"
    title="Add Bin">
    <template #body>
      <add-location-bin-form
        ref="offcanvasForm"
        @close="handleOffcanvasFormClose"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted"
        @addBay="handleAddBay"
        @addAisle="handleAddAisle"
        @addShelf="handleAddShelf"
        @addShelfLevel="handleAddShelfLevel"
        @addBinCategory="handleAddBinCategory" />
    </template>
    <template v-if="offcanvasForm?.isLoading || offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
  <router-view :parent="offcanvas" />
</template>
