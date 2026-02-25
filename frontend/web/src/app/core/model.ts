import { get } from 'lodash'
import { DateTime } from '@/app/core/date_time'

export class Model {
  protected _data: object

  constructor(data: object | null | undefined) {
    this._data = data && typeof data === 'object' ? data : {}
  }

  protected get(key: string, def: any): any {
    const value = get(this._data, key)
    return value ? value : def
  }

  protected exists(key: string): boolean {
    return get(this._data, key)
  }

  protected getString(key: string, def?: string): string {
    return `${this.get(key, def ? def : '')}`
  }

  protected getFloat(key: string, def?: number): number {
    return parseFloat(`${this.get(key, def ? def : 0)}`)
  }

  protected getNumber(key: string, def?: number): number {
    return parseInt(`${this.get(key, def ? def : 0)}`)
  }

  protected getObject(key: string, def?: any): any {
    return this.get(key, def ? def : {}) as any
  }

  protected getArray(key: string, def?: any[]): any[] {
    return this.get(key, def ? def : []) as Array<any>
  }

  protected getBoolean(key: string, def?: boolean): boolean {
    return this.get(key, def ? def : false)
  }

  protected getDateTime(key: string, def?: number): DateTime {
    return new DateTime(this.get(key, def))
  }

  protected getLatitude(key: string, def?: number): number {
    const value = this.getString(key)
    if (!value) {
      return def ? def : 0
    }
    const parts = value.split(',')
    if (parts.length == 0) {
      return def ? def : 0
    }
    return parts[0] ? parseFloat(parts[0]) : 0
  }

  protected getLongitude(key: string, def?: number) {
    const value = this.getString(key)
    if (!value) {
      return def ? def : 0
    }
    const parts = value.split(',')
    if (parts.length < 2) {
      return def ? def : 0
    }
    return parts[1] ? parseFloat(parts[1]) : 0
  }

  protected getMoneyStringFromCents(key: string, def?: number): string {
    return (parseInt(`${this.get(key, def ? def : 0)}`) / 100).toFixed(2)
  }

  protected getMoneyFloatFromCents(key: string, def?: number): number {
    return parseFloat(`${this.get(key, def ? def : 0)}`) / 100
  }

  protected isNull(key: string): boolean {
    const value = this.get(key, null)
    return !value
  }

  protected hasValue(key: string): boolean {
    const value = this.get(key, null)
    return value != null || value != undefined
  }

  public getData(): any {
    return this._data
  }

  public toString(): string {
    const object = JSON.parse(JSON.stringify(this))
    if (object['_data']) {
      delete object['_data']
    }
    return JSON.stringify(object)
  }
}

export class DatabaseModel extends Model {
  public readonly id: string
  public readonly deletedAt: DateTime
  public readonly createdAt: DateTime
  public readonly updatedAt: DateTime

  constructor(data: object | null | undefined) {
    super(data)
    this.id = this.getString('id')
    this.deletedAt = this.getDateTime('deleted_at')
    this.createdAt = this.getDateTime('created_at')
    this.updatedAt = this.getDateTime('updated_at')
  }
}
