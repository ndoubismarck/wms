<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useApp } from '@/app/app'
import { useOrdersStore } from '@/stores/orders_store'
import type { OrderModel, OrderPriority, OrderStatus } from '@/app/models/order_model'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueFormDatepickerInput from '@/views/shared/components/VueForm/VueFormDatepickerInput.vue'
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import { DateTime } from '@/app/core/date_time'

defineOptions({
  name: 'AdminOperationsOrdersPage',
})

const app = useApp()
const ordersStore = useOrdersStore()

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const selectedStatus = ref<'all' | OrderStatus>('all')
const selectedPriority = ref<'all' | OrderPriority>('all')

const orderOffcanvas = ref<IOffcanvas | null>(null)
const orderFormMode = ref<'create' | 'edit'>('create')
const editingOrder = ref<OrderModel | null>(null)
const orderFormKey = ref<number>(0)

const orders = computed(() => ordersStore.get().value)
const summary = computed(() => ordersStore.getSummary().value)

const statusOptions: Array<{ label: string; value: 'all' | OrderStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'New', value: 'new' },
  { label: 'Picking', value: 'picking' },
  { label: 'Packed', value: 'packed' },
  { label: 'Shipped', value: 'shipped' },
  { label: 'Backorder', value: 'backorder' },
  { label: 'Cancelled', value: 'cancelled' },
]

const priorityOptions: Array<{ label: string; value: 'all' | OrderPriority }> = [
  { label: 'All Priorities', value: 'all' },
  { label: 'Low', value: 'low' },
  { label: 'Normal', value: 'normal' },
  { label: 'High', value: 'high' },
  { label: 'Urgent', value: 'urgent' },
]

const statusSelectOptions = computed<ISelectInputOption[]>(() => {
  return statusOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const prioritySelectOptions = computed<ISelectInputOption[]>(() => {
  return priorityOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const orderFormInitial = computed(() => {
  if (!editingOrder.value || orderFormMode.value == 'create') {
    return {
      order_number: '',
      customer_name: '',
      channel: 'manual',
      status: 'new',
      priority: 'normal',
      items_count: '0',
      total_amount: '0',
      due_at: '',
      notes: '',
    }
  }
  return {
    order_number: editingOrder.value.orderNumber,
    customer_name: editingOrder.value.customerName,
    channel: editingOrder.value.channel,
    status: editingOrder.value.status,
    priority: editingOrder.value.priority,
    items_count: `${editingOrder.value.itemsCount}`,
    total_amount: `${editingOrder.value.totalAmount}`,
    due_at: toDateTimeLocal(editingOrder.value.dueAt),
    notes: editingOrder.value.notes,
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

const dueAtText = (row: OrderModel): string => {
  return row.dueAt ? new DateTime(row.dueAt).toDateTimeString() : 'N/A'
}

const getOrderTableColumns = (row: OrderModel): ITableColumn[] => {
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
      name: 'Due At',
      value: {
        text: dueAtText(row),
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

const openCreateOrder = () => {
  orderFormMode.value = 'create'
  editingOrder.value = null
  orderFormKey.value++
  orderOffcanvas.value?.show()
}

const openEditOrder = (row: OrderModel) => {
  orderFormMode.value = 'edit'
  editingOrder.value = row
  orderFormKey.value++
  orderOffcanvas.value?.show()
}

const orderFormTitle = computed(() => {
  return orderFormMode.value == 'create' ? 'Create Order' : 'Edit Order'
})

const submitOrderForm = async (data: VueFormData) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const payload: Record<string, unknown> = {
      order_number: `${data.order_number ?? ''}`.trim(),
      customer_name: `${data.customer_name ?? ''}`.trim(),
      channel: `${data.channel ?? ''}`.trim(),
      status: `${data.status ?? 'new'}`,
      priority: `${data.priority ?? 'normal'}`,
      items_count: Math.max(0, Number(data.items_count ?? 0)),
      total_amount: Math.max(0, Number(data.total_amount ?? 0)),
      notes: `${data.notes ?? ''}`.trim(),
    }
    if (data.due_at) {
      payload.due_at = new Date(`${data.due_at}`).toISOString()
    }

    if (orderFormMode.value == 'create') {
      const result = await app.services.operations.addOrder(payload)
      if (!result.success) {
        return
      }
    } else if (editingOrder.value) {
      const result = await app.services.operations.updateOrder(editingOrder.value.id, payload)
      if (!result.success) {
        return
      }
    }

    orderOffcanvas.value?.hide()
    await loadOrders(false)
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteOrder = async (row: OrderModel) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.operations.deleteOrderById(row.id)
    if (result.success) {
      await loadOrders(false)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const loadOrders = async (withLoader: boolean = false): Promise<void> => {
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
    if (selectedPriority.value != 'all') {
      payload.priority = selectedPriority.value
    }
    const result = await app.services.operations.getOrders(payload)
    if (result.success) {
      ordersStore.setAll(result.orders)
      ordersStore.setSummary(result.summary)
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
  selectedPriority.value = 'all'
}

watch([selectedStatus, selectedPriority], async () => {
  await loadOrders(false)
})

onMounted(async () => {
  await loadOrders(true)
})
</script>

<template>
  <div class="operations-orders-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Orders</h5>
          <p class="text-body-secondary mb-0">Manage warehouse order queue with create/edit backend records.</p>
        </div>
        <div class="d-flex gap-2 mt-4 mt-md-0">
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="resetFilters">
            <i class="bx bx-reset me-2"></i>
            Reset Filters
          </button>
          <button class="btn btn-primary" :disabled="isSubmitting" @click="openCreateOrder">
            <i class="bx bx-plus me-2"></i>
            Create Order
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
        <div class="card-header d-flex flex-column flex-lg-row justify-content-between align-items-lg-center gap-3">
          <h5 class="mb-0">Order Queue</h5>
          <div class="d-flex flex-column flex-md-row align-items-md-center gap-2">
            <div class="d-flex align-items-center gap-2">
              <label class="form-label mb-0 text-body-secondary">Status</label>
              <select v-model="selectedStatus" class="form-select">
                <option v-for="option in statusOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
            <div class="d-flex align-items-center gap-2">
              <label class="form-label mb-0 text-body-secondary">Priority</label>
              <select v-model="selectedPriority" class="form-select">
                <option v-for="option in priorityOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="orders" :columns="getOrderTableColumns" :checkboxes="true">
            <template #actions="{ row }">
              <div class="d-flex gap-2">
                <button type="button" class="btn btn-sm btn-outline-primary" @click="openEditOrder(row)">Edit</button>
                <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteOrder(row)">Delete</button>
              </div>
            </template>
            <template #noResults>
              <p class="py-6">No orders found for selected filters.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>

  <offcanvas ref="orderOffcanvas" :show="false" :overflow="isSubmitting" :title="orderFormTitle">
    <template #body>
      <div class="card">
        <div class="card-body">
          <vue-form :key="orderFormKey" @submit="submitOrderForm">
            <div class="mb-6">
              <vue-form-input
                name="order_number"
                type="text"
                label="Order Number"
                placeholder="SO-10458"
                validation="required"
                :value="orderFormInitial.order_number"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="customer_name"
                type="text"
                label="Customer Name"
                placeholder="Northfield Retail"
                validation="required"
                :value="orderFormInitial.customer_name"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="channel"
                type="text"
                label="Channel"
                placeholder="B2B Portal"
                :value="orderFormInitial.channel"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="status"
                label="Status"
                placeholder="Select status"
                :options="statusSelectOptions"
                :value="orderFormInitial.status"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="priority"
                label="Priority"
                placeholder="Select priority"
                :options="prioritySelectOptions"
                :value="orderFormInitial.priority"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="items_count"
                type="number"
                label="Items Count"
                placeholder="10"
                :value="orderFormInitial.items_count"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="total_amount"
                type="number"
                label="Total Amount"
                placeholder="1200.00"
                :value="orderFormInitial.total_amount"
                input-mode="decimal"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-datepicker-input
                name="due_at"
                label="Due At"
                mode="datetime-local"
                placeholder="Select date"
                :value="orderFormInitial.due_at"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-textarea-input
                name="notes"
                label="Notes"
                placeholder="Optional notes"
                :value="orderFormInitial.notes"
                :disabled="isSubmitting" />
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              {{ orderFormMode == 'create' ? 'Create Order' : 'Update Order' }}
            </button>
          </vue-form>
        </div>
      </div>
    </template>
  </offcanvas>
</template>

<style scoped></style>
