package showcase_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	showcaserecord_convert "resuming/backend-api/showcase/convert"
	"resuming/database"
	"resuming/service"
)

func RetrieveShowCaseRecordData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request().Context()
		retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve showcase record data."})
				}
				showcase, scErr := database.Queries.FindShowcaseRecordByUserId(ctx, user.ID)
				if scErr != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve showcase record data."})
				}
				if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, &showcase); syncErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve showcase record data."})
				}
				retrieved_data, err = service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
				if err != nil {
					return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve showcase record data."})
				}
			} else {
				return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve showcase record data."})
			}
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process showcase record data."})
		}

		response_data := echo.Map{
			"name":           showcaserecord_convert.ToStringSlice(data["name"]),
			"email":          showcaserecord_convert.ToStringSlice(data["email"]),
			"phone_number":   showcaserecord_convert.ToStringSlice(data["phone_number"]),
			"address":        showcaserecord_convert.ToStringSlice(data["address"]),
			"social_media":   showcaserecord_convert.ToStringSlice(data["social_media"]),
			"job_experience": showcaserecord_convert.ToJSON(data["job_experience"]),
			"education":      showcaserecord_convert.ToJSON(data["education"]),
			"skill":          showcaserecord_convert.ToStringSlice(data["skill"]),
			"certificate":    showcaserecord_convert.ToJSON(data["certificate"]),
			"language":       showcaserecord_convert.ToStringSlice(data["language"]),
			"project":        showcaserecord_convert.ToJSON(data["project"]),
		}
		c.Set("response_data", response_data)
		return nil
	}
}
