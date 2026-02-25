import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ProductModel } from '@/app/models/product_model'

export interface IProductsStore {
  get(): Ref<ProductModel[]>

  set(value: ProductModel | ProductModel[]): void

  getById(id: string): ProductModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useProductsStore = defineStore('products', (): IProductsStore => {
  const data = ref<ProductModel[]>([])

  const get = (): Ref<ProductModel[]> => {
    return data as Ref<ProductModel[]>
  }

  const set = (value: ProductModel | ProductModel[]): void => {
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

  const getById = (id: string): ProductModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as ProductModel
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
