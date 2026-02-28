import { DatabaseModel } from '@/app/core/model'
import { DateTime } from '@/app/core/date_time'

export class TaskModel extends DatabaseModel {
  public readonly locationId: string
  public readonly title: string
  public readonly description: string
  public readonly priority: string
  public readonly status: string
  public readonly scheduledByType: string
  public readonly scheduledByUserId: string
  public readonly hasDueAt: boolean
  public readonly dueAt: DateTime
  public readonly userIds: string[]
  public readonly teamIds: string[]

  constructor(data: any) {
    super(data)
    this.locationId = this.getString('location_id')
    this.title = this.getString('title')
    this.description = this.getString('description')
    this.priority = this.getString('priority')
    this.status = this.getString('status')
    this.scheduledByType = this.getString('scheduled_by_type')
    this.scheduledByUserId = this.getString('scheduled_by_user_id')
    this.hasDueAt = this.exists('due_at')
    this.dueAt = this.getDateTime('due_at')
    this.userIds = this.getArray('user_ids')
    this.teamIds = this.getArray('team_ids')
  }
}
