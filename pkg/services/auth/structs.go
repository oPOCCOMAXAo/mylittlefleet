package auth

import (
	"slices"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Audience string

type AudienceList []Audience

func (l AudienceList) Contains(other Audience) bool {
	for _, item := range l {
		if item == other {
			return true
		}
	}

	return false
}

func (l AudienceList) ContainsAll(other AudienceList) bool {
	for _, a := range other {
		if !l.Contains(a) {
			return false
		}
	}

	return true
}

func (l *AudienceList) Fix() {
	set := map[Audience]struct{}{}
	for _, a := range *l {
		set[a] = struct{}{}
	}

	*l = make(AudienceList, 0, len(set))
	for a := range set {
		*l = append(*l, a)
	}

	slices.SortFunc(*l, func(l, r Audience) int {
		return strings.Compare(string(l), string(r))
	})
}

var _ jwt.Claims = (*TokenClaims)(nil)

type TokenClaims struct {
	Issuer   string       `json:"iss"`           // Issuer - current service.
	IssuedAt int64        `json:"iat"`           // IssuedAt - unix timestamp. Should be set by default in token service.
	Subject  string       `json:"sub"`           // Subject - entity ID/user ID.
	Audience AudienceList `json:"aud"`           // Audience - service that should accept this token.
	ExpireAt *int64       `json:"exp,omitempty"` // ExpireAt - unix timestamp.
}

func (*TokenClaims) Valid() error {
	return nil
}

func (claims *TokenClaims) Clone() *TokenClaims {
	res := *claims

	return &res
}

func (claims *TokenClaims) GetIntSubject() (int64, error) {
	res, err := strconv.ParseInt(claims.Subject, 10, 64)

	return res, errors.WithStack(err)
}

type TokenClaimsOption func(*TokenClaims)

func WithAudience(aud Audience) TokenClaimsOption {
	return func(claims *TokenClaims) {
		claims.Audience = append(claims.Audience, aud)
	}
}

func WithEntityID(entityID string) TokenClaimsOption {
	return func(claims *TokenClaims) {
		claims.Subject = entityID
	}
}

func WithIntEntityID(entityID int64) TokenClaimsOption {
	return func(claims *TokenClaims) {
		claims.Subject = strconv.FormatInt(entityID, 10)
	}
}

func WithMaxAge(maxAge int64) TokenClaimsOption {
	return func(claims *TokenClaims) {
		claims.ExpireAt = lo.ToPtr(claims.IssuedAt + maxAge)
	}
}

type UserData struct {
	UserID int64
}
