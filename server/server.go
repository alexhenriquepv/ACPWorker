package main

import (
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
)

var (
	acp *Acp
)

type Coordinate []float64
type Coordinates []Coordinate

type Feature struct {
	Attributes interface{} `json:"attributes"`
	Geometry struct {
		GeometryType string `json:"type"`
		Coordinates Coordinates `json:"coordinates"`
	}
}

type Collection struct {
	Features []Feature `json:"features"`
}

type TaskArgs struct {
	Dataset string
	Collection string
	Method string
}

type ReqHandler struct {}

type FeaturesResponse struct {
	Message string `json:"message"`
	Data Collection `json:"data"`
	Error error `json:"error"`
}

func (r *Feature) AddCoordinates(c ...Coordinate) {
	r.Geometry.Coordinates = append(r.Geometry.Coordinates, c...)
}

func (r *FeaturesResponse) SetData(body Collection) {
	r.Data = body
}

func (r *FeaturesResponse) GetError() error {
	return r.Error
}

func (r *FeaturesResponse) GetBody() interface{} {
	return r.Data
}

func (h *ReqHandler) addJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func (h *ReqHandler) encodeResponse(w http.ResponseWriter, response *FeaturesResponse) {
	var (
		err error
	)
	if err = response.GetError(); err != nil {
		log.Print(err)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(response.GetBody()); err != nil {
		log.Printf("Error Encode response body: %s", err)
		return
	}
}

func (h *ReqHandler) GetAll (w http.ResponseWriter, r *http.Request) {

	var data TaskArgs
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	// nome do metodo a ser executado
	acp.PutString(data.Method)

	// nome do dataset
	acp.PutString(data.Dataset)

	// nome da collection
	acp.PutString(data.Collection)

	// espera o resultado, 0 - positivo, 1 - negativo
	status := acp.GetUbyte()

	var fResponse = &FeaturesResponse{}
	
	if status != 0 {
		log.Printf("ACP status fail")
		message := acp.GetString()
		fResponse.Message = message
	} else {
		no_records := acp.GetUint()
		no_fields := acp.GetUint()
		collection := Collection{}
		
		for i := 0; i < no_records; i++ {
			var feature Feature
			fields := make(map[string]interface{})
			log.Printf("enviando %d de %d", i + 1, no_records)
			for j := 0; j < no_fields; j++ {
				log.Printf("wait field_type")
				field_type := acp.GetString()
				log.Printf("wait key")
				key := acp.GetString()
				log.Printf("start field %s - %s, %d de %d", key, field_type, j + 1, no_fields)
				if field_type == "chain" {
					feature.Geometry.GeometryType = field_type
					no_segments := acp.GetUint()
					no_coords := no_segments * 2
					for k := 0; k < no_coords; k++ {
						coordinate := Coordinate{acp.GetFloat(), acp.GetFloat()}
						feature.AddCoordinates(coordinate)
					}
				} else if field_type == "point" {
					feature.Geometry.GeometryType = field_type
					coordinate := Coordinate{acp.GetFloat(), acp.GetFloat()}
					feature.AddCoordinates(coordinate)
				} else if field_type == "area" {
					feature.Geometry.GeometryType = field_type
					log.Printf("wait no_coords")
					no_coords := acp.GetUint()
					// sobra := acp.GetString()
					log.Printf("no_coords: %d", no_coords)
					for k := 0; k < no_coords; k++ {
						coordinate := Coordinate{acp.GetFloat(), acp.GetFloat()}
						feature.AddCoordinates(coordinate)
					}

					log.Printf("mounted polygon")
				} else {
					fields[key] = acp.GetString()
				}

				feature.Attributes = fields
				log.Printf("end field")
			}
			collection.Features = append(collection.Features, feature)
			log.Printf("end record")
		}
		log.Printf("end collection")
		fResponse.Message = "Data received ok"
		fResponse.Data = collection

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		h.encodeResponse(w, fResponse)
	}	
}

func main () {

	acp = NewAcp("w1")

	if err := acp.Connect("w1", 0, 1); err != nil {
		log.Panicf("ACP Connection error: %v\n", err)
	}

	handler := &ReqHandler{}

	r := mux.NewRouter()
	r.HandleFunc("/rpc", handler.GetAll).Methods("POST")

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