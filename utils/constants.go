package utils

// General Constants
const (
	ServerAddress                      = ":8080"
	MaxConsecutiveRequestErrorAccepted = 10
	InfuraUrl                          = "https://mainnet.infura.io/v3/"
	InfuraApiKey                       = "<REPLACE-BY-YOUR-INFURA-KEY>"
)

// Endpoints
const (
	LivenessEndpoint  = "/livez"
	ReadinessEndpoint = "/readyz"
	MetricsEndpoint   = "/metrics"
	GasPriceEndpoint  = "/eth/gasprice"
)

// HTTP Error Codes
const (
	Web3ClientCreationError = 600
	Web3RpcRequestError     = 601
	JsonEncodingError       = 700
)
