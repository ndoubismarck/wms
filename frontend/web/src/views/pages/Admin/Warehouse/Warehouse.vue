<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useApp } from '@/app/app.ts'
import { useSetupStore } from '@/stores/setup_store.ts'
import { DateTime } from '@/app/core/date_time.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import Offcanvas, { type IOffcanvas } from '@/views/shared/components/Offcanvas.vue'
import type { LocationAisleModel } from '@/app/models/location_aisle_model.ts'
import type { LocationBayModel } from '@/app/models/location_bay_model.ts'
import type { LocationShelfModel } from '@/app/models/location_shelf_model.ts'
import type { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model.ts'
import type { LocationBinModel } from '@/app/models/location_bin_model.ts'

defineOptions({
  name: 'AdminWarehouseLocationsPage',
})

const app = useApp()
const setupStore = useSetupStore()

const setup = computed(() => setupStore.get().value)

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const lastUpdatedAt = ref<string>('')

const aisles = ref<LocationAisleModel[]>([])
const bays = ref<LocationBayModel[]>([])
const shelves = ref<LocationShelfModel[]>([])
const shelfLevels = ref<LocationShelfLevelModel[]>([])
const bins = ref<LocationBinModel[]>([])

const selectedAisle = ref<LocationAisleModel | null>(null)
const selectedBay = ref<LocationBayModel | null>(null)
const selectedShelf = ref<LocationShelfModel | null>(null)
const selectedShelfLevel = ref<LocationShelfLevelModel | null>(null)

const baysOffcanvas = ref<IOffcanvas | null>(null)
const shelvesOffcanvas = ref<IOffcanvas | null>(null)
const levelsOffcanvas = ref<IOffcanvas | null>(null)
const binsOffcanvas = ref<IOffcanvas | null>(null)

const locationRequestParams = computed(() => {
  const params: Record<string, string> = {}
  if (setup.value?.location?.id) {
    params.location_id = setup.value.location.id
  }
  return params
})

const selectedAisleBays = computed(() => {
  if (!selectedAisle.value) {
    return []
  }
  return bays.value.filter((bay) => bay.aisleId == selectedAisle.value?.id)
})

const selectedBayShelves = computed(() => {
  if (!selectedBay.value) {
    return []
  }
  return shelves.value.filter((shelf) => shelf.bayId == selectedBay.value?.id)
})

const selectedShelfLevels = computed(() => {
  if (!selectedShelf.value) {
    return []
  }
  return shelfLevels.value.filter((level) => level.shelfId == selectedShelf.value?.id)
})

const selectedLevelBins = computed(() => {
  if (!selectedShelfLevel.value) {
    return []
  }
  return bins.value.filter((bin) => bin.shelfLevelId == selectedShelfLevel.value?.id)
})

const totals = computed(() => {
  return {
    aisles: aisles.value.length,
    bays: bays.value.length,
    shelves: shelves.value.length,
    levels: shelfLevels.value.length,
    bins: bins.value.length,
  }
})

const fetchAllAisles = async (): Promise<LocationAisleModel[]> => {
  const result: LocationAisleModel[] = []
  let page = 1
  while (true) {
    const response = await app.services.locations.getAisles({
      ...locationRequestParams.value,
      page,
      limit: 200,
      order: 'created_at.asc',
    })
    if (!response.success) {
      break
    }
    result.push(...response.aisles)
    if (!response.pagination.hasNext) {
      break
    }
    page = response.pagination.next
  }
  return result
}

const fetchAllBays = async (): Promise<LocationBayModel[]> => {
  const result: LocationBayModel[] = []
  let page = 1
  while (true) {
    const response = await app.services.locations.getBays({
      ...locationRequestParams.value,
      page,
      limit: 200,
      order: 'created_at.asc',
    })
    if (!response.success) {
      break
    }
    result.push(...response.bays)
    if (!response.pagination.hasNext) {
      break
    }
    page = response.pagination.next
  }
  return result
}

const fetchAllShelves = async (): Promise<LocationShelfModel[]> => {
  const result: LocationShelfModel[] = []
  let page = 1
  while (true) {
    const response = await app.services.locations.getShelves({
      ...locationRequestParams.value,
      page,
      limit: 200,
      order: 'created_at.asc',
    })
    if (!response.success) {
      break
    }
    result.push(...response.shelves)
    if (!response.pagination.hasNext) {
      break
    }
    page = response.pagination.next
  }
  return result
}

const fetchAllLevels = async (): Promise<LocationShelfLevelModel[]> => {
  const result: LocationShelfLevelModel[] = []
  let page = 1
  while (true) {
    const response = await app.services.locations.getShelfLevels({
      ...locationRequestParams.value,
      page,
      limit: 200,
      order: 'created_at.asc',
    })
    if (!response.success) {
      break
    }
    result.push(...response.shelfLevels)
    if (!response.pagination.hasNext) {
      break
    }
    page = response.pagination.next
  }
  return result
}

const fetchAllBins = async (): Promise<LocationBinModel[]> => {
  const result: LocationBinModel[] = []
  let page = 1
  while (true) {
    const response = await app.services.locations.getBins({
      ...locationRequestParams.value,
      page,
      limit: 200,
      order: 'created_at.asc',
    })
    if (!response.success) {
      break
    }
    result.push(...response.bins)
    if (!response.pagination.hasNext) {
      break
    }
    page = response.pagination.next
  }
  return result
}

const loadLocationHierarchy = async (withLoader: boolean = false): Promise<void> => {
  if (isRefreshing.value) {
    return
  }
  if (withLoader) {
    isLoading.value = true
  }
  isRefreshing.value = true
  try {
    const [newAisles, newBays, newShelves, newLevels, newBins] = await Promise.all([
      fetchAllAisles(),
      fetchAllBays(),
      fetchAllShelves(),
      fetchAllLevels(),
      fetchAllBins(),
    ])

    aisles.value = newAisles
    bays.value = newBays
    shelves.value = newShelves
    shelfLevels.value = newLevels
    bins.value = newBins
    lastUpdatedAt.value = new DateTime().toDateTimeString()
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isRefreshing.value = false
    if (withLoader) {
      isLoading.value = false
    }
  }
}

const getShelvesByBayId = (bayId: string) => {
  return shelves.value.filter((shelf) => shelf.bayId == bayId)
}

const getLevelsByShelfId = (shelfId: string) => {
  return shelfLevels.value.filter((level) => level.shelfId == shelfId)
}

const getBinsByLevelId = (levelId: string) => {
  return bins.value.filter((bin) => bin.shelfLevelId == levelId)
}

const getBayCountForAisle = (aisleId: string): number => {
  return bays.value.filter((bay) => bay.aisleId == aisleId).length
}

const getShelfCountForAisle = (aisleId: string): number => {
  const bayIds = new Set(bays.value.filter((bay) => bay.aisleId == aisleId).map((bay) => bay.id))
  return shelves.value.filter((shelf) => bayIds.has(shelf.bayId)).length
}

const getLevelCountForAisle = (aisleId: string): number => {
  const bayIds = new Set(bays.value.filter((bay) => bay.aisleId == aisleId).map((bay) => bay.id))
  const shelfIds = new Set(shelves.value.filter((shelf) => bayIds.has(shelf.bayId)).map((shelf) => shelf.id))
  return shelfLevels.value.filter((level) => shelfIds.has(level.shelfId)).length
}

const getBinCountForAisle = (aisleId: string): number => {
  const bayIds = new Set(bays.value.filter((bay) => bay.aisleId == aisleId).map((bay) => bay.id))
  const shelfIds = new Set(shelves.value.filter((shelf) => bayIds.has(shelf.bayId)).map((shelf) => shelf.id))
  const levelIds = new Set(shelfLevels.value.filter((level) => shelfIds.has(level.shelfId)).map((level) => level.id))
  return bins.value.filter((bin) => levelIds.has(bin.shelfLevelId)).length
}

const getLevelCountForBay = (bayId: string): number => {
  const shelfIds = new Set(getShelvesByBayId(bayId).map((shelf) => shelf.id))
  return shelfLevels.value.filter((level) => shelfIds.has(level.shelfId)).length
}

const getBinCountForBay = (bayId: string): number => {
  const shelfIds = new Set(getShelvesByBayId(bayId).map((shelf) => shelf.id))
  const levelIds = new Set(shelfLevels.value.filter((level) => shelfIds.has(level.shelfId)).map((level) => level.id))
  return bins.value.filter((bin) => levelIds.has(bin.shelfLevelId)).length
}

const getBinCountForShelf = (shelfId: string): number => {
  const levelIds = new Set(getLevelsByShelfId(shelfId).map((level) => level.id))
  return bins.value.filter((bin) => levelIds.has(bin.shelfLevelId)).length
}

const getAisleColumns = (row: LocationAisleModel): ITableColumn[] => {
  return [
    {
      name: 'Aisle',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Bays',
      value: {
        text: `${getBayCountForAisle(row.id)}`,
      },
    },
    {
      name: 'Shelves',
      value: {
        text: `${getShelfCountForAisle(row.id)}`,
      },
    },
    {
      name: 'Levels',
      value: {
        text: `${getLevelCountForAisle(row.id)}`,
      },
    },
    {
      name: 'Bins',
      value: {
        text: `${getBinCountForAisle(row.id)}`,
      },
    },
  ]
}

const getBayColumns = (row: LocationBayModel): ITableColumn[] => {
  return [
    {
      name: 'Bay',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Shelves',
      value: {
        text: `${getShelvesByBayId(row.id).length}`,
      },
    },
    {
      name: 'Levels',
      value: {
        text: `${getLevelCountForBay(row.id)}`,
      },
    },
    {
      name: 'Bins',
      value: {
        text: `${getBinCountForBay(row.id)}`,
      },
    },
  ]
}

const getShelfColumns = (row: LocationShelfModel): ITableColumn[] => {
  return [
    {
      name: 'Shelf',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Levels',
      value: {
        text: `${getLevelsByShelfId(row.id).length}`,
      },
    },
    {
      name: 'Bins',
      value: {
        text: `${getBinCountForShelf(row.id)}`,
      },
    },
  ]
}

const getLevelColumns = (row: LocationShelfLevelModel): ITableColumn[] => {
  return [
    {
      name: 'Shelf Level',
      value: {
        text: row.name,
      },
    },
    {
      name: 'Bins',
      value: {
        text: `${getBinsByLevelId(row.id).length}`,
      },
    },
  ]
}

const getBinColumns = (row: LocationBinModel): ITableColumn[] => {
  return [
    {
      name: 'Bin',
      value: {
        text: row.name,
      },
    },
  ]
}

const openAisleBays = (row: LocationAisleModel) => {
  selectedAisle.value = row
  selectedBay.value = null
  selectedShelf.value = null
  selectedShelfLevel.value = null
  baysOffcanvas.value?.show()
}

const openBayShelves = (row: LocationBayModel) => {
  selectedBay.value = row
  selectedShelf.value = null
  selectedShelfLevel.value = null
  shelvesOffcanvas.value?.show()
}

const openShelfLevels = (row: LocationShelfModel) => {
  selectedShelf.value = row
  selectedShelfLevel.value = null
  levelsOffcanvas.value?.show()
}

const openLevelBins = (row: LocationShelfLevelModel) => {
  selectedShelfLevel.value = row
  binsOffcanvas.value?.show()
}

const handleHideBaysOffcanvas = () => {
  shelvesOffcanvas.value?.hide()
  levelsOffcanvas.value?.hide()
  binsOffcanvas.value?.hide()
  selectedAisle.value = null
  selectedBay.value = null
  selectedShelf.value = null
  selectedShelfLevel.value = null
}

const handleHideShelvesOffcanvas = () => {
  levelsOffcanvas.value?.hide()
  binsOffcanvas.value?.hide()
  selectedBay.value = null
  selectedShelf.value = null
  selectedShelfLevel.value = null
}

const handleHideLevelsOffcanvas = () => {
  binsOffcanvas.value?.hide()
  selectedShelf.value = null
  selectedShelfLevel.value = null
}

const handleHideBinsOffcanvas = () => {
  selectedShelfLevel.value = null
}

onMounted(async () => {
  await loadLocationHierarchy(true)
})
</script>

<template>
  <div class="warehouse-locations-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Locations</h5>
          <p class="text-body-secondary mb-0">
            Explore warehouse structure from aisle to bin for
            <strong>{{ setup?.location?.name || 'current location' }}</strong>.
          </p>
        </div>
        <div class="d-flex align-items-center mt-4 mt-md-0">
          <small v-if="lastUpdatedAt" class="text-body-secondary me-3">Last updated {{ lastUpdatedAt }}</small>
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="loadLocationHierarchy(false)">
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
        <div class="col-12 col-md-6 col-xl-2 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Aisles</p>
              <h4 class="mb-0">{{ totals.aisles }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-2 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Bays</p>
              <h4 class="mb-0">{{ totals.bays }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-2 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Shelves</p>
              <h4 class="mb-0">{{ totals.shelves }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-2 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Levels</p>
              <h4 class="mb-0">{{ totals.levels }}</h4>
            </div>
          </div>
        </div>
        <div class="col-12 col-md-6 col-xl-2 mb-6">
          <div class="card h-100">
            <div class="card-body">
              <p class="text-body-secondary mb-1">Bins</p>
              <h4 class="mb-0">{{ totals.bins }}</h4>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Aisle Overview</h5>
          <span class="badge bg-label-primary">{{ aisles.length }}</span>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="aisles" :columns="getAisleColumns" :checkboxes="true">
            <template #actions="{ row }">
              <button type="button" class="btn btn-sm btn-outline-primary" @click="openAisleBays(row)">
                View Bays
              </button>
            </template>
            <template #noResults>
              <p class="py-6">No aisles found.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </div>

  <offcanvas
    ref="baysOffcanvas"
    :show="false"
    @hide="handleHideBaysOffcanvas"
    title="Aisle Bays">
    <template #body>
      <div class="card mb-0">
        <div class="card-header">
          <h6 class="mb-0">
            {{ selectedAisle ? `Aisle ${selectedAisle.name}` : 'Aisle' }}
          </h6>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="selectedAisleBays" :columns="getBayColumns" :checkboxes="true">
            <template #actions="{ row }">
              <button type="button" class="btn btn-sm btn-outline-primary" @click="openBayShelves(row)">
                View Shelves
              </button>
            </template>
            <template #noResults>
              <p class="py-6">No bays found for this aisle.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </offcanvas>

  <offcanvas
    ref="shelvesOffcanvas"
    :show="false"
    :parent="baysOffcanvas || undefined"
    @hide="handleHideShelvesOffcanvas"
    title="Bay Shelves">
    <template #body>
      <div class="card mb-0">
        <div class="card-header">
          <h6 class="mb-0">
            {{ selectedBay ? `Bay ${selectedBay.name}` : 'Bay' }}
          </h6>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="selectedBayShelves" :columns="getShelfColumns" :checkboxes="true">
            <template #actions="{ row }">
              <button type="button" class="btn btn-sm btn-outline-primary" @click="openShelfLevels(row)">
                View Levels
              </button>
            </template>
            <template #noResults>
              <p class="py-6">No shelves found for this bay.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </offcanvas>

  <offcanvas
    ref="levelsOffcanvas"
    :show="false"
    :parent="shelvesOffcanvas || undefined"
    @hide="handleHideLevelsOffcanvas"
    title="Shelf Levels">
    <template #body>
      <div class="card mb-0">
        <div class="card-header">
          <h6 class="mb-0">
            {{ selectedShelf ? `Shelf ${selectedShelf.name}` : 'Shelf' }}
          </h6>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="selectedShelfLevels" :columns="getLevelColumns" :checkboxes="true">
            <template #actions="{ row }">
              <button type="button" class="btn btn-sm btn-outline-primary" @click="openLevelBins(row)">
                View Bins
              </button>
            </template>
            <template #noResults>
              <p class="py-6">No levels found for this shelf.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </offcanvas>

  <offcanvas
    ref="binsOffcanvas"
    :show="false"
    :parent="levelsOffcanvas || undefined"
    @hide="handleHideBinsOffcanvas"
    title="Level Bins">
    <template #body>
      <div class="card mb-0">
        <div class="card-header">
          <h6 class="mb-0">
            {{ selectedShelfLevel ? `Level ${selectedShelfLevel.name}` : 'Level' }}
          </h6>
        </div>
        <div class="card-body p-0">
          <vue-table :rows="selectedLevelBins" :columns="getBinColumns" :checkboxes="true">
            <template #noResults>
              <p class="py-6">No bins found for this level.</p>
            </template>
          </vue-table>
        </div>
      </div>
    </template>
  </offcanvas>
</template>

<style scoped></style>
