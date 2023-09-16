package utils

import (
	"net/http"

	"go-instagram-clone/pkg/e"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

func AuthorizeUser(c echo.Context) (uuid.UUID, error) {
	userIDStr := c.Param("userId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, c.JSON(http.StatusBadRequest, err)
	}
	userClaims := c.Get("userClaims").(*CustomClaims)
	authenticatedUserID := userClaims.UserID

	if userID != authenticatedUserID {
		return uuid.UUID{}, c.JSON(http.StatusForbidden, e.ErrorResponse{Error: e.ErrForbidden})
	}
	return userID, nil
}
