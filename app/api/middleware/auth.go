package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"p2-individual-project/util/response"
)

// JWTAuth middleware buat ngejaga endpoint yang butuh login.
// dia baca token dari header Authorization, validasi, lalu titipin user_id ke context biar controller bisa ambil.
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ambil header Authorization, formatnya "Bearer <token>"
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, "token tidak ada", nil)
			}

			// pisahin "Bearer" dari token-nya
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, http.StatusUnauthorized, "format token salah", nil)
			}
			tokenString := parts[1]

			// parse dan validasi token pakai secret yang sama waktu bikin
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				return response.Error(c, http.StatusUnauthorized, "token tidak valid", nil)
			}

			// ambil isi token, simpan user_id ke context
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return response.Error(c, http.StatusUnauthorized, "token tidak valid", nil)
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				return response.Error(c, http.StatusUnauthorized, "token tidak valid", nil)
			}

			// titipin user_id ke context, controller ambil pakai c.Get("user_id")
			c.Set("user_id", userID)

			return next(c)
		}
	}
}
