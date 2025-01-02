package rmq_client

import (
	"errors"
)

var (
	ErrConnectionClosed = errors.New("rmq_rpc client - Client - RemoteCall - Connection closed")
)
