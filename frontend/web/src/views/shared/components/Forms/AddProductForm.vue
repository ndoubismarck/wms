<script setup lang="ts">
import VueFormSelectInput, {
  type ISelectInputOption
} from '@/views/shared/components/VueForm/VueFormSelectInput.vue'
import VueForm, {type VueFormData} from '@/views/shared/components/VueForm/VueForm.vue'
import {computed, onBeforeUnmount, onMounted, ref, type UnwrapRef, watch} from 'vue'
import {useApp} from '@/app/app'
import {useLocationsBinsStore} from '@/stores/location_bins_store.ts'
import {EAlertMessageType, type IAlertMessage} from '@/app/types.ts'
import type {ProductModel} from '@/app/models/product_model.ts'
import {useLocationsAislesStore} from '@/stores/location_aisles_store.ts'
import {useLocationsShelvesStore} from '@/stores/location_shelves_store.ts'
import {useLocationsShelfLevelsStore} from '@/stores/location_shelf_levels_store.ts'
import {useLocationsBaysStore} from '@/stores/location_bays_store.ts'
import AlertMessage from '@/views/shared/components/AlertMessage.vue'
import VueFormInput from '@/views/shared/components/VueForm/VueFormInput.vue'
import VueFormHtmlInput from '@/views/shared/components/VueForm/VueFormHtmlInput.vue'
import VueAccordionWizard from '@/views/shared/components/VueAccordionWizard/VueAccordionWizard.vue'
import VueAccordionWizardItem
  from '@/views/shared/components/VueAccordionWizard/VueAccordionWizardItem.vue'
import VueFormFileInput from '@/views/shared/components/VueForm/VueFormFileInput.vue'
import {useProductBrandsStore} from '@/stores/product_brands_store.ts'
import {
  ProductBrandCreatedEvent,
  ProductCategoryCreatedEvent,
  ProductSubcategoryCreatedEvent,
} from '@/views/shared/types/events.ts'
import {ProductBrandModel} from '@/app/models/product_brand_model.ts'
import {useProductCategoriesStore} from '@/stores/product_categories_store.ts'
import {useProductSubcategoriesStore} from '@/stores/product_subcategories_store.ts'
import type {ProductCategoryModel} from '@/app/models/product_category_model.ts'
import type {ProductSubcategoryModel} from '@/app/models/product_subcategory_model.ts'
import {useProductAttributesStore} from "@/stores/product_attributes_store.ts";
import type {ProductAttributeModel} from '@/app/models/product_attribute_model.ts'
import VueFormTextareaInput from "@/views/shared/components/VueForm/VueFormTextareaInput.vue";

const app = useApp()
const productBrandsStore = useProductBrandsStore()
const productCategoryStore = useProductCategoriesStore()
const productSubcategoryStore = useProductSubcategoriesStore()
const productAttributesStore = useProductAttributesStore()
const locationAislesStore = useLocationsAislesStore()
const locationBaysStore = useLocationsBaysStore()
const locationShelvesStore = useLocationsShelvesStore()
const locationShelfLevelsStore = useLocationsShelfLevelsStore()
const locationBinsStore = useLocationsBinsStore()

const selectedProductBrandId = ref<string | null>()
const selectedProductCategoryId = ref<string | null>()
const selectedProductSubcategoryId = ref<string | null>()
const selectedLocationAisleId = ref<string | null>()
const selectedLocationBayId = ref<string | null>()
const selectedLocationShelfId = ref<string | null>()
const selectedLocationShelfLevelId = ref<string | null>()
const selectedLocationBinId = ref<string | null>()

const productBrands = computed(() => productBrandsStore.get().value)
const productAttributes = computed(() => productAttributesStore.get().value)
const locationAisles = computed(() => locationAislesStore.get().value)
const productCategories = computed(
  () =>
    productCategoryStore.get({
      brandId: computed(() => selectedProductBrandId.value).value,
    }).value,
)
const productSubcategories = computed(
  () =>
    productSubcategoryStore.get({
      categoryId: computed(() => selectedProductCategoryId.value).value,
    }).value,
)
const locationBays = computed(() => locationBaysStore.get({aisleId: selectedLocationAisleId.value}).value)
const locationShelves = computed(() => locationShelvesStore.get({bayId: selectedLocationBayId.value}).value)
const locationShelfLevels = computed(
  () => locationShelfLevelsStore.get({shelfId: selectedLocationShelfId.value}).value,
)
const locationBins = computed(() => locationBinsStore.get({shelfLevelId: selectedLocationShelfLevelId.value}).value)

const form = ref<typeof VueForm | null>(null)
const alert = ref<IAlertMessage | null>(null)

const isLoading = ref<boolean>(true)
const isSubmitting = ref<boolean>(false)

const formData = ref<VueFormData>()
const accordionWizardRef = ref<typeof VueAccordionWizard>()

const productBrandOptions = ref<ISelectInputOption[]>([])
const productCategoryOptions = ref<ISelectInputOption[]>([])
const productSubcategoryOptions = ref<ISelectInputOption[]>([])
const locationAisleOptions = ref<ISelectInputOption[]>([])
const locationBayOptions = ref<ISelectInputOption[]>([])
const locationShelfOptions = ref<ISelectInputOption[]>([])
const locationShelfLevelOptions = ref<ISelectInputOption[]>([])
const locationBinOptions = ref<ISelectInputOption[]>([])

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

const getAttributeID = (item: ProductAttributeModel): string => {
  return item.id || item.attributeId
}

const getAttributeInputName = (item: ProductAttributeModel): string => {
  const attributeID = getAttributeID(item)
  const dataTypeByFieldType: Record<string, string> = {
    text: 'string',
    textarea: 'string',
    number: 'float',
    date: 'date',
    select: 'option',
    radio: 'option',
    checkbox: 'bool',
    multi_check: 'options',
  }
  const dataType = dataTypeByFieldType[item.fieldType] || 'string'
  return `attributes.${attributeID}.${dataType}`
}

const getAttributeValidation = (item: ProductAttributeModel): string | undefined => {
  if (!item.validationSchema || item.validationSchema.length == 0) {
    return undefined
  }
  return item.validationSchema.join('|')
}

const getAttributeSelectOptions = (item: ProductAttributeModel): ISelectInputOption[] => {
  if (item.options && item.options.length > 0) {
    return item.options
      .map((option) => {
        const value = option.optionId || option.id
        if (!value) {
          return null
        }
        return {
          label: option.label || value,
          value: value,
        }
      })
      .filter((option): option is ISelectInputOption => option != null)
  }

  if (item.fieldType == 'checkbox') {
    return [
      {
        label: 'Yes',
        value: true,
      },
      {
        label: 'No',
        value: false,
      },
    ]
  }

  return []
}

watch(
  productBrands,
  () => {
    loadOptions(productBrands, productBrandOptions, selectedProductBrandId)
  },
  {deep: true},
)

watch(
  productCategories,
  () => {
    loadOptions(productCategories, productCategoryOptions, selectedProductCategoryId)
  },
  {deep: true},
)

watch(
  productSubcategories,
  () => {
    loadOptions(productSubcategories, productSubcategoryOptions, selectedProductSubcategoryId)
  },
  {deep: true},
)

watch(
  locationAisles,
  () => {
    loadOptions(locationAisles, locationAisleOptions, selectedLocationAisleId)
  },
  {deep: true},
)

watch(
  locationBays,
  () => {
    loadOptions(locationBays, locationBayOptions, selectedLocationBayId)
  },
  {deep: true},
)

watch(
  locationShelves,
  () => {
    loadOptions(locationShelves, locationShelfOptions, selectedLocationShelfId)
  },
  {deep: true},
)

watch(
  locationShelfLevels,
  () => {
    loadOptions(locationShelfLevels, locationShelfLevelOptions, selectedLocationShelfLevelId)
  },
  {deep: true},
)

watch(
  locationBins,
  () => {
    loadOptions(locationBins, locationBinOptions, selectedLocationBinId)
  },
  {deep: true},
)

watch(
  selectedProductBrandId,
  async (newValue, oldValue) => {
    if (selectedProductBrandId.value) {
      if (newValue != oldValue) {
        productCategoryOptions.value = []
        selectedProductCategoryId.value = null
      }
      await getProductCategories()
    } else {
      productCategoryOptions.value = []
      selectedProductCategoryId.value = null
    }
  },
  {deep: true},
)

watch(
  selectedProductCategoryId,
  async (newValue, oldValue) => {
    if (selectedProductCategoryId.value) {
      if (newValue != oldValue) {
        productSubcategoryOptions.value = []
        selectedProductSubcategoryId.value = null
      }
      await getProductSubcategories()
    } else {
      productSubcategoryOptions.value = []
      selectedProductSubcategoryId.value = null
    }
  },
  {deep: true},
)

watch(
  selectedLocationAisleId,
  async (newValue, oldValue) => {
    if (selectedLocationAisleId.value) {
      if (newValue != oldValue) {
        locationBayOptions.value = []
        selectedLocationBayId.value = null
      }
      await getLocationBays()
    } else {
      locationBayOptions.value = []
      selectedLocationBayId.value = null
    }
  },
  {deep: true},
)

watch(
  selectedLocationBayId,
  async (newValue, oldValue) => {
    if (selectedLocationBayId.value) {
      if (newValue != oldValue) {
        locationShelfOptions.value = []
        selectedLocationShelfId.value = null
      }
      await getLocationShelves()
    } else {
      locationShelfOptions.value = []
      selectedLocationShelfId.value = null
    }
  },
  {deep: true},
)

watch(
  selectedLocationShelfId,
  async (newValue, oldValue) => {
    if (selectedLocationShelfId.value) {
      if (newValue != oldValue) {
        locationShelfLevelOptions.value = []
        selectedLocationShelfLevelId.value = null
      }
      await getLocationShelfLevels()
    } else {
      locationShelfLevelOptions.value = []
      selectedLocationShelfLevelId.value = null
    }
  },
  {deep: true},
)

watch(
  selectedLocationShelfLevelId,
  async (newValue, oldValue) => {
    if (selectedLocationShelfLevelId.value) {
      if (newValue != oldValue) {
        locationBinOptions.value = []
        selectedLocationBinId.value = null
      }
      await getLocationBins()
    } else {
      locationBinOptions.value = []
      selectedLocationBinId.value = null
    }
  },
  {deep: true},
)

const reset = () => {
  if (form.value) {
    form.value.reset()
  }
}

const loadOptions = (source: UnwrapRef<any>, options: UnwrapRef<any>, selected: UnwrapRef<any>) => {
  if (source.value) {
    source.value.forEach((val: any) => {
      const index = options.value.findIndex((value: any) => value.value == val.id)
      const item = {
        label: val.name,
        value: val.id,
        selected: val.id == selected.value
      }
      if (index < 0) {
        options.value.push(item)
      } else {
        options.value[index] = item
      }
    })
  }
}

const handleNext = async (data: VueFormData) => {
  Object.keys(data).forEach((key) => {
    if (!formData.value) {
      formData.value = {}
    }
    formData.value[key] = data[key]
  })
  if (accordionWizardRef.value) {
    if (accordionWizardRef.value.hasNext()) {
      accordionWizardRef.value.next()
    }
  }
  console.log(data)
}

const handleSubmit = async (data: VueFormData) => {
  return
  emits('submit')
  isSubmitting.value = true
  let success = false
  let payload: ProductModel | undefined
  try {
    const result = await app.services.products.add(data)
    await app.helpers.async.sleep(2000)
    if (result.success) {
      success = true
      payload = result.product
      form.value?.reset()
      alert.value = {
        type: EAlertMessageType.Success,
        body: {
          text: app.helpers.i18n.message('Product successfully added.'),
        },
      }
    } else {
      alert.value = {
        type: EAlertMessageType.Error,
        body: {
          text: app.helpers.i18n.message('Product could not be added.'),
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
}

const getProductBrands = async (search?: string) => {
  try {
    const result = await app.services.products.getBrands({
      page: 1,
      limit: 50,
      ...getSearchParams(search),
    })
    if (result.success) {
      productBrandsStore.set(result.brands)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getProductCategories = async (search?: string) => {
  try {
    const result = await app.services.products.getCategories({
      page: 1,
      limit: 50,
      brand_id: selectedProductBrandId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      productCategoryStore.set(result.categories)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getProductSubcategories = async (search?: string) => {
  try {
    const result = await app.services.products.getSubcategories({
      page: 1,
      limit: 50,
      category_id: selectedProductCategoryId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      productSubcategoryStore.set(result.subcategories)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getProductAttributes = async () => {
  try {
    const result = await app.services.products.getAttributes({
      page: 1,
      limit: 50,
    })
    if (result.success) {
      productAttributesStore.set(result.attributes)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getLocationAisles = async (search?: string) => {
  try {
    const result = await app.services.locations.getAisles({
      page: 1,
      limit: 50,
      ...getSearchParams(search),
    })
    if (result.success) {
      locationAislesStore.set(result.aisles)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getLocationBays = async (search?: string) => {
  try {
    const result = await app.services.locations.getBays({
      page: 1,
      limit: 50,
      aisle_id: selectedLocationAisleId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      locationBaysStore.set(result.bays)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getLocationShelves = async (search?: string) => {
  try {
    const result = await app.services.locations.getShelves({
      page: 1,
      limit: 50,
      bay_id: selectedLocationBayId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      locationShelvesStore.set(result.shelves)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getLocationShelfLevels = async (search?: string) => {
  try {
    const result = await app.services.locations.getShelfLevels({
      page: 1,
      limit: 50,
      shelf_id: selectedLocationShelfId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      locationShelfLevelsStore.set(result.shelfLevels)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const getLocationBins = async (search?: string) => {
  try {
    const result = await app.services.locations.getBins({
      page: 1,
      limit: 50,
      shelf_level_id: selectedLocationShelfLevelId.value,
      ...getSearchParams(search),
    })
    if (result.success) {
      locationBinsStore.set(result.bins)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const handleProductBrandInputChange = async (value: ISelectInputOption) => {
  selectedProductBrandId.value = value?.value
}

const handleProductCategoryInputChange = async (value: ISelectInputOption) => {
  selectedProductCategoryId.value = value?.value
}

const handleLocationBayInputChange = async (value: ISelectInputOption) => {
  selectedLocationBayId.value = value?.value
}

const handleLocationAisleInputChange = async (value: ISelectInputOption) => {
  selectedLocationAisleId.value = value?.value
}

const handleLocationShelfInputChange = async (value: ISelectInputOption) => {
  selectedLocationShelfId.value = value?.value
}

const handleLocationShelfLevelInputChange = async (value: ISelectInputOption) => {
  selectedLocationShelfLevelId.value = value?.value
}

const handleProductBrandSearch = (search: string) => {
  runDebouncedSearch('product-brand', async () => {
    await getProductBrands(search)
  })
}

const handleProductCategorySearch = (search: string) => {
  if (!selectedProductBrandId.value) {
    return
  }
  runDebouncedSearch('product-category', async () => {
    await getProductCategories(search)
  })
}

const handleProductSubcategorySearch = (search: string) => {
  if (!selectedProductCategoryId.value) {
    return
  }
  runDebouncedSearch('product-subcategory', async () => {
    await getProductSubcategories(search)
  })
}

const handleLocationAisleSearch = (search: string) => {
  runDebouncedSearch('location-aisle', async () => {
    await getLocationAisles(search)
  })
}

const handleLocationBaySearch = (search: string) => {
  if (!selectedLocationAisleId.value) {
    return
  }
  runDebouncedSearch('location-bay', async () => {
    await getLocationBays(search)
  })
}

const handleLocationShelfSearch = (search: string) => {
  if (!selectedLocationBayId.value) {
    return
  }
  runDebouncedSearch('location-shelf', async () => {
    await getLocationShelves(search)
  })
}

const handleLocationShelfLevelSearch = (search: string) => {
  if (!selectedLocationShelfId.value) {
    return
  }
  runDebouncedSearch('location-shelf-level', async () => {
    await getLocationShelfLevels(search)
  })
}

const handleLocationBinSearch = (search: string) => {
  if (!selectedLocationShelfLevelId.value) {
    return
  }
  runDebouncedSearch('location-bin', async () => {
    await getLocationBins(search)
  })
}

app.events.on(ProductBrandCreatedEvent, async (data: ProductBrandModel) => {
  selectedProductBrandId.value = data.id
  productBrandsStore.set(data)
})

app.events.on(ProductCategoryCreatedEvent, async (data: ProductCategoryModel) => {
  selectedProductCategoryId.value = data.id
  productCategoryStore.set(data)
})

app.events.on(ProductSubcategoryCreatedEvent, async (data: ProductSubcategoryModel) => {
  selectedProductSubcategoryId.value = data.id
  productSubcategoryStore.set(data)
})

defineExpose({
  reset,
  isLoading,
  isSubmitting,
})

const emits = defineEmits([
  'submit',
  'submitted',
  'addProductBrand',
  'addProductCategory',
  'addProductSubcategory',
  'addLocationAisle',
  'addLocationBay',
  'addLocationShelf',
  'addLocationShelfLevel',
  'addLocationBin',
])

onMounted(async () => {
  isLoading.value = true
  await getProductBrands()
  await getProductAttributes()
  await getLocationAisles()
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
</script>

<template>
  <alert-message v-if="alert" :params="alert"/>
  <vue-accordion-wizard ref="accordionWizardRef" v-slot="{ prev }">
    <vue-accordion-wizard-item title="Product Info">
      <vue-form
        ref="form"
        @submit="handleNext"
        :keepValues="true"
        :validationSchema="{
          name: 'required',
          brand_id: 'required',
          category_id: 'required',
          subcategory_id: 'required',
          description: 'required',
        }">
        <div class="mb-6">
          <vue-form-select-input
            :options="productBrandOptions"
            :disabled="isSubmitting || productBrandOptions.length == 0"
            @change="handleProductBrandInputChange"
            @search="handleProductBrandSearch"
            name="brand_id"
            label="Brand"
            placeholder="Select Brand">
            <template #addon-end>
              <button type="button" @click="emits('addProductBrand')"
                      class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="productCategoryOptions"
            :disabled="isSubmitting || !selectedProductBrandId || productCategoryOptions.length == 0"
            @change="handleProductCategoryInputChange"
            @search="handleProductCategorySearch"
            name="category_id"
            label="Category"
            placeholder="Select Category">
            <template #addon-end>
              <button
                :disabled="!selectedProductBrandId"
                @click="emits('addProductCategory', selectedProductBrandId)"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="productSubcategoryOptions"
            :disabled="isSubmitting || !selectedProductCategoryId || productSubcategoryOptions.length == 0"
            @search="handleProductSubcategorySearch"
            name="subcategory_id"
            type="text"
            label="Subcategory"
            placeholder="Select Subcategory">
            <template #addon-end>
              <button
                :disabled="!selectedProductCategoryId"
                @click="emits('addProductSubcategory', selectedProductCategoryId)"
                type="button"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-input
            :disabled="isSubmitting"
            name="name"
            type="text"
            label="Name"
            placeholder="Enter Product Name"/>
        </div>
        <div class="mb-6">
          <vue-form-html-input
            :disabled="isSubmitting"
            name="description"
            label="Description"
            placeholder="Product Description"/>
        </div>
        <button type="submit" class="btn btn-outline-primary">Next</button>
      </vue-form>
    </vue-accordion-wizard-item>
    <vue-accordion-wizard-item title="Product Details">
      <vue-form ref="form" @submit="handleNext">
        <div v-for="item in productAttributes" :key="getAttributeID(item)" class="mb-6">
          <template v-if="item.fieldType === 'textarea'">
            <vue-form-textarea-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :disabled="isSubmitting"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"/>
          </template>
          <template v-else-if="item.fieldType === 'number'">
            <vue-form-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :disabled="isSubmitting"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"
              type="number"/>
          </template>
          <template v-else-if="item.fieldType === 'date'">
            <vue-form-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :disabled="isSubmitting"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"
              type="date"/>
          </template>
          <template
            v-else-if="
              item.fieldType === 'select' ||
              item.fieldType === 'radio' ||
              item.fieldType === 'checkbox'
            ">
            <vue-form-select-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :disabled="isSubmitting || getAttributeSelectOptions(item).length == 0"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"
              :options="getAttributeSelectOptions(item)"/>
          </template>
          <template v-else-if="item.fieldType === 'multi_check'">
            <vue-form-select-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :isMulti="true"
              :disabled="isSubmitting || getAttributeSelectOptions(item).length == 0"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"
              :options="getAttributeSelectOptions(item)"/>
          </template>
          <template v-else>
            <vue-form-input
              :name="getAttributeInputName(item)"
              :label="item.label"
              :disabled="isSubmitting"
              :validation="getAttributeValidation(item)"
              :placeholder="item.placeholder"
              type="text"/>
          </template>
        </div>
        <div class="d-flex align-content-start">
          <button @click="prev" :disabled="isSubmitting" type="button"
                  class="btn btn-outline-primary">Prev
          </button>
          <button :disabled="isSubmitting" type="submit" class="btn btn-outline-primary ms-6">Next
          </button>
        </div>
      </vue-form>
    </vue-accordion-wizard-item>
    <vue-accordion-wizard-item title="Product Media">
      <vue-form ref="form" @submit="handleNext">
        <div class="mb-6">
          <vue-form-file-input
            :disabled="isSubmitting"
            name="variant.images"
            label="Images only, 1–10 files, up to 20 MB"
            :maxFileSize="2e7"
            :maxNumberOfFiles="10"
            :allowedFileTypes="['image/png', 'image/jpeg']"/>
        </div>
        <div class="d-flex align-content-start">
          <button @click="prev" type="button" class="btn btn-outline-primary">Prev</button>
          <button type="submit" class="btn btn-outline-primary ms-6">Next</button>
        </div>
      </vue-form>
    </vue-accordion-wizard-item>
    <vue-accordion-wizard-item title="Product Location">
      <vue-form ref="form" @submit="handleNext" :keepValues="true">
        <div class="mb-6">
          <vue-form-select-input
            :options="locationAisleOptions"
            :disabled="isSubmitting || locationAisleOptions.length == 0"
            @change="handleLocationAisleInputChange"
            @search="handleLocationAisleSearch"
            name="aisle_id"
            type="text"
            label="Aisle"
            placeholder="Select Aisle">
            <template #addon-end>
              <button type="button" @click="emits('addLocationAisle')"
                      class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="locationBayOptions"
            :disabled="isSubmitting || locationBayOptions.length == 0"
            @change="handleLocationBayInputChange"
            @search="handleLocationBaySearch"
            name="bay_id"
            type="text"
            label="Bay"
            placeholder="Select Bay">
            <template #addon-end>
              <button
                type="button"
                @click="emits('addLocationBay', selectedLocationAisleId)"
                :disabled="isSubmitting || !selectedLocationAisleId"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="locationShelfOptions"
            :disabled="isSubmitting || locationShelfOptions.length == 0"
            @change="handleLocationShelfInputChange"
            @search="handleLocationShelfSearch"
            name="shelf_id"
            type="text"
            label="Shelf"
            placeholder="Select Shelf">
            <template #addon-end>
              <button
                type="button"
                @click="emits('addLocationShelf', selectedLocationBayId)"
                :disabled="isSubmitting || !selectedLocationBayId"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="locationShelfLevelOptions"
            :disabled="isSubmitting || locationShelfLevelOptions.length == 0"
            @change="handleLocationShelfLevelInputChange"
            @search="handleLocationShelfLevelSearch"
            name="shelf_level_id"
            type="text"
            label="Shelf Level"
            placeholder="Select Shelf Level">
            <template #addon-end>
              <button
                type="button"
                @click="emits('addLocationShelfLevel', selectedLocationShelfId)"
                :disabled="isSubmitting || !selectedLocationShelfId"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="mb-6">
          <vue-form-select-input
            :options="locationBinOptions"
            :disabled="isSubmitting || locationBinOptions.length == 0"
            @search="handleLocationBinSearch"
            name="bin_id"
            type="text"
            label="Bin"
            placeholder="Product Bin">
            <template #addon-end>
              <button
                type="button"
                @click="emits('addLocationBin', selectedLocationShelfLevelId)"
                :disabled="isSubmitting || !selectedLocationShelfLevelId"
                class="btn btn-outline-primary ms-2">
                <i class="bx bx-plus icon-md"></i>
              </button>
            </template>
          </vue-form-select-input>
        </div>
        <div class="d-flex align-content-start">
          <button @click="prev" type="button" class="btn btn-outline-primary">Prev</button>
          <button type="submit" class="btn btn-outline-primary ms-6">Next</button>
        </div>
      </vue-form>
    </vue-accordion-wizard-item>
    <vue-accordion-wizard-item title="Review and Submit">
      <button @click="prev" type="button" class="btn btn-outline-primary">Prev</button>
      <button @click="handleSubmit" type="button" class="btn btn-primary ms-6">Submit</button>
    </vue-accordion-wizard-item>
  </vue-accordion-wizard>
</template>

<style scoped></style>
