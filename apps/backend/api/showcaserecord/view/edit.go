package showcaserecord_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	typing "resuming/api/showcaserecord/typing"
	util "resuming/api/showcaserecord/util"
	validator "resuming/api/showcaserecord/validator"
)

func EditShowCaseRecordData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		flag := c.Param("type-of-data")

		switch flag {
		case "name":
			var request typing.NameSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.NameSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "email":
			var request typing.EmailSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.EmailSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "phone-number":
			var request typing.PhoneNumberSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.PhoneNumberSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "address":
			var request typing.AddressSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.AddressSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "social-media":
			var request typing.SocialMediaSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.SocialMediaSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "job-experience":
			var request typing.JobExperienceSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.JobExperienceSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "education":
			var request typing.EducationSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.EducationSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "skill":
			var request typing.SkillSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.SkillSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "certificate":
			var request typing.CertificateSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.CertificateSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "language":
			var request typing.LanguageSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.LanguageSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "project":
			var request typing.ProjectSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := util.ValidateData[typing.ProjectSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			if validated_request.Index == nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": "Index is required."})
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to determine category of request."})
		}
	}
}
