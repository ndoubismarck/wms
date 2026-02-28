<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useApp } from '@/app/app'
import { useSuppliersStore } from '@/stores/suppliers_store'
import type { SupplierModel, SupplierStatus, SupplierType } from '@/app/models/supplier_model'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'

defineOptions({
  name: 'AdminSuppliersPage',
})

const app = useApp()
const suppliersStore = useSuppliersStore()

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const selectedStatus = ref<'all' | SupplierStatus>('all')

const supplierOffcanvas = ref<IOffcanvas | null>(null)
const supplierFormMode = ref<'create' | 'edit'>('create')
const editingSupplier = ref<SupplierModel | null>(null)
const supplierFormKey = ref<number>(0)

const suppliers = computed(() => suppliersStore.get().value)
const summary = computed(() => suppliersStore.getSummary().value)

const statusOptions: Array<{ label: string; value: 'all' | SupplierStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Approved', value: 'approved' },
  { label: 'Probation', value: 'probation' },
  { label: 'Blocked', value: 'blocked' },
]

const typeSelectOptions: ISelectInputOption[] = [
  { label: 'Raw Material', value: 'raw_material' },
  { label: 'Packaging', value: 'packaging' },
  { label: 'Finished Goods', value: 'finished_goods' },
  { label: 'Service', value: 'service' },
]

const statusSelectOptions = computed<ISelectInputOption[]>(() => {
  return statusOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const supplierFormInitial = computed(() => {
  if (!editingSupplier.value || supplierFormMode.value == 'create') {
    return {
      name: '',
      supplier_type: 'service',
      status: 'approved',
      contact_name: '',
      email: '',
      lead_time_days: '0',
      on_time_rate: '0',
      open_purchase_orders: '0',
      city: '',
    }
  }
  return {
    name: editingSupplier.value.name,
    supplier_type: editingSupplier.value.supplierType,
    status: editingSupplier.value.status,
    contact_name: editingSupplier.value.contactName,
    email: editingSupplier.value.email,
    lead_time_days: `${editingSupplier.value.leadTimeDays}`,
    on_time_rate: `${editingSupplier.value.onTimeRate}`,
    open_purchase_orders: `${editingSupplier.value.openPurchaseOrders}`,
    city: editingSupplier.value.city,
  }
})

const statusBadgeClass = (status: SupplierStatus): string => {
  if (status == 'approved') {
    return 'bg-success-subtle text-success'
  }
  if (status == 'probation') {
    return 'bg-warning-subtle text-warning'
  }
  return 'bg-danger-subtle text-danger'
}

const typeBadgeClass = (supplierType: SupplierType): string => {
  if (supplierType == 'raw_material') {
    return 'bg-primary-subtle text-primary'
  }
  if (supplierType == 'packaging') {
    return 'bg-info-subtle text-info'
  }
  if (supplierType == 'finished_goods') {
    return 'bg-label-secondary text-secondary'
  }
  return 'bg-label-dark text-dark'
}

const renderStatus = (status: SupplierStatus): string => {
  return `<span class="badge text-capitalize ${statusBadgeClass(status)}">${status}</span>`
}

const renderType = (supplierType: SupplierType): string => {
  return `<span class="badge text-capitalize ${typeBadgeClass(supplierType)}">${supplierType.replace('_', ' ')}</span>`
}

const getSuppliersTableColumns = (row: SupplierModel): ITableColumn[] => {
  return [
    {
      name: 'Supplier',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Contact',
      value: {
        text: `${row.contactName} | ${row.email}`,
      },
    },
    {
      name: 'Type',
      value: {
        html: renderType(row.supplierType),
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Lead Time',
      value: {
        text: `${row.leadTimeDays} days`,
      },
    },
    {
      name: 'On-Time Rate',
      value: {
        text: `${row.onTimeRate.toFixed(1)}%`,
      },
    },
    {
      name: 'Open POs',
      value: {
        text: `${row.openPurchaseOrders}`,
      },
    },
    {
      name: 'City',
      value: {
        text: row.city,
      },
    },
    {
      name: 'Updated',
      value: {
        text: row.updatedAt.toDateTimeString(),
      },
    },
  ]
}

const openCreateSupplier = () => {
  supplierFormMode.value = 'create'
  editingSupplier.value = null
  supplierFormKey.value++
  supplierOffcanvas.value?.show()
}

const openEditSupplier = (row: SupplierModel) => {
  supplierFormMode.value = 'edit'
  editingSupplier.value = row
  supplierFormKey.value++
  supplierOffcanvas.value?.show()
}

const supplierFormTitle = computed(() => {
  return supplierFormMode.value == 'create' ? 'Create Supplier' : 'Edit Supplier'
})

const submitSupplierForm = async (data: VueFormData) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const payload: Record<string, unknown> = {
      name: `${data.name ?? ''}`.trim(),
      supplier_type: `${data.supplier_type ?? 'service'}`,
      status: `${data.status ?? 'approved'}`,
      contact_name: `${data.contact_name ?? ''}`.trim(),
      email: `${data.email ?? ''}`.trim().toLowerCase(),
      lead_time_days: Math.max(0, Number(data.lead_time_days ?? 0)),
      on_time_rate: Math.min(100, Math.max(0, Number(data.on_time_rate ?? 0))),
      open_purchase_orders: Math.max(0, Number(data.open_purchase_orders ?? 0)),
      city: `${data.city ?? ''}`.trim(),
    }

    if (supplierFormMode.value == 'create') {
      const result = await app.services.operations.addSupplier(payload)
      if (!result.success) {
        return
      }
    } else if (editingSupplier.value) {
      const result = await app.services.operations.updateSupplier(editingSupplier.value.id, payload)
      if (!result.success) {
        return
      }
    }

    supplierOffcanvas.value?.hide()
    await loadSuppliers(false)
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteSupplier = async (row: SupplierModel) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.operations.deleteSupplierById(row.id)
    if (result.success) {
      await loadSuppliers(false)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const loadSuppliers = async (withLoader: boolean = false): Promise<void> => {
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
    if (selectedStatus.value != 'all') {
      payload.status = selectedStatus.value
    }
    const result = await app.services.operations.getSuppliers(payload)
    if (result.success) {
      suppliersStore.setAll(result.suppliers)
      suppliersStore.setSummary(result.summary)
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
  selectedStatus.value = 'all'
}

watch(selectedStatus, async () => {
  await loadSuppliers(false)
})

onMounted(async () => {
  await loadSuppliers(true)
})
</script>

<template>
  <div class="operations-suppliers-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Suppliers</h5>
          <p class="text-body-secondary mb-0">Manage suppliers with create/edit backend records.</p>
        </div>
        <div class="d-flex gap-2 mt-4 mt-md-0">
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="resetFilters">
            <i class="bx bx-reset me-2"></i>
            Reset Filters
          </button>
          <button class="btn btn-primary" :disabled="isSubmitting" @click="openCreateSupplier">
            <i class="bx bx-plus me-2"></i>
            Create Supplier
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
              <p class="text-body-secondary mb-1">Suppliers</p>
              <h4 class="mb-0">{{ summary.totalSuppliers }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Approved</p>
              <h4 class="mb-0 text-success">{{ summary.approvedCount }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Avg Lead Time</p>
              <h4 class="mb-0">{{ summary.avgLeadTime.toFixed(1) }} days</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Open POs</p>
              <h4 class="mb-0 text-primary">{{ summary.totalOpenPOs }}</h4>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
          <h5 class="mb-0">Supplier Register</h5>
          <div class="d-flex align-items-center gap-2">
            <label class="form-label mb-0 text-body-secondary">Status</label>
            <select v-model="selectedStatus" class="form-select">
              <option v-for="option in statusOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="suppliers" :columns="getSuppliersTableColumns" :checkboxes="true">
            <template #actions="{ row }">
              <div class="d-flex gap-2">
                <button type="button" class="btn btn-sm btn-outline-primary" @click="openEditSupplier(row)">Edit</button>
                <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteSupplier(row)">Delete</button>
              </div>
            </template>
            <template #noResults>
              <p class="py-6">No suppliers found for selected filters.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>

  <offcanvas ref="supplierOffcanvas" :show="false" :overflow="isSubmitting" :title="supplierFormTitle">
    <template #body>
      <div class="card">
        <div class="card-body">
          <vue-form :key="supplierFormKey" @submit="submitSupplierForm">
            <div class="mb-6">
              <vue-form-input
                name="name"
                type="text"
                label="Name"
                placeholder="Baltic Packaging Group"
                validation="required"
                :value="supplierFormInitial.name"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="contact_name"
                type="text"
                label="Contact Name"
                placeholder="Mantas Petrauskas"
                validation="required"
                :value="supplierFormInitial.contact_name"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="email"
                type="email"
                label="Email"
                placeholder="contact@balticpackaging.com"
                validation="required|email"
                :value="supplierFormInitial.email"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="city"
                type="text"
                label="City"
                placeholder="Kaunas"
                validation="required"
                :value="supplierFormInitial.city"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="supplier_type"
                label="Supplier Type"
                placeholder="Select type"
                :options="typeSelectOptions"
                :value="supplierFormInitial.supplier_type"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="status"
                label="Status"
                placeholder="Select status"
                :options="statusSelectOptions"
                :value="supplierFormInitial.status"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="lead_time_days"
                type="number"
                label="Lead Time (Days)"
                placeholder="10"
                validation="required"
                :value="supplierFormInitial.lead_time_days"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="on_time_rate"
                type="number"
                label="On-Time Rate (%)"
                placeholder="95"
                input-mode="decimal"
                validation="required"
                :value="supplierFormInitial.on_time_rate"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="open_purchase_orders"
                type="number"
                label="Open Purchase Orders"
                placeholder="0"
                validation="required"
                :value="supplierFormInitial.open_purchase_orders"
                :disabled="isSubmitting" />
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              {{ supplierFormMode == 'create' ? 'Create Supplier' : 'Update Supplier' }}
            </button>
          </vue-form>
        </div>
      </div>
    </template>
  </offcanvas>
</template>

<style scoped></style>
