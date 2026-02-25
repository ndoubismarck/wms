import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import { PaginationModel } from '@/app/models/pagination_model.ts'
import { UserModel } from '@/app/models/user_model'
import type { IApiResult } from '@/app/core/client'

export class UsersService extends AppService {
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

  public async getMany(params: any): Promise<IGetManyResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/users`, payload)
    const response = await this.api.private.get(path, true)
    const responseData = response.body.data['users']
    const resultModels: UserModel[] = []
    if (Array.isArray(responseData)) {
      responseData.forEach((val) => {
        resultModels.push(new UserModel(val))
      })
    }
    return {
      users: resultModels,
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async add(data: any): Promise<IAddResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/users`)
    const response = await this.api.private.post(path, payload, true)
    return {
      user: new UserModel(response.body.data['user']),
      success: response.body.code == 'success',
    }
  }

  public async update(id: string, data: any): Promise<IUpdateResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/users/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      user: new UserModel(response.body.data['user']),
      success: response.body.code == 'success',
    }
  }

  public async deleteById(id: string, params?: any): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/users/${id}`)
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

interface IAddResult extends IApiResult {
  user: UserModel
}

interface IUpdateResult extends IApiResult {
  user: UserModel
}

interface IDeleteByIdResult extends IApiResult {}

interface IGetManyResult extends IApiResult {
  users: UserModel[]
  pagination: PaginationModel
}
