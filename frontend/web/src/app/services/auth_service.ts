import { AppContext } from '@/app/core/context'
import { AppService } from '@/app/core/service'
import { UserModel } from '@/app/models/user_model'
import type { IApiResult } from '@/app/core/client'
import { HttpResponseEvent } from '@/app/core/events'

export class AuthService extends AppService {
  private onUnAuthorizedCallback?: (response: Response) => void

  constructor(ctx: AppContext) {
    super(ctx)
    ctx.events.on(HttpResponseEvent, (response: Response) => {
      if (response.status == 401) {
        if (this.onUnAuthorizedCallback) {
          this.onUnAuthorizedCallback(response)
        }
      }
    })
  }

  public async login(body: any): Promise<ILoginResult> {
    const response = await this.api.private.post('/auth/login', body)
    const token = response.body['data']['token']
    const success = response.body.code == 'success'
    if (success) {
      await this.api.private.saveAuthToken(token)
    }
    return {
      user: new UserModel(response.body['data']['user']),
      success: success,
    }
  }

  public async getMe(): Promise<IGetUserResult> {
    let path = `/me`
    const response = await this.api.private.get(path, true)
    return {
      user: new UserModel(response.body['data']['user']),
      success: response.body.code == 'success',
    }
  }

  public async logout(): Promise<boolean> {
    return this.api.private.deleteAuthToken()
  }

  public onUnAuthorized(callback: (response: Response) => void) {
    this.onUnAuthorizedCallback = callback
  }
}

interface IGetUserResult extends IApiResult {
  user: UserModel
}

interface ILoginResult extends IApiResult {
  user: UserModel
}
