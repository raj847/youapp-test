package service

import (
	"context"
	"errors"
	"net"
	"net/mail"
	"strings"
	"unicode"
	"youapp/entity"
	"youapp/repository"
	"youapp/utils"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserPasswordDontMatch = errors.New("password not match")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserInvalid           = errors.New("username invalid")
	ErrPasswordInvalid       = errors.New("password invalid")
	ErrEmailInvalid          = errors.New("email invalid")
)

func (s *UserService) AddAnggota(ctx context.Context, userReq entity.UserRegister) (entity.User, error) {
	existingUser, err := s.userRepository.GetUserByUsername(ctx, userReq.Username)
	if err != nil {
		return entity.User{}, err
	}

	if existingUser.ID != 0 {
		return entity.User{}, ErrUserAlreadyExists
	}

	if len(userReq.Username) < 3 || len(userReq.Username) > 11 {
		return entity.User{}, ErrUserInvalid
	}

	for _, v := range userReq.Username {
		if string(v) == strings.ToUpper(string(v)) {
			return entity.User{}, ErrUserInvalid
		}
	}

	_, err = mail.ParseAddress(userReq.Email)
	if err != nil {
		return entity.User{}, ErrEmailInvalid
	}

	domain := strings.Split(userReq.Email, "@")
	_, err = net.LookupMX(domain[1])
	if err != nil {
		return entity.User{}, ErrEmailInvalid
	}

	validPassword := validatePassword(userReq.Password)
	if !validPassword {
		return entity.User{}, ErrPasswordInvalid
	}

	anggota := entity.User{
		Username: userReq.Username,
		Password: userReq.Password,
		Email:    userReq.Email,
	}

	hashedPassword, err := utils.HashPassword(anggota.Password)
	if err != nil {
		return entity.User{}, err
	}

	anggota.Password = hashedPassword

	newUser, err := s.userRepository.Register(ctx, anggota)
	if err != nil {
		return entity.User{}, err
	}
	return newUser, nil
}

func validatePassword(password string) bool {
	var lower, upper, symbol bool
	moreThan := len(password) > 8

	for _, char := range password {
		if unicode.IsLower(char) {
			lower = true
			continue
		}

		if unicode.IsUpper(char) {
			upper = true
			continue
		}

		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			symbol = true
			continue
		}
	}

	return moreThan && lower && upper && symbol
}

func (s *UserService) LoginUser(ctx context.Context, cashierReq entity.UserLogin) (user entity.User, err error) {
	existingUser, err := s.userRepository.GetUserByUsername(ctx, cashierReq.Username)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}

	if utils.CheckPassword(cashierReq.Password, existingUser.Password) != nil {
		return entity.User{}, ErrUserPasswordDontMatch
	}

	return existingUser, nil
}
