package middlewares

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
)

const secretKey = "123e4567-e89b-12d3-a456-42661417"

type Middleware struct {
	Config *config.ServerConfig
}

func (m *Middleware) SetCookie(c *gin.Context) {
	var userId string
	var err error

	authCookie, _ := c.Request.Cookie("user")
	if authCookie == nil {
		userId = uuid.NewString()
		newCookie, err := m.encrypt(userId)
		if err != nil {
			log.Println(err)
			return
		}

		authCookie = &http.Cookie{
			Name:     "user",
			Value:    newCookie,
			Path:     "/",
			Domain:   m.Config.Address,
			MaxAge:   3600,
			Secure:   false,
			HttpOnly: false,
		}
		c.Request.AddCookie(authCookie)
	} else {
		userId, err = m.decrypt([]byte(authCookie.Value))
		if err != nil {
			log.Println(err)
			return
		}
	}

	c.SetCookie(
		"user",
		authCookie.Value,
		3600,
		"/",
		m.Config.Address,
		false,
		false,
	)

	c.Set("userId", userId)

	c.Next()
}

func (m *Middleware) encrypt(value string) (string, error) {
	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, err := helpers.GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return "", err
	}

	// зашифровываем
	dst := aesgcm.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(dst), nil
}

func (m *Middleware) decrypt(value []byte) (string, error) {
	var decodedValue []byte
	decodedValue, _ = hex.DecodeString(string(value))

	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, cipherText := decodedValue[:aesgcm.NonceSize()], decodedValue[aesgcm.NonceSize():]

	// расшифровываем
	if err != nil {
		return "", err
	}
	userId, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(userId), nil
}
