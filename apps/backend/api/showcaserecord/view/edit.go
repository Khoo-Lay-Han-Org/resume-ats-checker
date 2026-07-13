package showcaserecord_view

import (
	"net/http"

	"github.com/gin-gonic/gin"
	typing "resuming/api/showcaserecord/typing"
	util "resuming/api/showcaserecord/util"
	validator "resuming/api/showcaserecord/validator"
)

func EditShowCaseRecordData() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		flag := c.Param("type-of-data")

		switch flag {
		case "name":
			var request typing.NameSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.NameSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "email":
			var request typing.EmailSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.EmailSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "phone-number":
			var request typing.PhoneNumberSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.PhoneNumberSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "address":
			var request typing.AddressSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.AddressSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "social-media":
			var request typing.SocialMediaSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.SocialMediaSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "job-experience":
			var request typing.JobExperienceSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.JobExperienceSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "education":
			var request typing.EducationSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.EducationSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "skill":
			var request typing.SkillSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.SkillSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "certificate":
			var request typing.CertificateSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.CertificateSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "language":
			var request typing.LanguageSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.LanguageSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		case "project":
			var request typing.ProjectSection
			if err := c.ShouldBindJSON(&request); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
				return
			}

			validated_request, err := util.ValidateData[typing.ProjectSection](request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			err = validator.CheckTone(validated_request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if validated_request.Index == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Index is required."})
				return
			}
			err = util.EditShowCaseRecordData(validated_request, *validated_request.Index, public_user_id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
				return
			}

			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to determine category of request."})
			return
		}
	}
}
