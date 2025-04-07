package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/service"
	"github.com/barcek2281/adv2/auth/internal/store"
	"github.com/barcek2281/adv2/auth/models"
	"github.com/barcek2281/adv2/auth/utils"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	db           *store.MongoDB
	config       *config.Config
	adminService *service.AdminService
}

func NewAdminHandler(db *store.MongoDB, config *config.Config) *AdminHandler {
	return &AdminHandler{
		db:           db,
		config:       config,
		adminService: service.NewAdminrService(db, config),
	}
}

func (u *AdminHandler) Register() http.HandlerFunc {
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
		admin := models.Admin{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		msg, cookie, err := u.adminService.Register(&admin)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
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

func (u *AdminHandler) Login() http.HandlerFunc {
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

		admin, err := u.db.AdminRepo().Login(req.Email)
		if err != nil {
			utils.Response(w, r, http.StatusBadRequest, Res{Msg: "incorrect email or password"})
			log.Printf("login, no email: %v", err)
			return
		}

		if !admin.IsCorrectPassword(req.Password) {
			utils.Response(w, r, http.StatusBadRequest, Res{Msg: "incorrect email or password"})
			log.Printf("login, incorrect password")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"_id":  admin.ID,
				"role": admin.Role(),
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
