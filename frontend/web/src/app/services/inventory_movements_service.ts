import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import type { IApiResult } from '@/app/core/client'
import { PaginationModel } from '@/app/models/pagination_model'
import { InventoryMovementModel } from '@/app/models/inventory_movement_model'

type Payload = Record<string, unknown>

export class InventoryMovementsService extends AppService {
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

  public async getMany(params?: Payload): Promise<IGetManyResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/inventory/movements`, payload)
    const response = await this.api.private.get(path, true)
    return {
      movements: this.mapArray(response.body.data['movements'], (value) => new InventoryMovementModel(value)),
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
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
}

interface IGetManyResult extends IApiResult {
  movements: InventoryMovementModel[]
  pagination: PaginationModel
}
