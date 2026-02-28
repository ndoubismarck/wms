import { DatabaseModel } from '@/app/core/model'

export type OrderStatus = 'new' | 'picking' | 'packed' | 'shipped' | 'backorder' | 'cancelled'

export type OrderPriority = 'low' | 'normal' | 'high' | 'urgent'

export interface OrdersSummaryModel {
  totalOrders: number
  openOrders: number
  urgentOrders: number
  totalOrderValue: number
}

export class OrderModel extends DatabaseModel {
  public readonly orderNumber: string
  public readonly customerName: string
  public readonly channel: string
  public readonly status: OrderStatus
  public readonly priority: OrderPriority
  public readonly itemsCount: number
  public readonly totalAmount: number
  public readonly dueAt: string
  public readonly notes: string

  constructor(data: unknown) {
    super((typeof data == 'object' && data != null ? data : {}) as object)
    this.orderNumber = this.getString('order_number')
    this.customerName = this.getString('customer_name')
    this.channel = this.getString('channel')
    this.status = this.getString('status') as OrderStatus
    this.priority = this.getString('priority') as OrderPriority
    this.itemsCount = this.getNumber('items_count')
    this.totalAmount = this.getFloat('total_amount')
    this.dueAt = this.getString('due_at')
    this.notes = this.getString('notes')
  }
}
