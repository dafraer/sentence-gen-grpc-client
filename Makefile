generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/*.proto
run:
	go run cmd/main.go debug
push:
	git add . && git commit -m "$(m)" && git push
build-darwin:
	fyne package -os darwin -icon ../media/darwin_logo.png -src ./cmd -name Sengen --release
	hdiutil create -volname "Sengen" -srcfolder Sengen.app -ov -format UDZO Sengen.dmg
	rm -rf ./Sengen.app
	mv ./Sengen.dmg ./builds/darwin/Sengen.dmg
build-windows:
	fyne package -os windows -icon ../media/logo.png -src ./cmd -name Sengen -app-id com.kamil.sengen --release
	mv ./Sengen.exe ./builds/windows/Sengen.exe
build-linux:
	fyne package -os linux -icon ../media/logo.png -src ./cmd -name Sengen  --release
	mv ./Sengen.tar.xz ./builds/linux/Sengen.tar.xz
