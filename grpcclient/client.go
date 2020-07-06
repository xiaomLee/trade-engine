package apiClient

import "github.com/xiaomLee/grpc-end/client"

var (
	gRpcMapPool *client.MapPool
	servers     map[string]string
)

// use defaultDialFunc
func init() {
	client.InitClient(nil)
}
