package rpc_client

import (
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

func NewAuthClient(addr string) (intrvproto.AuthServiceClient, error) {
	conn, err := newConnection(addr)
	if err != nil {
		return nil, err
	}
	return intrvproto.NewAuthServiceClient(conn), nil
}
