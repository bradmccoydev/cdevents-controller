package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Info godoc
// @Summary Runtime information
// @Description Receive CDEvent
// @Tags HTTP API
// @Accept json
// @Produce json
// @Success 200 {object} api.RuntimeResponse
// @Router /api/cdevent/{key} [post]
func (s *Server) cdEventHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "cdEventHandler")
	defer span.End()

	key := mux.Vars(r)["key"]
	body, err := io.ReadAll(r.Body)
	fmt.Printf("key: %s, body: %s", key, body)
	if err != nil {
		s.ErrorResponse(w, r, span, "reading the request body failed", http.StatusBadRequest)
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb-headless:27017/cdevents?ssl=false"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("cdevents")
	collection := db.Collection("items")

	cdevent := CDEvent{
		SubjectID:        "1",
		SubjectSource:    "2",
		SubjectType:      "3",
		SubjectContent:   "4",
		ContextVersion:   "5",
		ContextID:        "6",
		ContextSource:    "7",
		ContextType:      "8",
		ContextTimestamp: "9",
	}

	_, err = collection.InsertOne(ctx, cdevent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Item inserted successfully! %s", key)

	s.JSONResponse(w, r, body)
}

type CDEvent struct {
	SubjectID        string `json:"subjectId"`
	SubjectSource    string `json:"subjectSource"`
	SubjectType      string `json:"subjectType"`
	SubjectContent   string `json:"subjectContent"`
	ContextVersion   string `json:"contextVersion"`
	ContextID        string `json:"contextId"`
	ContextSource    string `json:"contextSource"`
	ContextType      string `json:"contextType"`
	ContextTimestamp string `json:"contextTimestamp"`
}
