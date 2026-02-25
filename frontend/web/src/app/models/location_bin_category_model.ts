import { DatabaseModel } from '@/app/core/model'

export class LocationBinCategoryModel extends DatabaseModel {
  public readonly name: string

  constructor(data: any) {
    super(data)
    this.name = this.getString('name')
  }
}
