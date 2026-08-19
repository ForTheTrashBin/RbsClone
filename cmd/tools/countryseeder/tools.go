package tools

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/ForTheTrashBin/RbsClone/internal/serverconfig"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeCountry(logger *slog.Logger, dbPool *pgxpool.Pool) {

	wd, err := os.Getwd()

	logger.Info("Arbeitsverzeichnis", "Path", wd)

	thefile := filepath.Join(wd, "internal/tools/countries.csv")

	file, err := os.Open(thefile)

	if err != nil {

		logger.Error("Fehler beim Öffnen der CSV-Datei", "Fatal", err)

		return
	}

	defer file.Close()

	// 2. CSV-Reader konfigurieren (wichtig: Semikolon als Trennzeichen)

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	// 3. Kopfzeile (Header) überspringen

	_, err = reader.Read()

	if err != nil {

		logger.Error("Fehler beim Lesen der Kopfzeile", "error", err)

		return
	}

	// 4. Ihre bestehende Datenbankverbindung nutzen
	// (Hier wird vorausgesetzt, dass Sie 'dbConn' als *sql.DB oder pgx-Verbindung bereits haben)

	//-------------------------------------------------------------------------

	baseQueries := rbsdb.New(dbPool)

	logger.Info("Starte Import der Länder-Stammdaten...")

	count := 0

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	// 5. Zeile für Zeile einlesen und in die DB schreiben
	for {
		record, err := reader.Read()

		if err == io.EOF {

			break // Ende der Datei erreicht
		}

		if err != nil {

			logger.Error("Fehler beim Lesen einer Zeile (Übersprungen)", "error", err)

			continue
		}

		// Daten aus den CSV-Spalten extrahieren
		iso2 := record[0]
		// iso3 := record[1]
		landDe := record[2]
		// waehrung := record[3]

		// IBAN-Länge von String zu int/int32 konvertieren

		ibanLaengeInt, err := strconv.Atoi(record[4])

		if err != nil {

			logger.Error("Ungültige IBAN-Länge. Setze auf 0.", "Land", landDe, "error", err)

			ibanLaengeInt = 0
		}

		ibanLaenge := int16(ibanLaengeInt)

		// 6. Aufruf der von sqlc generierten Insert-Methode

		_, err = baseQueries.InsertCountry(ctx, rbsdb.InsertCountryParams{

			Shortcode:  iso2,
			Name:       landDe,
			Ibanlength: pgtype.Int2{Int16: ibanLaenge, Valid: ibanLaenge > 0},
		})
		if err != nil {

			logger.Error("Fehler beim Einfügen", "Land", landDe, "error", err)

			continue
		}

		count++
	}

	logger.Info("Import erfolgreich abgeschlossen!", "Count", count)
}

//-----------------------------------------------------------------------------
// The main entry-point of this app
//-----------------------------------------------------------------------------

func main() {

	var config serverconfig.Config

	//-------------------------------------------------------------------------
	// Read content of .env-file into the 'local' environment of this proccess
	//-------------------------------------------------------------------------

	if err := config.InitGoDotEnv(); err != nil {

		panic(fmt.Errorf("Error from 'godotenv': %w", err))
	}

	//-------------------------------------------------------------------------
	// Parse the content of the 'local' environmet into the Config-struct
	//-------------------------------------------------------------------------

	if err := config.InitCaarlos0(); err != nil {

		panic(fmt.Errorf("Error from 'env': %w", err))
	}

	//-------------------------------------------------------------------------
	// Create a new structured logger that writes JSON to stdout.
	//-------------------------------------------------------------------------

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel})) // Level: parseLogLevel(config.LogLevel)}))

	logger.Info("Application config", "DB_USER", config.DB.User)
	logger.Info("Application config", "DB_PASSWORD", config.DB.Password)
	logger.Info("Application config", "DB_HOST", config.DB.Host)
	logger.Info("Application config", "DB_PORT", config.DB.Port)
	logger.Info("Application config", "DB_DATABASE", config.DB.Database)

	logger.Info("Application config", "RT_HTTPPORT", config.GetHTTPPort())
	logger.Info("Application config", "RT_HTTPSPORT", config.GetHTTPSPort())

	//-------------------------------------------------------------------------
	// Connect to database
	//-------------------------------------------------------------------------

	// Create a pgxpool.Config for manual configuration

	pgxconfig, err := pgxpool.ParseConfig("")

	if err != nil {

		logger.Error("Can't parse postgres-config", "error", err)

		return
	}

	// Create a pgxpool.Config for manual configuration

	pgxconfig.ConnConfig.User = config.DB.User
	pgxconfig.ConnConfig.Password = config.DB.Password
	pgxconfig.ConnConfig.Host = config.DB.Host
	pgxconfig.ConnConfig.Port = config.DB.Port
	pgxconfig.ConnConfig.Database = config.DB.Database

	pgxconfig.ConnConfig.ConnectTimeout = 5 * time.Second

	// Connection pool tuning

	pgxconfig.MinConns = 5
	pgxconfig.MaxConns = 25
	pgxconfig.MaxConnLifetime = 1 * time.Hour
	pgxconfig.MaxConnIdleTime = 15 * time.Minute
	pgxconfig.HealthCheckPeriod = 1 * time.Minute

	//-------------------------------------------------------------------------
	// Connect to postgres database
	//-------------------------------------------------------------------------

	ctxDBConnect, cancelDBConnectTimeout := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancelDBConnectTimeout()

	//-------------------------------------------------------------------------

	dbPool, err := pgxpool.NewWithConfig(ctxDBConnect, pgxconfig)

	if err != nil {

		logger.Error("Can't open database", "error", err)

		return
	}

	defer dbPool.Close()

	//-------------------------------------------------------------------------

	if err := dbPool.Ping(ctxDBConnect); err != nil {

		logger.Error("Can't ping database", "error", err)

		return
	}

	//-------------------------------------------------------------------------

	cancelDBConnectTimeout()

	//-------------------------------------------------------------------------

	initializeCountry(logger, dbPool)
}
