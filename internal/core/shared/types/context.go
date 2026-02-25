package types

type IContext interface {
	App() App

	Cache() ICache

	Hooks() IHooks

	Events() IEvents

	Config() IConfig

	Logger() ILogger

	Server() IServer

	DotEnv() IDotEnv

	Database() IDatabase
}
