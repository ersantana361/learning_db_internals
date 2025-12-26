module github.com/ersantana/distributed-systems-learning/apps/api

go 1.25.5

require (
	github.com/ersantana/distributed-systems-learning/packages/protocol v0.0.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
)

replace github.com/ersantana/distributed-systems-learning/packages/protocol => ../../packages/protocol
