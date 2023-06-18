package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"youapp/entity"
	"youapp/service"

	"github.com/dgrijalva/jwt-go"
)

type UserAPI struct {
	userService *service.UserService
}

func NewUserAPI(
	userService *service.UserService,
) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

func (u *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Username == "" || user.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("register data is empty"))
		return
	}

	eUser, err := u.userService.AddAnggota(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrEmailInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"message": "register success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (u *UserAPI) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userReq entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if userReq.Username == "" || userReq.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("username or password is empty"))
		return
	}

	eUser, err := u.userService.LoginUser(r.Context(), userReq)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserPasswordDontMatch) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	expiresAt := time.Now().Add(5 * time.Hour)
	claims := entity.Claims{
		UserID: eUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	response := map[string]any{
		"data":        eUser,
		"tokenCookie": tokenString,
	}

	WriteJSON(w, http.StatusOK, response)
}
