import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ProductBrandModel } from '@/app/models/product_brand_model'

export interface IProductBrandsStore {
  get(): Ref<ProductBrandModel[]>

  set(value: ProductBrandModel | ProductBrandModel[]): void

  getById(id: string): ProductBrandModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useProductBrandsStore = defineStore('product_brands', (): IProductBrandsStore => {
  const data = ref<ProductBrandModel[]>([])

  const get = (): Ref<ProductBrandModel[]> => {
    return data as Ref<ProductBrandModel[]>
  }

  const set = (value: ProductBrandModel | ProductBrandModel[]): void => {
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

  const getById = (id: string): ProductBrandModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as ProductBrandModel
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
