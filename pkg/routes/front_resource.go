package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/data"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/utils"
)

// SendRequest godoc
// @Summary builds and send request message to backend
// @Description builds and send request message to backend for processing
// @Tags backend
// @Accept json
// @Param message body data.Request true "Message to process"
// @Success 202 {object} data.Response
// @Failure 400 {object} utils.HTTPError
//@Router /send-request [post]
func (e *Endpoints) SendRequest(c *gin.Context) {
	var request data.Request
	if err := c.BindJSON(&request); err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	} else {
		if request.Text == "" {
			utils.NewError(c, http.StatusBadRequest, err)
			return
		}
		log.Printf("Sending message %v  to backend", request)
		go sendMessage(request, e.DBConn)
		c.JSON(http.StatusAccepted, nil)

	}
}
