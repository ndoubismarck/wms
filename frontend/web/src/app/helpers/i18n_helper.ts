import { AppContext } from '@/app/core/context.ts'

export class I18nHelper {
  private ctx: AppContext

  constructor(ctx: AppContext) {
    this.ctx = ctx
  }

  public message(message: string): string {
    return message
  }

  public messageFromError(error: any): string {
    return 'An unexpected error has occurred.'
  }
}
