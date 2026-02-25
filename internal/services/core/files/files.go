package files

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/files"
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
	s := &Service{
		ctx:       ctx,
		providers: providers,
	}
	s.ctx.Hooks().OnStarted(s.start)
	return s
}

func (s *Service) start() error {
	go func() {
		if err := s.providers.Files().Cleanup(); err != nil {
			s.ctx.Logger().Error(err)
		}
		for {
			time.Sleep(time.Minute * 1)
			if err := s.providers.Files().Cleanup(); err != nil {
				s.ctx.Logger().Error(err)
			}
		}
	}()
	return nil
}

func (s *Service) Upload(data UploadData) (*UploadResult, error) {
	validation, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UploadResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: validation,
		}, nil
	}
	result, err := s.providers.Files().Upload(files.UploadFileData{
		File: data.File,
	})
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UploadResultPayload{
			File: *result,
		},
	}, nil
}
