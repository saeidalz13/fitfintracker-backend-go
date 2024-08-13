package handlers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	sqlc "github.com/saeidalz13/LifeStyle2/lifeStyleBack/db/sqlc"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
)

func extractEmailFromClaim(cookie string, tokenManager token.TokenManager) (string, error) {
	if cookie == "" {
		return "", fmt.Errorf(cn.ErrsFitFin.CookiePasetoName)
	}

	payload, err := tokenManager.VerifyToken(cookie)
	if err != nil {
		return "", fmt.Errorf(cn.ErrsFitFin.CookiePasetoValue)
	}
	return payload.Email, nil
}

func isContentTypeJson(ftx *fiber.Ctx) bool {
	contentType := ftx.Get("Content-Type")
	return contentType == "application/json"
}

func fetchUserFromToken(cookie string, ctx context.Context, q *sqlc.Queries, tokenManager token.TokenManager) (*sqlc.User, error) {
	userEmail, err := extractEmailFromClaim(cookie, tokenManager)
	if err != nil {
		return nil, err
	}
	user, err := q.SelectUser(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func validateEmail(email string) error {
	regex, err := regexp.Compile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if err != nil {
		return err
	}
	if !regex.MatchString(email) {
		return fmt.Errorf("user-provided email does NOT have proper email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < cn.PASSWORD_MIN_LEN {
		return fmt.Errorf("password must be a minimum of %d characters", cn.PASSWORD_MIN_LEN)
	}

	uppercaseRegex, err := regexp.Compile(`[A-Z]`)
	if err != nil {
		return err
	}

	digitRegex, err := regexp.Compile(`\d`)
	if err != nil {
		return err
	}

	if !uppercaseRegex.MatchString(password) || !digitRegex.MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter and one digit")
	}

	return nil
}
