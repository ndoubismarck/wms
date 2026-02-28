<script setup lang="ts">
import {computed, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useApp} from '@/app/app'
import type {TaskModel} from '@/app/models/task_model'
import type {UserModel} from '@/app/models/user_model'
import type {TeamModel} from '@/app/models/team_model'
import VueTable, {type ITableColumn} from '@/views/shared/components/VueTable.vue'
import Preloader from '@/views/shared/components/Preloader.vue'

const app = useApp()
const route = useRoute()
const router = useRouter()

const tasks = ref<TaskModel[]>([])
const users = ref<UserModel[]>([])
const teams = ref<TeamModel[]>([])

const page = ref<number>(1)
const hasNext = ref<boolean>(true)

const isLoading = ref<boolean>(true)
const isFetching = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)

const checkedTaskIds = ref<string[]>([])
const showAssignTaskForm = ref<boolean>(false)

const assignTaskForm = ref<{
  userIds: string[]
  teamIds: string[]
}>({
  userIds: [],
  teamIds: [],
})

type TaskFilterKey = 'pending-active' | 'pending' | 'active' | 'all'

const taskFilters: Record<
  TaskFilterKey,
  {
    label: string
    statuses: string[] | null
  }
> = {
  'pending-active': {
    label: 'Pending & Active',
    statuses: ['pending', 'in_progress'],
  },
  pending: {
    label: 'Pending',
    statuses: ['pending'],
  },
  active: {
    label: 'Active',
    statuses: ['in_progress'],
  },
  all: {
    label: 'All',
    statuses: null,
  },
}

const activeFilter = ref<TaskFilterKey>('pending-active')

const filteredTasks = computed(() => {
  const statuses = taskFilters[activeFilter.value].statuses
  if (!statuses) {
    return tasks.value
  }
  return tasks.value.filter((task) => statuses.includes(task.status))
})

const activeFilterLabel = computed(() => taskFilters[activeFilter.value].label)

const mergeTasks = (items: TaskModel[]) => {
  const next = [...tasks.value]
  items.forEach((item) => {
    const index = next.findIndex((val) => val.id == item.id)
    if (index == -1) {
      next.push(item)
      return
    }
    next[index] = item
  })
  tasks.value = next.sort((a, b) => b.createdAt.toUnix() - a.createdAt.toUnix())
}

const getTasks = async (reset: boolean = false) => {
  if (reset) {
    page.value = 1
    hasNext.value = true
    tasks.value = []
    checkedTaskIds.value = []
  }
  if (!hasNext.value) {
    return
  }
  const result = await app.services.tasks.getMany({
    page: page.value,
    limit: 25,
    order: 'created_at.desc',
  })
  if (!result.success) {
    return
  }
  mergeTasks(result.tasks)
  hasNext.value = !!result.pagination?.hasNext
  if (result.pagination?.hasNext) {
    page.value = result.pagination.next
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

const refresh = async () => {
  await Promise.all([getTasks(true), loadUsers(), loadTeams()])
}

const openAssignTask = () => {
  if (checkedTaskIds.value.length == 0) {
    return
  }
  const selectedTasks = tasks.value.filter((task) => checkedTaskIds.value.includes(task.id))
  if (selectedTasks.length == 1) {
    assignTaskForm.value = {
      userIds: [...selectedTasks[0]?.userIds],
      teamIds: [...selectedTasks[0]?.teamIds],
    }
  } else {
    assignTaskForm.value = {
      userIds: [],
      teamIds: [],
    }
  }
  showAssignTaskForm.value = true
}

const openAddTask = async () => {
  showAssignTaskForm.value = false
  await router.push({
    name: 'admin:tasks:add',
  })
}

const setTaskFilter = async (filterKey: TaskFilterKey) => {
  activeFilter.value = filterKey
  await getTasks(true)
}

const assignTasks = async () => {
  if (isSubmitting.value || checkedTaskIds.value.length == 0) {
    return
  }
  isSubmitting.value = true
  try {
    for (const taskID of checkedTaskIds.value) {
      await app.services.tasks.update(taskID, {
        user_ids: assignTaskForm.value.userIds,
        team_ids: assignTaskForm.value.teamIds,
      })
    }
    showAssignTaskForm.value = false
    await getTasks(true)
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteTask = async (taskId: string) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.tasks.deleteById(taskId)
    if (result.success) {
      tasks.value = tasks.value.filter((task) => task.id != taskId)
      checkedTaskIds.value = checkedTaskIds.value.filter((id) => id != taskId)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const renderStatus = (value?: string): string => {
  if (value == 'pending') {
    return `<span class="badge text-warning bg-warning-subtle text-capitalize">pending</span>`
  }
  if (value == 'in_progress') {
    return `<span class="badge text-primary bg-primary-subtle text-capitalize">active</span>`
  }
  if (value == 'completed') {
    return `<span class="badge text-success bg-success-subtle text-capitalize">completed</span>`
  }
  return `<span class="badge text-danger bg-danger-subtle text-capitalize">cancelled</span>`
}

const renderPriority = (value?: string): string => {
  if (value == 'urgent') {
    return `<span class="badge text-danger bg-danger-subtle text-capitalize">urgent</span>`
  }
  if (value == 'high') {
    return `<span class="badge text-warning bg-warning-subtle text-capitalize">high</span>`
  }
  if (value == 'normal') {
    return `<span class="badge text-info bg-info-subtle text-capitalize">normal</span>`
  }
  return `<span class="badge text-secondary bg-secondary-subtle text-capitalize">low</span>`
}

const getTableColumns = (row: TaskModel): ITableColumn[] => {
  return [
    {
      name: 'Task',
      value: {
        text: row.title,
      },
    },
    {
      name: 'Description',
      value: {
        text: row.description || '-',
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Priority',
      value: {
        html: renderPriority(row.priority),
      },
    },
    {
      name: 'Assignees',
      value: {
        html: `<span class="badge bg-label-info me-1">Users ${row.userIds.length}</span><span class="badge bg-label-secondary">Teams ${row.teamIds.length}</span>`,
      },
    },
    {
      name: 'Due',
      value: {
        text: row.hasDueAt ? row.dueAt.toDateTimeString() : '-',
      },
    },
    {
      name: 'Date Created',
      value: {
        text: row.createdAt.toDateTimeString(),
      },
    },
  ]
}

const handleTableChecked = (ids: string[]) => {
  checkedTaskIds.value = ids
}

const handleTableScrollEnd = async () => {
  if (isFetching.value || !hasNext.value) {
    return
  }
  isFetching.value = true
  await getTasks(false)
  await app.helpers.async.sleep(500)
  isFetching.value = false
}

const handleAssignSingleTask = (task: TaskModel) => {
  checkedTaskIds.value = [task.id]
  openAssignTask()
}

watch(
  () => route.name,
  async (name, oldName) => {
    if (name == 'admin:tasks' && oldName == 'admin:tasks:add') {
      await getTasks(true)
    }
  },
)

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
  <div class="card">
    <template v-if="isLoading">
      <div class="card-body position-relative h-px-300">
        <preloader :overlay="true"/>
      </div>
    </template>
    <template v-else>
      <div
        class="card-header d-flex justify-content-start align-items-center flex-wrap gap-2 border-bottom">
        <button
          class="btn btn-outline-primary d-flex justify-content-center"
          :disabled="isSubmitting"
          @click="openAddTask">
          <i class="bx bx-plus"></i>
          <span class="d-none d-sm-inline ms-2">Add Task</span>
        </button>

        <button
          class="btn btn-outline-primary d-flex justify-content-center"
          :disabled="checkedTaskIds.length == 0 || isSubmitting"
          @click="openAssignTask">
          <i class="bx bx-user-plus"></i>
          <span class="d-none d-sm-inline ms-2">Assign Task</span>
        </button>

        <div class="dropdown">
          <button
            class="btn btn-outline-primary d-flex justify-content-center"
            data-bs-toggle="dropdown"
            data-bs-boundary="viewport"
            aria-expanded="false">
            <i class="bx bx-filter-alt"></i>
            <span class="d-none d-sm-inline ms-2">Filter: {{ activeFilterLabel }}</span>
          </button>
          <ul class="dropdown-menu dropdown-menu-end">
            <li v-for="(filter, key) in taskFilters" :key="key">
              <a href="#" class="dropdown-item"
                 @click.prevent="setTaskFilter(key as TaskFilterKey)">
                {{ filter.label }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div v-if="showAssignTaskForm" class="card-body border-bottom">
        <h6 class="mb-3">Assign Task</h6>
        <small class="text-body-secondary d-block mb-3">Selected tasks: {{
            checkedTaskIds.length
          }}</small>
        <div class="row g-3">
          <div class="col-12 col-md-6">
            <label class="form-label">Users</label>
            <select v-model="assignTaskForm.userIds" class="form-select" size="6" multiple>
              <option v-for="user in users" :key="user.id" :value="user.id">{{
                  user.fullName
                }}
              </option>
            </select>
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Teams</label>
            <select v-model="assignTaskForm.teamIds" class="form-select" size="6" multiple>
              <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
            </select>
          </div>
          <div class="col-12 d-flex gap-2">
            <button class="btn btn-primary" :disabled="isSubmitting" @click="assignTasks">Apply
              Assignment
            </button>
            <button class="btn btn-outline-secondary" :disabled="isSubmitting"
                    @click="showAssignTaskForm = false">Cancel
            </button>
          </div>
        </div>
      </div>

      <div class="card-body p-0">
        <vue-table
          :rows="filteredTasks"
          :columns="getTableColumns"
          :checkboxes="true"
          @checked="handleTableChecked"
          @scrollend="handleTableScrollEnd">
          <template #actions="{ row }">
            <div class="dropdown">
              <button class="btn btn-lg text-primary border-0" type="button"
                      data-bs-toggle="dropdown" data-bs-boundary="viewport">
                <i class="bx bx-dots-vertical"></i>
              </button>
              <ul class="dropdown-menu dropdown-menu-end">
                <li>
                  <a href="#" class="dropdown-item" @click.prevent="handleAssignSingleTask(row)">
                    Assign Task
                  </a>
                </li>
                <li>
                  <a href="#" class="dropdown-item text-danger" @click.prevent="deleteTask(row.id)">
                    Delete
                  </a>
                </li>
              </ul>
            </div>
          </template>
          <template #noResults>
            <div class="no-results d-flex justify-content-center align-items-center">
              <p class="py-6 mb-0">No tasks found for the selected filter.</p>
            </div>
          </template>
          <template #preloader>
            <div v-if="isFetching" class="position-relative h-auto py-6 my-6">
              <preloader :overlay="true"/>
            </div>
          </template>
        </vue-table>
      </div>
    </template>
  </div>
  <router-view/>
</template>
