package portfolio_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/tool"
)

func RetrievePortfolioData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve session data."})
		}

		public_user_id := retrieved_public_user_id.(string)

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

		data_detail := data["detail"].(map[string]any)

		template_id := data["template_id"].(float64)
		name := data_detail["name"].([]any)
		email := data_detail["email"].([]any)
		phone_number := data_detail["phone_number"].([]any)
		address := data_detail["address"].([]any)
		social_media := data_detail["social_media"].([]any)
		job_experience := data_detail["job_experience"].([]any)
		education := data_detail["education"].([]any)
		skill := data_detail["skill"].([]any)
		certificate := data_detail["certificate"].([]any)
		language := data_detail["language"].([]any)
		project := data_detail["project"].([]any)

		response_data := echo.Map{
			"template_id":    template_id,
			"name":           name,
			"email":          email,
			"phone_number":   phone_number,
			"address":        address,
			"social_media":   social_media,
			"job_experience": job_experience,
			"education":      education,
			"skill":          skill,
			"certificate":    certificate,
			"language":       language,
			"project":        project,
		}
		c.Set("response_data", response_data)

		return nil
	}
}
