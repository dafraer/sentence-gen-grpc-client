generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
run:
	go run cmd/main.go debug
push:
	git add . && git commit -m "$(m)" && git push
build darwin:
	 fyne package -os darwin -icon ../media/darwin_logo.png -src ./cmd -name sengen