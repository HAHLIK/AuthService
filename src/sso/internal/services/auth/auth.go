package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/HAHLIK/AuthService/sso/internal/domain/models"
	"github.com/HAHLIK/AuthService/sso/internal/lib/jwt"
	"github.com/HAHLIK/AuthService/sso/internal/lib/logger"
	"github.com/HAHLIK/AuthService/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

var ErrInvalidCredentails = errors.New("invalid credentails")

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int64) (token string, err error) {
	const op string = "Auth.Login"

	defer func() {
		if err != nil {
			err = fmt.Errorf("%s : %w", op, err)
		}
	}()

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("app_id", appId),
	)

	log.Info("Attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", logger.Err(err))

			return "", ErrInvalidCredentails
		}

		a.log.Error("failed to get user", logger.Err(err))

		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid", logger.Err(err))

		return "", ErrInvalidCredentails
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		return "", err
	}

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")

		return "", err
	}

	log.Info("user logged in successfully")

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (id int64, err error) {
	const op string = "Auth.RegisterNewUser"

	defer func() {
		if err != nil {
			err = fmt.Errorf("%s : %w", op, err)
		}
	}()

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed generate password hash", logger.Err(err))

		return -1, err
	}

	id, err = a.userSaver.SaveUser(ctx, email, passwordHash)
	if err != nil {
		log.Error("failed to save user", logger.Err(err))

		return -1, err
	}

	log.Info("User registred")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op string = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userId),
	)

	log.Info("Checking user admin status")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("can't isAdmin : %e", err)
	}

	log.Info("User admin status checked", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
