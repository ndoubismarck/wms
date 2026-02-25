import Cookies from 'js-cookie'
import { AppContext } from '@/app/core/context'
import { HttpResponseEvent } from '@/app/core/events'

export class PrivateApiClient {
  private readonly ctx: AppContext

  constructor(ctx: AppContext) {
    this.ctx = ctx
  }

  public url(path: string): string {
    if (import.meta.env.DEV || !import.meta.env.PROD) {
      let host = window.location.hostname
      if (!host || host.length == 0) {
        host = 'localhost'
      }
      return `http://${host}:8080/app/v1${path}`
    }
    return `/app/v1${path}`
  }

  public path(path: string, params?: any) {
    if (!params) {
      return path
    }
    params = new URLSearchParams(params).toString()
    return `${path}?${params}`
  }

  public async get(path: string, auth: boolean = false): Promise<IApiResponse> {
    const url = this.url(path)
    const headers = await this.getHeaders(auth)
    const response = await fetch(url, {
      method: 'get',
      headers: headers,
    })
    return this.handleJsonResponse(response)
  }

  public async put(path: string, body: any, auth: boolean = false): Promise<IApiResponse> {
    const url = this.url(path)
    const headers = await this.getHeaders(auth)
    const response = await fetch(url, {
      body: JSON.stringify(body),
      method: 'put',
      headers: headers,
    })
    return this.handleJsonResponse(response)
  }

  public async post(path: string, body: any, auth: boolean = false): Promise<IApiResponse> {
    const url = this.url(path)
    const headers = await this.getHeaders(auth)
    const response = await fetch(url, {
      body: JSON.stringify(body),
      method: 'post',
      headers: headers,
    })
    return this.handleJsonResponse(response)
  }

  public async delete(path: string, auth: boolean = false): Promise<IApiResponse> {
    const url = this.url(path)
    const headers = await this.getHeaders(auth)
    const response = await fetch(url, {
      method: 'delete',
      headers: headers,
    })
    return this.handleJsonResponse(response)
  }

  public async onEvent(path: string, name: string, callback: any) {
    const url = this.url(path)
    const token = await this.getAuthToken()
    const eventSource = new EventSource(this.url(`${url}?authToken=${token}`))
    eventSource.onopen = () => {}
    eventSource.addEventListener(name, (event: MessageEvent) => {
      callback(event.data)
    })
    eventSource.onerror = (event) => {}
  }

  public async getHeaders(auth: boolean = false): Promise<Record<string, string>> {
    const headers: Record<string, string> = {}
    if (auth) {
      const token = await this.getAuthToken()
      headers['Authorization'] = `Bearer ${token}`
    }
    return headers
  }

  public async deleteAuthToken(): Promise<boolean> {
    Cookies.remove('auth_token')
    return (await this.getAuthToken()) == null
  }

  public async saveAuthToken(token: string) {
    const payload = this.decodeAuthToken(token)
    if (!payload || !payload.expiresAt) {
      Cookies.set('auth_token', token, {
        sameSite: 'Strict',
      })
    } else {
      Cookies.set('auth_token', token, {
        expires: payload.expiresAt,
        sameSite: 'Strict',
      })
    }
  }

  public async getAuthToken(): Promise<string | null> {
    const token = Cookies.get('auth_token')
    if (!token || token.trim().length == 0) {
      return null
    }
    return token.trim()
  }

  public async hasAuthToken(): Promise<boolean> {
    return (await this.getAuthToken()) != null
  }

  private decodeAuthToken(token: string | null): IJWTToken | null {
    if (token == null) {
      return null
    }
    let val = token.split('.')[1]
    let data = {
      id: null,
      iat: null,
      exp: null,
    }
    if (!val) {
      return null
    }
    data = JSON.parse(atob(val))

    const id = data.id
    const exp = data.exp
    const iat = data.iat
    if (!id || id <= 0) {
      return null
    }
    return {
      id: id,
      issuedAt: iat ? new Date(iat * 1000) : null,
      expiresAt: exp ? new Date(exp * 1000) : null,
    }
  }

  private async handleJsonResponse(response: Response): Promise<IApiResponse> {
    const jsonResponse = await response.json()
    this.ctx.events.emit(HttpResponseEvent, response)
    if (!jsonResponse) {
      return {
        body: {
          code: 'failed',
          data: {},
        },
        status: response.status,
      }
    }
    return {
      body: {
        code: jsonResponse.code ? jsonResponse.code : 'failed',
        data: jsonResponse.data ? jsonResponse.data : {},
      },
      status: response.status,
      pagination: jsonResponse.pagination,
    }
  }
}

interface IJWTToken {
  id: string
  issuedAt: Date | null
  expiresAt: Date | null
}

interface IApiResponse {
  body: {
    code: string
    data: any
  }
  status: number
  pagination?: any
}

export interface IApiResult {
  success: boolean
}
