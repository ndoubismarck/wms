import { Model } from '@/app/core/model'

export class PaginationModel extends Model {
  public readonly page: number
  public readonly prev: number
  public readonly next: number
  public readonly limit: number
  public readonly order: string

  public readonly hasPrev: boolean
  public readonly hasNext: boolean

  public readonly totalPages: number
  public readonly totalRecords: number

  constructor(data: any) {
    super(data)

    this.page = this.getNumber('page')
    this.prev = this.getNumber('prev')
    this.next = this.getNumber('next')
    this.limit = this.getNumber('limit')
    this.order = this.getString('order')

    this.hasPrev = this.getBoolean('has_prev')
    this.hasNext = this.getBoolean('has_next')

    this.totalPages = this.getNumber('total_pages')
    this.totalRecords = this.getNumber('total_records')
  }

  public get isFirstPage(): boolean {
    return this.page === 1
  }

  public get isLastPage(): boolean {
    return !this.hasNext
  }
}
