export enum EPreloaderSize {
  SM = 'sm',
  MD = 'md',
  LG = 'lg',
}

export interface IAlertMessage {
  type: EAlertMessageType
  title?: IAlertMessageTitle
  body: IAlertMessageBody
}

export interface IAlertMessageTitle {
  html?: string
  text?: string
}

export interface IAlertMessageBody {
  html?: string
  text?: string
}

export enum EAlertMessageType {
  Info = 'info',
  Error = 'error',
  Success = 'success',
  Warning = 'warn',
  Progress = 'progress',
}
