<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRoute, useRouter } from 'vue-router'
import { ref, toRef } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import AddLocationBinForm from '@/views/shared/components/Forms/AddLocationBinForm.vue'
import type { LocationBinModel } from '@/app/models/location_bin_model.ts'
import { useApp } from '@/app/app.ts'
import { LocationBinCreatedEvent } from '@/views/shared/types/events.ts'

const app = useApp()
const route = useRoute()
const router = useRouter()
const shelfLevelId = route.params?.shelfLevelId as string

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<any | null>()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

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
    app.events.emit(LocationBinCreatedEvent, data)
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
        :shelfLevelId="shelfLevelId"
        @close="handleOffcanvasFormClose"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted" />
    </template>
    <template v-if="offcanvasForm?.isLoading || offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
  <router-view :parent="offcanvas" />
</template>
