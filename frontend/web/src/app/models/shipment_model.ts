import { DatabaseModel } from '@/app/core/model'

export type ShipmentStatus = 'draft' | 'ready' | 'in_transit' | 'delivered' | 'exception'

export interface ShipmentsSummaryModel {
  totalShipments: number
  inTransit: number
  delivered: number
  exceptions: number
  totalPackages: number
}

export class ShipmentModel extends DatabaseModel {
  public readonly shipmentNumber: string
  public readonly orderNumber: string
  public readonly carrier: string
  public readonly service: string
  public readonly status: ShipmentStatus
  public readonly packages: number
  public readonly trackingCode: string
  public readonly eta: string
  public readonly destination: string

  constructor(data: unknown) {
    super((typeof data == 'object' && data != null ? data : {}) as object)
    this.shipmentNumber = this.getString('shipment_number')
    this.orderNumber = this.getString('order_number')
    this.carrier = this.getString('carrier')
    this.service = this.getString('service')
    this.status = this.getString('status') as ShipmentStatus
    this.packages = this.getNumber('packages')
    this.trackingCode = this.getString('tracking_code')
    this.eta = this.getString('eta')
    this.destination = this.getString('destination')
  }
}
