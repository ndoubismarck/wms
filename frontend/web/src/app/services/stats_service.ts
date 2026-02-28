import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import type { IApiResult } from '@/app/core/client'
import { StatsModel } from '@/app/models/stats_model'

export class StatsService extends AppService {
  constructor(ctx: AppContext) {
    super(ctx)
  }

  public async get(lid: string): Promise<IGetSetupResult> {
    const path = `/${lid}/stats`
    const response = await this.api.private.get(path, true)
    return {
      stats: new StatsModel(response.body['data']['stats']),
      success: response.body.code == 'success',
    }
  }
}

interface IGetSetupResult extends IApiResult {
  stats: StatsModel
}
