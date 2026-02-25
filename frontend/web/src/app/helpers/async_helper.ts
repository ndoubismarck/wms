import { AppContext } from '@/app/core/context'

export class AsyncHelper {
  private ctx: AppContext

  constructor(ctx: AppContext) {
    this.ctx = ctx
  }

  public async sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }
}
