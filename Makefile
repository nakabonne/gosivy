
.PHONY: test
test:
	go test -race -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...

.PHONY: mockgen
mockgen:
	mockgen -source=./diagnoser/gui/widgets.go -destination=./diagnoser/gui/widgets.mock.go -package=gui
