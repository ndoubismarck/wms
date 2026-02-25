import { DatabaseModel } from '@/app/core/model'

export class LocationModel extends DatabaseModel {
  public readonly code: string
  public readonly name: string
  public readonly description: string

  constructor(data: any) {
    super(data)
    this.code = this.getString('code')
    this.name = this.getString('name')
    this.description = this.getString('description')
  }
}
