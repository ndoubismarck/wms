<script setup lang="ts">
import { computed, ref } from 'vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminOperationsOrdersPage',
})

type OrderStatus = 'new' | 'picking' | 'packed' | 'shipped' | 'backorder' | 'cancelled'
type OrderPriority = 'low' | 'normal' | 'high' | 'urgent'

interface OrderRecord {
  id: string
  orderNumber: string
  customerName: string
  channel: string
  status: OrderStatus
  priority: OrderPriority
  itemsCount: number
  totalAmount: number
  dueBy: string
  createdAt: string
}

const selectedStatus = ref<'all' | OrderStatus>('all')

const statusOptions: Array<{ label: string; value: 'all' | OrderStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'New', value: 'new' },
  { label: 'Picking', value: 'picking' },
  { label: 'Packed', value: 'packed' },
  { label: 'Shipped', value: 'shipped' },
  { label: 'Backorder', value: 'backorder' },
  { label: 'Cancelled', value: 'cancelled' },
]

const orders = ref<OrderRecord[]>([
  {
    id: 'ord-0001',
    orderNumber: 'SO-10458',
    customerName: 'Northfield Retail',
    channel: 'B2B Portal',
    status: 'new',
    priority: 'high',
    itemsCount: 14,
    totalAmount: 4820.5,
    dueBy: 'Mar 01, 2026 09:00 AM',
    createdAt: 'Feb 28, 2026 08:12 AM',
  },
  {
    id: 'ord-0002',
    orderNumber: 'SO-10457',
    customerName: 'Bluebird Stores',
    channel: 'EDI',
    status: 'picking',
    priority: 'urgent',
    itemsCount: 26,
    totalAmount: 12240,
    dueBy: 'Feb 28, 2026 02:00 PM',
    createdAt: 'Feb 28, 2026 07:18 AM',
  },
  {
    id: 'ord-0003',
    orderNumber: 'SO-10456',
    customerName: 'Horizon Market',
    channel: 'Marketplace',
    status: 'packed',
    priority: 'normal',
    itemsCount: 8,
    totalAmount: 1990.75,
    dueBy: 'Mar 01, 2026 10:30 AM',
    createdAt: 'Feb 27, 2026 05:44 PM',
  },
  {
    id: 'ord-0004',
    orderNumber: 'SO-10455',
    customerName: 'Packline Distribution',
    channel: 'B2B Portal',
    status: 'shipped',
    priority: 'normal',
    itemsCount: 32,
    totalAmount: 16420.2,
    dueBy: 'Feb 27, 2026 06:00 PM',
    createdAt: 'Feb 27, 2026 11:11 AM',
  },
  {
    id: 'ord-0005',
    orderNumber: 'SO-10454',
    customerName: 'Vertex Supplies',
    channel: 'EDI',
    status: 'backorder',
    priority: 'high',
    itemsCount: 18,
    totalAmount: 7140.6,
    dueBy: 'Mar 02, 2026 03:00 PM',
    createdAt: 'Feb 27, 2026 10:02 AM',
  },
  {
    id: 'ord-0006',
    orderNumber: 'SO-10453',
    customerName: 'Prime Goods Co.',
    channel: 'Marketplace',
    status: 'cancelled',
    priority: 'low',
    itemsCount: 4,
    totalAmount: 640,
    dueBy: 'N/A',
    createdAt: 'Feb 26, 2026 04:50 PM',
  },
])

const filteredOrders = computed(() => {
  if (selectedStatus.value == 'all') {
    return orders.value
  }
  return orders.value.filter((order) => order.status == selectedStatus.value)
})

const summary = computed(() => {
  const source = filteredOrders.value
  let totalOrderValue = 0
  let openOrders = 0
  let urgentOrders = 0

  source.forEach((order) => {
    totalOrderValue += order.totalAmount
    if (order.status != 'shipped' && order.status != 'cancelled') {
      openOrders++
    }
    if (order.priority == 'urgent') {
      urgentOrders++
    }
  })

  return {
    totalOrders: source.length,
    totalOrderValue,
    openOrders,
    urgentOrders,
  }
})

const statusBadgeClass = (status: OrderStatus): string => {
  if (status == 'new') {
    return 'bg-label-primary text-primary'
  }
  if (status == 'picking') {
    return 'bg-label-warning text-warning'
  }
  if (status == 'packed') {
    return 'bg-label-info text-info'
  }
  if (status == 'shipped') {
    return 'bg-label-success text-success'
  }
  if (status == 'backorder') {
    return 'bg-label-danger text-danger'
  }
  return 'bg-label-secondary text-secondary'
}

const priorityBadgeClass = (priority: OrderPriority): string => {
  if (priority == 'urgent') {
    return 'bg-danger-subtle text-danger'
  }
  if (priority == 'high') {
    return 'bg-warning-subtle text-warning'
  }
  if (priority == 'normal') {
    return 'bg-info-subtle text-info'
  }
  return 'bg-secondary-subtle text-secondary'
}

const renderStatus = (status: OrderStatus): string => {
  return `<span class="badge text-capitalize ${statusBadgeClass(status)}">${status}</span>`
}

const renderPriority = (priority: OrderPriority): string => {
  return `<span class="badge text-capitalize ${priorityBadgeClass(priority)}">${priority}</span>`
}

const getOrderTableColumns = (row: OrderRecord): ITableColumn[] => {
  return [
    {
      name: 'Order',
      value: {
        text: row.orderNumber,
      },
    },
    {
      name: 'Customer',
      value: {
        text: row.customerName,
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
      name: 'Items',
      value: {
        text: `${row.itemsCount}`,
      },
    },
    {
      name: 'Value',
      value: {
        text: `$${row.totalAmount.toFixed(2)}`,
      },
    },
    {
      name: 'Due By',
      value: {
        text: row.dueBy,
      },
    },
    {
      name: 'Created',
      value: {
        text: row.createdAt,
      },
    },
  ]
}
</script>

<template>
  <div class="operations-orders-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Orders</h5>
          <p class="text-body-secondary mb-0">
            Track intake, fulfillment stage, and order value concentration across channels.
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
            <p class="text-body-secondary mb-1">Orders</p>
            <h4 class="mb-0">{{ summary.totalOrders }}</h4>
          </div>
        </div>
      </div>
      <div class="col-12 col-md-6 col-xl-3 mb-6">
        <div class="card h-100">
          <div class="card-body">
            <p class="text-body-secondary mb-1">Open Orders</p>
            <h4 class="mb-0 text-warning">{{ summary.openOrders }}</h4>
          </div>
        </div>
      </div>
      <div class="col-12 col-md-6 col-xl-3 mb-6">
        <div class="card h-100">
          <div class="card-body">
            <p class="text-body-secondary mb-1">Urgent Orders</p>
            <h4 class="mb-0 text-danger">{{ summary.urgentOrders }}</h4>
          </div>
        </div>
      </div>
      <div class="col-12 col-md-6 col-xl-3 mb-6">
        <div class="card h-100">
          <div class="card-body">
            <p class="text-body-secondary mb-1">Order Value</p>
            <h4 class="mb-0">${{ summary.totalOrderValue.toFixed(2) }}</h4>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
        <h5 class="mb-0">Order Queue</h5>
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
        <vue-table :rows="filteredOrders" :columns="getOrderTableColumns" :checkboxes="true">
          <template #actions="{ row }">
            <button type="button" class="btn btn-sm btn-outline-primary" :title="`Open ${row.orderNumber}`">
              View
            </button>
          </template>
          <template #noResults>
            <p class="py-6">No orders found for the selected status.</p>
          </template>
        </vue-table>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
