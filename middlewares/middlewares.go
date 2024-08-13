package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
)

func IsLoggedIn(tokenManager token.TokenManager) fiber.Handler {
	return func(ftx *fiber.Ctx) error {
		cookie := ftx.Cookies("paseto")
		if cookie == "" {
			log.Println(cn.ErrsFitFin.CookiePasetoName)
			return ftx.SendStatus(fiber.StatusUnauthorized)
		}
	
		_, err := tokenManager.VerifyToken(cookie)
		if err != nil {
			log.Println(cn.ErrsFitFin.CookiePasetoValue)
			return ftx.SendStatus(fiber.StatusUnauthorized)
		}
	
		return ftx.Next()
	}
}

