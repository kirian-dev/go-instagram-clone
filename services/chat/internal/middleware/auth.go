package middleware

import (
	"errors"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			mw.log.Infof("header auth %s", bearerHeader)
			if bearerHeader == "" {
				return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrAuthorizationHeaderNotSet})
			}
			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.log.Error("Auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrUnauthorized})
				}

				token := headerParts[1]

				claims, err := utils.ValidateJWTToken(token, mw.cfg)
				if err != nil {
					if err == jwt.ErrTokenExpired {
						mw.log.Error(e.ErrTokenExpired)
						return c.JSON(http.StatusUnauthorized, errors.New(e.ErrTokenExpired))
					}
					mw.log.Error("validateJWTToken", zap.String("JWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrUnauthorized})
				}
				c.Set("userClaims", claims)
				return next(c)
			}
			return nil
		}
	}
}

func (mw *MiddlewareManager) AdminAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("userClaims").(*utils.CustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrUnauthorized})
			}

			// Check if the user has the admin role
			if userClaims.Role != "admin" {
				return c.JSON(http.StatusForbidden, e.ErrorResponse{Error: e.ErrForbidden})
			}

			return next(c)
		}
	}
}
