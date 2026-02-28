import { Events } from '@/app/core/events'
import { AppContext } from '@/app/core/context'
import { AuthService } from '@/app/services/auth_service'
import { ProductsService } from '@/app/services/products_service'
import { LocationsService } from '@/app/services/locations_service'
import { LoggerService } from '@/app/services/logger_service'
import { I18nHelper } from '@/app/helpers/i18n_helper.ts'
import { AsyncHelper } from '@/app/helpers/async_helper'
import { SetupService } from '@/app/services/setup_service'
import { StatsService } from '@/app/services/stats_service'
import { VeeValidateHelper } from '@/app/helpers/vee-validate.ts'
import type { PrivateApiClient } from '@/app/core/client.ts'
import { UsersService } from '@/app/services/users_service'
import { TeamsService } from '@/app/services/teams_service'
import { TasksService } from '@/app/services/tasks_service'

class App {
  public readonly api: IApiClient
  public readonly events: Events
  public readonly helpers: IHelpers
  public readonly services: IServices

  constructor() {
    const ctx = new AppContext()
    this.api = {
      private: ctx.privateApiClient,
    }
    this.events = ctx.events
    this.helpers = {
      i18n: new I18nHelper(ctx),
      async: new AsyncHelper(ctx),
      veeValidate: new VeeValidateHelper(ctx),
    }
    this.services = {
      auth: new AuthService(ctx),
      setup: new SetupService(ctx),
      stats: new StatsService(ctx),
      logger: new LoggerService(ctx),
      products: new ProductsService(ctx),
      locations: new LocationsService(ctx),
      users: new UsersService(ctx),
      teams: new TeamsService(ctx),
      tasks: new TasksService(ctx),
    }
  }
}

const app = new App()
export const useApp = (): App => {
  return app
}

interface IHelpers {
  i18n: I18nHelper
  async: AsyncHelper
  veeValidate: VeeValidateHelper
}

interface IApiClient {
  private: PrivateApiClient
}

interface IServices {
  auth: AuthService
  setup: SetupService
  stats: StatsService
  logger: LoggerService
  products: ProductsService
  locations: LocationsService
  users: UsersService
  teams: TeamsService
  tasks: TasksService
}
