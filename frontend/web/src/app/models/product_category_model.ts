import { DatabaseModel } from '@/app/core/model'

export class ProductCategoryModel extends DatabaseModel {
  public readonly name: string
  public readonly brandId: string
  public readonly description: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
    this.brandId = this.getString('brand_id')
    this.description = this.getString('description')
  }
}
