// Branch feature/login
package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fazaalexander/GenuineID/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func CreateToken(user *models.User, id uint, username string, email string, role string) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetClaims(token string) (*jwtCustomClaims, error) {
	token = strings.Replace(token, "Bearer ", "", 1)

	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"error": "Invalid token",
		})
	}

	claims := jwtToken.Claims.(*jwtCustomClaims)
	return claims, nil
}

// Verify token to ensure that user role is "customer"
func CustomerTokenVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")[7:]
		jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Unauthorized",
			})
		}

		claims, ok := jwtToken.Claims.(*jwtCustomClaims)
		log.Println(claims)
		if !ok || claims.Role != "customer" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "Forbidden Access",
			})
		}

		// lanjutkan ke handler/controller apabila token dan role user valid
		return next(c)
	}
}

// Verify token to ensure that user role is "seller"
func SellerTokenVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")[7:]
		jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Unauthorized",
			})
		}

		claims, ok := jwtToken.Claims.(*jwtCustomClaims)
		if !ok || claims.Role != "seller" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "Forbidden Access",
			})
		}

		// lanjutkan ke handler/controller apabila token dan role user valid
		return next(c)
	}
}

// Verify token to ensure that user role is "admin"
func AdminTokenVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")[7:]
		jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Unauthorized",
			})
		}

		claims, ok := jwtToken.Claims.(*jwtCustomClaims)
		if !ok || claims.Role != "admin" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "Forbidden Access",
			})
		}

		// lanjutkan ke handler/controller apabila token dan role user valid
		return next(c)
	}
}
