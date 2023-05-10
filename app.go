package main

import (
	"KilnBackendHomeExercice/utils"
	"encoding/json"
	"github.com/onrik/ethrpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
)

type Counter int32

var (
	isReady                 = false
	isLive                  = true
	failedRequests  Counter = 0
	succeedRequests Counter = 0
)

func main() {
	// Kubernetes endpoints
	http.HandleFunc(utils.ReadinessEndpoint, handleReady)
	http.HandleFunc(utils.LivenessEndpoint, handleLive)

	// Prometheus metric endpoint
	http.Handle(utils.MetricsEndpoint, promhttp.Handler())

	// Gas Price endpoint
	http.HandleFunc(utils.GasPriceEndpoint, handleGasPrice)

	// Launch the web app on port 8080
	isReady = true
	err := http.ListenAndServe(utils.ServerAddress, nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

// Readiness HTTP endpoint handler for Kubernetes to be used as a probe
func handleReady(w http.ResponseWriter, r *http.Request) {
	if isReady {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// Liveness HTTP endpoint handler for Kubernetes to be used as a probe
func handleLive(w http.ResponseWriter, r *http.Request) {
	if isLive {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// Gas price HTTP endpoint
func handleGasPrice(w http.ResponseWriter, r *http.Request) {

	client := ethrpc.New(utils.InfuraUrl + utils.InfuraApiKey)
	// Checking the Web3 RPC Client
	_, err := client.Web3ClientVersion()
	if err != nil {
		failedRequests.handleFailed(w, err, utils.Web3ClientCreationError)
		return
	}

	// Doing the GesPrice RPC request and check no error happen
	gasPrice, err := client.EthGasPrice()
	if err != nil {
		failedRequests.handleFailed(w, err, utils.Web3RpcRequestError)
		return
	}

	jsonResult := map[string]interface{}{
		"gasPrice": "0x" + gasPrice.Text(16),
	}

	// Jsonify the Map and check no error happen during the Json Encoding
	err = json.NewEncoder(w).Encode(jsonResult)
	if err != nil {
		failedRequests.handleFailed(w, err, utils.JsonEncodingError)
		return
	}
	succeedRequests.handleSucceed(&failedRequests)
}

// Handle failed HTTP request
func (c *Counter) handleFailed(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
	atomic.AddInt32((*int32)(c), 1)

	nbFailedRequests := atomic.LoadInt32((*int32)(c))
	log.Info("Request failed ! Number failed requests since last success: ", nbFailedRequests)

	// Let's say to K8s to launch the restart policy
	if nbFailedRequests > utils.MaxConsecutiveRequestErrorAccepted {
		isLive = false
	}
}

// Handle succeed HTTP request
func (c *Counter) handleSucceed(failedRequests *Counter) {
	atomic.AddInt32((*int32)(c), 1)
	log.Info("Request succeed ! Number succeed requests: ", atomic.LoadInt32((*int32)(c)))

	atomic.StoreInt32((*int32)(failedRequests), 0)
}
