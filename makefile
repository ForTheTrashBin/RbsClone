# Vor einem (erneuten) Aufruf von sqlc generate MUSS, MUSS, MUSS das 'generated'-Verzeichnis von Hand gelöscht werden

#------------------------------------------------------------------------------

sqlc:
	rm -rf internal/rbsdb/*
	sqlc generate -f sqlc.yaml
	
#------------------------------------------------------------------------------

goose-validate:
	goose -dir ./db/migrations validate -v

goose-status:
	goose -dir ./db/migrations postgres "host=skylax-dkt-01-docker port=5432 user=postgres password=eiterbatzen123 dbname=my_database sslmode=disable" status

goose-up:
	goose -dir ./db/migrations postgres "host=skylax-dkt-01-docker port=5432 user=postgres password=eiterbatzen123 dbname=my_database sslmode=disable" up

goose-down:
	goose -dir ./db/migrations postgres "host=skylax-dkt-01-docker port=5432 user=postgres password=eiterbatzen123 dbname=my_database sslmode=disable" down

#------------------------------------------------------------------------------

clean:
	go clean -cache -testcache -modcache
	go mod tidy

# skylax-dkt-01-docker:5432/my_database