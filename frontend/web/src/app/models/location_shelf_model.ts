import { DatabaseModel } from '@/app/core/model.ts'
import { LocationBayModel } from '@/app/models/location_bay_model.ts'

export class LocationShelfModel extends DatabaseModel {
  public readonly name: string
  public readonly bay: LocationBayModel
  public readonly bayId: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
    this.bay = new LocationBayModel(this.getObject('bay'))
    this.bayId = this.getString('bay_id')
  }
}
