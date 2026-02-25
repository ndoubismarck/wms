package types

type (
	Event         map[string][]EventListener
	EventListener func(data ...any)
)

type IEvents interface {
	On(name string, listener ...EventListener)
	Emit(name string, data ...any)
}

const (
	EventTypeDatabaseMigrationComplete = "event.database.migration.complete"
)
