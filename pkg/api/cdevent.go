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
	"gopkg.in/yaml.v3"
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

	var cdevent CDEvent
	if err := yaml.Unmarshal(body, &cdevent); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(s.config.MongodbURL))
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

	result, err := collection.InsertOne(ctx, cdevent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	s.JSONResponse(w, r, body)
}

type CDEvent struct {
	Context struct {
		Version   string    `json:"version"`
		ID        string    `json:"id"`
		Source    string    `json:"source"`
		Type      string    `json:"type"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"context"`
	Subject struct {
		ID      string `json:"id"`
		Source  string `json:"source"`
		Type    string `json:"type"`
		Content struct {
		} `json:"content"`
	} `json:"subject"`
}
