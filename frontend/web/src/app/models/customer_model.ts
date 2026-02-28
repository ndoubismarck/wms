import { DatabaseModel } from '@/app/core/model'

export type CustomerTier = 'standard' | 'growth' | 'enterprise'
export type CustomerStatus = 'active' | 'paused' | 'churn_risk'

export interface CustomersSummaryModel {
  totalCustomers: number
  activeCount: number
  enterpriseCount: number
  avgLtv: number
}

export class CustomerModel extends DatabaseModel {
  public readonly name: string
  public readonly email: string
  public readonly phone: string
  public readonly city: string
  public readonly tier: CustomerTier
  public readonly status: CustomerStatus
  public readonly totalOrders: number
  public readonly lifetimeValue: number
  public readonly lastOrderAt: string

  constructor(data: unknown) {
    super((typeof data == 'object' && data != null ? data : {}) as object)
    this.name = this.getString('name')
    this.email = this.getString('email')
    this.phone = this.getString('phone')
    this.city = this.getString('city')
    this.tier = this.getString('tier') as CustomerTier
    this.status = this.getString('status') as CustomerStatus
    this.totalOrders = this.getNumber('total_orders')
    this.lifetimeValue = this.getFloat('lifetime_value')
    this.lastOrderAt = this.getString('last_order_at')
  }
}
