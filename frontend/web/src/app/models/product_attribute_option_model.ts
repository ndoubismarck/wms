import { DatabaseModel } from '@/app/core/model'

export class ProductAttributeOptionModel extends DatabaseModel {
  public readonly optionId: string
  public readonly label: string
  public readonly productAttributeId: string
  public readonly attributeId: string

  constructor(data: any) {
    const payload =
      typeof data === 'string'
        ? {
            id: data,
            option_id: data,
            label: data,
          }
        : data
    super(payload)
    this.optionId = this.getString('option_id') || this.getString('id')
    this.label = this.getString('label') || this.getString('name') || this.optionId
    this.productAttributeId =
      this.getString('product_attribute_id') || this.getString('attribute_id') || this.getString('attribute_value_id')
    this.attributeId = this.productAttributeId
  }
}
