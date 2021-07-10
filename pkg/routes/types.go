package routes

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

//Endpoints is a marker for defining REST API routes
type Endpoints struct {
	DBConn *sql.DB
}

//NewEndpoints provides handle to the REST URI defined in this API
func NewEndpoints() *Endpoints {
	return &Endpoints{}
}

//Request the worker request to process
type Request struct {
	//Text is any text to process
	Text string `json:"text"`
	//Uppercase change the Text to uppercase
	Uppercase bool `json:"upperCase,omitempty"`
	//Reverse   reverse Text
	Reverse bool `json:"reverse,omitempty"`
	//SleepMillis add some sleep to processing
	SleepMillis int `json:"sleepMillis,omitempty"`
}

//CloudWorker the Cloud Worker info
type CloudWorker struct {
	//RequestId the request id
	RequestId string `json:"requestId"`
	//WorkerId the worker id
	WorkerId string `json:"workerId"`
	//Cloud the cloud which processed the request
	Cloud string `json:"cloud"`
	//Response the processed text with all applied transformations
	Response string `json:"response"`
	//Response the number of requests processed by the cloud
	RequestsProcessed int `json:"requestsProcessed"`
	//LastProcessedTimestamp is the last time when this worker processed the request
	LastProcessedTimestamp time.Time `json:"lastProcessedTimestamp"`
}

//CloudWorkers represents the rows of each CloudWorker
type CloudWorkers []CloudWorker

//CloudWorkerRequest holds the number of requests processed by each cloud
type CloudWorkerRequest struct {
	//Cloud the cloud which processed the request
	Cloud string `json:"cloud"`
	//Response the total number of requests processed by the cloud
	RequestsProcessed int `json:"requestsProcessed"`
}

//CloudWorkerRequests represents the rows of each CloudWorker
type CloudWorkerRequests []CloudWorkerRequest

//Message handles the message that needs to be processed
type Message struct {
	//RequestId the unique request id
	RequestId string  `json:"requestId"`
	Request   Request `json:"request"`
}

//Response is the processed Request
type Response struct {
	//RequestId the request id
	RequestId string `json:"requestId"`
	//WorkerId the worker id
	WorkerId string `json:"workerId"`
	//CloudId the cloud which processed the request
	CloudId string `json:"cloudId"`
	//Text the processed text with all applied transformations
	Text string `json:"text"`
}

func (r *Request) String() string {
	return fmt.Sprintf("Request{text=%s, uppercase=%s, reverse=%s}",
		r.Text, strconv.FormatBool(r.Uppercase), strconv.FormatBool(r.Reverse))
}

func (r *Response) String() string {
	return fmt.Sprintf("Response{requestId=%s,workerId=%s,text=%s, cloud=%s}",
		r.RequestId, r.WorkerId, r.CloudId, r.Text)
}

//func (r *Data) String() string {
//	return fmt.Sprintf("Response{requestId=%s,workerid=%s,cloudId=%s,Text=%s}",
//		r.RequestId, r.WorkerId, r.Cloud, r.Response)
//}

func (m *Message) String() string {
	return fmt.Sprintf("Message{request=%s,requestId=%s}",
		m.Request.String(), m.RequestId)
}
