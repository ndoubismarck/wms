package services

import (
	"server/internal/core/shared/types"
	"server/internal/services/core/auth"
	"server/internal/services/core/files"
	"server/internal/services/core/inventory"
	"server/internal/services/core/locations"
	"server/internal/services/core/operations"
	"server/internal/services/core/products"
	"server/internal/services/core/setup"
	"server/internal/services/core/stats"
	"server/internal/services/core/tasks"
	"server/internal/services/core/teams"
	"server/internal/services/core/users"
	"server/internal/services/providers"
)

type Services struct {
	auth       *auth.Service
	setup      *setup.Service
	stats      *stats.Service
	files      *files.Service
	products   *products.Service
	locations  *locations.Service
	inventory  *inventory.Service
	users      *users.Service
	teams      *teams.Service
	tasks      *tasks.Service
	operations *operations.Service
}

func New(ctx types.IContext, providers *providers.Providers) *Services {
	return &Services{
		auth:       auth.New(ctx, providers),
		setup:      setup.New(ctx, providers),
		stats:      stats.New(ctx, providers),
		files:      files.New(ctx, providers),
		products:   products.New(ctx, providers),
		locations:  locations.New(ctx, providers),
		inventory:  inventory.New(ctx, providers),
		users:      users.New(ctx, providers),
		teams:      teams.New(ctx, providers),
		tasks:      tasks.New(ctx, providers),
		operations: operations.New(ctx, providers),
	}
}

func (s *Services) Auth() *auth.Service {
	return s.auth
}

func (s *Services) Setup() *setup.Service {
	return s.setup
}

func (s *Services) Stats() *stats.Service {
	return s.stats
}

func (s *Services) Files() *files.Service {
	return s.files
}

func (s *Services) Products() *products.Service {
	return s.products
}

func (s *Services) Locations() *locations.Service {
	return s.locations
}

func (s *Services) Inventory() *inventory.Service {
	return s.inventory
}

func (s *Services) Users() *users.Service {
	return s.users
}

func (s *Services) Teams() *teams.Service {
	return s.teams
}

func (s *Services) Tasks() *tasks.Service {
	return s.tasks
}

func (s *Services) Operations() *operations.Service {
	return s.operations
}
