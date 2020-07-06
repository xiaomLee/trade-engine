package grpchandler

import grpc_end "github.com/xiaomLee/grpc-end"

func SayHi(c *grpc_end.GRpcContext) {
	name := c.StringParamDefault("name", "")
	c.SuccessResponse("hi " + name)
}
