package showcaserecord_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	util "resuming/api/showcaserecord/util"
	"resuming/database"
	"resuming/tool"
)

func RetrieveShowCaseRecordData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
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
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
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
			"name":           util.ToStringSlice(data["name"]),
			"email":          util.ToStringSlice(data["email"]),
			"phone_number":   util.ToStringSlice(data["phone_number"]),
			"address":        util.ToStringSlice(data["address"]),
			"social_media":   util.ToStringSlice(data["social_media"]),
			"job_experience": util.ToJSON(data["job_experience"]),
			"education":      util.ToJSON(data["education"]),
			"skill":          util.ToStringSlice(data["skill"]),
			"certificate":    util.ToJSON(data["certificate"]),
			"language":       util.ToStringSlice(data["language"]),
			"project":        util.ToJSON(data["project"]),
		}
		c.Set("response_data", response_data)
		return nil
	}
}
