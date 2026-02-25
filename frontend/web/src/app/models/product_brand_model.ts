import { DatabaseModel } from '@/app/core/model'

export class ProductBrandModel extends DatabaseModel {
  public readonly locationId: string
  public readonly name: string

  constructor(data: any) {
    super(data)
    this.locationId = this.getString('location_id')
    this.name = this.getString('name')
  }
}
