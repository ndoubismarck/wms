<script setup lang="ts">
import {computed, onMounted, ref} from 'vue'
import {useApp} from '@/app/app'
import VueForm, {type VueFormData} from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import VueFormSelectInput, {
  type ISelectInputOption
} from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueFormDatepickerInput from '@/views/shared/components/VueForm/VueFormDatepickerInput.vue'
import VueFormMultiSelectInput from '@/views/shared/components/VueForm/VueFormMultiSelectInput.vue'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import {EAlertMessageType, type IAlertMessage} from '@/app/types.ts'
import type {TaskModel} from '@/app/models/task_model.ts'
import type {UserModel} from '@/app/models/user_model.ts'
import type {TeamModel} from '@/app/models/team_model.ts'

const app = useApp()

const emits = defineEmits<{
  (e: 'close'): void
  (e: 'submit'): void
  (e: 'submitted', success: boolean, payload?: TaskModel): void
}>()

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isLoading = ref<boolean>(true)
const isSubmitting = ref<boolean>(false)

const users = ref<UserModel[]>([])
const teams = ref<TeamModel[]>([])

const priorityOptions = ref<ISelectInputOption[]>([
  {label: 'Low', value: 'low'},
  {label: 'Normal', value: 'normal'},
  {label: 'High', value: 'high'},
  {label: 'Urgent', value: 'urgent'},
])

const userOptions = computed<ISelectInputOption[]>(() => {
  return users.value.map((user) => ({
    label: user.fullName,
    value: user.id,
  }))
})

const teamOptions = computed<ISelectInputOption[]>(() => {
  return teams.value.map((team) => ({
    label: team.name,
    value: team.id,
  }))
})

const normalizeIDs = (value: any): string[] => {
  if (!value) {
    return []
  }
  if (Array.isArray(value)) {
    return Array.from(new Set(value.map((item) => `${item}`.trim()).filter((item) => item.length > 0)))
  }
  const parsed = `${value}`.trim()
  return parsed.length > 0 ? [parsed] : []
}

const loadUsers = async () => {
  const result = await app.services.users.getMany({
    page: 1,
    limit: 200,
    order: 'created_at.desc',
  })
  if (result.success) {
    users.value = result.users
  }
}

const loadTeams = async () => {
  const result = await app.services.teams.getMany({
    page: 1,
    limit: 200,
    order: 'created_at.desc',
  })
  if (result.success) {
    teams.value = result.teams
  }
}

const reset = () => {
  form.value?.reset()
  alert.value = null
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')

  isSubmitting.value = true
  let success = false
  let payload: TaskModel | undefined

  try {
    const dueAt = data.due_at ? new Date(data.due_at).toISOString() : undefined

    const result = await app.services.tasks.add({
      title: `${data.title ?? ''}`.trim(),
      description: `${data.description ?? ''}`.trim(),
      priority: `${data.priority ?? 'normal'}`,
      due_at: dueAt,
      user_ids: normalizeIDs(data.user_ids),
      team_ids: normalizeIDs(data.team_ids),
    })

    if (result.success) {
      success = true
      payload = result.task
      reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Task successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Task could not be added.'),
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

onMounted(async () => {
  isLoading.value = true
  try {
    await Promise.all([loadUsers(), loadTeams()])
  } catch (ex) {
    app.services.logger.error(ex)
  }
  isLoading.value = false
})

defineExpose({
  reset,
  isLoading,
  isSubmitting,
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert"/>

        <div class="mb-6">
          <vue-form-input
            name="title"
            type="text"
            label="Title"
            placeholder="Pick items for outgoing order"
            validation="required"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="mb-6">
          <vue-form-textarea-input
            name="description"
            label="Description"
            placeholder="Task details"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="mb-6">
          <vue-form-select-input
            name="priority"
            label="Priority"
            placeholder="Select priority"
            :options="priorityOptions"
            :value="'normal'"
            validation="required"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="mb-6">
          <vue-form-datepicker-input
            name="due_at"
            label="Due Date or Time"
            mode="datetime-local"
            placeholder="Select due date and time"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="mb-6">
          <vue-form-multi-select-input
            name="user_ids"
            label="Assign Users"
            placeholder="Search users"
            :options="userOptions"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="mb-6">
          <vue-form-multi-select-input
            name="team_ids"
            label="Assign Teams"
            placeholder="Search teams"
            :options="teamOptions"
            :disabled="isLoading || isSubmitting"/>
        </div>

        <div class="d-flex gap-2">
          <button type="submit" :disabled="isLoading || isSubmitting" class="btn btn-primary">
            Submit
          </button>
        </div>
      </vue-form>
    </div>
  </div>
</template>
