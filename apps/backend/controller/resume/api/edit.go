package resume_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	dto "resuming/controller/resume/dto"
	resume_find "resuming/controller/resume/find"
	validator "resuming/controller/resume/validator"
	"resuming/service"
	"resuming/systemconfig"
)

func ChooseTemplate() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request dto.ChooseTemplateRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		template_id, err := validator.ValidateTemplateID(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		retrieved_data, err := resume_find.GetResumeData(public_user_id)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve resume data."})
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume data."})
		}

		data["template_id"] = template_id

		serialised_data, err := json.Marshal(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process resume data."})
		}

		err = service.Valkey.Do(
			c.Request().Context(),
			service.Valkey.B().Set().
				Key(public_user_id+":resume_data").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to save resume data."})
		}

		return nil
	}
}
