import { DatabaseModel } from '@/app/core/model'

export type InventoryMovementType =
  | 'receive'
  | 'ship'
  | 'adjust'
  | 'damage'
  | 'reserve'
  | 'release'
  | 'transfer_out'
  | 'transfer_in'

export class InventoryMovementModel extends DatabaseModel {
  public readonly inventoryId: string
  public readonly referenceId: string
  public readonly movementType: InventoryMovementType
  public readonly quantity: number
  public readonly reason: string

  constructor(data: unknown) {
    super((typeof data == 'object' && data != null ? data : {}) as object)
    this.inventoryId = this.getString('inventory_id')
    this.referenceId = this.getString('reference_id')
    this.movementType = this.getString('movement_type') as InventoryMovementType
    this.quantity = this.getNumber('quantity')
    this.reason = this.getString('reason')
  }
}
