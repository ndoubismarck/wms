package dotenv

import (
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/fileutil"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type dotEnv struct {
	ctx  types.IContext
	data *types.DotEnv
}

func New(ctx types.IContext) types.IDotEnv {
	return &dotEnv{
		ctx: ctx,
	}
}

func (e *dotEnv) Get() (*types.DotEnv, error) {
	if e.data != nil {
		return e.data, nil
	}
	if fileutil.Exists(".env") {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}
	data := types.DotEnv{}
	if err := env.Parse(&data); err != nil {
		return nil, err
	}
	e.data = &data
	return e.data, nil
}
