package response

type Response struct {
	Status  int
	Message string
}

type ResponseStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}
