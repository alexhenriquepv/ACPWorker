package main

import (
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"time"
)

var(
	acp *Acp
	cacheData = cache.New(2 * time.Minute, 3 * time.Minute)
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
	name string
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

	cacheData.Set(response.Data.name, response, cache.DefaultExpiration)
}

func (h *ReqHandler) GetAll (w http.ResponseWriter, r *http.Request) {

	var fResponse = &FeaturesResponse{}
	var data TaskArgs
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	// verifica a existencia dos dados no cache
	content, found := cacheData.Get(data.Collection)
	if found {
		log.Println("content in cache")
		fResponse = content.(*FeaturesResponse)
	} else {
		log.Println("content not cached")

		// nome do metodo a ser executado
		acp.PutString(data.Method)

		// nome do dataset
		acp.PutString(data.Dataset)

		// nome da collection
		acp.PutString(data.Collection)

		// espera o resultado, 0 - positivo, 1 - negativo
		status := acp.GetUbyte()

		if status != 0 {
			log.Printf("ACP status fail")
			message := acp.GetString()
			fResponse.Message = message
			w.WriteHeader(http.StatusBadRequest)
		} else {
			no_records := acp.GetUint()
			no_fields := acp.GetUint()
			collection := Collection{ name: data.Collection }
			
			for i := 0; i < no_records; i++ {
				var feature Feature
				fields := make(map[string]interface{})
				for j := 0; j < no_fields; j++ {
					field_type := acp.GetString()
					key := acp.GetString()
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
						no_coords := acp.GetUint()
						for k := 0; k < no_coords; k++ {
							coordinate := Coordinate{acp.GetFloat(), acp.GetFloat()}
							feature.AddCoordinates(coordinate)
						}
					} else {
						fields[key] = acp.GetString()
					}

					feature.Attributes = fields
				}
				collection.Features = append(collection.Features, feature)
			}
			
			fResponse.Message = "Data received ok"
			fResponse.Data = collection
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	h.encodeResponse(w, fResponse)
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