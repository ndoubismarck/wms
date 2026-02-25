import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import type { SetupModel } from '@/app/models/setup_model'

export interface ISetupStore {
  get(): Ref<SetupModel | undefined>

  set(value: SetupModel): void
}

export const useSetupStore = defineStore('setup', (): ISetupStore => {
  const data = ref<SetupModel>()

  const get = (): Ref<SetupModel | undefined> => {
    return data as Ref<SetupModel | undefined>
  }

  const set = (value: SetupModel): void => {
    data.value = value
  }

  return { set, get }
})
