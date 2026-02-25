package stats

import (
	"math"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database"
	"server/internal/services/providers/database/queries"
	"strings"
	"time"
)

type Service struct {
	ctx types.IContext
	tmp struct {
		admin *entities.User
	}
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func (s *Service) Get(data GetData) (*GetResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	inventory, err := s.getInventoryStats(data, query)
	if err != nil {
		return nil, err
	}
	payload := GetResultPayload{}
	if inventory != nil {
		payload.Stats.Inventory = *inventory
	}
	return &GetResult{
		Code:    types.ServiceResultCodeSuccess,
		Payload: payload,
	}, nil
}

func (s *Service) getInventoryStats(data GetData, query *database.Query) (*inventoryStats, error) {
	now := time.Now().UTC()
	calc := func(current, previous uint64) (string, float64) {
		if previous == 0 {
			if current > 0 {
				return "up", 100.0
			}
			return "stable", 0.0
		}
		diff := float64(current) - float64(previous)
		percent := (diff / float64(previous)) * 100
		switch {
		case diff > 0:
			return "up", percent
		case diff < 0:
			return "down", percent
		default:
			return "stable", 0
		}
	}
	weekCurrent := now.AddDate(0, 0, -7)
	weekPrevious := now.AddDate(0, 0, -14)
	monthCurrent := now.AddDate(0, -1, 0)
	monthPrevious := now.AddDate(0, -2, 0)
	yearCurrent := now.AddDate(-1, 0, 0)
	yearPrevious := now.AddDate(-2, 0, 0)

	locationID := strings.TrimSpace(data.LocationID)
	var locationIDPtr *string
	if locationID != "" {
		locationIDPtr = &locationID
	}

	type totals struct {
		products  uint64
		inventory uint64
		stock     uint64
	}

	countProducts := func(from, to *time.Time) (uint64, error) {
		return query.Products().Count(queries.ProductParams{
			LocationID:    locationIDPtr,
			FromUpdatedAt: from,
			ToUpdatedAt:   to,
		})
	}

	sumInventory := func(from, to *time.Time) (totals, error) {
		rows, err := query.Inventory().GenerateSummary(queries.InventoryParams{
			LocationID:    locationIDPtr,
			FromUpdatedAt: from,
			ToUpdatedAt:   to,
		})
		if err != nil {
			return totals{}, err
		}
		var quantityOnHand uint64
		for _, row := range rows {
			quantityOnHand += row.QuantityOnHand
		}
		totalVariants := uint64(len(rows))
		// Requirement: stock = total_variants * quantity_on_hand
		stock := uint64(0)
		if totalVariants > 0 && quantityOnHand > 0 {
			if quantityOnHand > (math.MaxUint64 / totalVariants) {
				stock = math.MaxUint64
			} else {
				stock = totalVariants * quantityOnHand
			}
		}
		productCount, err := countProducts(from, to)
		if err != nil {
			return totals{}, err
		}
		return totals{
			products:  productCount,
			inventory: totalVariants,
			stock:     stock,
		}, nil
	}

	currentTotals, err := sumInventory(nil, nil)
	if err != nil {
		return nil, err
	}

	currentWeekTotals, err := sumInventory(&weekCurrent, &now)
	if err != nil {
		return nil, err
	}
	previousWeekTotals, err := sumInventory(&weekPrevious, &weekCurrent)
	if err != nil {
		return nil, err
	}

	currentMonthTotals, err := sumInventory(&monthCurrent, &now)
	if err != nil {
		return nil, err
	}
	previousMonthTotals, err := sumInventory(&monthPrevious, &monthCurrent)
	if err != nil {
		return nil, err
	}

	currentYearTotals, err := sumInventory(&yearCurrent, &now)
	if err != nil {
		return nil, err
	}
	previousYearTotals, err := sumInventory(&yearPrevious, &yearCurrent)
	if err != nil {
		return nil, err
	}

	stats := inventoryStats{}

	stats.Stock.Total = currentTotals.stock
	stats.Stock.WeekTrend.Direction, stats.Stock.WeekTrend.Change = calc(currentWeekTotals.stock, previousWeekTotals.stock)
	stats.Stock.MonthTrend.Direction, stats.Stock.MonthTrend.Change = calc(currentMonthTotals.stock, previousMonthTotals.stock)
	stats.Stock.AnnualTrend.Direction, stats.Stock.AnnualTrend.Change = calc(currentYearTotals.stock, previousYearTotals.stock)

	stats.Products.Total = currentTotals.products
	stats.Products.WeekTrend.Direction, stats.Products.WeekTrend.Change = calc(currentWeekTotals.products, previousWeekTotals.products)
	stats.Products.MonthTrend.Direction, stats.Products.MonthTrend.Change = calc(currentMonthTotals.products, previousMonthTotals.products)
	stats.Products.AnnualTrend.Direction, stats.Products.AnnualTrend.Change = calc(currentYearTotals.products, previousYearTotals.products)

	stats.Inventory.Total = currentTotals.inventory
	stats.Inventory.WeekTrend.Direction, stats.Inventory.WeekTrend.Change = calc(currentWeekTotals.inventory, previousWeekTotals.inventory)
	stats.Inventory.MonthTrend.Direction, stats.Inventory.MonthTrend.Change = calc(currentMonthTotals.inventory, previousMonthTotals.inventory)
	stats.Inventory.AnnualTrend.Direction, stats.Inventory.AnnualTrend.Change = calc(currentYearTotals.inventory, previousYearTotals.inventory)

	return &stats, nil
}
