package config

type ApiEnvConfig struct {
	Port        string
	Env         string
	Host        string
	AuthService string
}

const DEV_ENV = "DEV"
const STAGE_ENV = "STAGE"
const PROD_ENV = "PROD"
