<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useApp } from '@/app/app'
import { useCustomersStore } from '@/stores/customers_store'
import type { CustomerModel, CustomerStatus, CustomerTier } from '@/app/models/customer_model'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueFormDatepickerInput from '@/views/shared/components/VueForm/VueFormDatepickerInput.vue'
import { DateTime } from '@/app/core/date_time'

defineOptions({
  name: 'AdminCustomersPage',
})

const app = useApp()
const customersStore = useCustomersStore()

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const selectedTier = ref<'all' | CustomerTier>('all')
const selectedStatus = ref<'all' | CustomerStatus>('all')

const customerOffcanvas = ref<IOffcanvas | null>(null)
const customerFormMode = ref<'create' | 'edit'>('create')
const editingCustomer = ref<CustomerModel | null>(null)
const customerFormKey = ref<number>(0)

const customers = computed(() => customersStore.get().value)
const summary = computed(() => customersStore.getSummary().value)

const tierOptions: Array<{ label: string; value: 'all' | CustomerTier }> = [
  { label: 'All Tiers', value: 'all' },
  { label: 'Standard', value: 'standard' },
  { label: 'Growth', value: 'growth' },
  { label: 'Enterprise', value: 'enterprise' },
]

const statusOptions: Array<{ label: string; value: 'all' | CustomerStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Active', value: 'active' },
  { label: 'Paused', value: 'paused' },
  { label: 'Churn Risk', value: 'churn_risk' },
]

const tierSelectOptions = computed<ISelectInputOption[]>(() => {
  return tierOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const statusSelectOptions = computed<ISelectInputOption[]>(() => {
  return statusOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const customerFormInitial = computed(() => {
  if (!editingCustomer.value || customerFormMode.value == 'create') {
    return {
      name: '',
      email: '',
      phone: '',
      city: '',
      tier: 'standard',
      status: 'active',
      total_orders: '0',
      lifetime_value: '0',
      last_order_at: '',
    }
  }
  return {
    name: editingCustomer.value.name,
    email: editingCustomer.value.email,
    phone: editingCustomer.value.phone,
    city: editingCustomer.value.city,
    tier: editingCustomer.value.tier,
    status: editingCustomer.value.status,
    total_orders: `${editingCustomer.value.totalOrders}`,
    lifetime_value: `${editingCustomer.value.lifetimeValue}`,
    last_order_at: toDateTimeLocal(editingCustomer.value.lastOrderAt),
  }
})

const tierBadgeClass = (tier: CustomerTier): string => {
  if (tier == 'enterprise') {
    return 'bg-primary-subtle text-primary'
  }
  if (tier == 'growth') {
    return 'bg-info-subtle text-info'
  }
  return 'bg-secondary-subtle text-secondary'
}

const statusBadgeClass = (status: CustomerStatus): string => {
  if (status == 'active') {
    return 'bg-success-subtle text-success'
  }
  if (status == 'paused') {
    return 'bg-warning-subtle text-warning'
  }
  return 'bg-danger-subtle text-danger'
}

const renderTier = (tier: CustomerTier): string => {
  return `<span class="badge text-capitalize ${tierBadgeClass(tier)}">${tier}</span>`
}

const renderStatus = (status: CustomerStatus): string => {
  return `<span class="badge text-capitalize ${statusBadgeClass(status)}">${status.replace('_', ' ')}</span>`
}

const lastOrderText = (row: CustomerModel): string => {
  return row.lastOrderAt ? new DateTime(row.lastOrderAt).toDateTimeString() : 'N/A'
}

const getCustomersTableColumns = (row: CustomerModel): ITableColumn[] => {
  return [
    {
      name: 'Customer',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Contact',
      value: {
        text: `${row.email} | ${row.phone || 'N/A'}`,
      },
    },
    {
      name: 'Tier',
      value: {
        html: renderTier(row.tier),
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Orders',
      value: {
        text: `${row.totalOrders}`,
      },
    },
    {
      name: 'LTV',
      value: {
        text: `$${row.lifetimeValue.toFixed(2)}`,
      },
    },
    {
      name: 'Last Order',
      value: {
        text: lastOrderText(row),
      },
    },
    {
      name: 'City',
      value: {
        text: row.city,
      },
    },
  ]
}

const toDateTimeLocal = (value: string): string => {
  if (!value) {
    return ''
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 16)
}

const openCreateCustomer = () => {
  customerFormMode.value = 'create'
  editingCustomer.value = null
  customerFormKey.value++
  customerOffcanvas.value?.show()
}

const openEditCustomer = (row: CustomerModel) => {
  customerFormMode.value = 'edit'
  editingCustomer.value = row
  customerFormKey.value++
  customerOffcanvas.value?.show()
}

const customerFormTitle = computed(() => {
  return customerFormMode.value == 'create' ? 'Create Customer' : 'Edit Customer'
})

const submitCustomerForm = async (data: VueFormData) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const payload: Record<string, unknown> = {
      name: `${data.name ?? ''}`.trim(),
      email: `${data.email ?? ''}`.trim(),
      phone: `${data.phone ?? ''}`.trim(),
      city: `${data.city ?? ''}`.trim(),
      tier: `${data.tier ?? 'standard'}`,
      status: `${data.status ?? 'active'}`,
      total_orders: Math.max(0, Number(data.total_orders ?? 0)),
      lifetime_value: Math.max(0, Number(data.lifetime_value ?? 0)),
    }
    if (data.last_order_at) {
      payload.last_order_at = new Date(`${data.last_order_at}`).toISOString()
    }

    if (customerFormMode.value == 'create') {
      const result = await app.services.operations.addCustomer(payload)
      if (!result.success) {
        return
      }
    } else if (editingCustomer.value) {
      const result = await app.services.operations.updateCustomer(editingCustomer.value.id, payload)
      if (!result.success) {
        return
      }
    }

    customerOffcanvas.value?.hide()
    await loadCustomers(false)
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteCustomer = async (row: CustomerModel) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.operations.deleteCustomerById(row.id)
    if (result.success) {
      await loadCustomers(false)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const loadCustomers = async (withLoader: boolean = false): Promise<void> => {
  if (isRefreshing.value) {
    return
  }
  if (withLoader) {
    isLoading.value = true
  }
  isRefreshing.value = true
  try {
    const payload: Record<string, unknown> = {
      page: 1,
      limit: 200,
      order: 'updated_at.desc',
    }
    if (selectedTier.value != 'all') {
      payload.tier = selectedTier.value
    }
    if (selectedStatus.value != 'all') {
      payload.status = selectedStatus.value
    }
    const result = await app.services.operations.getCustomers(payload)
    if (result.success) {
      customersStore.setAll(result.customers)
      customersStore.setSummary(result.summary)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isRefreshing.value = false
    if (withLoader) {
      isLoading.value = false
    }
  }
}

const resetFilters = () => {
  selectedTier.value = 'all'
  selectedStatus.value = 'all'
}

watch([selectedTier, selectedStatus], async () => {
  await loadCustomers(false)
})

onMounted(async () => {
  await loadCustomers(true)
})
</script>

<template>
  <div class="operations-customers-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Customers</h5>
          <p class="text-body-secondary mb-0">Manage customer accounts with create/edit backend records.</p>
        </div>
        <div class="d-flex gap-2 mt-4 mt-md-0">
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="resetFilters">
            <i class="bx bx-reset me-2"></i>
            Reset Filters
          </button>
          <button class="btn btn-primary" :disabled="isSubmitting" @click="openCreateCustomer">
            <i class="bx bx-plus me-2"></i>
            Create Customer
          </button>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="card">
      <div class="card-body position-relative h-px-300">
        <preloader :overlay="true" />
      </div>
    </div>

    <template v-else>
      <div class="row">
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Customers</p>
              <h4 class="mb-0">{{ summary.totalCustomers }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Active</p>
              <h4 class="mb-0 text-success">{{ summary.activeCount }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Enterprise</p>
              <h4 class="mb-0 text-primary">{{ summary.enterpriseCount }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Avg LTV</p>
              <h4 class="mb-0">${{ summary.avgLtv.toFixed(2) }}</h4>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex flex-column flex-lg-row justify-content-between align-items-lg-center gap-3">
          <h5 class="mb-0">Customer Directory</h5>
          <div class="d-flex flex-column flex-md-row align-items-md-center gap-2">
            <div class="d-flex align-items-center gap-2">
              <label class="form-label mb-0 text-body-secondary">Tier</label>
              <select v-model="selectedTier" class="form-select">
                <option v-for="option in tierOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
            <div class="d-flex align-items-center gap-2">
              <label class="form-label mb-0 text-body-secondary">Status</label>
              <select v-model="selectedStatus" class="form-select">
                <option v-for="option in statusOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="customers" :columns="getCustomersTableColumns" :checkboxes="true">
            <template #actions="{ row }">
              <div class="d-flex gap-2">
                <button type="button" class="btn btn-sm btn-outline-primary" @click="openEditCustomer(row)">Edit</button>
                <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteCustomer(row)">Delete</button>
              </div>
            </template>
            <template #noResults>
              <p class="py-6">No customers found for selected filters.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>

  <offcanvas ref="customerOffcanvas" :show="false" :overflow="isSubmitting" :title="customerFormTitle">
    <template #body>
      <div class="card">
        <div class="card-body">
          <vue-form :key="customerFormKey" @submit="submitCustomerForm">
            <div class="mb-6">
              <vue-form-input
                name="name"
                type="text"
                label="Name"
                placeholder="Northfield Retail"
                validation="required"
                :value="customerFormInitial.name"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="email"
                type="email"
                label="Email"
                placeholder="ops@northfieldretail.com"
                validation="required|email"
                :value="customerFormInitial.email"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="phone"
                type="text"
                label="Phone"
                placeholder="+370 612 33001"
                validation="required"
                :value="customerFormInitial.phone"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="city"
                type="text"
                label="City"
                placeholder="Vilnius"
                validation="required"
                :value="customerFormInitial.city"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="tier"
                label="Tier"
                placeholder="Select tier"
                :options="tierSelectOptions"
                :value="customerFormInitial.tier"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="status"
                label="Status"
                placeholder="Select status"
                :options="statusSelectOptions"
                :value="customerFormInitial.status"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="total_orders"
                type="number"
                label="Total Orders"
                placeholder="0"
                :value="customerFormInitial.total_orders"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="lifetime_value"
                type="number"
                label="Lifetime Value"
                placeholder="0"
                :value="customerFormInitial.lifetime_value"
                input-mode="decimal"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-datepicker-input
                name="last_order_at"
                label="Last Order At"
                mode="datetime-local"
                placeholder="Select date"
                :value="customerFormInitial.last_order_at"
                :disabled="isSubmitting" />
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              {{ customerFormMode == 'create' ? 'Create Customer' : 'Update Customer' }}
            </button>
          </vue-form>
        </div>
      </div>
    </template>
  </offcanvas>
</template>

<style scoped></style>
