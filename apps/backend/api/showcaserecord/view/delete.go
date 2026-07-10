package showcaserecord_view

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	typing "resuming/api/showcaserecord/typing"
	util "resuming/api/showcaserecord/util"
	validator "resuming/api/showcaserecord/validator"
)

func DeleteShowCaseRecordData() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.SpecificPortoflioDataRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		polished_request, err := validator.ValidateSpecificPortfolioDataRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		index, err := strconv.Atoi(polished_request.Index)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid index format."})
			return
		}

		err = util.DeleteShowCaseRecordData(polished_request.SectionTitle, index, public_user_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Operation failed."})
			return
		}

	}
}
