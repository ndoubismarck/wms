import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { OrderModel, type OrdersSummaryModel } from '@/app/models/order_model'
import { emptyOrdersSummary } from '@/app/services/operations_service'

export interface IOrdersStore {
  get(): Ref<OrderModel[]>

  set(value: OrderModel | OrderModel[]): void

  setAll(value: OrderModel[]): void

  getSummary(): Ref<OrdersSummaryModel>

  setSummary(value: OrdersSummaryModel): void

  deleteAll(): void
}

export const useOrdersStore = defineStore('orders', (): IOrdersStore => {
  const data = ref<OrderModel[]>([])
  const summary = ref<OrdersSummaryModel>(emptyOrdersSummary())

  const get = (): Ref<OrderModel[]> => {
    return data as Ref<OrderModel[]>
  }

  const set = (value: OrderModel | OrderModel[]): void => {
    if (Array.isArray(value)) {
      value.forEach((item) => {
        const index = data.value.findIndex((current) => current.id == item.id)
        if (index < 0) {
          data.value.push(item)
        } else {
          data.value[index] = item
        }
      })
    } else {
      const index = data.value.findIndex((current) => current.id == value.id)
      if (index < 0) {
        data.value.push(value)
      } else {
        data.value[index] = value
      }
    }
    data.value = data.value.sort((a, b) => b.createdAt.toUnix() - a.createdAt.toUnix())
  }

  const setAll = (value: OrderModel[]): void => {
    data.value = [...value].sort((a, b) => b.createdAt.toUnix() - a.createdAt.toUnix())
  }

  const getSummary = (): Ref<OrdersSummaryModel> => {
    return summary as Ref<OrdersSummaryModel>
  }

  const setSummary = (value: OrdersSummaryModel): void => {
    summary.value = value
  }

  const deleteAll = (): void => {
    data.value = []
    summary.value = emptyOrdersSummary()
  }

  return { set, setAll, get, getSummary, setSummary, deleteAll }
})
