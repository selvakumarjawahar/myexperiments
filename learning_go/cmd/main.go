package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)
var server string
var timed time.Duration
var DBPort int
var DBName string
var DBUser string
var DBPassword string
var DBHost string

func main() {
	var rootCmd = &cobra.Command{Use: "ncc3-metrics-service", Run: func(c *cobra.Command, args []string) {}}
	rootCmd.Flags().StringVarP(&server, "server", "", "http://localhost", "Netrounds server URL")
	rootCmd.Flags().DurationVarP(&timed, "duration", "", 10*time.Second, "format:- string of decimal numbers with unit suffix,Valid time units are ns,us,ms,s,m,h")
	rootCmd.Flags().StringVarP(
		&DBHost, "db-host", "", "localhost", "The name of the database host")
	rootCmd.Flags().IntVarP(
		&DBPort, "db-port", "", 5432, "The port the database connection")
	rootCmd.Flags().StringVarP(
		&DBName, "db-name", "", "", "The name of the postgres database")
	//rootCmd.MarkFlagRequired("db-name")
	rootCmd.Flags().StringVarP(
		&DBUser, "db-user", "", "", "The database user")
	//rootCmd.MarkFlagRequired("db-user")
	rootCmd.Flags().StringVarP(
		&DBPassword, "db-password", "", "", "The database password")

	viper.BindPFlags(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Quit on help.
	rootCmd.Flags().Visit(func(f *flag.Flag) {
		if f.Name == "help" {
			os.Exit(0)
		}
	})

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		DBHost, DBPort, DBUser, DBPassword, DBName)
	//connStr += " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic().Err(err).Msg("Failed opening DB connection")
	}
	defer db.Close()
	// Open doesn't open a connection. Check that DB is reachable:
	err = db.Ping()
	if err != nil {
		log.Panic().Err(err).Msg("DB unreachable")
	}

	runTestQueries(db)

}

func runTestQueries(db *sql.DB) {
	rows, err := db.Query("SELECT max(time), measurement_id FROM $1  GROUP BY measurement_id","monitor_metrics_http")
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return
	}
	defer rows.Close()
	var maxTime time.Time
	var measurementID int32
	for rows.Next() {
		err := rows.Scan(&maxTime, &measurementID)
		if err != nil {
			log.Debug().Msgf("Error scanning rows %s \n ", err)
			continue
		}
		fmt.Printf("max time = %v MID = %d \n",maxTime,measurementID)
	}
}