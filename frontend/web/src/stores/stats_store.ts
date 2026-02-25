import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import type { StatsModel } from '@/app/models/stats_model'

export interface IStatsStore {
  get(): Ref<StatsModel | undefined>

  set(value: StatsModel): void
}

export const useStatsStore = defineStore('stats', (): IStatsStore => {
  const data = ref<StatsModel>()

  const get = (): Ref<StatsModel | undefined> => {
    return data as Ref<StatsModel | undefined>
  }

  const set = (value: StatsModel): void => {
    data.value = value
  }

  return { set, get }
})
