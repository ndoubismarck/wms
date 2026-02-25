import {DatabaseModel} from '@/app/core/model'
import {ProductAttributeOptionModel} from '@/app/models/product_attribute_option_model'

export class ProductAttributeModel extends DatabaseModel {
  public readonly attributeId: string
  public readonly attributeOptionId: string
  public readonly label: string
  public readonly placeholder: string
  public readonly appliesTo: string
  public readonly dataType: string
  public readonly fieldType: string
  public readonly validationSchema: string[]
  public readonly options: ProductAttributeOptionModel[]

  constructor(data: any) {
    super(data)
    this.attributeId =
      this.getString('product_attribute_id') || this.getString('attribute_id') || this.getString('id')
    this.attributeOptionId =
      this.getString('product_attribute_option_id') || this.getString('attribute_option_id')
    this.label =
      this.getString('label') ||
      this.getString('product_attribute_definition.label') ||
      this.getString('definition.label')
    this.placeholder =
      this.getString('placeholder') ||
      this.getString('product_attribute_definition.placeholder') ||
      this.getString('definition.placeholder')
    this.appliesTo =
      this.getString('applies_to') ||
      this.getString('product_attribute_definition.applies_to') ||
      this.getString('definition.applies_to')
    this.dataType =
      this.getString('data_type') ||
      this.getString('product_attribute_definition.data_type') ||
      this.getString('definition.data_type')
    this.fieldType =
      this.getString('field_type') ||
      this.getString('product_attribute_definition.field_type') ||
      this.getString('definition.field_type')
    const validationRules = this.getArray('validation_rules')
    const validationRulesFromAttributeDefinition = this.getArray('product_attribute_definition.validation_rules')
    const validationRulesFromDefinition = this.getArray('definition.validation_rules')
    const validationSchema = this.getArray('validation_schema')
    const validationSchemaFromDefinition = this.getArray('definition.validation_schema')
    this.validationSchema =
      validationRules.length > 0
        ? validationRules
        : validationRulesFromAttributeDefinition.length > 0
          ? validationRulesFromAttributeDefinition
          : validationRulesFromDefinition.length > 0
            ? validationRulesFromDefinition
            : validationSchema.length > 0
              ? validationSchema
              : validationSchemaFromDefinition
    const dataOptions = this.getArray('options')
    const resolvedOptions =
      dataOptions.length > 0
        ? dataOptions
        : this.getArray('product_attribute_definition.options').length > 0
          ? this.getArray('product_attribute_definition.options')
          : this.getArray('definition.options')
    this.options = resolvedOptions.map((value) =>
      new ProductAttributeOptionModel(
        typeof value === 'string'
          ? {
              id: value,
              option_id: value,
              label: value,
            }
          : value,
      ),
    )
  }
}
