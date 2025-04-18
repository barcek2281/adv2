package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/store"
	"github.com/barcek2281/adv2/auth/models"
	"github.com/golang-jwt/jwt"
)

var (
	ErrIncorrectData = errors.New("incorrect data")
	ErrEncrypt       = errors.New("cannot create a password")
	ErrEmailUsed     = errors.New("email is used")
	ErrRequestEmail  = errors.New("cannot send email")
	ErrJWT           = errors.New("cannot create JWT token")
)

type UserService struct {
	db     *store.MongoDB
	config *config.Config
}

func NewUserService(db *store.MongoDB, config *config.Config) *UserService {
	return &UserService{
		db:     db,
		config: config,
	}
}

func (us *UserService) Register(user *models.User) (string, *http.Cookie, error) {
	// validation part
	if !user.IsValid() {
		log.Printf("incorrect user property")
		return "", nil, ErrIncorrectData
	}
	if err := user.CryptPassword(); err != nil {
		log.Printf("cannot encrypt user: %v", err)
		return "", nil, ErrEncrypt
	}

	user.RegisterAt = time.Now()

	// create part
	if err := us.db.UserRepo().Register(user); err != nil {
		log.Printf("cannot save user into db: %v", err)
		return "", nil, err
	}

	// sending email part
	type Req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	req := Req{
		Username: user.Username,
		Email:    user.Email,
	}
	resB2, _ := json.Marshal(&req)
	resp, err := http.Post(us.config.MailURI+"/mail/register", "application/json", bytes.NewBuffer(resB2))

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("%v", err)
		return "", nil, ErrRequestEmail
	}

	// jwt token part
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"_id":  user.ID,
			"role": user.Role(),
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(us.config.SECRET))
	if err != nil {
		log.Printf("cannot save cookie: %v", err)
		return "", nil, ErrJWT
	}

	cookie := http.Cookie{Name: "token", Value: tokenString}

	return "user register successfully", &cookie, nil
}
