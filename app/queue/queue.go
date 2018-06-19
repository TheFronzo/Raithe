package queue

import (
	"github.com/catmullet/Raithe/app/store"
	"github.com/gin-gonic/gin/json"
	"github.com/catmullet/Raithe/app/auth/model"
)

type Message struct {
	Queue   string              `json:"queue"`
	Message interface{}         `json:"message"`
	Token   model.SecurityToken `json:"security_token"`
}

type PopRequest struct {
	Queue string              `json:"queue"`
	Token model.SecurityToken `json:"security_token"`
}

type PushResponse struct {
	Success bool `json:"success"`
}

type PopResponse struct {
	Queue string `json:"queue"`
	Body interface{} `json:"body"`
}

func PushToQueue(msg Message) error {
	byteSlice, _ := json.Marshal(msg)

	err := store.Set(msg.Queue, byteSlice)
	return err
}

func GetFromQueue(queue string) ([]byte, error) {
	return store.Get(queue)
}