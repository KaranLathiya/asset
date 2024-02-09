package response

import (
	"asset/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

func MessageShow(code int, errMessage string, w http.ResponseWriter) {
	var message model.Message
	message.Code = code
	message.Message = errMessage
	user_data, _ := json.MarshalIndent(message, "", "  ")
	w.WriteHeader(code)
	w.Write(user_data)
}

func DatabaseErrorShow(err error) (string, int) {
	if dbErr, ok := err.(*pq.Error); ok { // For PostgreSQL database driver (pq)
		// Access PostgreSQL-specific error fields
		// errCode,_ :=  strconv.Atoi(dbErr.Code)
		errCode := dbErr.Code
		// errMessage := errCode.Name()
		// errDetail := dbErr.Detail
		// Handle the PostgreSQL-specific error
		// fmt.Println(errCode)
		// fmt.Println(errDetail)
		// fmt.Println(errMessage)
		switch errCode {
		case "23502":
			// not-null constraint violation
			return "Some required data was left out", 400

		case "23503":
			// foreign key violation
			return "This record can't be changed because another record refers to it", 409

		case "23505":
			// unique constraint violation
			return "This record contains duplicated data that conflicts with what is already in the database", 409

		case "23514":
			// check constraint violation
			return "This record contains inconsistent or out-of-range data", 400

		}
	}
	fmt.Println(err.Error())
	return "Internal server error", 500
}
