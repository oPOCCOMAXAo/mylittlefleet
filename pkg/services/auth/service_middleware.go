package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/mylittlefleet/pkg/models"
)

func (s *Service) MiddlewareAuth(ctx *gin.Context) {
	if !s.hasUsers {
		ctx.Redirect(http.StatusFound, s.setupPath)

		return
	}

	res, err := ctx.Cookie(Cookie)
	if err != nil || res == "" {
		ctx.Redirect(http.StatusFound, s.authPath)

		return
	}

	claims, err := s.token.Validate(res)
	if err != nil {
		ctx.Redirect(http.StatusFound, s.authPath)

		return
	}

	uid, err := claims.GetIntSubject()
	if err != nil {
		ctx.Redirect(http.StatusFound, s.authPath)

		return
	}

	CtxValue.Set(ctx, claims)
	CtxUserID.Set(ctx, uid)
}

func (s *Service) SetupSession(
	ctx *gin.Context,
	user *models.User,
) error {
	token, _, err := s.token.SignNew(
		WithIntEntityID(user.ID),
		WithAudience(AudAuth),
		WithMaxAge(TokenMaxAge),
	)
	if err != nil {
		return err
	}

	ctx.SetCookie(Cookie, token, TokenMaxAge, "/", "", false, true)

	return nil
}

func (s *Service) ClearSession(
	ctx *gin.Context,
) error {
	ctx.SetCookie(Cookie, "", -1, "/", "", false, true)

	return nil
}
