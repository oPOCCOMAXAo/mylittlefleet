package auth

import (
	"context"
	"time"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/auth/repo"
	"github.com/opoccomaxao/mylittlefleet/pkg/utils/envutils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	config Config
	repo   *repo.Repo
	token  *TokenService

	// timestamp for tokens when full info should be reloaded instead of just extending expiration.
	refreshReloadBefore int64

	hasUsers bool

	authPath  string
	setupPath string
}

type Config struct {
	TokenKey    envutils.HexBytes `env:"TOKEN_KEY,required,notEmpty"`                            // hex encoded key.
	TokenIssuer string            `env:"TOKEN_ISSUER"                envDefault:"mylittlefleet"` // token issuer.
}

func NewService(
	config Config,
	repo *repo.Repo,
) *Service {
	return &Service{
		config: config,
		repo:   repo,
		token:  NewTokenService(config.TokenIssuer, config.TokenKey),

		authPath:  "/login",
		setupPath: "/setup",
	}
}

func (s *Service) OnStart(ctx context.Context) error {
	err := s.initHasUsers(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) initHasUsers(ctx context.Context) error {
	totalUsers, err := s.repo.GetTotalUsers(ctx)
	if err != nil {
		return err
	}

	s.hasUsers = totalUsers > 0

	return nil
}

// When some data is changed, we need to reload full info from the database.
// This method sets the timestamp when full info should be reloaded.
func (s *Service) MarkFullReloadAfterNow() {
	s.refreshReloadBefore = time.Now().Unix()
}

func (s *Service) RefreshAuth(
	ctx context.Context,
	claims *TokenClaims,
) (string, *TokenClaims, error) {
	if claims.IssuedAt >= s.refreshReloadBefore {
		return s.token.SignClaims(claims,
			WithMaxAge(TokenMaxAge),
		)
	}

	userID, err := claims.GetIntSubject()
	if err != nil {
		return "", nil, err
	}

	data, err := s.GetUserData(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	return s.token.SignNew(
		WithIntEntityID(data.UserID),
		WithAudience(AudAuth),
		WithMaxAge(TokenMaxAge),
	)
}

func (s *Service) GetUserData(
	ctx context.Context,
	userID int64,
) (*UserData, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserData{
		UserID: user.ID,
	}, nil
}

func (s *Service) HasUsers() bool {
	return s.hasUsers
}

func (s *Service) HashPassword(password string) (string, error) {
	const additionalCost = 4

	res, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost+additionalCost,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(res), nil
}

type CreateUserParams struct {
	Username string
	Password string
}

func (s *Service) CreateUser(
	ctx context.Context,
	params CreateUserParams,
) (*models.User, error) {
	var err error

	user := models.User{
		Login: params.Username,
	}

	user.Password, err = s.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	s.hasUsers = true

	return &user, nil
}

//nolint:revive
type AuthUserParams struct {
	Username string
	Password string
}

func (s *Service) AuthUser(
	ctx context.Context,
	params AuthUserParams,
) (*models.User, error) {
	user, err := s.repo.GetUserByLogin(ctx, params.Username)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(models.ErrInvalidAuth)
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return nil, errors.WithStack(models.ErrInvalidAuth)
	}

	return user, nil
}
