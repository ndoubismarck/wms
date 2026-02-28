import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import type { IApiResult } from '@/app/core/client'
import { PaginationModel } from '@/app/models/pagination_model'
import { TeamModel } from '@/app/models/team_model'

export class TeamsService extends AppService {
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
    const path = this.api.private.path(`/${resolvedLocationID}/teams`, payload)
    const response = await this.api.private.get(path, true)
    return {
      teams: this.mapArray(response.body.data['teams'], (value) => new TeamModel(value)),
      success: response.body.code == 'success',
      pagination: new PaginationModel(response.pagination),
    }
  }

  public async getOne(id: string, params?: any): Promise<IGetOneResult> {
    const { locationId, payload } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams/${id}`, payload)
    const response = await this.api.private.get(path, true)
    return {
      team: new TeamModel(response.body.data['team']),
      success: response.body.code == 'success',
    }
  }

  public async add(data: any): Promise<IAddResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams`)
    const response = await this.api.private.post(path, payload, true)
    return {
      team: new TeamModel(response.body.data['team']),
      success: response.body.code == 'success',
    }
  }

  public async update(id: string, data: any): Promise<IUpdateResult> {
    const { locationId, payload } = this.extractPayload(data)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams/${id}`)
    const response = await this.api.private.put(path, payload, true)
    return {
      team: new TeamModel(response.body.data['team']),
      success: response.body.code == 'success',
    }
  }

  public async deleteById(id: string, params?: any): Promise<IDeleteByIdResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams/${id}`)
    const response = await this.api.private.delete(path, true)
    return {
      success: response.body.code == 'success',
    }
  }

  public async getMembers(teamID: string, params?: any): Promise<IGetMembersResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams/${teamID}/users`)
    const response = await this.api.private.get(path, true)
    const userIDs = response.body.data['user_ids']
    return {
      userIDs: Array.isArray(userIDs) ? userIDs : [],
      success: response.body.code == 'success',
    }
  }

  public async setMembers(teamID: string, userIDs: string[], params?: any): Promise<ISetMembersResult> {
    const { locationId } = this.extractPayload(params)
    const resolvedLocationID = await this.resolveLocationID(locationId)
    const path = this.api.private.path(`/${resolvedLocationID}/teams/${teamID}/users`)
    const response = await this.api.private.put(path, { user_ids: userIDs }, true)
    const ids = response.body.data['user_ids']
    return {
      userIDs: Array.isArray(ids) ? ids : [],
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

interface IGetManyResult extends IApiResult {
  teams: TeamModel[]
  pagination: PaginationModel
}

interface IGetOneResult extends IApiResult {
  team: TeamModel
}

interface IAddResult extends IApiResult {
  team: TeamModel
}

interface IUpdateResult extends IApiResult {
  team: TeamModel
}

interface IDeleteByIdResult extends IApiResult {}

interface IGetMembersResult extends IApiResult {
  userIDs: string[]
}

interface ISetMembersResult extends IApiResult {
  userIDs: string[]
}
