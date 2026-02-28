<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useApp } from '@/app/app.ts'
import { DateTime } from '@/app/core/date_time.ts'
import type { InventoryMovementModel, InventoryMovementType } from '@/app/models/inventory_movement_model.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import { useSetupStore } from '@/stores/setup_store.ts'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'

defineOptions({
  name: 'AdminInventoryMovementsPage',
})

const app = useApp()
const setupStore = useSetupStore()

const setup = computed(() => setupStore.get().value)

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const lastUpdatedAt = ref<string>('')
const selectedMovementType = ref<'all' | InventoryMovementType>('all')

const movements = ref<InventoryMovementModel[]>([])

const movementTypeOptions: Array<{
  label: string
  value: 'all' | InventoryMovementType
}> = [
  { label: 'All Types', value: 'all' },
  { label: 'Receive', value: 'receive' },
  { label: 'Ship', value: 'ship' },
  { label: 'Adjust', value: 'adjust' },
  { label: 'Damage', value: 'damage' },
  { label: 'Reserve', value: 'reserve' },
  { label: 'Release', value: 'release' },
  { label: 'Transfer Out', value: 'transfer_out' },
  { label: 'Transfer In', value: 'transfer_in' },
]

const incomingTypes: InventoryMovementType[] = ['receive', 'transfer_in', 'release']
const outgoingTypes: InventoryMovementType[] = ['ship', 'transfer_out', 'reserve']
const adjustmentTypes: InventoryMovementType[] = ['adjust', 'damage']

const filteredMovements = computed(() => {
  if (selectedMovementType.value == 'all') {
    return movements.value
  }
  return movements.value.filter((movement) => movement.movementType == selectedMovementType.value)
})

const summary = computed(() => {
  let totalQuantity = 0
  let incomingQuantity = 0
  let outgoingQuantity = 0
  let adjustmentQuantity = 0

  filteredMovements.value.forEach((movement) => {
    totalQuantity += movement.quantity
    if (incomingTypes.includes(movement.movementType)) {
      incomingQuantity += movement.quantity
    }
    if (outgoingTypes.includes(movement.movementType)) {
      outgoingQuantity += movement.quantity
    }
    if (adjustmentTypes.includes(movement.movementType)) {
      adjustmentQuantity += movement.quantity
    }
  })

  return {
    totalRecords: filteredMovements.value.length,
    totalQuantity,
    incomingQuantity,
    outgoingQuantity,
    adjustmentQuantity,
  }
})

const movementTypeLabel = (value: InventoryMovementType): string => {
  const option = movementTypeOptions.find((item) => item.value == value)
  return option?.label ?? value
}

const movementBadgeClass = (value: InventoryMovementType): string => {
  if (value == 'receive' || value == 'transfer_in') {
    return 'bg-success-subtle text-success'
  }
  if (value == 'ship' || value == 'transfer_out') {
    return 'bg-primary-subtle text-primary'
  }
  if (value == 'adjust') {
    return 'bg-info-subtle text-info'
  }
  if (value == 'damage') {
    return 'bg-danger-subtle text-danger'
  }
  if (value == 'reserve' || value == 'release') {
    return 'bg-warning-subtle text-warning'
  }
  return 'bg-secondary-subtle text-secondary'
}

const renderMovementType = (value: InventoryMovementType): string => {
  return `<span class="badge text-capitalize ${movementBadgeClass(value)}">${movementTypeLabel(value)}</span>`
}

const getTableColumns = (row: InventoryMovementModel): ITableColumn[] => {
  return [
    {
      name: 'Date',
      value: {
        text: row.createdAt.toDateTimeString(),
      },
    },
    {
      name: 'Type',
      value: {
        html: renderMovementType(row.movementType),
      },
    },
    {
      name: 'Quantity',
      value: {
        text: `${row.quantity}`,
      },
    },
    {
      name: 'Reference',
      value: {
        text: row.referenceId,
      },
    },
    {
      name: 'Inventory',
      value: {
        text: row.inventoryId,
      },
    },
    {
      name: 'Reason',
      value: {
        text: row.reason || 'No reason provided',
      },
    },
  ]
}

const loadMovements = async (withLoader: boolean = false): Promise<void> => {
  if (isRefreshing.value) {
    return
  }
  if (withLoader) {
    isLoading.value = true
  }
  isRefreshing.value = true
  try {
    const result = await app.services.inventoryMovements.getMany({
      page: 1,
      limit: 200,
      order: 'created_at.desc',
    })
    if (result.success) {
      movements.value = result.movements
      lastUpdatedAt.value = new DateTime().toDateTimeString()
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

onMounted(async () => {
  await loadMovements(true)
})
</script>

<template>
  <div class="inventory-movements-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Inventory Movements</h5>
          <p class="text-body-secondary mb-0">
            Review inbound, outbound, and adjustments for
            <strong>{{ setup?.location?.name || 'current location' }}</strong>.
          </p>
        </div>
        <div class="d-flex align-items-center mt-4 mt-md-0">
          <small v-if="lastUpdatedAt" class="text-body-secondary me-3">Last updated {{ lastUpdatedAt }}</small>
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="loadMovements(false)">
            <i class="bx bx-refresh me-2" />
            Refresh
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
              <p class="text-body-secondary mb-1">Records</p>
              <h4 class="mb-0">{{ summary.totalRecords }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Incoming Qty</p>
              <h4 class="mb-0 text-success">{{ summary.incomingQuantity }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Outgoing Qty</p>
              <h4 class="mb-0 text-primary">{{ summary.outgoingQuantity }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Adjustment Qty</p>
              <h4 class="mb-0 text-info">{{ summary.adjustmentQuantity }}</h4>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
          <h5 class="mb-0">Movement Records</h5>
          <div class="d-flex align-items-center gap-2">
            <label class="form-label mb-0 text-body-secondary">Filter Type</label>
            <select v-model="selectedMovementType" class="form-select">
              <option v-for="option in movementTypeOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="filteredMovements" :columns="getTableColumns" :checkboxes="true">
            <template #noResults>
              <p class="py-6">No movement records found.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped></style>
