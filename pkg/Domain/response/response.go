package response

import "encoding/json"

type Response struct {
	Status  int
	Message string
}

type ResponseStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

func (r Response) BuildResponse(data interface{}) ([]byte, error) {
	rs := ResponseStruct{Message: r.Message, Data: data}
	return json.Marshal(rs)
}
