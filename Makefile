generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/*.proto

push:
	git add . && git commit -m "$(m)" && git push

build-ankiweb:
	cd myaddon && zip -r ../myaddon.ankiaddon *
