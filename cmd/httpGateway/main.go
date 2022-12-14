package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	gw "homework-1/pkg/api/v1"
	"io/ioutil"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	// init swagger
	err := mux.HandlePath("GET", "/swagger.json", swaggerHandler)
	if err != nil {
		return errors.Wrap(err, "Can't init swagger handler")
	}

	// init grpc gateway
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = gw.RegisterApiServiceHandlerFromEndpoint(ctx, mux, config.ProxyApiServiceAddress, opts)
	if err != nil {
		return errors.Wrap(err, "Can't init grpc gateway")
	}

	return http.ListenAndServe(config.HTTPGatewayServiceAddress, mux)
}

func swaggerHandler(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	fileBytes, err := ioutil.ReadFile("pkg/api/v1/api.swagger.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(fileBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}
