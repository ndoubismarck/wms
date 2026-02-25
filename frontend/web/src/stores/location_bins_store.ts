import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationBinModel } from '@/app/models/location_bin_model'

interface IGetParams {
  shelfLevelId?: string | null
}

export interface ILocationBinsStore {
  get(params?: IGetParams): Ref<LocationBinModel[]>

  set(value: LocationBinModel | LocationBinModel[]): void

  getById(id: string): LocationBinModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsBinsStore = defineStore('location_bins', (): ILocationBinsStore => {
  const data = ref<LocationBinModel[]>([])
  const filteredData = ref<LocationBinModel[]>([])

  const get = (params?: IGetParams): Ref<LocationBinModel[]> => {
    if (params && params.shelfLevelId) {
      filteredData.value = data.value.filter((item) => item.shelfLevelId == params.shelfLevelId)
      return filteredData as Ref<LocationBinModel[]>
    }
    return data as Ref<LocationBinModel[]>
  }

  const set = (value: LocationBinModel | LocationBinModel[]): void => {
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

  const getById = (id: string): LocationBinModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationBinModel
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
})
