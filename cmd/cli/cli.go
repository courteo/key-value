package cli

import (
	"bufio"
	"context"
	"fmt"
	"github.com/courteo/key-value/internal/database"
	"github.com/courteo/key-value/internal/database/compute"
	"github.com/courteo/key-value/internal/database/storage"
	in_memory "github.com/courteo/key-value/internal/database/storage/engine/in-memory"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	osSignCh := make(chan os.Signal, 1)

	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	flag := true

	l, _ := zap.NewProduction()

	database := getDatabase(l)

	for flag {
		select {
		case sig := <-osSignCh:
			fmt.Printf("Received signal %s\n", sig.String())

			flag = false

			break
		default:
			in := bufio.NewReader(os.Stdin)

			query, err := in.ReadString('\n')
			if err != nil {
				fmt.Printf("Error input: %s\n", err.Error())
				continue
			}

			query = query[:len(query)-1]

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

			value, err := database.HandleQuery(ctx, query)
			if err != nil {
				fmt.Printf("Error executing query: %s\n", err.Error())
				continue
			}

			if value != "" {
				fmt.Println(value)
				continue
			}
		}
	}
}

func getDatabase(l *zap.Logger) *database.Database {
	computer := compute.NewComputer(l)

	engine := in_memory.New(l)
	storage := storage.New(l, engine)

	database := database.New(l, computer, storage)

	return database
}
