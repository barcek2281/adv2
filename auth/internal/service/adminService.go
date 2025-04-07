package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/barcek2281/adv2/auth/internal/config"
	"github.com/barcek2281/adv2/auth/internal/store"
	"github.com/barcek2281/adv2/auth/models"
	"github.com/golang-jwt/jwt"
)

type AdminService struct {
	db     *store.MongoDB
	config *config.Config
}

func NewAdminrService(db *store.MongoDB, config *config.Config) *AdminService {
	return &AdminService{
		db:     db,
		config: config,
	}
}

func (us *AdminService) Register(admin *models.Admin) (string, *http.Cookie, error) {
	// validation part
	if !admin.IsValid() {
		log.Printf("incorrect user property")
		return "", nil, fmt.Errorf("invalid")
	}
	if err := admin.CryptPassword(); err != nil {
		log.Printf("cannot encrypt user: %v", err)
		return "", nil, err
	}

	admin.RegisterAt = time.Now()

	// create part
	if err := us.db.AdminRepo().Register(admin); err != nil {
		log.Printf("cannot save user into db: %v", err)
		return "", nil, err
	}

	// sending email part
	type Req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	req := Req{
		Username: admin.Username,
		Email:    admin.Email,
	}
	resB2, _ := json.Marshal(&req)
	resp, err := http.Post(us.config.MailURI+"/mail/register", "application/json", bytes.NewBuffer(resB2))

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("%v", err)
		return "", nil, err
	}

	// jwt token part
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"_id":  admin.ID,
			"role": admin.Role(),
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(us.config.SECRET))
	if err != nil {
		log.Printf("cannot save cookie: %v", err)
		return "", nil, err
	}

	cookie := http.Cookie{Name: "token", Value: tokenString}

	return "admin register successfully", &cookie, nil
}
