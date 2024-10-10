run:
	go run .

build:
	go build -o bin/todoapp .

build-prod:
	go build -ldflags="-s -w" -o ./bin/todoapp -tags prod

clean:
	rm -rf bin/

# Only for releases on Github Actions
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
	mkdir bin
	sh scripts/update_version.sh
	fyne package -os darwin --release -exe ftodo.app --tags prod
	mv ftodo.app bin/
	sh scripts/restore_version.sh