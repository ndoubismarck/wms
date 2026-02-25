import { AppContext } from '@/app/core/context'
import { PrivateApiClient } from '@/app/core/client'

export class AppService {
  protected readonly ctx: AppContext
  protected readonly api: IApiClient

  constructor(ctx: AppContext) {
    this.ctx = ctx
    this.api = {
      private: new PrivateApiClient(ctx),
    }
  }

  protected mapArray<T>(data: any, mapper: (value: any) => T): T[] {
    const result: T[] = []
    if (Array.isArray(data)) {
      data.forEach((val) => {
        result.push(mapper(val))
      })
    }
    return result
  }
}

interface IApiClient {
  private: PrivateApiClient
}
