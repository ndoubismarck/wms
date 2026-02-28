import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ShipmentModel, type ShipmentsSummaryModel } from '@/app/models/shipment_model'
import { emptyShipmentsSummary } from '@/app/services/operations_service'

export interface IShipmentsStore {
  get(): Ref<ShipmentModel[]>

  setAll(value: ShipmentModel[]): void

  getSummary(): Ref<ShipmentsSummaryModel>

  setSummary(value: ShipmentsSummaryModel): void

  deleteAll(): void
}

export const useShipmentsStore = defineStore('shipments', (): IShipmentsStore => {
  const data = ref<ShipmentModel[]>([])
  const summary = ref<ShipmentsSummaryModel>(emptyShipmentsSummary())

  const get = (): Ref<ShipmentModel[]> => {
    return data as Ref<ShipmentModel[]>
  }

  const setAll = (value: ShipmentModel[]): void => {
    data.value = [...value].sort((a, b) => b.updatedAt.toUnix() - a.updatedAt.toUnix())
  }

  const getSummary = (): Ref<ShipmentsSummaryModel> => {
    return summary as Ref<ShipmentsSummaryModel>
  }

  const setSummary = (value: ShipmentsSummaryModel): void => {
    summary.value = value
  }

  const deleteAll = (): void => {
    data.value = []
    summary.value = emptyShipmentsSummary()
  }

  return { setAll, get, getSummary, setSummary, deleteAll }
})
