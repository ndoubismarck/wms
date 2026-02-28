import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { SupplierModel, type SuppliersSummaryModel } from '@/app/models/supplier_model'
import { emptySuppliersSummary } from '@/app/services/operations_service'

export interface ISuppliersStore {
  get(): Ref<SupplierModel[]>

  setAll(value: SupplierModel[]): void

  getSummary(): Ref<SuppliersSummaryModel>

  setSummary(value: SuppliersSummaryModel): void

  deleteAll(): void
}

export const useSuppliersStore = defineStore('suppliers', (): ISuppliersStore => {
  const data = ref<SupplierModel[]>([])
  const summary = ref<SuppliersSummaryModel>(emptySuppliersSummary())

  const get = (): Ref<SupplierModel[]> => {
    return data as Ref<SupplierModel[]>
  }

  const setAll = (value: SupplierModel[]): void => {
    data.value = [...value].sort((a, b) => b.openPurchaseOrders - a.openPurchaseOrders)
  }

  const getSummary = (): Ref<SuppliersSummaryModel> => {
    return summary as Ref<SuppliersSummaryModel>
  }

  const setSummary = (value: SuppliersSummaryModel): void => {
    summary.value = value
  }

  const deleteAll = (): void => {
    data.value = []
    summary.value = emptySuppliersSummary()
  }

  return { setAll, get, getSummary, setSummary, deleteAll }
})
