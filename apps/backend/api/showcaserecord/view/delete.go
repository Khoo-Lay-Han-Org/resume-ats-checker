package showcaserecord_view

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	typing "resuming/api/showcaserecord/typing"
	util "resuming/api/showcaserecord/util"
	validator "resuming/api/showcaserecord/validator"
)

func DeleteShowCaseRecordData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.SpecificPortoflioDataRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		polished_request, err := validator.ValidateSpecificPortfolioDataRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		index, err := strconv.Atoi(polished_request.Index)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid index format."})
		}

		err = util.DeleteShowCaseRecordData(polished_request.SectionTitle, index, public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
		}

		return nil
	}
}
