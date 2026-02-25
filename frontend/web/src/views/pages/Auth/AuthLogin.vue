<script setup lang="ts">
import { ref } from 'vue'
import { useApp } from '@/app/app.ts'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueForm, { type VueFormActions, type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { useRouter } from 'vue-router'

const app = useApp()
const router = useRouter()

const rememberMe = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const alertMessage = ref<IAlertMessage | null>(null)

const handleSubmit = async (data: VueFormData, actions: VueFormActions) => {
  try {
    isSubmitting.value = true
    const result = await app.services.auth.login({
      username: data.username,
      password: data.password,
      remember: rememberMe.value,
    })
    await app.helpers.async.sleep(2000)
    isSubmitting.value = false
    if (result.success) {
      alertMessage.value = {
        type: EAlertMessageType.Success,
        body: {
          text: 'Login successful.',
        },
      }
    } else {
      alertMessage.value = {
        type: EAlertMessageType.Error,
        body: {
          text: 'Login failed. Invalid username or password.',
        },
      }
    }
    if (result.success) {
      await app.helpers.async.sleep(3500)
      window.location.href = router.resolve({ name: 'admin:home' }).href
    }
  } catch (ex) {
    app.services.logger.error(ex)
    await app.helpers.async.sleep(2000)
    isSubmitting.value = false
    alertMessage.value = {
      type: EAlertMessageType.Error,
      body: {
        text: 'Login failed. An unexpected error has occurred.',
      },
    }
  }
}
</script>

<template>
  <div class="card px-sm-6 px-0">
    <div class="card-body position-relative">
      <preloader v-if="isSubmitting" :overlay="true" />
      <div class="app-brand justify-content-center d-none">
        <router-link :to="{ name: 'admin:home' }" class="app-brand-link gap-2">
          <img src="/assets/img/app/app-icon-dark.svg" alt="" />
        </router-link>
      </div>
      <h4 class="mb-1">Welcome! 👋</h4>
      <p class="mb-6">Please sign-in to your account</p>

      <vue-form @submit="handleSubmit" class="mb-6">
        <alert-message v-if="alertMessage" :params="alertMessage" />
        <div class="mb-6">
          <vue-form-input name="username" type="text" label="Email" placeholder="Email" />
        </div>
        <div class="mb-6 form-password-toggle">
          <vue-form-input name="password" type="password" label="Password" placeholder="Password" />
        </div>
        <div class="mb-8">
          <div class="d-flex justify-content-between">
            <div class="form-check mb-0">
              <input class="form-check-input" type="checkbox" v-model="rememberMe" />
              <label class="form-check-label" for="remember-me"> Remember Me </label>
            </div>
            <a href="https://fiverr.com/inbox/ndoubismarck" target="_blank">
              <span>Forgot Password?</span>
            </a>
          </div>
        </div>
        <div class="mb-6">
          <button class="btn btn-primary d-grid w-100" type="submit">Login</button>
        </div>
      </vue-form>

      <p class="text-center">
        <span>Need help?</span>
        <span>&nbsp;</span>
        <a href="https://fiverr.com/inbox/ndoubismarck" target="_blank">
          <span>Get In Touch</span>
        </a>
      </p>
    </div>
  </div>
</template>
