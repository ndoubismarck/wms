import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import type { IApiResult } from '@/app/core/client'
import { PaginationModel } from '@/app/models/pagination_model'
import { OrderModel, type OrdersSummaryModel } from '@/app/models/order_model'
import { ShipmentModel, type ShipmentsSummaryModel } from '@/app/models/shipment_model'
import { CustomerModel, type CustomersSummaryModel } from '@/app/models/customer_model'
import { SupplierModel, type SuppliersSummaryModel } from '@/app/models/supplier_model'

type Payload = Record<string, unknown>

const emptyOrdersSummary = (): OrdersSummaryModel => ({
  totalOrders: 0,
  openOrders: 0,
  urgentOrders: 0,
  totalOrderValue: 0,
})

const emptyShipmentsSummary = (): ShipmentsSummaryModel => ({
  totalShipments: 0,
  inTransit: 0,
  delivered: 0,
  exceptions: 0,
  totalPackages: 0,
})

const emptyCustomersSummary = (): CustomersSummaryModel => ({
  totalCustomers: 0,
  activeCount: 0,
  enterpriseCount: 0,
  avgLtv: 0,
})

const emptySuppliersSummary = (): SuppliersSummaryModel => ({
  totalSuppliers: 0,
  approvedCount: 0,
  avgLeadTime: 0,
  totalOpenPOs: 0,
})

export class OperationsService extends AppService {
  private currentLocationId: string | null = null

  constructor(ctx: AppContext) {
    super(ctx)
  }

  public setCurrentLocationId(locationId: string | null): void {
    this.currentLocationId = locationId
  }

  public getCurrentLocationId(): string | null {
    return this.currentLocationId
  }

  public async getOrders(params?: Payload): Promise<IGetOrdersResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/orders`, payload)
    const response = await this.api.private.get(path, true)
    const summary = this.getObject(response.body.data['summary'])
    return {
      orders: this.mapArray(response.body.data['orders'], (value) => new OrderModel(value)),
      summary: {
        totalOrders: Number(summary['total_orders'] ?? 0),
        openOrders: Number(summary['open_orders'] ?? 0),
        urgentOrders: Number(summary['urgent_orders'] ?? 0),
        totalOrderValue: Number(summary['total_order_value'] ?? 0),
      },
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addOrder(data: Payload): Promise<IAddOrderResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/orders`)
    const response = await this.api.private.post(path, payload, true)
    return {
      order: new OrderModel(response.body.data['order']),
      success: response.body.code == 'success',
    }
  }

  public async updateOrder(id: string, data: Payload): Promise<IUpdateOrderResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/orders/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      order: new OrderModel(response.body.data['order']),
      success: response.body.code == 'success',
    }
  }

  public async deleteOrderById(id: string, params?: Payload): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/orders/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getShipments(params?: Payload): Promise<IGetShipmentsResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/shipments`, payload)
    const response = await this.api.private.get(path, true)
    const summary = this.getObject(response.body.data['summary'])
    return {
      shipments: this.mapArray(response.body.data['shipments'], (value) => new ShipmentModel(value)),
      summary: {
        totalShipments: Number(summary['total_shipments'] ?? 0),
        inTransit: Number(summary['in_transit'] ?? 0),
        delivered: Number(summary['delivered'] ?? 0),
        exceptions: Number(summary['exceptions'] ?? 0),
        totalPackages: Number(summary['total_packages'] ?? 0),
      },
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addShipment(data: Payload): Promise<IAddShipmentResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/shipments`)
    const response = await this.api.private.post(path, payload, true)
    return {
      shipment: new ShipmentModel(response.body.data['shipment']),
      success: response.body.code == 'success',
    }
  }

  public async updateShipment(id: string, data: Payload): Promise<IUpdateShipmentResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/shipments/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      shipment: new ShipmentModel(response.body.data['shipment']),
      success: response.body.code == 'success',
    }
  }

  public async deleteShipmentById(id: string, params?: Payload): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/shipments/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getCustomers(params?: Payload): Promise<IGetCustomersResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/customers`, payload)
    const response = await this.api.private.get(path, true)
    const summary = this.getObject(response.body.data['summary'])
    return {
      customers: this.mapArray(response.body.data['customers'], (value) => new CustomerModel(value)),
      summary: {
        totalCustomers: Number(summary['total_customers'] ?? 0),
        activeCount: Number(summary['active_count'] ?? 0),
        enterpriseCount: Number(summary['enterprise_count'] ?? 0),
        avgLtv: Number(summary['avg_ltv'] ?? 0),
      },
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addCustomer(data: Payload): Promise<IAddCustomerResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/customers`)
    const response = await this.api.private.post(path, payload, true)
    return {
      customer: new CustomerModel(response.body.data['customer']),
      success: response.body.code == 'success',
    }
  }

  public async updateCustomer(id: string, data: Payload): Promise<IUpdateCustomerResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/customers/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      customer: new CustomerModel(response.body.data['customer']),
      success: response.body.code == 'success',
    }
  }

  public async deleteCustomerById(id: string, params?: Payload): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/customers/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getSuppliers(params?: Payload): Promise<IGetSuppliersResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/suppliers`, payload)
    const response = await this.api.private.get(path, true)
    const summary = this.getObject(response.body.data['summary'])
    return {
      suppliers: this.mapArray(response.body.data['suppliers'], (value) => new SupplierModel(value)),
      summary: {
        totalSuppliers: Number(summary['total_suppliers'] ?? 0),
        approvedCount: Number(summary['approved_count'] ?? 0),
        avgLeadTime: Number(summary['avg_lead_time'] ?? 0),
        totalOpenPOs: Number(summary['total_open_pos'] ?? 0),
      },
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addSupplier(data: Payload): Promise<IAddSupplierResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/suppliers`)
    const response = await this.api.private.post(path, payload, true)
    return {
      supplier: new SupplierModel(response.body.data['supplier']),
      success: response.body.code == 'success',
    }
  }

  public async updateSupplier(id: string, data: Payload): Promise<IUpdateSupplierResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/suppliers/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      supplier: new SupplierModel(response.body.data['supplier']),
      success: response.body.code == 'success',
    }
  }

  public async deleteSupplierById(id: string, params?: Payload): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/operations/suppliers/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  private async resolveLocationID(locationId?: string | null): Promise<string> {
    if (locationId && locationId.trim().length > 0) {
      this.currentLocationId = locationId.trim()
      return this.currentLocationId
    }
    if (this.currentLocationId && this.currentLocationId.trim().length > 0) {
      return this.currentLocationId
    }
    const path = this.api.private.path('/locations', {
      page: 1,
      limit: 1,
    })
    const response = await this.api.private.get(path, true)
    if (response.body.code != 'success') {
      throw new Error('No locations found')
    }
    const locations = response.body.data['locations']
    if (!Array.isArray(locations) || locations.length == 0) {
      throw new Error('No locations found')
    }
    const defaultLocationID = `${locations[0]?.id ?? ''}`.trim()
    if (defaultLocationID.length == 0) {
      throw new Error('No locations found')
    }
    this.currentLocationId = defaultLocationID
    return this.currentLocationId
  }

  private extractPayload(payload?: Payload): {
    locationId?: string | null
    payload: Payload
  } {
    if (!payload || typeof payload != 'object') {
      return {
        payload: {},
      }
    }
    const copy: Payload = { ...payload }
    const locationField = copy['location_id'] ?? copy['locationId'] ?? null
    const locationId = typeof locationField == 'string' ? locationField : null
    delete copy['location_id']
    delete copy['locationId']
    return {
      locationId,
      payload: copy,
    }
  }

  private getObject(value: unknown): Record<string, unknown> {
    if (typeof value == 'object' && value != null) {
      return value as Record<string, unknown>
    }
    return {}
  }
}

interface IGetOrdersResult extends IApiResult {
  orders: OrderModel[]
  summary: OrdersSummaryModel
  pagination: PaginationModel
}

interface IAddOrderResult extends IApiResult {
  order: OrderModel
}

interface IUpdateOrderResult extends IApiResult {
  order: OrderModel
}

interface IGetShipmentsResult extends IApiResult {
  shipments: ShipmentModel[]
  summary: ShipmentsSummaryModel
  pagination: PaginationModel
}

interface IAddShipmentResult extends IApiResult {
  shipment: ShipmentModel
}

interface IUpdateShipmentResult extends IApiResult {
  shipment: ShipmentModel
}

interface IGetCustomersResult extends IApiResult {
  customers: CustomerModel[]
  summary: CustomersSummaryModel
  pagination: PaginationModel
}

interface IAddCustomerResult extends IApiResult {
  customer: CustomerModel
}

interface IUpdateCustomerResult extends IApiResult {
  customer: CustomerModel
}

interface IGetSuppliersResult extends IApiResult {
  suppliers: SupplierModel[]
  summary: SuppliersSummaryModel
  pagination: PaginationModel
}

interface IAddSupplierResult extends IApiResult {
  supplier: SupplierModel
}

interface IUpdateSupplierResult extends IApiResult {
  supplier: SupplierModel
}

interface IDeleteByIdResult extends IApiResult {}

export { emptyOrdersSummary, emptyShipmentsSummary, emptyCustomersSummary, emptySuppliersSummary }
