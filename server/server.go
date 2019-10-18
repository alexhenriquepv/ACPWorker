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

type BodyElement map[string]interface{}
type Body []BodyElement

type TaskArgs struct {
	Dataset string
	Collection string
	Method string
}

type Task struct {
	Message string
	Data Body
}

type TaskService struct {}

func (h *TaskService) GetAll (req *http.Request, args *TaskArgs, res *Task) error {
	
	// nome do metodo a ser executado
	acp.PutString(args.Method)

	// nome do dataset
	acp.PutString(args.Dataset)

	// nome da collection
	acp.PutString(args.Collection)

	// espera o resultado, 0 - positivo, 1 - negativo
	status := acp.GetUbyte()
	
	if status != 0 {
		log.Printf("ACP status fail")
		message := acp.GetString()
		res.Message = message
	} else {
		no_records := acp.GetUint()
		no_fields := acp.GetUint()
		body := make(Body, 0, no_records)
		
		for i := 0; i < no_records; i++ {
			bodyEl := make(BodyElement)

			for j := 0; j < no_fields; j++ {
				key := acp.GetString()
				val := acp.GetString()
				bodyEl[key] = val
			}

			body = append(body, bodyEl)
		}
		
		res.Message = "Data received ok"
		res.Data = body
	}
	
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