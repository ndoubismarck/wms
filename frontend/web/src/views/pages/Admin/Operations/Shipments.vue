<script setup lang="ts">
import { computed, ref } from 'vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminOperationsShipmentsPage',
})

type ShipmentStatus = 'draft' | 'ready' | 'in_transit' | 'delivered' | 'exception'

interface ShipmentRecord {
  id: string
  shipmentNumber: string
  orderNumber: string
  carrier: string
  service: string
  status: ShipmentStatus
  packages: number
  trackingCode: string
  eta: string
  destination: string
}

const selectedStatus = ref<'all' | ShipmentStatus>('all')

const statusOptions: Array<{ label: string; value: 'all' | ShipmentStatus }> = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Draft', value: 'draft' },
  { label: 'Ready', value: 'ready' },
  { label: 'In Transit', value: 'in_transit' },
  { label: 'Delivered', value: 'delivered' },
  { label: 'Exception', value: 'exception' },
]

const shipments = ref<ShipmentRecord[]>([
  {
    id: 'shp-001',
    shipmentNumber: 'SHP-90015',
    orderNumber: 'SO-10458',
    carrier: 'DHL',
    service: 'Express',
    status: 'ready',
    packages: 3,
    trackingCode: 'DH-20485881',
    eta: 'Mar 01, 2026',
    destination: 'Vilnius',
  },
  {
    id: 'shp-002',
    shipmentNumber: 'SHP-90014',
    orderNumber: 'SO-10457',
    carrier: 'FedEx',
    service: 'Priority Overnight',
    status: 'in_transit',
    packages: 6,
    trackingCode: 'FD-08441177',
    eta: 'Feb 28, 2026',
    destination: 'Kaunas',
  },
  {
    id: 'shp-003',
    shipmentNumber: 'SHP-90013',
    orderNumber: 'SO-10455',
    carrier: 'DPD',
    service: 'Ground',
    status: 'exception',
    packages: 2,
    trackingCode: 'DP-44813751',
    eta: 'Mar 02, 2026',
    destination: 'Klaipeda',
  },
  {
    id: 'shp-004',
    shipmentNumber: 'SHP-90012',
    orderNumber: 'SO-10452',
    carrier: 'UPS',
    service: 'Standard',
    status: 'delivered',
    packages: 5,
    trackingCode: 'UP-77128388',
    eta: 'Feb 27, 2026',
    destination: 'Panevezys',
  },
  {
    id: 'shp-005',
    shipmentNumber: 'SHP-90011',
    orderNumber: 'SO-10450',
    carrier: 'DHL',
    service: 'Economy',
    status: 'draft',
    packages: 1,
    trackingCode: 'Pending',
    eta: 'TBD',
    destination: 'Siauliai',
  },
])

const filteredShipments = computed(() => {
  if (selectedStatus.value == 'all') {
    return shipments.value
  }
  return shipments.value.filter((shipment) => shipment.status == selectedStatus.value)
})

const summary = computed(() => {
  let inTransit = 0
  let delivered = 0
  let exceptions = 0
  let totalPackages = 0
  filteredShipments.value.forEach((shipment) => {
    totalPackages += shipment.packages
    if (shipment.status == 'in_transit') {
      inTransit++
    }
    if (shipment.status == 'delivered') {
      delivered++
    }
    if (shipment.status == 'exception') {
      exceptions++
    }
  })
  return {
    totalShipments: filteredShipments.value.length,
    inTransit,
    delivered,
    exceptions,
    totalPackages,
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

const getShipmentTableColumns = (row: ShipmentRecord): ITableColumn[] => {
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
        text: row.eta,
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
</script>

<template>
  <div class="operations-shipments-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Shipments</h5>
          <p class="text-body-secondary mb-0">
            Track packing handover, carrier execution, and delivery issues in one place.
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
        <vue-table :rows="filteredShipments" :columns="getShipmentTableColumns" :checkboxes="true">
          <template #actions="{ row }">
            <button type="button" class="btn btn-sm btn-outline-primary" :title="`Track ${row.shipmentNumber}`">
              Track
            </button>
          </template>
          <template #noResults>
            <p class="py-6">No shipments found for the selected status.</p>
          </template>
        </vue-table>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
