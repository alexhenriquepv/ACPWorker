package main

import (
	"net/http"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
)

var (
	acp *Acp
)

type TaskArgs struct {
	Dataset string
	Collection string
}

type Task struct {
	Message string
}

type TaskService struct {}

func (h *TaskService) Show (req *http.Request, args *TaskArgs, res *Task) error {
	
	// nome do metodo a ser executado
	acp.PutString("get_all")

	// nome do dataset
	acp.PutString(args.Dataset)

	// nome da collection
	acp.PutString(args.Collection)

	// espera o resultado, 0 - positivo, 1 - negativo
	status := acp.GetUbyte()
	
	// espera a mensagem de resposta
	message := acp.GetString()

	if status != 0 {
		log.Printf("ACP status fail")
	}

	res.Message = message
	return nil
}

func main () {

	acp = NewAcp("w1")

	if err := acp.Connect("w1", 0, 1); err != nil {
		log.Panicf("ACP Connection error: %v\n", err)
	}

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(TaskService), "")
	http.Handle("/rpc", s)
	
	e := http.ListenAndServe("10.84.125.24:4001", nil)

	if e == nil {
		log.Printf("server ok")
	} else {
		log.Fatal("server fail")
	}
}