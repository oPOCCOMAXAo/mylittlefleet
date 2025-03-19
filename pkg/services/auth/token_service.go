package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/pkg/errors"
)

type TokenService struct {
	issuer string
	key    []byte
	method jwt.SigningMethod
	now    func() time.Time
}

func NewTokenService(
	issuer string,
	key []byte,
) *TokenService {
	return &TokenService{
		issuer: issuer,
		key:    key,
		method: jwt.SigningMethodHS256,
		now:    time.Now,
	}
}

func (s *TokenService) SetNowFunc(now func() time.Time) {
	if now == nil {
		return
	}

	s.now = now
}

func (s *TokenService) resetClaims(claims *TokenClaims) {
	claims.Issuer = s.issuer
	claims.IssuedAt = s.now().Unix()
	claims.ExpireAt = nil
}

func (s *TokenService) SignNew(
	options ...TokenClaimsOption,
) (string, *TokenClaims, error) {
	return s.SignClaims(&TokenClaims{}, options...)
}

func (s *TokenService) SignClaims(
	claims *TokenClaims,
	options ...TokenClaimsOption,
) (string, *TokenClaims, error) {
	claims = claims.Clone()
	s.resetClaims(claims)

	for _, option := range options {
		option(claims)
	}

	if claims.Subject == "" {
		return "", nil, errors.Wrapf(models.ErrInvalidAuth, "empty subject")
	}

	if len(claims.Audience) == 0 {
		return "", nil, errors.Wrapf(models.ErrInvalidAuth, "empty audience")
	}

	claims.Audience.Fix()

	token, err := jwt.
		NewWithClaims(s.method, claims).
		SignedString(s.key)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return token, claims, nil
}

func (s *TokenService) Validate(
	token string,
) (*TokenClaims, error) {
	var claims TokenClaims

	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		if claims.Issuer != s.issuer {
			return nil, errors.WithMessage(models.ErrInvalidAuth, "invalid issuer")
		}

		if claims.ExpireAt != nil && *claims.ExpireAt < s.now().Unix() {
			return nil, errors.WithMessage(models.ErrInvalidAuth, "expired")
		}

		return s.key, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &claims, nil
}
