package auth

import "github.com/opoccomaxao/mylittlefleet/pkg/utils/ginutils"

const (
	CtxValue  ginutils.CtxValuePointer[TokenClaims] = "__auth"
	CtxUserID ginutils.CtxValue[int64]              = "__user_id"

	Cookie string = "auth"

	CookieMaxAge = 60 * 60 * 24 * 7
	TokenMaxAge  = 3600

	AudAuth Audience = "auth"
)
