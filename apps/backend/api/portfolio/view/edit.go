package portfolio_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/portfolio/typing"
	validator "resuming/api/portfolio/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ChooseTemplate() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChooseTemplateRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		template_id, err := validator.ValidateTemplateID(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve portfolio data."})
				}
				portfolio, dbErr := database.Queries.FindPortfolioByUserId(ctx, user.ID)
				if dbErr != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve portfolio data."})
				}
				if syncErr := database.SyncIndividualPortfolioDataSessionStore(public_user_id, &portfolio); syncErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve portfolio data."})
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
				if err != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve portfolio data."})
				}
			} else {
				return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve portfolio data."})
			}
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process portfolio data."})
		}

		data["template_id"] = template_id

		serialised_data, err := json.Marshal(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process portfolio data."})
		}

		err = tool.Valkey.Do(
			c.Request().Context(),
			tool.Valkey.B().Set().
				Key(public_user_id+":portfolio_data").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to save portfolio data."})
		}

		return nil
	}
}
