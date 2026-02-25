package auth

import (
	"fmt"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/securityutil"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Service struct {
	ctx       types.IContext
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func (s *Service) Login(data LoginData) (*LoginResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &LoginResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &LoginResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	usersQuery := query.Users()
	result, exists, err := usersQuery.FindOneByEmailAddress(data.Username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &LoginResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	if !securityutil.CheckPasswordHash(data.Password, result.Password) {
		return &LoginResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	token, err := s.generateJWTAuthToken(result.ID, !data.Remember)
	if err != nil {
		return nil, err
	}
	loginAt := time.Now().UTC()
	if _, _, err = usersQuery.UpdateByID(result.ID, entities.User{
		LastLoginAt:        &loginAt,
		LastLoginIPAddress: &data.IPAddress,
	}); err != nil {
		return nil, err
	}
	result, exists, err = usersQuery.FindOneByID(result.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &LoginResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	return &LoginResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: &LoginResultPayload{
			User:  result,
			Token: token,
		},
	}, nil
}

func (s *Service) JWTAuth(data JWTAuthData) (*JWTAuthResult, error) {
	conf, err := s.ctx.Config().Get()
	if err != nil {
		return nil, err
	}
	key := []byte(conf.Security.Key)
	token, err := jwt.ParseWithClaims(data.Token, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return &JWTAuthResult{
			Code: types.ServiceResultCodeFailed,
		}, nil
	}
	claims, ok := token.Claims.(*jwtClaims)
	if ok && token.Valid {
		query, ok := s.providers.Database().Query()
		if !ok {
			return &JWTAuthResult{
				Code: types.ServiceResultCodeFailed,
			}, nil
		}
		usersQuery := query.Users()
		result, exists, err := usersQuery.FindOneByID(claims.ID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return &JWTAuthResult{
				Code: types.ServiceResultCodeFailed,
			}, nil
		}
		if claims.ExpiresAt != nil {
			if claims.ExpiresAt.UTC().Before(time.Now().UTC()) {
				return &JWTAuthResult{
					Code: types.ServiceResultCodeFailed,
				}, nil
			}
		}
		return &JWTAuthResult{
			Code: types.ServiceResultCodeSuccess,
			Payload: &JWTAuthResultPayload{
				Token: data.Token,
				User:  result,
			},
		}, nil
	}
	return &JWTAuthResult{
		Code: types.ServiceResultCodeFailed,
	}, nil
}

func (s *Service) generateJWTAuthToken(userId string, sessionToken bool) (string, error) {
	conf, err := s.ctx.Config().Get()
	if err != nil {
		return "", err
	}
	iat := time.Now().UTC()
	token := jwt.New(jwt.SigningMethodHS256)
	if sessionToken {
		token.Claims = &jwtClaims{
			ID: userId,
			RegisteredClaims: &jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(iat),
			},
		}
	} else {
		token.Claims = &jwtClaims{
			ID: userId,
			RegisteredClaims: &jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(iat),
				ExpiresAt: jwt.NewNumericDate(iat.Add(time.Hour * 24 * 30)),
			},
		}
	}
	tkn, err := token.SignedString([]byte(conf.Security.Key))
	if err != nil {
		return "", err
	}
	return tkn, nil
}
