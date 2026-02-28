<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import AddTaskForm from '@/views/shared/components/Forms/AddTaskForm.vue'
import { ref, toRef } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<any | null>()

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

const handleOffcanvasFormSubmitted = () => {}
</script>

<template>
  <offcanvas
    ref="offcanvas"
    :show="true"
    :parent="parent"
    @hide="handleHideOffcanvas"
    :overflow="offcanvasForm?.isSubmitting"
    title="Add Task">
    <template #body>
      <add-task-form
        ref="offcanvasForm"
        @close="handleOffcanvasFormClose"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted" />
    </template>
    <template v-if="offcanvasForm?.isLoading || offcanvasForm?.isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
</template>
