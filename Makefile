run:
	go run .

build:
	go build -o bin/todoapp .

build-prod:
	go build -ldflags="-s -w" -o ./bin/todoapp -tags prod

clean:
	rm -rf bin/

# For releases on Github Actions
generate-textfiles:
	sh scripts/generate_license-readme.sh

package-linux:
	sh scripts/update_version.sh
	fyne package -os linux --release -exe ftodo --tags prod
	sh scripts/restore_version.sh

package-windows:
	sh scripts/update_version.sh
	fyne package -os windows --release -exe bin/ftodo.exe --tags prod --appID com.emarifer.ftodo
	sh scripts/restore_version.sh

package-darwin:
	sh scripts/update_version.sh
	fyne package -os darwin --release --tags prod
	sh scripts/restore_version.sh