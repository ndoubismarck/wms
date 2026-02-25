import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationBinCategoryModel } from '@/app/models/location_bin_category_model'

interface IGetParams {}

export interface ILocationBinCategoriesStore {
  get(params?: IGetParams): Ref<LocationBinCategoryModel[]>

  set(value: LocationBinCategoryModel | LocationBinCategoryModel[]): void

  getById(id: string): LocationBinCategoryModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsBinCategoriesStore = defineStore(
  'location_bin_categories',
  (): ILocationBinCategoriesStore => {
    const data = ref<LocationBinCategoryModel[]>([])
    const filteredData = ref<LocationBinCategoryModel[]>([])

    const get = (params?: IGetParams): Ref<LocationBinCategoryModel[]> => {
      return data as Ref<LocationBinCategoryModel[]>
    }

    const set = (value: LocationBinCategoryModel | LocationBinCategoryModel[]): void => {
      if (Array.isArray(value)) {
        value.forEach((item) => {
          const index = data.value.findIndex((value) => value.id == item.id)
          if (index < 0) {
            data.value.push(item)
          } else {
            data.value[index] = item
          }
        })
      } else {
        const item = value
        const index = data.value.findIndex((value) => value.id == item.id)
        if (index < 0) {
          data.value.push(value)
        } else {
          data.value[index] = value
        }
      }
      data.value = data.value.sort((a, b) => b.createdAt.toUnix() - a.createdAt.toUnix())
    }

    const getById = (id: string): LocationBinCategoryModel | null => {
      const index = data.value.findIndex((value) => value.id == id)
      if (index >= 0) {
        return data.value[index] as LocationBinCategoryModel
      }
      return null
    }

    const deleteById = (id: string): void => {
      const index = data.value.findIndex((value) => value.id == id)
      if (index >= 0) {
        data.value.splice(index, 1)
      }
    }

    const deleteAll = (): void => {
      data.value = []
    }

    return { set, get, getById, deleteById, deleteAll }
  },
)
