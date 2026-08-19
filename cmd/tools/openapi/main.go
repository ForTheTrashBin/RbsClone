package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func main() {

	if err := run(); err != nil {

		log.Fatal(err)
	}

	log.Println("Creation of output file successful")
}

func run() error {

	inUrl := "https://www.rbsclone.de:8443/openapi.yaml"

	outName := "openapi.yaml"

	//-------------------------------------------------------------------------
	// Create a http-request and copy the content
	//-------------------------------------------------------------------------

	response, err := http.Get(inUrl)

	if err != nil {

		return fmt.Errorf("Can't create http-response-object url: %s, error: %w", inUrl, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		return fmt.Errorf("Http-response-object returning none OK-status, status: %d", response.StatusCode)
	}

	//-------------------------------------------------------------------------
	// Evaluate current tool position to place result
	//-------------------------------------------------------------------------

	_, goFileName, _, ok := runtime.Caller(0)

	if !ok {

		return fmt.Errorf("Can't evaluate file position")
	}

	outPath := filepath.Join(filepath.Dir(goFileName), outName)

	//-------------------------------------------------------------------------
	// Create the output file
	//-------------------------------------------------------------------------

	outFile, err := os.Create(outPath)

	if err != nil {

		return fmt.Errorf("Can't create http-response-object, name: %s, error: %w", outPath, err)
	}

	defer outFile.Close()

	//-------------------------------------------------------------------------

	if _, err := io.Copy(outFile, response.Body); err != nil {

		return fmt.Errorf("Can't copy response body, error: %w", err)
	}

	return nil
}
