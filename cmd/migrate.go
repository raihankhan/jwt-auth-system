package cmd

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/raihankhan/jwt-auth-system/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  "Commands to manage database migrations (up, down).",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("down")
	},
}

// SetupMigrateCommands adds the migrate subcommand and its children to the rootCmd.
func SetupMigrateCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
}

func runMigrations(direction string) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbSourceURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)

	m, err := migrate.New("file://migrations", dbSourceURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatalf("Migration 'up' failed: %v", err)
			} else {
				fmt.Println("No migrations to apply.")
			}
			return
		}
		fmt.Println("Migrations applied successfully (up).")
	} else if direction == "down" {
		if err := m.Steps(-1); err != nil {
			log.Fatalf("Migration 'down' failed: %v", err)
			return
		}
		fmt.Println("Migration rolled back successfully (down).")
	} else {
		log.Fatalf("Invalid migration direction: %s", direction)
	}
}
