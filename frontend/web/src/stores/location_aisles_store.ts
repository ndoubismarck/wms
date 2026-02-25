import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationAisleModel } from '@/app/models/location_aisle_model'

export interface ILocationAislesStore {
  get(): Ref<LocationAisleModel[]>

  set(value: LocationAisleModel | LocationAisleModel[]): void

  getById(id: string): LocationAisleModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsAislesStore = defineStore('location_aisles', (): ILocationAislesStore => {
  const data = ref<LocationAisleModel[]>([])

  const get = (): Ref<LocationAisleModel[]> => {
    return data as Ref<LocationAisleModel[]>
  }

  const set = (value: LocationAisleModel | LocationAisleModel[]): void => {
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

  const getById = (id: string): LocationAisleModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationAisleModel
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
