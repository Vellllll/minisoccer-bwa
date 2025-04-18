package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	errConstant "user-service/constants/error"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.InterfaceRepositoryRegistry
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

type InterfaceUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

func NewUserService(repository repositories.InterfaceRepositoryRegistry) InterfaceUserService {
	return &UserService{repository: repository}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		logrus.Errorf("username is not found")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logrus.Errorf("password incorrect")
		return nil, err
	}

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) IsUsenameExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) IsEmailExist(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	if user != nil {
		return true
	}

	return false
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.IsUsenameExist(ctx, req.Username) {
		return nil, errConstant.ErrUsernameExist
	}

	if u.IsEmailExist(ctx, req.Email) {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	user, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:            req.Name,
		Username:        req.Username,
		Password:        string(hashedPassword),
		ConfirmPassword: req.ConfirmPassword,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		RoleID:          constants.Customer,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Username:    user.Username,
			Name:        user.Name,
			Email:       user.Email,
			Role:        user.Role.Code,
			PhoneNumber: user.PhoneNumber,
		},
	}

	return response, nil
}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      *dto.UserResponse
	)

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	isUsernameExist := u.IsUsenameExist(ctx, req.Username)
	if isUsernameExist && user.Username != req.Username {
		checkUsername, err = u.repository.GetUser().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}

		if checkUsername != nil {
			return nil, errConstant.ErrUsernameExist
		}
	}

	isEmailExist := u.IsEmailExist(ctx, req.Email)
	if isEmailExist && user.Email != req.Email {
		checkEmail, err = u.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errConstant.ErrEmailExist
		}
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	password = string(hashedPassword)

	userResult, err = u.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}, uuid)

	if err != nil {
		return nil, err
	}

	data = &dto.UserResponse{
		UUID:        user.UUID,
		Name:        userResult.Name,
		Username:    userResult.Username,
		PhoneNumber: userResult.PhoneNumber,
		Email:       userResult.Email,
	}

	return data, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      *dto.UserResponse
	)

	data = &dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		PhoneNumber: userLogin.Username,
		Email:       userLogin.Email,
		Role:        userLogin.Role,
	}

	return data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		logrus.Errorf("user not found")
		return nil, err
	}

	data := *&dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}

	return &data, nil
}
