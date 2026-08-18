package tools

import (
	"context"
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeCountry(logger *slog.Logger, dbPool *pgxpool.Pool) {

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
