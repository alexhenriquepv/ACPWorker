package main

import (
	"net/http"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
)

var (
	acp *Acp
)

type TaskArgs struct {
	Dataset string
	Collection string
	Method string
}

type Task struct {
	Message string
}

type TaskService struct {}

func (h *TaskService) Show (req *http.Request, args *TaskArgs, res *Task) error {
	
	// nome do metodo a ser executado
	acp.PutString(args.Method)

	// nome do dataset
	acp.PutString(args.Dataset)

	// nome da collection
	acp.PutString(args.Collection)

	// espera o resultado, 0 - positivo, 1 - negativo
	status := acp.GetUbyte()
	
	// espera a primeira mensagem de resposta (nomes dos campos)
	message := acp.GetString()
	
	// espera a segunda mensagem de resposta (valores dos campos)
	// values := acp.GetString()
	// log.Printf("values %s", values)
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
	r := mux.NewRouter()
	r.Handle("/rpc", s)

	ct := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	or := handlers.AllowedOrigins([]string{"*"})
	mt := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	e := http.ListenAndServe(":4001", handlers.CORS(ct, or, mt)(r))

	if e == nil {
		log.Printf("server ok")
	} else {
		acp.PutString("server_fail")
		log.Fatal("server fail")
	}
}