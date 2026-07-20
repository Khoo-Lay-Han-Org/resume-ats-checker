package client_support_view

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	typing "resuming/api/client-support/typing"
	validator "resuming/api/client-support/validator"
	"resuming/tool"
)

func ClientReportOtherClient() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_reporting_public_user_id := c.Get("public_user_id")
		if retrieved_reporting_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve session id."})
		}

		reporting_public_user_id := retrieved_reporting_public_user_id.(string)

		var request typing.ClientReportRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to retrieve request."})
		}

		polished_request, err := validator.ValidateClientReportRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		report_type := polished_request.ReportType
		target_public_user_id := polished_request.TargetClientPublicUserId

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_report_logs").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to find user."})
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user data."})
		}

		new_report_log := map[string]any{
			"public_id":                uuid.New(),
			"reporting_public_user_id": reporting_public_user_id,
			"target_public_user_id":    target_public_user_id,
			"type":                     report_type,
		}

		data = append(data, new_report_log)

		serialised_data, err := json.Marshal(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to serialise new report log"})
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_report_logs").
				Value(string(serialised_data)).
				Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store client report logs"})
		}

		return nil
	}
}

func AIReport() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
