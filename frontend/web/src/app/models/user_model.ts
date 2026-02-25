import { DatabaseModel } from '@/app/core/model'

export class UserModel extends DatabaseModel {
  public readonly role: string
  public readonly firstName: string
  public readonly lastName: string
  public readonly fullName: string
  public readonly emailAddress: string

  constructor(data: any) {
    super(data)
    this.role = this.getString('role')
    this.firstName = this.getString('first_name')
    this.lastName = this.getString('last_name')
    this.fullName = `${this.getString('first_name')} ${this.getString('last_name')}`
    this.emailAddress = this.getString('email_address')
  }
}
