package main

import (
	"context"
	"fmt"
	"github.com/Insulince/jnet-api/untitled/cmd/api/database"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet-api/untitled/cmd/api/rest"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	trainingexecutor "github.com/Insulince/jnet-api/untitled/cmd/api/training-executor"
	"github.com/Insulince/jnet/pkg/network"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	workerCount = 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := Run(); err != nil {
		panic(errors.Wrap(err, filepath.Base(os.Args[0])))
	}
}

func Run() error {
	fmt.Println("jnet-api is starting")

	fmt.Println("creating context")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errs := make(chan error)

	fmt.Println("connecting to mongo server")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return errors.Wrap(err, "connecting to mongo server")
	}
	defer func() { _ = client.Disconnect(ctx) }()

	fmt.Println("creating mongo")
	m, err := database.NewMongo(ctx, client)
	if err != nil {
		return errors.Wrap(err, "new mongo")
	}

	fmt.Println("creating training queue")
	tq := make(networks.TrainingQueue)

	fmt.Println("creating service")
	nt := network.NewProtoTranslator(network.WithCompression())
	s := service.New(m, nt, tq)

	fmt.Println("creating and executing trainer")
	te := trainingexecutor.New(s, workerCount)
	go te.Execute(ctx, tq, errs)

	fmt.Println("setting up web socket infrastructure")
	var wsUpgrader websocket.Upgrader
	wsUpgrader.ReadBufferSize = 1024
	wsUpgrader.WriteBufferSize = 1024
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	fmt.Println("creating router")
	r := rest.NewRouter(s, wsUpgrader)

	var c cors.Options
	c.AllowedOrigins = []string{"*"}
	c.AllowedMethods = []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
	c.AllowedHeaders = []string{"*"}
	c.AllowCredentials = false

	h := cors.New(c).Handler(r)

	go func() {
		fmt.Println("jnet-api serving")
		err = http.ListenAndServe(":8080", h)
		if err != nil {
			errs <- errors.Wrap(err, "listen and serve")
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errs:
		return err
	}
}
