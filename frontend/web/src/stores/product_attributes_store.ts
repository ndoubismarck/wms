import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ProductAttributeModel } from '@/app/models/product_attribute_model'

export interface IProductAttributesStore {
  get(): Ref<ProductAttributeModel[]>

  set(value: ProductAttributeModel | ProductAttributeModel[]): void

  getById(id: string): ProductAttributeModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useProductAttributesStore = defineStore('product_attributes', (): IProductAttributesStore => {
  const data = ref<ProductAttributeModel[]>([])

  const get = (): Ref<ProductAttributeModel[]> => {
    return data as Ref<ProductAttributeModel[]>
  }

  const set = (value: ProductAttributeModel | ProductAttributeModel[]): void => {
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

  const getById = (id: string): ProductAttributeModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as ProductAttributeModel
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
