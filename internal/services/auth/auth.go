package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/iriskin77/grpc_go/internal/domain/models"
	"github.com/iriskin77/grpc_go/internal/jwt"
	"github.com/iriskin77/grpc_go/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Auth struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
	log         *slog.Logger
}

// New returns a new instance of the Auth service
func New(userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration, logger *slog.Logger) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
		log:         logger,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {

	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attemping to login user")

	user, err := a.usrProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err.Error())
			return "", fmt.Errorf(op, err.Error())
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Invalid credentials", err.Error())
		return "", fmt.Errorf(op, err.Error())
	}

	app, err := a.appProvider.App(ctx, appID)

	if err != nil {
		return "", fmt.Errorf(op, err.Error())
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf(op, err.Error())
	}

	return token, nil

}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {

	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error("failed to save user", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)

	}

	return id, nil

}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf(op, err.Error())
	}

	return isAdmin, nil
}
