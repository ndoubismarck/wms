<script setup lang="ts">
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import { EPreloaderSize } from '@/app/types.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import AddStaffForm from '@/views/shared/components/Forms/AddStaffForm.vue'
import { ref, toRef } from 'vue'
import { useRouter } from 'vue-router'
import { useUsersStore } from '@/stores/users_store.ts'
import type { UserModel } from '@/app/models/user_model.ts'

const usersStore = useUsersStore()
const router = useRouter()

const props = defineProps<{
  parent?: IOffcanvas
}>()

const parent = toRef(props, 'parent')

const offcanvas = ref<IOffcanvas>()
const offcanvasForm = ref<{
  reset: () => void
  isSubmitting: boolean
} | null>(null)

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

const handleOffcanvasFormSubmitted = (success: boolean, user?: UserModel) => {
  if (success && user) {
    usersStore.set(user)
  }
}
</script>

<template>
  <Offcanvas
    ref="offcanvas"
    :show="true"
    :parent="parent"
    @hide="handleHideOffcanvas"
    :overflow="offcanvasForm?.isSubmitting"
    title="Add Staff User">
    <template #body>
      <AddStaffForm
        ref="offcanvasForm"
        @close="handleOffcanvasFormClose"
        @submit="handleOffcanvasFormSubmit"
        @submitted="handleOffcanvasFormSubmitted" />
    </template>
    <template v-if="offcanvasForm?.isSubmitting" #overlay>
      <Preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </Offcanvas>
</template>
