package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/data"
	"github.com/swaggo/swag/example/celler/httputil"
)

const ()

// Responses godoc
// @Summary Retrieves all responses processed by the backend
// @Description Retrieves all responses processed by the backend sorted by timestamp
// @Tags worker
// @Produce json
// @Success 200 {object} data.CloudWorkers  "Processed response data"
// @Router /workers/all [get]
func (e *Endpoints) Responses(c *gin.Context) {
	if rows, err := e.DBConn.Query(data.DMLALLRESPONSESCLOUDWORKERS); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	} else {
		var cws data.CloudWorkers
		for rows.Next() {
			var cw data.CloudWorker
			if err := rows.Scan(&cw.RequestId, &cw.WorkerId, &cw.Cloud, &cw.Response, &cw.LastProcessedTimestamp); err != nil {
				log.Printf("Error while retriving row data %s", err)
			} else {
				cws = append(cws, cw)
			}

		}
		c.JSON(http.StatusOK, cws)
	}
}

// CloudWorkerRequests godoc
// @Summary Cloud Workers and the total number of messages processed by them
// @Description List of all the Cloud Workers and total number of messages processed by them
// @Tags worker
// @Produce json
// @Success 200 {object} data.CloudWorkerRequests "The total number of requests processed by each cloud"
// @Router /workers/cloud [get]
func (e *Endpoints) CloudWorkerRequests(c *gin.Context) {
	if rows, err := e.DBConn.Query(data.DMLREQUESTSPROCESSEDBYCLOUDWORKER); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	} else {
		var cwrs data.CloudWorkerRequests
		for rows.Next() {
			var cwr data.CloudWorkerRequest
			if err := rows.Scan(&cwr.Cloud, &cwr.RequestsProcessed); err != nil {
				log.Printf("Error while retriving row data %s", err)
			} else {
				cwrs = append(cwrs, cwr)
			}

		}
		c.JSON(http.StatusOK, cwrs)
	}
}
