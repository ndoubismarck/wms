import { configure, defineRule } from 'vee-validate'
import * as rules from '@vee-validate/rules'
import { localize } from '@vee-validate/i18n'

import { AppContext } from '@/app/core/context.ts'

export class VeeValidateHelper {
  private ctx: AppContext

  constructor(ctx: AppContext) {
    this.ctx = ctx
    Object.keys(rules).forEach((rule: string) => {
      const validator = rules[rule as keyof typeof rules] as ValidatorFn
      if (typeof validator === 'function') {
        defineRule(rule, validator)
      }
    })

    configure({
      generateMessage: localize('en', {
        messages: {
          email: 'Please enter a valid email address',
          required: 'This field is required',
        },
      }),
    })
  }

  public errorMessage(message?: string): string | undefined {
    if (message) {
      message = message.trim()
      message = message.charAt(0).toUpperCase() + message.slice(1)
    }
    return message
  }
}

type ValidatorFn = (
  value: unknown,
  params:
    | [string]
    | {
        locale?: string
      },
) => boolean
