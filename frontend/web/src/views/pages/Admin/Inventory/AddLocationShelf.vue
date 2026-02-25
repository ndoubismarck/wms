<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { useRoute, useRouter } from 'vue-router'
import { ref, toRef } from 'vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useApp } from '@/app/app.ts'
import type { LocationShelfModel } from '@/app/models/location_shelf_model.ts'
import AddLocationShelfForm from '@/views/shared/components/Forms/AddLocationShelfForm.vue'

const app = useApp()
const route = useRoute()
const router = useRouter()

const bayId = route.params?.bayId as string

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

const handleOffcanvasFormSubmitted = async (success: boolean, data?: LocationShelfModel) => {
  if (success && data) {
    app.events.emit('forms.location.add.shelf.submitted', data)
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
    title="Add Shelf">
    <template #body>
      <add-location-shelf-form
        ref="offcanvasForm"
        :bayId="bayId"
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
