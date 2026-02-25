import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ProductSubcategoryModel } from '@/app/models/product_subcategory_model'

interface IGetParams {
  categoryId?: string | null
}

export interface IProductSubcategoriesStore {
  get(params?: IGetParams): Ref<ProductSubcategoryModel[]>

  set(value: ProductSubcategoryModel | ProductSubcategoryModel[]): void

  getById(id: string): ProductSubcategoryModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useProductSubcategoriesStore = defineStore('product_subcategories', (): IProductSubcategoriesStore => {
  const data = ref<ProductSubcategoryModel[]>([])
  const filteredData = ref<ProductSubcategoryModel[]>([])

  const get = (params?: IGetParams): Ref<ProductSubcategoryModel[]> => {
    if (params && params.categoryId) {
      filteredData.value = data.value.filter((item) => item.categoryId == params.categoryId)
      return filteredData as Ref<ProductSubcategoryModel[]>
    }
    return data as Ref<ProductSubcategoryModel[]>
  }

  const set = (value: ProductSubcategoryModel | ProductSubcategoryModel[]): void => {
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

  const getById = (id: string): ProductSubcategoryModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as ProductSubcategoryModel
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
