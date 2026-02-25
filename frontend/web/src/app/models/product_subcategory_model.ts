import { DatabaseModel } from '@/app/core/model'

export class ProductSubcategoryModel extends DatabaseModel {
  public readonly categoryId: string
  public readonly code: string
  public readonly name: string
  public readonly description: string

  constructor(data: any) {
    super(data)
    this.categoryId = this.getString('category_id')
    this.code = this.getString('code')
    this.name = this.getString('name')
    this.description = this.getString('description')
  }
}
