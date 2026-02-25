import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { LocationBayModel } from '@/app/models/location_bay_model'

interface IGetParams {
  aisleId?: string | null
}

export interface ILocationBaysStore {
  get(params?: IGetParams): Ref<LocationBayModel[]>

  set(value: LocationBayModel | LocationBayModel[]): void

  getById(id: string): LocationBayModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useLocationsBaysStore = defineStore('location_bays', (): ILocationBaysStore => {
  const data = ref<LocationBayModel[]>([])
  const filteredData = ref<LocationBayModel[]>([])

  const get = (params?: IGetParams): Ref<LocationBayModel[]> => {
    if (params && params.aisleId) {
      filteredData.value = data.value.filter((item) => item.aisleId == params.aisleId)
      return filteredData as Ref<LocationBayModel[]>
    }
    return data as Ref<LocationBayModel[]>
  }

  const set = (value: LocationBayModel | LocationBayModel[]): void => {
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

  const getById = (id: string): LocationBayModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as LocationBayModel
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
