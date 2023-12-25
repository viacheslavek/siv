all:

.PHONY: run_server
run_server:
	cd ssh_testing/server && go run cmd/server.go

.PHONY: run_client
run_client:
	cd ssh_testing/client && go run cmd/client.go

.PHONY: build_visualizer
build_visualizer:
	"/Users/slavaruswarrior/Applications/CLion 2023.2.2.app/Contents/bin/cmake/mac/bin/cmake" --build /Users/slavaruswarrior/Documents/GitHub/siv/visualizer/cmake-build-debug --target ConsoleGraph -j 6

.PHONY: run_visualizer
run_visualizer:
	./visualizer/cmake-build-debug/ConsoleGraph < ./visualizer/visfifo
