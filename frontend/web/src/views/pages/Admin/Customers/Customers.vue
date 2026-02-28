<script setup lang="ts">
import { computed, ref } from 'vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminCustomersPage',
})

type CustomerTier = 'standard' | 'growth' | 'enterprise'
type CustomerStatus = 'active' | 'paused' | 'churn_risk'

interface CustomerRecord {
  id: string
  name: string
  email: string
  phone: string
  city: string
  tier: CustomerTier
  status: CustomerStatus
  totalOrders: number
  lifetimeValue: number
  lastOrderAt: string
}

const selectedTier = ref<'all' | CustomerTier>('all')

const tierOptions: Array<{ label: string; value: 'all' | CustomerTier }> = [
  { label: 'All Tiers', value: 'all' },
  { label: 'Standard', value: 'standard' },
  { label: 'Growth', value: 'growth' },
  { label: 'Enterprise', value: 'enterprise' },
]

const customers = ref<CustomerRecord[]>([
  {
    id: 'cus-001',
    name: 'Northfield Retail',
    email: 'ops@northfieldretail.com',
    phone: '+370 612 33001',
    city: 'Vilnius',
    tier: 'enterprise',
    status: 'active',
    totalOrders: 142,
    lifetimeValue: 234500.75,
    lastOrderAt: 'Feb 28, 2026 08:12 AM',
  },
  {
    id: 'cus-002',
    name: 'Bluebird Stores',
    email: 'logistics@bluebird.lt',
    phone: '+370 612 33002',
    city: 'Kaunas',
    tier: 'growth',
    status: 'active',
    totalOrders: 89,
    lifetimeValue: 113420.1,
    lastOrderAt: 'Feb 28, 2026 07:18 AM',
  },
  {
    id: 'cus-003',
    name: 'Horizon Market',
    email: 'warehouse@horizonmarket.com',
    phone: '+370 612 33003',
    city: 'Klaipeda',
    tier: 'growth',
    status: 'paused',
    totalOrders: 64,
    lifetimeValue: 74210.9,
    lastOrderAt: 'Feb 22, 2026 03:44 PM',
  },
  {
    id: 'cus-004',
    name: 'Packline Distribution',
    email: 'planning@packline.io',
    phone: '+370 612 33004',
    city: 'Panevezys',
    tier: 'enterprise',
    status: 'active',
    totalOrders: 188,
    lifetimeValue: 319880.55,
    lastOrderAt: 'Feb 27, 2026 11:11 AM',
  },
  {
    id: 'cus-005',
    name: 'Prime Goods Co.',
    email: 'hello@primegoods.co',
    phone: '+370 612 33005',
    city: 'Siauliai',
    tier: 'standard',
    status: 'churn_risk',
    totalOrders: 21,
    lifetimeValue: 17590.3,
    lastOrderAt: 'Jan 30, 2026 04:55 PM',
  },
])

const filteredCustomers = computed(() => {
  if (selectedTier.value == 'all') {
    return customers.value
  }
  return customers.value.filter((customer) => customer.tier == selectedTier.value)
})

const summary = computed(() => {
  let totalLtv = 0
  let activeCount = 0
  let enterpriseCount = 0
  filteredCustomers.value.forEach((customer) => {
    totalLtv += customer.lifetimeValue
    if (customer.status == 'active') {
      activeCount++
    }
    if (customer.tier == 'enterprise') {
      enterpriseCount++
    }
  })
  return {
    totalCustomers: filteredCustomers.value.length,
    activeCount,
    enterpriseCount,
    avgLtv: filteredCustomers.value.length > 0 ? totalLtv / filteredCustomers.value.length : 0,
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

const getCustomersTableColumns = (row: CustomerRecord): ITableColumn[] => {
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
        text: `${row.email} | ${row.phone}`,
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
        text: row.lastOrderAt,
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
  <div class="operations-customers-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Customers</h5>
          <p class="text-body-secondary mb-0">
            Monitor account health, commercial value, and service tier distribution.
          </p>
        </div>
        <button class="btn btn-outline-primary mt-4 mt-md-0" @click="selectedTier = 'all'">
          <i class="bx bx-reset me-2"></i>
          Reset Filters
        </button>
      </div>
    </div>

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
      <div class="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
        <h5 class="mb-0">Customer Directory</h5>
        <div class="d-flex align-items-center gap-2">
          <label class="form-label mb-0 text-body-secondary">Tier</label>
          <select v-model="selectedTier" class="form-select">
            <option v-for="option in tierOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>
      <div class="card-body p-0">
        <vue-table :rows="filteredCustomers" :columns="getCustomersTableColumns" :checkboxes="true">
          <template #actions="{ row }">
            <button type="button" class="btn btn-sm btn-outline-primary" :title="`Open ${row.name}`">
              Profile
            </button>
          </template>
          <template #noResults>
            <p class="py-6">No customers found for the selected tier.</p>
          </template>
        </vue-table>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
