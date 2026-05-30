module github.com/alekparkhomenko/investor/investor

replace github.com/alekparkhomenko/investor/platform => ../plantform

go 1.25.7

require (
	github.com/alekparkhomenko/investor/platform v0.0.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgx/v5 v5.9.2
	github.com/joho/godotenv v1.5.1
	github.com/pressly/goose/v3 v3.27.1
	github.com/swaggo/http-swagger/v2 v2.0.2
	go.uber.org/zap v1.27.1
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/edaniel30/loki-logger-go v0.6.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/rogpeppe/go-internal v1.15.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/swaggo/files/v2 v2.0.0 // indirect
	github.com/swaggo/swag v1.8.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/tools v0.43.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
