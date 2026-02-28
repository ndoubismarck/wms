import { DatabaseModel } from '@/app/core/model'

export type SupplierStatus = 'approved' | 'probation' | 'blocked'
export type SupplierType = 'raw_material' | 'packaging' | 'finished_goods' | 'service'

export interface SuppliersSummaryModel {
  totalSuppliers: number
  approvedCount: number
  avgLeadTime: number
  totalOpenPOs: number
}

export class SupplierModel extends DatabaseModel {
  public readonly name: string
  public readonly supplierType: SupplierType
  public readonly status: SupplierStatus
  public readonly contactName: string
  public readonly email: string
  public readonly leadTimeDays: number
  public readonly onTimeRate: number
  public readonly openPurchaseOrders: number
  public readonly city: string

  constructor(data: unknown) {
    super((typeof data == 'object' && data != null ? data : {}) as object)
    this.name = this.getString('name')
    this.supplierType = this.getString('supplier_type') as SupplierType
    this.status = this.getString('status') as SupplierStatus
    this.contactName = this.getString('contact_name')
    this.email = this.getString('email')
    this.leadTimeDays = this.getNumber('lead_time_days')
    this.onTimeRate = this.getFloat('on_time_rate')
    this.openPurchaseOrders = this.getNumber('open_purchase_orders')
    this.city = this.getString('city')
  }
}
