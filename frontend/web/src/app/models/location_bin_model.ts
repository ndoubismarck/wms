import { DatabaseModel } from '@/app/core/model.ts'
import { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model.ts'

export class LocationBinModel extends DatabaseModel {
  public readonly name: string
  public readonly shelfLevel: LocationShelfLevelModel
  public readonly shelfLevelId: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
    this.shelfLevel = new LocationShelfLevelModel(this.getObject('shelf_level'))
    this.shelfLevelId = this.getString('shelf_level_id')
  }
}
