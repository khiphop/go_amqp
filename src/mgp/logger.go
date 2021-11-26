package mgp

import (
	"log"
)

func Info(msg string, uuid string)  {
	content := uuid + " | " +msg
	log.Println(content)
}