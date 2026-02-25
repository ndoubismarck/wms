import type { AppContext } from '@/app/core/context'

export class Events {
  private readonly events: any

  constructor(ctx: AppContext) {
    this.events = {}
  }

  public on(name: string, callback: (data: any) => void) {
    if (!this.events[name]) {
      this.events[name] = []
    }
    this.events[name].push(callback)
  }

  public emit(name: string, data: any) {
    if (this.events[name]) {
      this.events[name].forEach((val: (data: any) => void) => {
        val(data)
      })
    }
  }
}

export const HttpResponseEvent = 'http.response'
