<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useApp } from '@/app/app'
import type { TeamModel } from '@/app/models/team_model'
import type { UserModel } from '@/app/models/user_model'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import { EPreloaderSize } from '@/app/types.ts'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import VueFormMultiSelectInput from '@/views/shared/components/VueForm/VueFormMultiSelectInput.vue'
import type { ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminTeamsPage',
})

const app = useApp()

const isLoading = ref<boolean>(true)
const isFetching = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const isLoadingUsers = ref<boolean>(false)

const teams = ref<TeamModel[]>([])
const users = ref<UserModel[]>([])

const createTeamOffcanvas = ref<IOffcanvas | null>(null)
const createTeamForm = ref<typeof VueForm | null>(null)

const userOptions = computed<ISelectInputOption[]>(() => {
  return users.value.map((user) => ({
    label: user.fullName,
    value: user.id,
  }))
})

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

const loadUsers = async () => {
  isLoadingUsers.value = true
  try {
    const result = await app.services.users.getMany({
      page: 1,
      limit: 200,
      order: 'created_at.desc',
    })
    if (result.success) {
      users.value = result.users
    }
  } finally {
    isLoadingUsers.value = false
  }
}

const refresh = async () => {
  if (isFetching.value) {
    return
  }
  isFetching.value = true
  try {
    await loadTeams()
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isFetching.value = false
  }
}

const resetForm = () => {
  createTeamForm.value?.reset()
}

const openCreateTeam = async () => {
  if (users.value.length == 0 && !isLoadingUsers.value) {
    try {
      await loadUsers()
    } catch (ex) {
      app.services.logger.error(ex)
    }
  }
  createTeamOffcanvas.value?.show()
}

const handleHideCreateTeam = () => {
  resetForm()
}

const normalizeIDs = (value: unknown): string[] => {
  if (!value) {
    return []
  }
  if (Array.isArray(value)) {
    return Array.from(new Set(value.map((item) => `${item}`.trim()).filter((item) => item.length > 0)))
  }
  const parsed = `${value}`.trim()
  return parsed.length > 0 ? [parsed] : []
}

const getTableColumns = (row: TeamModel): ITableColumn[] => {
  return [
    {
      name: 'Name',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Description',
      value: {
        text: row.description || 'No description',
      },
    },
    {
      name: 'Members',
      value: {
        text: `${row.userIds.length}`,
      },
    },
  ]
}

const createTeam = async (data: VueFormData) => {
  if (isSubmitting.value) {
    return
  }
  const name = `${data.name ?? ''}`.trim()
  if (!name) {
    return
  }

  isSubmitting.value = true
  try {
    const result = await app.services.teams.add({
      name,
      description: `${data.description ?? ''}`.trim(),
      user_ids: normalizeIDs(data.user_ids),
    })
    if (result.success) {
      resetForm()
      await loadTeams()
      createTeamOffcanvas.value?.hide()
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteTeam = async (teamId: string) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.teams.deleteById(teamId)
    if (result.success) {
      await loadTeams()
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    await loadTeams()
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="card">
    <template v-if="isLoading">
      <div class="card-body position-relative h-px-300">
        <preloader :overlay="true" />
      </div>
    </template>
    <template v-else>
      <div class="card-header d-flex justify-content-between align-items-center">
        <div class="d-flex justify-content-start">
          <button
            style="height: 40px !important"
            class="btn btn-outline-primary d-flex justify-content-center"
            :disabled="isSubmitting"
            @click="openCreateTeam">
            <i class="bx bx-plus" />
            <span class="d-none d-sm-inline ms-2">Create Team</span>
          </button>

          <button
            style="height: 40px !important"
            class="btn btn-outline-primary d-flex justify-content-center ms-4"
            :disabled="isFetching || isSubmitting"
            @click="refresh">
            <i class="bx bx-refresh" />
            <span class="d-none d-sm-inline ms-2">Refresh</span>
          </button>
        </div>

        <span class="badge bg-label-primary">{{ teams.length }}</span>
      </div>

      <div class="card-body p-0">
        <vue-table :rows="teams" :columns="getTableColumns" :checkboxes="true">
          <template #actions="{ row }">
            <button class="btn btn-sm btn-outline-danger" :disabled="isSubmitting" @click="deleteTeam(row.id)">
              Delete
            </button>
          </template>
          <template #noResults>
            <p class="py-6">No teams found.</p>
          </template>
        </vue-table>
      </div>
    </template>
  </div>

  <offcanvas
    ref="createTeamOffcanvas"
    :show="false"
    @hide="handleHideCreateTeam"
    :overflow="isSubmitting"
    title="Create Team">
    <template #body>
      <div class="card">
        <div class="card-body">
          <vue-form ref="createTeamForm" @submit="createTeam">
            <div class="mb-6">
              <vue-form-input
                name="name"
                type="text"
                label="Team Name"
                placeholder="Receiving Team"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-textarea-input
                name="description"
                label="Description"
                placeholder="What this team handles"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-multi-select-input
                name="user_ids"
                label="Members"
                placeholder="Search users"
                :options="userOptions"
                :disabled="isSubmitting || isLoadingUsers" />
              <small v-if="isLoadingUsers" class="text-body-secondary">Loading members...</small>
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">Create Team</button>
          </vue-form>
        </div>
      </div>
    </template>
    <template v-if="isSubmitting" #overlay>
      <preloader :size="EPreloaderSize.SM" :overlay="true" />
    </template>
  </offcanvas>
</template>
