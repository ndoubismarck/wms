import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { UserModel } from '@/app/models/user_model'

export interface IUsersStore {
  get(): Ref<UserModel[]>

  set(value: UserModel | UserModel[]): void

  getById(id: string): UserModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useUsersStore = defineStore('users', (): IUsersStore => {
  const data = ref<UserModel[]>([])

  const get = (): Ref<UserModel[]> => {
    return data as Ref<UserModel[]>
  }

  const set = (value: UserModel | UserModel[]): void => {
    if (Array.isArray(value)) {
      value.forEach((item) => {
        const index = data.value.findIndex((val) => val.id == item.id)
        if (index < 0) {
          data.value.push(item)
        } else {
          data.value[index] = item
        }
      })
    } else {
      const item = value
      const index = data.value.findIndex((val) => val.id == item.id)
      if (index < 0) {
        data.value.push(value)
      } else {
        data.value[index] = value
      }
    }
    data.value = data.value.sort((a, b) => b.createdAt.toUnix() - a.createdAt.toUnix())
  }

  const getById = (id: string): UserModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as UserModel
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
