import { AppService } from '@/app/core/service'
import { AppContext } from '@/app/core/context'
import type { IApiResult } from '@/app/core/client'
import { SetupModel } from '@/app/models/setup_model'

export class SetupService extends AppService {
  constructor(ctx: AppContext) {
    super(ctx)
  }

  public async get(): Promise<IGetSetupResult> {
    let path = `/setup`
    const response = await this.api.private.get(path, true)
    const model = new SetupModel(response.body['data']['setup'])
    if (model.location.id) {
      this.ctx.storage.set('location_id', model.location.id)
    }
    return {
      setup: model,
      success: response.body.code == 'success',
    }
  }
}

interface IGetSetupResult extends IApiResult {
  setup: SetupModel
}
