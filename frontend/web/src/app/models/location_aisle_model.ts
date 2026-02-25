import { DatabaseModel } from '@/app/core/model.ts'

export class LocationAisleModel extends DatabaseModel {
  public readonly name: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
  }
}
