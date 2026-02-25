import { DatabaseModel } from '@/app/core/model'

export class StatsModel extends DatabaseModel {
  public readonly inventory: {
    bins: Metric
    stock: Metric
    products: Metric
    productVariants: Metric
  }

  constructor(data: any) {
    super(data)
    this.inventory = {
      bins: {
        total: this.getNumber('inventory.bins.total'),
        week_trend: {
          change: this.getNumber('inventory.bins.week_trend.change'),
          direction: this.getString('inventory.bins.week_trend.direction') as TrendDirection,
        },
        month_trend: {
          change: this.getNumber('inventory.bins.month_trend.change'),
          direction: this.getString('inventory.bins.month_trend.direction') as TrendDirection,
        },
        annual_trend: {
          change: this.getNumber('inventory.bins.annual_trend.change'),
          direction: this.getString('inventory.bins.annual_trend.direction') as TrendDirection,
        },
      },
      stock: {
        total: this.getNumber('inventory.stock.total'),
        week_trend: {
          change: this.getNumber('inventory.stock.week_trend.change'),
          direction: this.getString('inventory.stock.week_trend.direction') as TrendDirection,
        },
        month_trend: {
          change: this.getNumber('inventory.stock.month_trend.change'),
          direction: this.getString('inventory.stock.month_trend.direction') as TrendDirection,
        },
        annual_trend: {
          change: this.getNumber('inventory.stock.annual_trend.change'),
          direction: this.getString('inventory.stock.annual_trend.direction') as TrendDirection,
        },
      },

      products: {
        total: this.getNumber('inventory.products.total'),
        week_trend: {
          change: this.getNumber('inventory.products.week_trend.change'),
          direction: this.getString('inventory.products.week_trend.direction') as TrendDirection,
        },
        month_trend: {
          change: this.getNumber('inventory.products.month_trend.change'),
          direction: this.getString('inventory.products.month_trend.direction') as TrendDirection,
        },
        annual_trend: {
          change: this.getNumber('inventory.products.annual_trend.change'),
          direction: this.getString('inventory.products.annual_trend.direction') as TrendDirection,
        },
      },

      productVariants: {
        total: this.getNumber('inventory.productVariants.total'),
        week_trend: {
          change: this.getNumber('inventory.productVariants.week_trend.change'),
          direction: this.getString('inventory.productVariants.week_trend.direction') as TrendDirection,
        },
        month_trend: {
          change: this.getNumber('inventory.productVariants.month_trend.change'),
          direction: this.getString('inventory.productVariants.month_trend.direction') as TrendDirection,
        },
        annual_trend: {
          change: this.getNumber('inventory.productVariants.annual_trend.change'),
          direction: this.getString('inventory.productVariants.annual_trend.direction') as TrendDirection,
        },
      },
    }
  }
}

export type TrendDirection = 'up' | 'down' | 'stable'

export interface Trend {
  change: number
  direction: TrendDirection
}

export interface Metric {
  total: number
  week_trend: Trend
  month_trend: Trend
  annual_trend: Trend
}
