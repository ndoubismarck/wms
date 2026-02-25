import { DatabaseModel } from '@/app/core/model'
import { LocationBinModel } from '@/app/models/location_bin_model.ts'
import { ProductBrandModel } from '@/app/models/product_brand_model'
import { ProductCategoryModel } from '@/app/models/product_category_model'
import { ProductSubcategoryModel } from '@/app/models/product_subcategory_model'
import { ProductAttributeModel } from '@/app/models/product_attribute_model'

export class ProductModel extends DatabaseModel {
  public readonly sku: string
  public readonly bin: LocationBinModel
  public readonly brandId: string
  public readonly categoryId: string
  public readonly subcategoryId: string
  public readonly name: string
  public readonly description: string
  public readonly status: string
  public readonly brand: ProductBrandModel
  public readonly category: ProductCategoryModel
  public readonly subcategory: ProductSubcategoryModel
  public readonly attributes: ProductAttributeModel[]

  constructor(data: any) {
    super(data)
    this.sku = this.getString('sku') || this.getString('variants.0.sku')
    this.bin = new LocationBinModel(this.getObject('bin'))
    this.brandId = this.getString('brand_id')
    this.categoryId = this.getString('category_id')
    this.subcategoryId = this.getString('subcategory_id')
    this.name = this.getString('name')
    this.description = this.getString('description')
    this.status = this.getString('status')
    this.brand = new ProductBrandModel(this.getObject('brand'))
    this.category = new ProductCategoryModel(this.getObject('category'))
    this.subcategory = new ProductSubcategoryModel(this.getObject('subcategory'))
    this.attributes = this.getArray('attributes').map((value) => new ProductAttributeModel(value))
  }
}
