package model

type Job struct {
	ID      int    `json:"id"`
	Type    string `json:type`
	Payload string `json:payload`
	Status  string `json:status`
}

type CreateJob struct {
	Type    string `json:type`
	Payload string `json:payload`
}
