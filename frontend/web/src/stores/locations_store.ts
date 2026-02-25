import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationModel } from '@/app/models/location_model'

export interface ILocationsStore {
  get(): Ref<LocationModel[]>

  set(value: LocationModel | LocationModel[]): void

  getById(id: string): LocationModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsStore = defineStore('locations', (): ILocationsStore => {
  const data = ref<LocationModel[]>([])

  const get = (): Ref<LocationModel[]> => {
    return data as Ref<LocationModel[]>
  }

  const set = (value: LocationModel | LocationModel[]): void => {
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

  const getById = (id: string): LocationModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationModel
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
