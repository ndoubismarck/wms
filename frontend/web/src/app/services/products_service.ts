import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import { PaginationModel } from '@/app/models/pagination_model.ts'
import { ProductModel } from '@/app/models/product_model.ts'
import { ProductBrandModel } from '@/app/models/product_brand_model'
import { ProductCategoryModel } from '@/app/models/product_category_model'
import { ProductSubcategoryModel } from '@/app/models/product_subcategory_model'
import { ProductAttributeModel } from '@/app/models/product_attribute_model'
import type { IApiResult } from '@/app/core/client'

export class ProductsService extends AppService {
  constructor(ctx: AppContext) {
    super(ctx)
  }

  public async add(data: any): Promise<IAddResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products`)
    const response = await this.api.private.post(path, data, true)
    return {
      product: new ProductModel(response.body.data['product']),
      success: response.body.code == 'success',
    }
  }

  public async update(id: string, data: any): Promise<IUpdateResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/${id}`)
    const response = await this.api.private.put(path, data, true)
    return {
      product: new ProductModel(response.body.data['product']),
      success: response.body.code == 'success',
    }
  }

  public async deleteById(id: string, params?: any): Promise<IDeleteByIdResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getOne(id: string, params?: any): Promise<IGetOneResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/${id}`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      product: new ProductModel(response.body.data['product']),
    }
  }

  public async getMany(params: any): Promise<IGetManyResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      products: this.mapArray(response.body.data['products'], (value) => new ProductModel(value)),
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addCategory(data: any): Promise<IAddCategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/categories`)
    const response = await this.api.private.post(path, data, true)
    return {
      success: response.body.code == 'success',
      category: new ProductCategoryModel(response.body.data['category']),
    }
  }

  public async updateCategory(id: string, data: any): Promise<IUpdateCategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/categories/${id}`)
    const response = await this.api.private.put(path, data, true)
    return {
      success: response.body.code == 'success',
      category: new ProductCategoryModel(response.body.data['category']),
    }
  }

  public async deleteCategory(id: string, params?: any): Promise<IDeleteCategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/categories/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getCategory(id: string, params?: any): Promise<IGetCategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/categories/${id}`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      category: new ProductCategoryModel(response.body.data['category']),
    }
  }

  public async getCategories(params: any): Promise<IGetCategoriesResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/categories`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      categories: this.mapArray(response.body.data['categories'], (value) => new ProductCategoryModel(value)),
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addSubcategory(data: any): Promise<IAddSubcategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/subcategories`)
    const response = await this.api.private.post(path, data, true)
    return {
      success: response.body.code == 'success',
      subcategory: new ProductSubcategoryModel(response.body.data['subcategory']),
    }
  }

  public async updateSubcategory(id: string, data: any): Promise<IUpdateSubcategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/subcategories/${id}`)
    const response = await this.api.private.put(path, data, true)
    return {
      success: response.body.code == 'success',
      subcategory: new ProductSubcategoryModel(response.body.data['subcategory']),
    }
  }

  public async deleteSubcategory(id: string): Promise<IDeleteSubcategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/subcategories/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getSubcategory(id: string, params?: any): Promise<IGetSubcategoryResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/subcategories/${id}`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      subcategory: new ProductSubcategoryModel(response.body.data['subcategory']),
    }
  }

  public async getSubcategories(params: any): Promise<IGetSubcategoriesResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/subcategories`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      subcategories: this.mapArray(response.body.data['subcategories'], (value) => new ProductSubcategoryModel(value)),
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addBrand(data: any): Promise<IAddBrandResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/brands`)
    const response = await this.api.private.post(path, data, true)
    return {
      success: response.body.code == 'success',
      brand: new ProductBrandModel(response.body.data['brand']),
    }
  }

  public async updateBrand(id: string, data: any): Promise<IUpdateBrandResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/brands/${id}`)
    const response = await this.api.private.put(path, data, true)
    return {
      success: response.body.code == 'success',
      brand: new ProductBrandModel(response.body.data['brand']),
    }
  }

  public async deleteBrand(id: string, params?: any): Promise<IDeleteBrandResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/brands/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getBrand(id: string, params?: any): Promise<IGetBrandResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/brands/${id}`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      brand: new ProductBrandModel(response.body.data['brand']),
    }
  }

  public async getBrands(params: any): Promise<IGetBrandsResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/brands`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      brands: this.mapArray(response.body.data['brands'], (value) => new ProductBrandModel(value)),
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addAttribute(data: any): Promise<IAddAttributeResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/attributes`)
    const response = await this.api.private.post(path, data, true)
    return {
      success: response.body.code == 'success',
      attribute: new ProductAttributeModel(response.body.data['attribute']),
    }
  }

  public async updateAttribute(id: string, data: any): Promise<IUpdateAttributeResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/attributes/${id}`)
    const response = await this.api.private.put(path, data, true)
    return {
      success: response.body.code == 'success',
      attribute: new ProductAttributeModel(response.body.data['attribute']),
    }
  }

  public async deleteAttribute(id: string, params?: any): Promise<IDeleteAttributeResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/attributes/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getAttribute(id: string, params?: any): Promise<IGetAttributeResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/attributes/${id}`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      attribute: new ProductAttributeModel(response.body.data['attribute']),
    }
  }

  public async getAttributes(params: any): Promise<IGetAttributesResult> {
    const lid = await this.ctx.storage.get('location_id')
    const path = this.api.private.path(`/${lid}/products/attributes`, params)
    const response = await this.api.private.get(path, true)
    return {
      success: response.body.code == 'success',
      attributes: this.mapArray(response.body.data['attributes'], (value) => new ProductAttributeModel(value)),
      pagination: new PaginationModel(response.pagination),
    }
  }
}

interface IAddResult extends IApiResult {
  product: ProductModel
}

interface IUpdateResult extends IApiResult {
  product: ProductModel
}

interface IDeleteByIdResult extends IApiResult {}

interface IGetOneResult extends IApiResult {
  product: ProductModel
}

interface IGetManyResult extends IApiResult {
  products: ProductModel[]
  pagination: PaginationModel
}

interface IAddCategoryResult extends IApiResult {
  category: ProductCategoryModel
}

interface IUpdateCategoryResult extends IApiResult {
  category: ProductCategoryModel
}

interface IDeleteCategoryResult extends IApiResult {}

interface IGetCategoryResult extends IApiResult {
  category: ProductCategoryModel
}

interface IGetCategoriesResult extends IApiResult {
  categories: ProductCategoryModel[]
  pagination: PaginationModel
}

interface IAddSubcategoryResult extends IApiResult {
  subcategory: ProductSubcategoryModel
}

interface IUpdateSubcategoryResult extends IApiResult {
  subcategory: ProductSubcategoryModel
}

interface IDeleteSubcategoryResult extends IApiResult {}

interface IGetSubcategoryResult extends IApiResult {
  subcategory: ProductSubcategoryModel
}

interface IGetSubcategoriesResult extends IApiResult {
  subcategories: ProductSubcategoryModel[]
  pagination: PaginationModel
}

interface IAddBrandResult extends IApiResult {
  brand: ProductBrandModel
}

interface IUpdateBrandResult extends IApiResult {
  brand: ProductBrandModel
}

interface IDeleteBrandResult extends IApiResult {}

interface IGetBrandResult extends IApiResult {
  brand: ProductBrandModel
}

interface IGetBrandsResult extends IApiResult {
  brands: ProductBrandModel[]
  pagination: PaginationModel
}

interface IAddAttributeResult extends IApiResult {
  attribute: ProductAttributeModel
}

interface IUpdateAttributeResult extends IApiResult {
  attribute: ProductAttributeModel
}

interface IDeleteAttributeResult extends IApiResult {}

interface IGetAttributeResult extends IApiResult {
  attribute: ProductAttributeModel
}

interface IGetAttributesResult extends IApiResult {
  attributes: ProductAttributeModel[]
  pagination: PaginationModel
}
