package server

import (
	"log"

	"github.com/ekr-paolo-carraro/todoTest/Todo-service/postgres"
)

func Init() {
	delegate, err := postgres.NewTodoDelegate()
	if err != nil {
		log.Fatal(err)
	} else {
		r, err := NewRouter(*delegate)
		if err != nil {
			log.Fatal(err)
			return
		}
		r.Run()
	}

}
