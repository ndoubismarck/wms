import { Events } from '@/app/core/events'
import { Storage } from '@/app/core/storage'
import { PrivateApiClient } from '@/app/core/client'

export class AppContext {
  public readonly events: Events
  public readonly storage: Storage
  public readonly privateApiClient: PrivateApiClient

  constructor() {
    this.events = new Events(this)
    this.storage = new Storage(this)
    this.privateApiClient = new PrivateApiClient(this)
  }
}
