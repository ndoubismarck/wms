<script setup lang="ts">
import { computed, ref } from 'vue'
import { useApp } from '@/app/app.ts'
import { type IAlertMessage } from '@/app/types.ts'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueForm, { type VueFormActions, type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { useRouter } from 'vue-router'
import VueFormButton from '@/views/shared/components/VueForm/VueFormButton.vue'
import { useSetupStore } from '@/stores/setup_store.ts'

const app = useApp()
const router = useRouter()

const isSubmitting = ref<boolean>(false)
const alertMessage = ref<IAlertMessage | null>(null)

const setup = computed(() => useSetupStore().get().value)

const handleSubmit = async (data: VueFormData, actions: VueFormActions) => {}
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
      <p class="mb-6">Please fill out the form below to get started</p>

      <vue-form @submit="handleSubmit" class="mb-6">
        <alert-message v-if="alertMessage" :params="alertMessage" />
        <div class="mb-6">
          <vue-form-input name="username" type="text" label="Email" placeholder="Email" />
        </div>
        <div class="mb-6 form-password-toggle">
          <vue-form-input name="password" type="password" label="Password" placeholder="Password" />
        </div>
        <div class="mb-6">
          <vue-form-button text="Get Started" class="btn-primary d-grid w-100" />
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
