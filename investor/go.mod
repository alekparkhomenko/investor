module github.com/alekparkhomenko/investor/investor

replace github.com/alekparkhomenko/investor/platform => ../plantform

go 1.25.7

require (
	github.com/alekparkhomenko/investor/platform v0.0.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/edaniel30/loki-logger-go v0.6.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
)
