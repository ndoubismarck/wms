<script setup lang="ts">
import { useApp } from '@/app/app'
import { computed, onBeforeMount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { EPreloaderSize } from '@/app/types'
import ResizeObserver from 'resize-observer-polyfill'
import SetupLayout from '@/views/layouts/SetupLayout.vue'
import AuthLayout from '@/views/layouts/AuthLayout.vue'
import AdminLayout from '@/views/layouts/AdminLayout.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useUserStore } from '@/stores/user_store'
import { useSetupStore } from '@/stores/setup_store'
import ErrorLayout from '@/views/layouts/ErrorLayout.vue'
import { Popover } from 'bootstrap'

const app = useApp()

const userStore = useUserStore()
const setupStore = useSetupStore()

const me = computed(() => userStore.get().value)

const route = useRoute()
const router = useRouter()

const isError = ref<boolean>(false)
const isLoading = ref<boolean>(true)
const isAuthPage = ref<boolean>(false)
const isAdminPage = ref<boolean | null>(null)
const isSetupPage = ref<boolean | null>(null)
const isValidSetup = ref<boolean>(false)
const isAuthenticated = ref<boolean | null>(null)

const getMe = async (): Promise<boolean> => {
  try {
    const result = await app.services.auth.getMe()
    if (result.success) {
      userStore.set(result.user)
      return result.user.id.length > 0
    }
  } catch (ex) {
    isError.value = true
    app.services.logger.error(ex)
  }
  return false
}

const getSetup = async (): Promise<boolean> => {
  try {
    const result = await app.services.setup.get()
    if (result.success) {
      setupStore.set(result.setup)
      return result.setup.admin && result.setup.database
    }
  } catch (ex) {
    isError.value = true
    app.services.logger.error(ex)
  }
  return false
}

const enforceSetup = async (init: boolean) => {
  if (isValidSetup.value != null) {
    if (!isValidSetup.value && !isSetupPage.value) {
      if (init) {
        await app.helpers.async.sleep(3000)
      }
      window.location.href = router.resolve({ name: 'setup:main' }).href
    }
  }
}

const enforceAuthentication = async (init: boolean) => {
  if (isValidSetup.value != null && isAuthenticated.value != null) {
    if (isValidSetup.value && !isAuthPage.value && !isAuthenticated.value) {
      if (init) {
        await app.helpers.async.sleep(3000)
      }
      window.location.href = router.resolve({ name: 'auth:login' }).href
    }
  }
}

const enforceAuthorization = async (init: boolean) => {
  if (me.value && isAuthenticated.value) {
    //todo: add authorization
  }
}

watch(route, () => {
  if (route.name) {
    const routeName = route.name as string
    isAuthPage.value = routeName.startsWith('auth:')
    isSetupPage.value = routeName.startsWith('setup:')
    isAdminPage.value = routeName.startsWith('admin:')
  }
  isError.value = !!(isSetupPage.value && isValidSetup.value)
})

watch(isValidSetup, async () => {
  await enforceSetup(false)
})

watch(isAuthenticated, async () => {
  await enforceAuthentication(false)
})

onMounted(async () => {
  isValidSetup.value = await getSetup()
  await enforceSetup(true)
  if (isValidSetup.value) {
    isAuthenticated.value = await getMe()
    await enforceAuthentication(true)
    await enforceAuthorization(true)
  }
  if (isSetupPage.value && isValidSetup.value) {
    isError.value = true
  }
  await app.helpers.async.sleep(3000)
  isLoading.value = false

  document.querySelectorAll('[data-bs-toggle="popover"]').forEach((el) => new Popover(el))
})

onBeforeMount(() => {
  isLoading.value = true
  if (route.name) {
    const routeName = route.name as string
    isAuthPage.value = routeName.startsWith('auth:')
    isSetupPage.value = routeName.startsWith('setup:')
  }
})

app.services.auth.onUnAuthorized((response: Response) => {
  isAuthenticated.value = false
  app.services.logger.warn('unauthorized:', response.url)
})

window.ResizeObserver = ResizeObserver
</script>

<template>
  <template v-if="isLoading">
    <preloader :size="EPreloaderSize.LG" />
  </template>
  <template v-else>
    <template v-if="isError">
      <error-layout>
        <router-view />
      </error-layout>
    </template>
    <template v-else>
      <template v-if="isValidSetup">
        <template v-if="isAuthPage || !isAuthenticated">
          <auth-layout>
            <router-view />
          </auth-layout>
        </template>
        <template v-else>
          <admin-layout>
            <router-view />
          </admin-layout>
        </template>
      </template>
      <template v-else>
        <setup-layout>
          <router-view />
        </setup-layout>
      </template>
    </template>
  </template>
</template>
