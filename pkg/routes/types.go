package routes

import (
	"database/sql"
)

//Endpoints is a marker for defining REST API routes
type Endpoints struct {
	DBConn *sql.DB
}

//NewEndpoints provides handle to the REST URI defined in this API
func NewEndpoints() *Endpoints {
	return &Endpoints{}
}
