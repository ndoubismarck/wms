import { AppContext } from '@/app/core/context'
import { AppService } from '@/app/core/service'

export class LoggerService extends AppService {
  constructor(ctx: AppContext) {
    super(ctx)
  }

  public log(...data: any[]) {
    console.log(...data)
  }

  public info(...data: any[]) {
    console.info(...data)
  }

  public warn(...data: any[]) {
    console.warn(...data)
  }

  public error(...data: any[]) {
    console.error(...data)
  }
}
