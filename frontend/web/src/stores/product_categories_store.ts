import { type Ref, ref } from 'vue'
import { defineStore } from 'pinia'
import { ProductCategoryModel } from '@/app/models/product_category_model'

interface IGetParams {
  brandId?: string | null
}

export interface IProductCategoriesStore {
  get(params?: IGetParams): Ref<ProductCategoryModel[]>

  set(value: ProductCategoryModel | ProductCategoryModel[]): void

  getById(id: string): ProductCategoryModel | null

  deleteById(id: string): void

  deleteAll(): void
}

export const useProductCategoriesStore = defineStore('product_categories', (): IProductCategoriesStore => {
  const data = ref<ProductCategoryModel[]>([])
  const filteredData = ref<ProductCategoryModel[]>([])

  const get = (params?: IGetParams): Ref<ProductCategoryModel[]> => {
    if (params && params.brandId) {
      filteredData.value = data.value.filter((item) => item.brandId == params.brandId)
      return filteredData as Ref<ProductCategoryModel[]>
    }
    return data as Ref<ProductCategoryModel[]>
  }

  const set = (value: ProductCategoryModel | ProductCategoryModel[]): void => {
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

  const getById = (id: string): ProductCategoryModel | null => {
    const index = data.value.findIndex((value) => value.id == id)
    if (index >= 0) {
      return data.value[index] as ProductCategoryModel
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
