<script setup lang="ts">
import { computed, ref } from 'vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminSuppliersPage',
})

type SupplierStatus = 'approved' | 'probation' | 'blocked'
type SupplierType = 'raw_material' | 'packaging' | 'finished_goods' | 'service'

interface SupplierRecord {
  id: string
  name: string
  supplierType: SupplierType
  status: SupplierStatus
  contactName: string
  email: string
  leadTimeDays: number
  onTimeRate: number
  openPurchaseOrders: number
  city: string
}

const selectedStatus = ref<'all' | SupplierStatus>('all')

const statusOptions: Array<{ label: string; value: 'all' | SupplierStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Approved', value: 'approved' },
  { label: 'Probation', value: 'probation' },
  { label: 'Blocked', value: 'blocked' },
]

const suppliers = ref<SupplierRecord[]>([
  {
    id: 'sup-001',
    name: 'Baltic Packaging UAB',
    supplierType: 'packaging',
    status: 'approved',
    contactName: 'Rasa Milte',
    email: 'orders@balticpack.lt',
    leadTimeDays: 6,
    onTimeRate: 97.4,
    openPurchaseOrders: 4,
    city: 'Vilnius',
  },
  {
    id: 'sup-002',
    name: 'Nord Components',
    supplierType: 'raw_material',
    status: 'approved',
    contactName: 'Jonas Petrauskas',
    email: 'supply@nordcomponents.eu',
    leadTimeDays: 9,
    onTimeRate: 93.2,
    openPurchaseOrders: 7,
    city: 'Kaunas',
  },
  {
    id: 'sup-003',
    name: 'Prime Wholesale',
    supplierType: 'finished_goods',
    status: 'probation',
    contactName: 'Milda Kazlaite',
    email: 'procurement@primewholesale.com',
    leadTimeDays: 14,
    onTimeRate: 84.1,
    openPurchaseOrders: 3,
    city: 'Klaipeda',
  },
  {
    id: 'sup-004',
    name: 'Ops Service Partners',
    supplierType: 'service',
    status: 'approved',
    contactName: 'Tomas Grigas',
    email: 'support@ospartners.io',
    leadTimeDays: 2,
    onTimeRate: 99.1,
    openPurchaseOrders: 1,
    city: 'Siauliai',
  },
  {
    id: 'sup-005',
    name: 'Legacy Distribution',
    supplierType: 'raw_material',
    status: 'blocked',
    contactName: 'Asta Vilke',
    email: 'contact@legacydist.net',
    leadTimeDays: 18,
    onTimeRate: 71.6,
    openPurchaseOrders: 0,
    city: 'Panevezys',
  },
])

const filteredSuppliers = computed(() => {
  if (selectedStatus.value == 'all') {
    return suppliers.value
  }
  return suppliers.value.filter((supplier) => supplier.status == selectedStatus.value)
})

const summary = computed(() => {
  let approvedCount = 0
  let totalLeadTimeDays = 0
  let totalOpenPOs = 0
  filteredSuppliers.value.forEach((supplier) => {
    totalLeadTimeDays += supplier.leadTimeDays
    totalOpenPOs += supplier.openPurchaseOrders
    if (supplier.status == 'approved') {
      approvedCount++
    }
  })
  return {
    totalSuppliers: filteredSuppliers.value.length,
    approvedCount,
    avgLeadTime:
      filteredSuppliers.value.length > 0 ? totalLeadTimeDays / filteredSuppliers.value.length : 0,
    totalOpenPOs,
  }
})

const supplierTypeLabel = (value: SupplierType): string => {
  return value.replace('_', ' ').replace('_', ' ')
}

const statusBadgeClass = (status: SupplierStatus): string => {
  if (status == 'approved') {
    return 'bg-success-subtle text-success'
  }
  if (status == 'probation') {
    return 'bg-warning-subtle text-warning'
  }
  return 'bg-danger-subtle text-danger'
}

const renderStatus = (status: SupplierStatus): string => {
  return `<span class="badge text-capitalize ${statusBadgeClass(status)}">${status}</span>`
}

const getSuppliersTableColumns = (row: SupplierRecord): ITableColumn[] => {
  return [
    {
      name: 'Supplier',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Category',
      value: {
        text: supplierTypeLabel(row.supplierType),
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Contact',
      value: {
        text: `${row.contactName} | ${row.email}`,
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
  ]
}
</script>

<template>
  <div class="operations-suppliers-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Suppliers</h5>
          <p class="text-body-secondary mb-0">
            Monitor supplier reliability, risk status, and purchasing load.
          </p>
        </div>
        <button class="btn btn-outline-primary mt-4 mt-md-0" @click="selectedStatus = 'all'">
          <i class="bx bx-reset me-2"></i>
          Reset Filters
        </button>
      </div>
    </div>

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
        <vue-table :rows="filteredSuppliers" :columns="getSuppliersTableColumns" :checkboxes="true">
          <template #actions="{ row }">
            <button type="button" class="btn btn-sm btn-outline-primary" :title="`Open ${row.name}`">
              View
            </button>
          </template>
          <template #noResults>
            <p class="py-6">No suppliers found for the selected status.</p>
          </template>
        </vue-table>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
