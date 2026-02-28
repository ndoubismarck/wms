<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useApp } from '@/app/app'
import type { TeamModel } from '@/app/models/team_model'
import type { UserModel } from '@/app/models/user_model'

const app = useApp()

const isLoading = ref<boolean>(true)
const isSubmitting = ref<boolean>(false)

const teams = ref<TeamModel[]>([])
const users = ref<UserModel[]>([])

const form = ref<{
  name: string
  description: string
  userIds: string[]
}>({
  name: '',
  description: '',
  userIds: [],
})

const selectedTeamId = ref<string>('')
const selectedTeamUserIds = ref<string[]>([])

const selectedTeam = computed(() => teams.value.find((team) => team.id == selectedTeamId.value) ?? null)

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
  const result = await app.services.users.getMany({
    page: 1,
    limit: 200,
    order: 'created_at.desc',
  })
  if (result.success) {
    users.value = result.users
  }
}

const refresh = async () => {
  await Promise.all([loadTeams(), loadUsers()])
}

const resetForm = () => {
  form.value = {
    name: '',
    description: '',
    userIds: [],
  }
}

const createTeam = async () => {
  if (isSubmitting.value) {
    return
  }
  const name = form.value.name.trim()
  if (!name) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.teams.add({
      name: name,
      description: form.value.description,
      user_ids: form.value.userIds,
    })
    if (result.success) {
      resetForm()
      await loadTeams()
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
      if (selectedTeamId.value == teamId) {
        selectedTeamId.value = ''
      }
      await loadTeams()
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const saveSelectedTeamAssignments = async () => {
  if (!selectedTeamId.value || isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.teams.update(selectedTeamId.value, {
      user_ids: selectedTeamUserIds.value,
    })
    if (result.success) {
      await loadTeams()
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

watch(selectedTeamId, () => {
  const team = selectedTeam.value
  selectedTeamUserIds.value = team ? [...team.userIds] : []
})

onMounted(async () => {
  isLoading.value = true
  try {
    await refresh()
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="row g-4">
    <div class="col-12 col-lg-5">
      <div class="card">
        <div class="card-header">
          <h5 class="mb-0">Create Team</h5>
        </div>
        <div class="card-body">
          <div class="mb-3">
            <label class="form-label">Team Name</label>
            <input v-model="form.name" type="text" class="form-control" placeholder="Receiving Team" />
          </div>
          <div class="mb-3">
            <label class="form-label">Description</label>
            <textarea
              v-model="form.description"
              class="form-control"
              rows="3"
              placeholder="What this team handles"></textarea>
          </div>
          <div class="mb-3">
            <label class="form-label">Members</label>
            <select v-model="form.userIds" class="form-select" size="6" multiple>
              <option v-for="user in users" :key="user.id" :value="user.id">{{ user.fullName }}</option>
            </select>
          </div>
          <button class="btn btn-primary" :disabled="isSubmitting" @click="createTeam">Create Team</button>
        </div>
      </div>
    </div>

    <div class="col-12 col-lg-7">
      <div class="card mb-4">
        <div class="card-header">
          <h5 class="mb-0">Manage Team Assignments</h5>
        </div>
        <div class="card-body">
          <div class="mb-3">
            <label class="form-label">Team</label>
            <select v-model="selectedTeamId" class="form-select">
              <option value="">Select team</option>
              <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
            </select>
          </div>
          <template v-if="selectedTeam">
            <div class="mb-3">
              <label class="form-label">Members</label>
              <select v-model="selectedTeamUserIds" class="form-select" size="6" multiple>
                <option v-for="user in users" :key="user.id" :value="user.id">{{ user.fullName }}</option>
              </select>
            </div>
            <button class="btn btn-outline-primary" :disabled="isSubmitting" @click="saveSelectedTeamAssignments">
              Save Team Assignments
            </button>
          </template>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Teams</h5>
          <span class="badge bg-label-primary">{{ teams.length }}</span>
        </div>
        <div class="table-responsive">
          <table class="table mb-0">
            <thead>
              <tr>
                <th>Name</th>
                <th>Members</th>
                <th class="text-end">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="isLoading">
                <td colspan="3" class="text-center py-4">Loading teams...</td>
              </tr>
              <tr v-else-if="teams.length == 0">
                <td colspan="3" class="text-center py-4">No teams found.</td>
              </tr>
              <tr v-for="team in teams" :key="team.id">
                <td>
                  <div class="fw-semibold">{{ team.name }}</div>
                  <small class="text-body-secondary">{{ team.description || 'No description' }}</small>
                </td>
                <td>{{ team.userIds.length }}</td>
                <td class="text-end">
                  <button class="btn btn-sm btn-outline-danger" :disabled="isSubmitting" @click="deleteTeam(team.id)">
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
