package main

import (
	"context"
	"github.com/Eydzhpee08/mittingroom/cmd/app"
	"github.com/gorilla/mux"
	"github.com/Eydzhpee08/mittingroom/pkg/users"
	"github.com/Eydzhpee08/mittingroom/pkg/security"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}

}

func execute(host string, port string, dsn string) (err error) {
	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool,error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		users.NewService,
		security.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			return err
		}
	}

	err = container.Invoke(func(server *app.Server){
		server.Init()
	})
	if err != nil {
		return err
	}
	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})

}
