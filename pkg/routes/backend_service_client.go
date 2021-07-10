package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const (
	DMLINSERTCLOUDWORKER = `INSERT INTO cloud_workers('requestId','workerId','cloud','response') 
VALUES(?,?,?,?);`
)

var (
	client            = resty.New()
	backendServiceUrl = os.Getenv("HYBRID_CLOUD_BACKEND_URL")
)

func sendMessage(request Request, db *sql.DB) {
	requestId := generateRequestId()
	rChan := make(chan *resty.Response)
	eChan := make(chan error)

	//send it to backend
	go func() {
		message := &Message{
			RequestId: requestId,
			Request:   request,
		}
		log.Printf("Sending message for request %s to backend %s", requestId, message)
		response, err := client.R().
			EnableTrace().
			SetBody(message).
			Post(fmt.Sprintf("%s/process", backendServiceUrl))
		if err != nil {
			eChan <- err
		} else {
			rChan <- response
		}
	}()

	go saveMessageToDB(requestId, rChan, eChan, db)
}

func saveMessageToDB(requestId string, rChan chan *resty.Response, eChan chan error, db *sql.DB) {
	log.Printf("Waiting to save response to Database for request %s", requestId)
	select {
	case res := <-rChan:
		log.Printf("Processing  response for request %s", requestId)
		var bResponse Response
		b := res.Body()
		err := json.Unmarshal(b, &bResponse)
		if err != nil {
			log.Printf("Error marshalling the response %s", err)
		} else {
			log.Printf("Saving response %s", bResponse)
			if tx, err := db.Begin(); err != nil {
				log.Printf("Unable to begin transaction %s", err)
				return
			} else {
				if stmt, err := db.Prepare(DMLINSERTCLOUDWORKER); err != nil {
					log.Printf("Error preparing statement %s", err)
				} else {
					if _, err := stmt.Exec(bResponse.RequestId, bResponse.WorkerId,
						bResponse.CloudId, bResponse.Text); err != nil {
						if e := tx.Rollback(); e != nil {
							log.Printf("Unable to rollback %s", err)
						}
						log.Printf("Error saving response for requestId %s,%s", requestId, err)
						return
					}
					log.Printf("Done Saving response message for requestId %s", requestId)
					tx.Commit()
				}
			}
		}
	case err := <-eChan:
		log.Printf("Error processing message for requestId %s, %s", requestId, err)
	}
}

func generateRequestId() string {
	return "request-" + uuid.New().String()[:4]
}
