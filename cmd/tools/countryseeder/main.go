package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/ForTheTrashBin/RbsClone/internal/serverconfig"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedCountries(logger *slog.Logger, dbPool *pgxpool.Pool) error {

	inputName := "countries.csv"

	//-------------------------------------------------------------------------
	// Evaluate current tool position to place result
	//-------------------------------------------------------------------------

	_, goFileName, _, ok := runtime.Caller(0)

	if !ok {

		return fmt.Errorf("Can't evaluate file position")
	}

	inputPath := filepath.Join(filepath.Dir(goFileName), inputName)

	file, err := os.Open(inputPath)

	if err != nil {

		return fmt.Errorf("Can't open csv file, error: %w", err)
	}

	defer file.Close()

	//-------------------------------------------------------------------------
	// Configure CSV-Reader (importand: Semikolon as delimiter)
	//-------------------------------------------------------------------------

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	//-------------------------------------------------------------------------
	// Skip headline
	//-------------------------------------------------------------------------

	_, err = reader.Read()

	if err != nil {

		return fmt.Errorf("Error reading headline, error: %w", err)
	}

	//-------------------------------------------------------------------------
	// Read CSV file and insert into database
	//-------------------------------------------------------------------------

	logger.Info("Start import of country data...")

	baseQueries := rbsdb.New(dbPool)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for {
		record, err := reader.Read()

		if err == io.EOF {

			break // End of file reached
		}

		if err != nil {

			logger.Error("Error reading a line (skiped)", "error", err)

			continue
		}

		//---------------------------------------------------------------------
		// Read data from CSV columns
		//---------------------------------------------------------------------

		shortcode := record[0]
		name := record[2]

		//---------------------------------------------------------------------
		// Convert IBAN length to int/int32
		//---------------------------------------------------------------------

		ibanLengthInt, err := strconv.Atoi(record[4])

		if err != nil {

			logger.Error("Invalid IBAN-length. Set to zero!", "country", shortcode)

			ibanLengthInt = 0
		}

		ibanLength := int16(ibanLengthInt)

		//---------------------------------------------------------------------
		// Call database via SQLC to insert record
		//---------------------------------------------------------------------

		_, err = baseQueries.InsertCountry(ctx, rbsdb.InsertCountryParams{

			Shortcode:  shortcode,
			Name:       name,
			Ibanlength: pgtype.Int2{Int16: ibanLength, Valid: ibanLength > 0},
		})

		if err != nil {

			return fmt.Errorf("Error insert record: %w", err)

			// logger.Error("Error insert record", "country", shortcode, "error", err)

			// continue
		}
	}

	logger.Info("Import of country data finished")

	return nil
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

	if err := seedCountries(logger, dbPool); err != nil {

		logger.Error("Error initializing countries", "error", err)
	}
}
