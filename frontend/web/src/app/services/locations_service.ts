import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import { PaginationModel } from '@/app/models/pagination_model.ts'
import { LocationModel } from '@/app/models/location_model'
import { LocationAisleModel } from '@/app/models/location_aisle_model.ts'
import { LocationBayModel } from '@/app/models/location_bay_model.ts'
import { LocationShelfModel } from '@/app/models/location_shelf_model.ts'
import { LocationShelfLevelModel } from '@/app/models/location_shelf_level_model.ts'
import { LocationBinModel } from '@/app/models/location_bin_model.ts'
import type { IApiResult } from '@/app/core/client'

export class LocationsService extends AppService {
  private currentLocationId: string | null = null

  constructor(ctx: AppContext) {
    super(ctx)
  }

  public async getMany(params: any): Promise<IGetManyResult> {
    const path = this.api.private.path('/locations', params)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['locations']
    const resultModels: LocationModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationModel(val))
      })
    }
    return {
      locations: resultModels,
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public setCurrentLocationId(locationId: string | null): void {
    this.currentLocationId = locationId
  }

  public getCurrentLocationId(): string | null {
    return this.currentLocationId
  }

  public async addAisle(data: any): Promise<IAddAisleResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/aisles`)
    const response = await this.api.private.post(path, payload, true)
    return {
      aisle: new LocationAisleModel(response.body.data['aisle']),
      success: response.body.code == 'success',
    }
  }

  public async getAisles(params: any): Promise<IGetAislesResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/aisles`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['aisles']
    const resultModels: LocationAisleModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationAisleModel(val))
      })
    }
    return {
      aisles: resultModels,
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addBay(data: any): Promise<IAddBayResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/bays`)
    const response = await this.api.private.post(path, payload, true)
    return {
      bay: new LocationBayModel(response.body.data['bay']),
      success: response.body.code == 'success',
    }
  }

  public async getBays(params: any): Promise<IGetBaysResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/bays`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['bays']
    const resultModels: LocationBayModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationBayModel(val))
      })
    }
    return {
      bays: resultModels,
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addShelf(data: any): Promise<IAddShelfResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/shelves`)
    const response = await this.api.private.post(path, payload, true)
    return {
      shelf: new LocationShelfModel(response.body.data['shelf']),
      success: response.body.code == 'success',
    }
  }

  public async getShelves(params: any): Promise<IGetShelvesResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/shelves`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['shelves']
    const resultModels: LocationShelfModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationShelfModel(val))
      })
    }
    return {
      shelves: resultModels,
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async addShelfLevel(data: any): Promise<IAddShelfLevelResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/shelves/levels`)
    const response = await this.api.private.post(path, payload, true)
    return {
      success: response.body.code == 'success',
      shelfLevel: new LocationShelfLevelModel(response.body.data['shelf_level']),
    }
  }

  public async getShelfLevels(params: any): Promise<IGetShelfLevelsResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/shelves/levels`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['shelf_levels']
    const resultModels: LocationShelfLevelModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationShelfLevelModel(val))
      })
    }
    return {
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
      shelfLevels: resultModels,
    }
  }

  public async addBin(data: any): Promise<IAddBinResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/bins`)
    const response = await this.api.private.post(path, payload, true)
    return {
      bin: new LocationBinModel(response.body.data['bin']),
      success: response.body.code == 'success',
    }
  }

  public async getBins(params: any): Promise<IGetBinsResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/locations/bins`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['bins']
    const resultModels: LocationBinModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new LocationBinModel(val))
      })
    }
    return {
      bins: resultModels,
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
    const result = await this.getMany({
      page: 1,
      limit: 1,
    })
    if (!result.success || result.locations.length == 0) {
      throw new Error('No locations found')
    }
    this.currentLocationId = result.locations[0]?.id ?? null
    if (!this.currentLocationId || this.currentLocationId.trim().length == 0) {
      throw new Error('No locations found')
    }
    return this.currentLocationId
  }

  private extractPayload(payload?: any): {
    locationId?: string | null
    payload: any
  } {
    if (!payload || typeof payload !== 'object') {
      return {
        payload: {},
      }
    }
    const copy = { ...payload }
    const locationId = copy['location_id'] ?? copy['locationId'] ?? null
    delete copy['location_id']
    delete copy['locationId']
    return {
      locationId,
      payload: copy,
    }
  }
}

interface IGetManyResult extends IApiResult {
  locations: LocationModel[]
  pagination: PaginationModel
}

interface IAddAisleResult extends IApiResult {
  aisle: LocationAisleModel
}

interface IGetAislesResult extends IApiResult {
  aisles: LocationAisleModel[]
  pagination: PaginationModel
}

interface IAddBayResult extends IApiResult {
  bay: LocationBayModel
}

interface IGetBaysResult extends IApiResult {
  bays: LocationBayModel[]
  pagination: PaginationModel
}

interface IAddShelfResult extends IApiResult {
  shelf: LocationShelfModel
}

interface IGetShelvesResult extends IApiResult {
  shelves: LocationShelfModel[]
  pagination: PaginationModel
}

interface IAddShelfLevelResult extends IApiResult {
  shelfLevel: LocationShelfLevelModel
}

interface IGetShelfLevelsResult extends IApiResult {
  pagination: PaginationModel
  shelfLevels: LocationShelfLevelModel[]
}

interface IAddBinResult extends IApiResult {
  bin: LocationBinModel
}

interface IGetBinsResult extends IApiResult {
  bins: LocationBinModel[]
  pagination: PaginationModel
}
