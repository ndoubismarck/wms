import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model'

interface IGetParams {
  shelfId?: string | null
}

export interface ILocationShelfLevelsStore {
  get(params?: IGetParams): Ref<LocationShelfLevelModel[]>

  set(value: LocationShelfLevelModel | LocationShelfLevelModel[]): void

  getById(id: string): LocationShelfLevelModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsShelfLevelsStore = defineStore('location_shelf_levels', (): ILocationShelfLevelsStore => {
  const data = ref<LocationShelfLevelModel[]>([])
  const filteredData = ref<LocationShelfLevelModel[]>([])

  const get = (params?: IGetParams): Ref<LocationShelfLevelModel[]> => {
    if (params && params.shelfId) {
      filteredData.value = data.value.filter((item) => item.shelfId == params.shelfId)
      return filteredData as Ref<LocationShelfLevelModel[]>
    }
    return data as Ref<LocationShelfLevelModel[]>
  }

  const set = (value: LocationShelfLevelModel | LocationShelfLevelModel[]): void => {
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

  const getById = (id: string): LocationShelfLevelModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationShelfLevelModel
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
