package showcase_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	crud "resuming/controller/showcase/crud"
	dto "resuming/controller/showcase/dto"
	validator "resuming/controller/showcase/validator"
)

func AddShowCaseRecordData() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		flag := c.Param("type-of-data")

		switch flag {
		case "name":
			var request dto.NameSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.NameSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "email":
			var request dto.EmailSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.EmailSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "phone-number":
			var request dto.PhoneNumberSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.PhoneNumberSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "address":
			var request dto.AddressSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.AddressSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "social-media":
			var request dto.SocialMediaSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.SocialMediaSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "job-experience":
			var request dto.JobExperienceSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.JobExperienceSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "education":
			var request dto.EducationSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.EducationSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "skill":
			var request dto.SkillSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.SkillSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "certificate":
			var request dto.CertificateSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.CertificateSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "language":
			var request dto.LanguageSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.LanguageSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		case "project":
			var request dto.ProjectSection
			if err := c.Bind(&request); err != nil {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
			}

			validated_request, err := crud.ValidateData[dto.ProjectSection](request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
			}

			err = crud.InsertShowCaseRecordData(validated_request, public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Operation failed."})
			}

			return nil
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to determine category of request."})
		}
	}
}
