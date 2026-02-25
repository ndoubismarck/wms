<script setup lang="ts">
import VueFormTextareaInput from '@/views/shared/components/VueForm/VueFormTextareaInput.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormSelectInput, { type ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import VueFormButton from '@/views/shared/components/VueForm/VueFormButton.vue'
import type { LocationAisleModel } from '@/app/models/location_aisle_model.ts'
import { useLocationsAislesStore } from '@/stores/location_aisles_store.ts'
import { useApp } from '@/app/app.ts'
import { useLocationsBaysStore } from '@/stores/location_bays_store.ts'
import type { LocationBayModel } from '@/app/models/location_bay_model.ts'
import VueForm, { type VueFormData } from '@/views/shared/components/VueForm/VueForm.vue'
import { useLocationsShelvesStore } from '@/stores/location_shelves_store.ts'
import { useLocationsShelfLevelsStore } from '@/stores/location_shelf_levels_store.ts'
import { useProductCategoriesStore } from '@/stores/product_categories_store.ts'
import type { LocationShelfModel } from '@/app/models/location_shelf_model.ts'
import type { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model.ts'
import type { ProductCategoryModel } from '@/app/models/product_category_model.ts'
import { EAlertMessageType, type IAlertMessage } from '@/app/types.ts'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import type { LocationBinModel } from '@/app/models/location_bin_model.ts'

const app = useApp()
const baysStore = useLocationsBaysStore()
const aislesStore = useLocationsAislesStore()
const shelvesStore = useLocationsShelvesStore()
const shelfLevelsStore = useLocationsShelfLevelsStore()
const binCategoriesStore = useProductCategoriesStore()

const emits = defineEmits([
  'close',
  'submit',
  'submitted',
  'addBay',
  'addAisle',
  'addShelf',
  'addShelfLevel',
  'addBinCategory',
])

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)
const isLoading = ref<boolean>(true)
const isSubmitting = ref<boolean>(false)

const selectedBayId = ref<string | null>(null)
const selectedAisleId = ref<string | null>(null)
const selectedShelfId = ref<string | null>(null)
const selectedShelfLevelId = ref<string | null>(null)
const selectedBinCategoryId = ref<string | null>(null)

const allAisles = computed(() => aislesStore.get().value)
const allBays = computed(() => baysStore.get({ aisleId: selectedAisleId.value }).value)
const allShelves = computed(() => shelvesStore.get({ bayId: selectedBayId.value }).value)
const allShelfLevels = computed(() => shelfLevelsStore.get({ shelfId: selectedShelfId.value }).value)
const allBinCategories = computed(() => binCategoriesStore.get().value)

const bayOptions = ref<ISelectInputOption[]>([])
const aisleOptions = ref<ISelectInputOption[]>([])
const shelfOptions = ref<ISelectInputOption[]>([])
const shelfLevelOptions = ref<ISelectInputOption[]>([])
const binCategoryOptions = ref<ISelectInputOption[]>([])

const SEARCH_DEBOUNCE_MS = 350
const searchTimers: Record<string, ReturnType<typeof setTimeout> | undefined> = {}

const getSearchParams = (search?: string) => {
  const trimmed = search?.trim()
  if (!trimmed) {
    return {}
  }
  return {
    search: trimmed,
  }
}

const runDebouncedSearch = (key: string, handler: () => Promise<void>) => {
  if (searchTimers[key]) {
    clearTimeout(searchTimers[key])
  }
  searchTimers[key] = setTimeout(() => {
    void handler()
  }, SEARCH_DEBOUNCE_MS)
}

watch(
  allBays,
  () => {
    loadBayOptions()
  },
  { deep: true },
)

watch(
  allAisles,
  () => {
    loadAisleOptions()
  },
  { deep: true },
)

watch(
  allShelves,
  () => {
    loadShelfOptions()
  },
  { deep: true },
)

watch(
  allShelfLevels,
  () => {
    loadShelfLevelOptions()
  },
  { deep: true },
)

watch(
  allBinCategories,
  () => {
    loadBinCategoryOptions()
  },
  { deep: true },
)

const reset = () => {
  if (form.value) {
    form.value.reset()
  }
}

const getBays = async (search?: string) => {
  if (selectedAisleId.value) {
    try {
      const result = await app.services.locations.getBays({
        page: 1,
        aisle_id: selectedAisleId.value,
        ...getSearchParams(search),
      })
      if (result.success) {
        baysStore.set(result.bays)
      }
    } catch (ex) {
      app.services.logger.error(ex)
    }
  }
}

const getAisles = async (search?: string) => {
  try {
    const result = await app.services.locations.getAisles({
      page: 1,
      ...getSearchParams(search),
    })
    if (result.success) {
      aislesStore.set(result.aisles)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getShelves = async (search?: string) => {
  if (selectedBayId.value) {
    try {
      const result = await app.services.locations.getShelves({
        page: 1,
        bay_id: selectedBayId.value,
        ...getSearchParams(search),
      })
      if (result.success) {
        shelvesStore.set(result.shelves)
      }
    } catch (ex) {
      app.services.logger.error(ex)
    }
  }
}

const getShelfLevels = async (search?: string) => {
  if (selectedShelfId.value) {
    try {
      const result = await app.services.locations.getShelfLevels({
        page: 1,
        shelf_id: selectedShelfId.value,
        ...getSearchParams(search),
      })
      if (result.success) {
        shelfLevelsStore.set(result.shelfLevels)
      }
    } catch (ex) {
      app.services.logger.error(ex)
    }
  }
}

const getBinCategories = async (search?: string) => {
  try {
    const result = await app.services.products.getCategories({
      page: 1,
      ...getSearchParams(search),
    })
    if (result.success) {
      binCategoriesStore.set(result.categories)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const loadAisleOptions = () => {
  if (allAisles.value) {
    allAisles.value.forEach((val) => {
      const index = aisleOptions.value.findIndex((value) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selectedAisleId.value,
      }
      if (index < 0) {
        aisleOptions.value.push(item)
      } else {
        aisleOptions.value[index] = item
      }
    })
  }
}

const loadBayOptions = () => {
  if (allBays.value) {
    allBays.value.forEach((val) => {
      const index = bayOptions.value.findIndex((value) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selectedBayId.value,
      }
      if (index < 0) {
        bayOptions.value.push(item)
      } else {
        bayOptions.value[index] = item
      }
    })
  }
}

const loadShelfOptions = () => {
  if (allShelves.value) {
    allShelves.value.forEach((val) => {
      const index = shelfOptions.value.findIndex((value) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selectedShelfId.value,
      }
      if (index < 0) {
        shelfOptions.value.push(item)
      } else {
        shelfOptions.value[index] = item
      }
    })
  }
}

const loadShelfLevelOptions = () => {
  if (allShelfLevels.value) {
    allShelfLevels.value.forEach((val) => {
      const index = shelfLevelOptions.value.findIndex((value) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selectedShelfLevelId.value,
      }
      if (index < 0) {
        shelfLevelOptions.value.push(item)
      } else {
        shelfLevelOptions.value[index] = item
      }
    })
  }
}

const loadBinCategoryOptions = () => {
  if (allBinCategories.value) {
    allBinCategories.value.forEach((val) => {
      const index = binCategoryOptions.value.findIndex((value) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selectedBinCategoryId.value,
      }
      if (index < 0) {
        binCategoryOptions.value.push(item)
      } else {
        binCategoryOptions.value[index] = item
      }
    })
  }
}

const handleSubmit = async (data: VueFormData) => {
  emits('submit')
  isSubmitting.value = true
  let success = false
  let payload: LocationBinModel | undefined
  try {
    const result = await app.services.locations.addBin(data)
    await app.helpers.async.sleep(2000)
    if (result.success) {
      success = true
      payload = result.bin
      form.value?.reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Bin successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Bin could not be added.'),
        },
      }
    }
  } catch (ex) {
    alert.value = {
      type: EAlertMessageType.Error,
      body: {
        text: app.helpers.i18n.messageFromError(ex),
      },
    }
    app.services.logger.error(ex)
  }
  isSubmitting.value = false
  emits('submitted', success, payload)
  await app.helpers.async.sleep(5000)
  if (success) {
    emits('close')
  }
}

const handleAddAisle = () => {
  emits('addAisle')
}

const handleAddBay = () => {
  if (selectedAisleId.value) {
    emits('addBay', selectedAisleId.value)
  }
}

const handleAddShelf = () => {
  emits('addShelf', selectedBayId.value)
}

const handleAddShelfLevel = () => {
  emits('addShelfLevel', selectedShelfId.value)
}

const handleAddBinCategory = () => {
  emits('addBinCategory')
}

const handleAisleInputChange = async (value?: ISelectInputOption) => {
  if (value) {
    selectedAisleId.value = value.value
    await getBays()
  } else {
    selectedBayId.value = null
    selectedAisleId.value = null
    selectedShelfId.value = null
    selectedShelfLevelId.value = null
    form.value?.resetField('bay_id')
  }
}

const handleBayInputChange = async (value?: ISelectInputOption) => {
  if (value) {
    selectedBayId.value = value.value
    await getShelves()
  } else {
    selectedBayId.value = null
    selectedShelfId.value = null
    selectedShelfLevelId.value = null
    form.value?.resetField('shelf_id')
  }
}

const handleShelfInputChange = async (value?: ISelectInputOption) => {
  if (value) {
    selectedShelfId.value = value.value
    await getShelfLevels()
  } else {
    selectedShelfId.value = null
    selectedShelfLevelId.value = null
    form.value?.resetField('shelf_level_id')
  }
}

const handleAisleSearch = (search: string) => {
  runDebouncedSearch('aisle', async () => {
    await getAisles(search)
  })
}

const handleBaySearch = (search: string) => {
  if (!selectedAisleId.value) {
    return
  }
  runDebouncedSearch('bay', async () => {
    await getBays(search)
  })
}

const handleShelfSearch = (search: string) => {
  if (!selectedBayId.value) {
    return
  }
  runDebouncedSearch('shelf', async () => {
    await getShelves(search)
  })
}

const handleShelfLevelSearch = (search: string) => {
  if (!selectedShelfId.value) {
    return
  }
  runDebouncedSearch('shelf-level', async () => {
    await getShelfLevels(search)
  })
}

const handleBinCategorySearch = (search: string) => {
  runDebouncedSearch('bin-category', async () => {
    await getBinCategories(search)
  })
}

onMounted(async () => {
  isLoading.value = true
  await getAisles()
  await getBinCategories()
  await app.helpers.async.sleep(3000)
  isLoading.value = false
})

onBeforeUnmount(() => {
  Object.values(searchTimers).forEach((timer) => {
    if (timer) {
      clearTimeout(timer)
    }
  })
})

defineExpose({
  reset,
  isLoading,
  isSubmitting,
})

app.events.on('forms.location.add.aisle.submitted', async (data: LocationAisleModel) => {
  selectedAisleId.value = data.id
  await app.helpers.async.sleep(10)
  aislesStore.set(data)
  await getBays()
})

app.events.on('forms.location.add.bay.submitted', async (data: LocationBayModel) => {
  selectedBayId.value = data.id
  await app.helpers.async.sleep(10)
  baysStore.set(data)
  await getShelves()
})

app.events.on('forms.location.add.shelf.submitted', async (data: LocationShelfModel) => {
  selectedShelfId.value = data.id
  await app.helpers.async.sleep(10)
  shelvesStore.set(data)
  await getShelfLevels()
})

app.events.on('forms.location.add.shelf.level.submitted', async (data: LocationShelfLevelModel) => {
  selectedShelfLevelId.value = data.id
  await app.helpers.async.sleep(10)
  shelfLevelsStore.set(data)
})

app.events.on('forms.location.add.bin.category.submitted', async (data: ProductCategoryModel) => {
  selectedBinCategoryId.value = data.id
  await app.helpers.async.sleep(10)
  binCategoriesStore.set(data)
})
</script>

<template>
  <div class="card">
    <div class="card-body">
      <vue-form ref="form" @submit="handleSubmit">
        <alert-message v-if="alert" :params="alert" />
        <div class="mb-6">
          <vue-form-select-input
            name="aisle_id"
            @change="handleAisleInputChange"
            @search="handleAisleSearch"
            :options="aisleOptions"
            :disabled="isSubmitting || aisleOptions.length == 0"
            label="Aisle"
            placeholder="Select Aisle">
            <template #addon-end>
              <button
                @click="handleAddAisle"
                :disabled="isSubmitting || aisleOptions.length == 0"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            name="bay_id"
            @change="handleBayInputChange"
            @search="handleBaySearch"
            :options="bayOptions"
            :disabled="isSubmitting || !selectedAisleId"
            label="Bay"
            placeholder="Select Bay">
            <template #addon-end>
              <button
                @click="handleAddBay"
                :disabled="isSubmitting || !selectedAisleId"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            name="shelf_id"
            @change="handleShelfInputChange"
            @search="handleShelfSearch"
            :options="shelfOptions"
            :disabled="isSubmitting || !selectedBayId"
            type="text"
            label="Shelf"
            placeholder="Select Shelf">
            <template #addon-end>
              <button
                @click="handleAddShelf"
                :disabled="isSubmitting || !selectedBayId"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            name="shelf_level_id"
            @search="handleShelfLevelSearch"
            :options="shelfLevelOptions"
            :disabled="isSubmitting || !selectedShelfId"
            label="Shelf Level"
            placeholder="Select Shelf Level">
            <template #addon-end>
              <button
                @click="handleAddShelfLevel"
                :disabled="isSubmitting || !selectedShelfId"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            name="category_id"
            @search="handleBinCategorySearch"
            :options="binCategoryOptions"
            :disabled="isSubmitting"
            label="Category"
            placeholder="Category">
            <template #addon-end>
              <button
                @click="handleAddBinCategory"
                :disabled="isSubmitting"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-input name="code" :disabled="isSubmitting" type="text" label="Code" placeholder="Code" />
        </div>
        <div class="mb-6">
          <vue-form-input name="name" :disabled="isSubmitting" type="text" label="Name" placeholder="Name" />
        </div>
        <div class="mb-6">
          <vue-form-textarea-input
            name="description"
            :disabled="isSubmitting"
            label="Description"
            placeholder="Enter description" />
        </div>
        <vue-form-button :disabled="isSubmitting" text="Submit" />
      </vue-form>
    </div>
  </div>
</template>

<style scoped></style>
