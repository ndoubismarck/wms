<script setup lang="ts">
import { ref } from 'vue'
import { useApp } from '@/app/app'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import type { UserModel } from '@/app/models/user_model.ts'

const app = useApp()

const emits = defineEmits<{
  (e: 'close'): void
  (e: 'submit'): void
  (e: 'submitted', success: boolean, payload?: UserModel): void
}>()

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isSubmitting = ref<boolean>(false)

const roleOptions = ref<ISelectInputOption[]>([
  {
    label: 'Admin',
    value: 'admin',
  },
  {
    label: 'Editor',
    value: 'editor',
  },
  {
    label: 'Subscriber',
    value: 'subscriber',
  },
])

const reset = () => {
  form.value?.reset()
  alert.value = null
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')
  isSubmitting.value = true

  let success = false
  let payload: UserModel | undefined

  try {
    const result = await app.services.users.add({
      first_name: `${data.first_name ?? ''}`.trim(),
      last_name: `${data.last_name ?? ''}`.trim(),
      email_address: `${data.email_address ?? ''}`.trim(),
      phone_number: `${data.phone_number ?? ''}`.trim(),
      password: `${data.password ?? ''}`,
      role: `${data.role ?? 'subscriber'}`,
    })

    if (result.success) {
      success = true
      payload = result.user
      reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Staff user successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Staff user could not be added.'),
        },
      }
    }
  } catch (ex) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.messageFromError(ex),
      },
    }
    app.services.logger.error(ex)
  }

  isSubmitting.value = false
  emits('submitted', success, payload)

  if (success) {
    await app.helpers.async.sleep(300)
    emits('close')
  }
}

defineExpose({
  reset,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert" />

        <div class="mb-6">
          <vue-form-input
            name="first_name"
            type="text"
            label="First Name"
            placeholder="John"
            validation="required"
            :disabled="isSubmitting" />
        </div>

        <div class="mb-6">
          <vue-form-input
            name="last_name"
            type="text"
            label="Last Name"
            placeholder="Doe"
            validation="required"
            :disabled="isSubmitting" />
        </div>

        <div class="mb-6">
          <vue-form-input
            name="email_address"
            type="email"
            label="Email Address"
            placeholder="john.doe@example.com"
            validation="required|email"
            :disabled="isSubmitting" />
        </div>

        <div class="mb-6">
          <vue-form-input
            name="phone_number"
            type="tel"
            label="Phone Number"
            placeholder="+1 555 000 0000"
            :disabled="isSubmitting" />
        </div>

        <div class="mb-6">
          <vue-form-select-input
            :options="roleOptions"
            :disabled="isSubmitting"
            :value="'subscriber'"
            name="role"
            label="Role"
            placeholder="Select Role"
            validation="required" />
        </div>

        <div class="mb-6">
          <vue-form-input
            name="password"
            type="password"
            label="Password"
            placeholder="Enter Password"
            validation="required"
            :disabled="isSubmitting" />
        </div>

        <button type="submit" :disabled="isSubmitting" class="btn btn-primary">Submit</button>
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
