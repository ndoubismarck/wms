import { Model } from '@/app/core/model'
import { LocationModel } from '@/app/models/location_model.ts'

export class SetupModel extends Model {
  public readonly admin: boolean
  public readonly database: boolean
  public readonly location: LocationModel

  constructor(data: any) {
    super(data)
    this.admin = this.getBoolean('admin')
    this.database = this.getBoolean('database')
    this.location = new LocationModel(this.getObject('location'))
  }
}
