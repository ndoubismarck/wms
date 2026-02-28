<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useApp } from '@/app/app'
import { useShipmentsStore } from '@/stores/shipments_store'
import type { ShipmentModel, ShipmentStatus } from '@/app/models/shipment_model'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import Preloader from '@/views/shared/components/Preloader.vue'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueFormDatepickerInput from '@/views/shared/components/VueForm/VueFormDatepickerInput.vue'
import { DateTime } from '@/app/core/date_time'

defineOptions({
  name: 'AdminOperationsShipmentsPage',
})

const app = useApp()
const shipmentsStore = useShipmentsStore()

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const isSubmitting = ref<boolean>(false)
const selectedStatus = ref<'all' | ShipmentStatus>('all')

const shipmentOffcanvas = ref<IOffcanvas | null>(null)
const shipmentFormMode = ref<'create' | 'edit'>('create')
const editingShipment = ref<ShipmentModel | null>(null)
const shipmentFormKey = ref<number>(0)

const shipments = computed(() => shipmentsStore.get().value)
const summary = computed(() => shipmentsStore.getSummary().value)

const statusOptions: Array<{ label: string; value: 'all' | ShipmentStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Draft', value: 'draft' },
  { label: 'Ready', value: 'ready' },
  { label: 'In Transit', value: 'in_transit' },
  { label: 'Delivered', value: 'delivered' },
  { label: 'Exception', value: 'exception' },
]

const statusSelectOptions = computed<ISelectInputOption[]>(() => {
  return statusOptions.filter((option) => option.value != 'all') as ISelectInputOption[]
})

const shipmentFormInitial = computed(() => {
  if (!editingShipment.value || shipmentFormMode.value == 'create') {
    return {
      shipment_number: '',
      order_number: '',
      carrier: '',
      service: '',
      status: 'draft',
      packages: '1',
      tracking_code: '',
      eta: '',
      destination: '',
    }
  }
  return {
    shipment_number: editingShipment.value.shipmentNumber,
    order_number: editingShipment.value.orderNumber,
    carrier: editingShipment.value.carrier,
    service: editingShipment.value.service,
    status: editingShipment.value.status,
    packages: `${editingShipment.value.packages}`,
    tracking_code: editingShipment.value.trackingCode,
    eta: toDateTimeLocal(editingShipment.value.eta),
    destination: editingShipment.value.destination,
  }
})

const statusBadgeClass = (status: ShipmentStatus): string => {
  if (status == 'draft') {
    return 'bg-label-secondary text-secondary'
  }
  if (status == 'ready') {
    return 'bg-label-info text-info'
  }
  if (status == 'in_transit') {
    return 'bg-label-primary text-primary'
  }
  if (status == 'delivered') {
    return 'bg-label-success text-success'
  }
  return 'bg-label-danger text-danger'
}

const renderStatus = (status: ShipmentStatus): string => {
  return `<span class="badge text-capitalize ${statusBadgeClass(status)}">${status.replace('_', ' ')}</span>`
}

const etaText = (row: ShipmentModel): string => {
  return row.eta ? new DateTime(row.eta).toDateTimeString() : 'N/A'
}

const getShipmentTableColumns = (row: ShipmentModel): ITableColumn[] => {
  return [
    {
      name: 'Shipment',
      value: {
        text: row.shipmentNumber,
      },
    },
    {
      name: 'Order',
      value: {
        text: row.orderNumber,
      },
    },
    {
      name: 'Carrier',
      value: {
        text: `${row.carrier} (${row.service})`,
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Packages',
      value: {
        text: `${row.packages}`,
      },
    },
    {
      name: 'Tracking',
      value: {
        text: row.trackingCode,
      },
    },
    {
      name: 'ETA',
      value: {
        text: etaText(row),
      },
    },
    {
      name: 'Destination',
      value: {
        text: row.destination,
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

const openCreateShipment = () => {
  shipmentFormMode.value = 'create'
  editingShipment.value = null
  shipmentFormKey.value++
  shipmentOffcanvas.value?.show()
}

const openEditShipment = (row: ShipmentModel) => {
  shipmentFormMode.value = 'edit'
  editingShipment.value = row
  shipmentFormKey.value++
  shipmentOffcanvas.value?.show()
}

const shipmentFormTitle = computed(() => {
  return shipmentFormMode.value == 'create' ? 'Create Shipment' : 'Edit Shipment'
})

const submitShipmentForm = async (data: VueFormData) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const payload: Record<string, unknown> = {
      shipment_number: `${data.shipment_number ?? ''}`.trim(),
      order_number: `${data.order_number ?? ''}`.trim(),
      carrier: `${data.carrier ?? ''}`.trim(),
      service: `${data.service ?? ''}`.trim(),
      status: `${data.status ?? 'draft'}`,
      packages: Math.max(0, Number(data.packages ?? 0)),
      tracking_code: `${data.tracking_code ?? ''}`.trim(),
      destination: `${data.destination ?? ''}`.trim(),
    }
    if (data.eta) {
      payload.eta = new Date(`${data.eta}`).toISOString()
    }

    if (shipmentFormMode.value == 'create') {
      const result = await app.services.operations.addShipment(payload)
      if (!result.success) {
        return
      }
    } else if (editingShipment.value) {
      const result = await app.services.operations.updateShipment(editingShipment.value.id, payload)
      if (!result.success) {
        return
      }
    }

    shipmentOffcanvas.value?.hide()
    await loadShipments(false)
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const deleteShipment = async (row: ShipmentModel) => {
  if (isSubmitting.value) {
    return
  }
  isSubmitting.value = true
  try {
    const result = await app.services.operations.deleteShipmentById(row.id)
    if (result.success) {
      await loadShipments(false)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isSubmitting.value = false
  }
}

const loadShipments = async (withLoader: boolean = false): Promise<void> => {
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
    const result = await app.services.operations.getShipments(payload)
    if (result.success) {
      shipmentsStore.setAll(result.shipments)
      shipmentsStore.setSummary(result.summary)
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
  await loadShipments(false)
})

onMounted(async () => {
  await loadShipments(true)
})
</script>

<template>
  <div class="operations-shipments-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Shipments</h5>
          <p class="text-body-secondary mb-0">Manage shipment lifecycle with create/edit backend records.</p>
        </div>
        <div class="d-flex gap-2 mt-4 mt-md-0">
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="resetFilters">
            <i class="bx bx-reset me-2"></i>
            Reset Filters
          </button>
          <button class="btn btn-primary" :disabled="isSubmitting" @click="openCreateShipment">
            <i class="bx bx-plus me-2"></i>
            Create Shipment
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
              <p class="text-body-secondary mb-1">Shipments</p>
              <h4 class="mb-0">{{ summary.totalShipments }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">In Transit</p>
              <h4 class="mb-0 text-primary">{{ summary.inTransit }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Delivered</p>
              <h4 class="mb-0 text-success">{{ summary.delivered }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Exceptions</p>
              <h4 class="mb-0 text-danger">{{ summary.exceptions }}</h4>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
          <h5 class="mb-0">Shipment List</h5>
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
          <vue-table :rows="shipments" :columns="getShipmentTableColumns" :checkboxes="true">
            <template #actions="{ row }">
              <div class="d-flex gap-2">
                <button type="button" class="btn btn-sm btn-outline-primary" @click="openEditShipment(row)">Edit</button>
                <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteShipment(row)">Delete</button>
              </div>
            </template>
            <template #noResults>
              <p class="py-6">No shipments found for selected filters.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>

  <offcanvas ref="shipmentOffcanvas" :show="false" :overflow="isSubmitting" :title="shipmentFormTitle">
    <template #body>
      <div class="card">
        <div class="card-body">
          <vue-form :key="shipmentFormKey" @submit="submitShipmentForm">
            <div class="mb-6">
              <vue-form-input
                name="shipment_number"
                type="text"
                label="Shipment Number"
                placeholder="SHP-90015"
                validation="required"
                :value="shipmentFormInitial.shipment_number"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="order_number"
                type="text"
                label="Order Number"
                placeholder="SO-10458"
                validation="required"
                :value="shipmentFormInitial.order_number"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="carrier"
                type="text"
                label="Carrier"
                placeholder="DHL"
                validation="required"
                :value="shipmentFormInitial.carrier"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="service"
                type="text"
                label="Service"
                placeholder="Express"
                validation="required"
                :value="shipmentFormInitial.service"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-select-input
                name="status"
                label="Status"
                placeholder="Select status"
                :options="statusSelectOptions"
                :value="shipmentFormInitial.status"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="packages"
                type="number"
                label="Packages"
                placeholder="1"
                :value="shipmentFormInitial.packages"
                validation="required"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="tracking_code"
                type="text"
                label="Tracking Code"
                placeholder="DH-20485881"
                validation="required"
                :value="shipmentFormInitial.tracking_code"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-datepicker-input
                name="eta"
                label="ETA"
                mode="datetime-local"
                placeholder="Select ETA"
                :value="shipmentFormInitial.eta"
                :disabled="isSubmitting" />
            </div>

            <div class="mb-6">
              <vue-form-input
                name="destination"
                type="text"
                label="Destination"
                placeholder="Vilnius"
                validation="required"
                :value="shipmentFormInitial.destination"
                :disabled="isSubmitting" />
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
              {{ shipmentFormMode == 'create' ? 'Create Shipment' : 'Update Shipment' }}
            </button>
          </vue-form>
        </div>
      </div>
    </template>
  </offcanvas>
</template>

<style scoped></style>
