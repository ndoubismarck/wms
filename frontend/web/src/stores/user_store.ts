import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import type { UserModel } from '@/app/models/user_model'

export interface IUserStore {
  get(): Ref<UserModel | undefined>

  set(value: UserModel): void
}

export const useUserStore = defineStore('user', (): IUserStore => {
  const data = ref<UserModel>()

  const get = (): Ref<UserModel | undefined> => {
    return data as Ref<UserModel | undefined>
  }

  const set = (value: UserModel): void => {
    data.value = value
  }

  return { set, get }
})
