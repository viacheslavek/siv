all:

.PHONY: run_server
run_server:
	cd ssh_testing/server && go run cmd/server.go

.PHONY: run_client
run_client:
	cd ssh_testing/client && go run cmd/client.go

.PHONY: build_visualizer
build_visualizer:
	echo "build visualizer"

.PHONY: run_visualizer
run_visualizer:
	echo "run visualizer"
