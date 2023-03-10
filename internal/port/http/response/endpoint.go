package response

import (
	"time"
)

type EndpointResponse struct {
	Id        string    `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type EndpointRequest struct {
	Status int       `json:"status"`
	Time   time.Time `json:"time"`
}
