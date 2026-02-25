package stats

import "server/internal/core/shared/types"

type Duration string

const (
	DurationWeek  Duration = "week"
	DurationMonth          = "month"
	DurationYear           = "year"
)

type (
	inventoryStats struct {
		Stock struct {
			Total     uint64 `json:"total"`
			WeekTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"week_trend"`
			MonthTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"month_trend"`
			AnnualTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"annual_trend"`
		} `json:"stock"`
		Products struct {
			Total     uint64 `json:"total"`
			WeekTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"week_trend"`
			MonthTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"month_trend"`
			AnnualTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"annual_trend"`
		} `json:"products"`
		Inventory struct {
			Total     uint64 `json:"total"`
			WeekTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"week_trend"`
			MonthTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"month_trend"`
			AnnualTrend struct {
				Change    float64 `json:"change"`
				Direction string  `json:"direction"`
			} `json:"annual_trend"`
		} `json:"inventory"`
	}
)

type (
	GetData struct {
		UserID     string `json:"-"`
		LocationID string `json:"-"`
	}
	GetResult struct {
		Code    types.ServiceResultCode
		Payload GetResultPayload
	}

	GetResultPayload struct {
		Stats struct {
			Inventory inventoryStats `json:"inventory"`
		} `json:"stats"`
	}
)
