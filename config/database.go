package config

import (
	"context"
	"fmt"
	"log"

	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/migrate"
	_ "github.com/lib/pq"
)

func ConnectDb(cfg Config) *ent.Client {
	dbUrlTemplate := "host=%s port=%s user=%s dbname=%s password=%s sslmode=disable"

	dbUrl := fmt.Sprintf(
		dbUrlTemplate,
		cfg.PostgresServer,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresDB,
		cfg.PostgresPassword,
	)

	client, err := ent.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
