package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"log"
	"net/http"
)

const (
	DMLALLRESPONSESCLOUDWORKERS       = `SELECT requestId,workerId,cloud,response,timestamp from cloud_workers ORDER BY timestamp desc;`
	DMLREQUESTSPROCESSEDBYCLOUDWORKER = `SELECT cw.cloud as cloud,sum(cw.requestsProcessed) as requestTotal 
  FROM cloud_workers cw GROUP BY cw.cloud;`
)

// Responses godoc
// @Summary Retrieves all responses processed by the backend
// @Description Retrieves all responses processed by the backend sorted by timestamp
// @Tags worker
// @Produce json
// @Success 200 {object} routes.CloudWorkers  "Processed response data"
// @Router /workers/all [get]
func (e *Endpoints) Responses(c *gin.Context) {
	if rows, err := e.DBConn.Query(DMLALLRESPONSESCLOUDWORKERS); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	} else {
		var cws CloudWorkers
		for rows.Next() {
			var cw CloudWorker
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
// @Success 200 {object} routes.CloudWorkerRequests "The total number of requests processed by each cloud"
// @Router /workers/cloud [get]
func (e *Endpoints) CloudWorkerRequests(c *gin.Context) {
	if rows, err := e.DBConn.Query(DMLREQUESTSPROCESSEDBYCLOUDWORKER); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	} else {
		var cwrs CloudWorkerRequests
		for rows.Next() {
			var cwr CloudWorkerRequest
			if err := rows.Scan(&cwr.Cloud, &cwr.RequestsProcessed); err != nil {
				log.Printf("Error while retriving row data %s", err)
			} else {
				cwrs = append(cwrs, cwr)
			}

		}
		c.JSON(http.StatusOK, cwrs)
	}
}
