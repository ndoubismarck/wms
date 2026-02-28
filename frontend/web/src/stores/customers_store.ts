import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { CustomerModel, type CustomersSummaryModel } from '@/app/models/customer_model'
import { emptyCustomersSummary } from '@/app/services/operations_service'

export interface ICustomersStore {
  get(): Ref<CustomerModel[]>

  setAll(value: CustomerModel[]): void

  getSummary(): Ref<CustomersSummaryModel>

  setSummary(value: CustomersSummaryModel): void

  deleteAll(): void
}

export const useCustomersStore = defineStore('customers', (): ICustomersStore => {
  const data = ref<CustomerModel[]>([])
  const summary = ref<CustomersSummaryModel>(emptyCustomersSummary())

  const get = (): Ref<CustomerModel[]> => {
    return data as Ref<CustomerModel[]>
  }

  const setAll = (value: CustomerModel[]): void => {
    data.value = [...value].sort((a, b) => b.totalOrders - a.totalOrders)
  }

  const getSummary = (): Ref<CustomersSummaryModel> => {
    return summary as Ref<CustomersSummaryModel>
  }

  const setSummary = (value: CustomersSummaryModel): void => {
    summary.value = value
  }

  const deleteAll = (): void => {
    data.value = []
    summary.value = emptyCustomersSummary()
  }

  return { setAll, get, getSummary, setSummary, deleteAll }
})
