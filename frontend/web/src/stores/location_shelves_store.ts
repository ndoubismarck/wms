import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationShelfModel } from '@/app/models/location_shelf_model'

interface IGetParams {
  bayId?: string | null
}

export interface ILocationShelvesStore {
  get(params?: IGetParams): Ref<LocationShelfModel[]>

  set(value: LocationShelfModel | LocationShelfModel[]): void

  getById(id: string): LocationShelfModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsShelvesStore = defineStore('location_shelves', (): ILocationShelvesStore => {
  const data = ref<LocationShelfModel[]>([])
  const filteredData = ref<LocationShelfModel[]>([])

  const get = (params?: IGetParams): Ref<LocationShelfModel[]> => {
    if (params && params.bayId) {
      filteredData.value = data.value.filter((item) => item.bayId == params.bayId)
      return filteredData as Ref<LocationShelfModel[]>
    }
    return data as Ref<LocationShelfModel[]>
  }

  const set = (value: LocationShelfModel | LocationShelfModel[]): void => {
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

  const getById = (id: string): LocationShelfModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationShelfModel
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
