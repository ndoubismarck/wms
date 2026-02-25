import Cookies from 'js-cookie'
import type { AppContext } from '@/app/core/context'

export class Storage {
  private readonly data: any

  constructor(ctx: AppContext) {
    this.data = {}
  }

  public set(key: string, value: string, options?: any) {
    this.data[key] = value
    Cookies.set(key, value)
  }

  public async get(key: string): Promise<string | null> {
    let value = this.data[key]
    if (value) {
      return value
    }
    value = Cookies.get(key)
    if (value) {
      return value
    }
    return null
  }
}
