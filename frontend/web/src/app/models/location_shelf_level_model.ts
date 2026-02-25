import { DatabaseModel } from '@/app/core/model.ts'
import { LocationShelfModel } from '@/app/models/location_shelf_model.ts'

export class LocationShelfLevelModel extends DatabaseModel {
  public readonly name: string
  public readonly shelf: LocationShelfModel
  public readonly shelfId: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
    this.shelf = new LocationShelfModel(this.getObject('shelf'))
    this.shelfId = this.getString('shelf_id')
  }
}
