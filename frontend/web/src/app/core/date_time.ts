import moment from 'moment'
import momentTimezone from 'moment-timezone'

export class DateTime {
  public readonly year: number
  public readonly month: number

  public static timeFormat: string = 'hh:mm A'
  public static dateFormat: string = 'MM/DD/YYYY'
  public static dateTimeFormat: string = 'MM/DD/YYYY hh:mm A'

  private readonly date: Date

  constructor(value?: string | number | Date | null | undefined, format?: string) {
    let date: Date = new Date()
    if (value) {
      if (typeof value === 'object') {
        date = momentTimezone(value).toDate()
      }
      if (typeof value === 'string') {
        if (format) {
          date = momentTimezone(value, format).toDate()
        } else {
          date = momentTimezone(new Date(value)).toDate()
        }
      }
      if (typeof value === 'number') {
        date = momentTimezone(new Date((value as number) * 1000)).toDate()
      }
    }
    this.date = date
    this.year = this.date.getFullYear()
    this.month = this.date.getMonth()
  }

  public unix(): number {
    if (!this.date) {
      return 0
    }
    return momentTimezone(this.date).unix()
  }

  public toDate(): Date {
    if (!this.date) {
      return new Date()
    }
    return this.date
  }

  public isPast(): boolean {
    return momentTimezone().diff(this.date) > 0
  }

  public isFuture(): boolean {
    return momentTimezone().diff(this.date) < 0
  }

  public moment(): momentTimezone.Moment {
    return momentTimezone(this.date)
  }

  public isToday(): boolean {
    const now = momentTimezone()
    const date = momentTimezone(this.date)
    return now.year() == date.year() && now.month() == date.month() && now.day() == date.day()
  }

  public timeAgo(unit: string | DateTimeUnit): number {
    return momentTimezone().diff(this.date, unit as moment.unitOfTime.Diff)
  }

  public format(format: string, timezone?: null | string): string {
    if (timezone) {
      return momentTimezone(this.date).tz(timezone).format(format)
    }
    return momentTimezone(this.date).format(format)
  }

  public formatted(): string {
    return this.toDateTimeString()
  }

  public toUnix(): number {
    return momentTimezone(this.date).unix()
  }

  public getHours(): number {
    return momentTimezone(this.date).get('hours')
  }

  public greeting(): string {
    const hh = parseInt(this.moment().format('HH'))
    return hh < 12 ? 'Good morning' : hh < 18 ? 'Good afternoon' : 'Good evening'
  }

  public toString(): string {
    return this.toISOString()
  }

  public toISOString(): string {
    return momentTimezone(this.date).toISOString(false)
  }

  public toDateString(): string {
    return this.format(DateTime.dateFormat)
  }

  public toTimeString(): string {
    return this.format(DateTime.timeFormat)
  }

  public toDateTimeString(): string {
    return this.format(DateTime.dateTimeFormat)
  }
}

export enum DateTimeUnit {
  days = 'days',
  hours = 'hours',
  weeks = 'weeks',
  years = 'years',
  months = 'months',
  minutes = 'minutes',
  seconds = 'seconds',
  microseconds = 'microseconds',
  milliseconds = 'milliseconds',
}
