package util

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	security "lib/security/punish"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// ReadRequestBody : Takes request body and converts it to JSON
func ReadRequestBody(body io.Reader, response http.ResponseWriter, req *http.Request, clientConn *mongo.Client) map[string]string {
	bodyByteArr, err := ioutil.ReadAll(body)
	if Must(err) {
		RejectResponse(response, clientConn, req)
		return nil
	}
	var retVal map[string]string
	errorr := json.Unmarshal(bodyByteArr, &retVal)
	if Must(errorr) {
		RejectResponse(response, clientConn, req)
		return nil
	}
	// Putting the IP Address of the Request in this result
	// since it'll help with rejection and banning logic.
	retVal["ipAddr"] = strings.Split(req.RemoteAddr, ":")[0]
	return retVal
}

// RejectResponse : Rejects invalid requests
func RejectResponse(response http.ResponseWriter, clientConn *mongo.Client, request *http.Request) {
	rejected := map[string]string{"result": "Rejected Request"}
	clientConn.Disconnect(context.Background())
	json.NewEncoder(response).Encode(rejected)
	security.Reject(strings.Split(request.RemoteAddr, ":")[0])
	return
}

// SendResponse : Sends a successful response
func SendResponse(response http.ResponseWriter, res map[string]interface{}, clientConn *mongo.Client) {
	clientConn.Disconnect(context.Background())
	json.NewEncoder(response).Encode(res)
	return
}
