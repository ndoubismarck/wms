import { DatabaseModel } from '@/app/core/model'

export class TeamModel extends DatabaseModel {
  public readonly locationId: string
  public readonly name: string
  public readonly description: string
  public readonly isSystem: boolean
  public readonly userIds: string[]

  constructor(data: any) {
    super(data)
    this.locationId = this.getString('location_id')
    this.name = this.getString('name')
    this.description = this.getString('description')
    this.isSystem = this.getBoolean('is_system')
    this.userIds = this.getArray('user_ids')
  }
}
