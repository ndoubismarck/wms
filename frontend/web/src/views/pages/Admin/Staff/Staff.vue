<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import { useUsersStore } from '@/stores/users_store.ts'
import { useApp } from '@/app/app.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import type { UserModel } from '@/app/models/user_model.ts'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'AdminStaffPage',
})

const app = useApp()
const router = useRouter()
const usersStore = useUsersStore()

const users = computed(() => usersStore.get().value)

const page = ref<number>(1)
const isLoading = ref<boolean>(false)
const isFetching = ref<boolean>(false)
const checkedUserIds = ref<string[]>([])

const getUsers = async () => {
  try {
    const result = await app.services.users.getMany({
      page: page.value,
      limit: 25,
    })
    if (result.success) {
      if (result.pagination.hasNext) {
        page.value = result.pagination.next
      }
      usersStore.set(result.users)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const renderRole = (value?: string): string => {
  if (value == 'admin') {
    return `<span class="badge text-success bg-success-subtle text-capitalize">${value}</span>`
  }
  if (value == 'editor') {
    return `<span class="badge text-warning bg-warning-subtle text-capitalize">${value}</span>`
  }
  return `<span class="badge text-secondary bg-secondary-subtle text-capitalize">${value}</span>`
}

const getTableColumns = (row: UserModel): ITableColumn[] => {
  return [
    {
      name: 'Name',
      value: {
        text: row.fullName,
      },
    },
    {
      name: 'Email',
      value: {
        text: row.emailAddress,
      },
    },
    {
      name: 'Role',
      value: {
        html: renderRole(row.role),
      },
    },
    {
      name: 'Date Updated',
      value: {
        text: row.updatedAt.toDateTimeString(),
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
  checkedUserIds.value = ids
}

const handleTableScrollEnd = async () => {
  if (!isFetching.value) {
    isFetching.value = true
    await getUsers()
    await app.helpers.async.sleep(3000)
    isFetching.value = false
  }
}

onMounted(async () => {
  isLoading.value = users.value.length == 0
  await getUsers()
  await app.helpers.async.sleep(3000)
  isLoading.value = false
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
      <div class="card-header d-flex justify-content-start">
        <button
          style="height: 40px !important"
          class="btn btn-outline-primary d-flex justify-content-center"
          :disabled="isFetching"
          @click="router.push({ name: 'admin:staff:add' })">
          <i class="bx bx-plus" />
          <span class="d-none d-sm-inline ms-2">Add Staff User</span>
        </button>

        <button
          style="height: 40px !important"
          class="btn btn-outline-primary d-flex justify-content-center ms-4"
          :disabled="isFetching"
          @click="handleTableScrollEnd">
          <i class="bx bx-refresh" />
          <span class="d-none d-sm-inline ms-2">Refresh</span>
        </button>

        <div class="dropdown ms-4">
          <button
            style="height: 40px !important"
            class="btn btn-outline-primary d-flex justify-content-center"
            data-bs-toggle="dropdown"
            data-bs-boundary="viewport"
            aria-expanded="false"
            :disabled="checkedUserIds.length == 0">
            <span class="d-none d-sm-inline me-2">With Selected</span>
            <i class="bx bx-chevron-down" />
          </button>
          <ul class="dropdown-menu">
            <li>
              <a href="#" class="dropdown-item text-primary">
                <div class="d-flex justify-content-start"><i class="bx bx-download me-2"></i><span>Export</span></div>
              </a>
            </li>
            <li>
              <a href="#" class="dropdown-item text-danger">
                <div class="d-flex justify-content-start"><i class="bx bx-trash me-2"></i><span>Delete</span></div>
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div class="card-body p-0" style="min-height: 300px">
        <vue-table
          :rows="users"
          :columns="getTableColumns"
          @checked="handleTableChecked"
          @scrollend="handleTableScrollEnd"
          :checkboxes="true">
          <template #actions="{ row }">
            <button type="button" class="btn btn-lg text-primary border-0" :title="row.emailAddress">
              <i class="bx bx-dots-vertical"></i>
            </button>
          </template>
          <template #preloader>
            <div v-if="isFetching" class="position-relative h-auto py-6 my-6">
              <preloader :overlay="true" />
            </div>
          </template>
        </vue-table>
      </div>
    </template>
  </div>
  <router-view />
</template>

<style scoped></style>
