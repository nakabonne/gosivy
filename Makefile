
.PHONY: test
test:
	go test -race -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: mockgen
mockgen:
	mockgen -source=./diagnoser/tui/widgets.go -destination=./diagnoser/tui/widgets.mock.go -package=tui
	mockgen -source=./diagnoser/diagnoser.go -destination=./diagnoser/diagnoser.mock.go -package=diagnoser
