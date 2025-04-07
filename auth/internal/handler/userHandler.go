package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/service"
	"github.com/barcek2281/adv2/auth/internal/store"
	"github.com/barcek2281/adv2/auth/models"
	"github.com/barcek2281/adv2/auth/utils"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	db          *store.MongoDB
	config      *config.Config
	userService *service.UserService
}

func NewUserHandler(db *store.MongoDB, config *config.Config) *UserHandler {
	return &UserHandler{
		db:          db,
		config:      config,
		userService: service.NewUserService(db, config),
	}
}

func (u *UserHandler) RegisterUser() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type Response struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("register, cannot read from json: %v", err)
			return
		}
		user := models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		msg, cookie, err := u.userService.Register(&user)
		if err != nil {

			utils.Response(w, r, http.StatusInternalServerError, Response{Msg: "server error"})

			return
		}

		http.SetCookie(w, cookie)
		res := Response{
			Msg: msg,
		}
		utils.Response(w, r, http.StatusCreated, res)
		log.Printf("handle users/register")
	}
}

func (u *UserHandler) Login() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type Res struct {
		Msg string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			log.Printf("register, cannot read from json: %v", err)
			return
		}

		user, err := u.db.UserRepo().Login(req.Email)
		if err != nil {
			utils.Response(w, r, http.StatusBadRequest, Res{Msg: "incorrect email or password"})
			log.Printf("login, no email: %v", err)
			return
		}

		if !user.IsCorrectPassword(req.Password) {
			utils.Response(w, r, http.StatusBadRequest, Res{Msg: "incorrect email or password"})
			log.Printf("login, incorrect password")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"_id":  user.ID,
				"role": user.Role(),
				"exp":  time.Now().Add(time.Hour * 24).Unix(),
			})

		tokenString, err := token.SignedString([]byte(u.config.SECRET))
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			log.Printf("cannot save cookie: %v", err)
			return
		}

		cookie := http.Cookie{Name: "token", Value: tokenString}

		http.SetCookie(w, &cookie)

		utils.Response(w, r, http.StatusAccepted, Res{Msg: "login successfully"})
		log.Printf("handle /users/login")
	}
}
