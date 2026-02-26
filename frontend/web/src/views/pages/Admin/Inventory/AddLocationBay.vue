<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRoute, useRouter } from 'vue-router'
import { ref, toRef } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useApp } from '@/app/app.ts'
import type { LocationBayModel } from '@/app/models/location_bay_model.ts'
import AddLocationBayForm from '@/views/shared/components/Forms/AddLocationBayForm.vue'
import { LocationBayCreatedEvent } from '@/views/shared/types/events.ts'

const app = useApp()
const route = useRoute()
const router = useRouter()

const aisleId = route.params?.aisleId as string

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

const handleOffcanvasFormSubmitted = async (success: boolean, data?: LocationBayModel) => {
  if (success && data) {
    app.events.emit(LocationBayCreatedEvent, data)
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
    title="Add Bay">
    <template #body>
      <add-location-bay-form
        ref="offcanvasForm"
        :aisleId="aisleId"
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
