package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONAppErrorReporter() gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("\n %v empty \n ",c.Errors)
		c.Next()
		fmt.Println("Req went through handler")
		if (len(c.Errors) != 0) {
			fmt.Printf("\n%v\n",c.Errors)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "error": "username already sign up"})
			return 
		}
		fmt.Println("SUCCESS")
	}
}
