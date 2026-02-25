import { DatabaseModel } from '@/app/core/model.ts'
import { LocationAisleModel } from '@/app/models/location_aisle_model.ts'

export class LocationBayModel extends DatabaseModel {
  public readonly name: string
  public readonly aisle: LocationAisleModel
  public readonly aisleId: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
    this.aisle = new LocationAisleModel(this.getObject('aisle'))
    this.aisleId = this.getString('aisle_id')
  }
}
