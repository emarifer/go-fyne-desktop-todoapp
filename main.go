package main

import "github.com/emarifer/go-fyne-desktop-todoapp/internal/app"

func main() {
	a := app.NewApp()
	defer a.Cleanup()

	a.Run()
}

/* REFERENCES:
https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang

https://stackoverflow.com/questions/71971679/button-action-for-a-specific-list-item-in-fyne

https://stackoverflow.com/questions/66896228/click-event-on-container
https://docs.fyne.io/extend/extending-widgets

Update a collection item given its ID:
https://github.com/ostafen/clover/blob/v2/examples/update/main.go#L32

Advanced Go Build Techniques:
https://dev.to/jacktt/go-build-in-advance-4o8n#iii-build-tags
*/

/* COMMANDS TO BUILD RELEASE:
git tag v1.0.3 && git push origin v1.0.3
go build -ldflags="-s -w" -o ./bin/todoapp -tags=prod
fyne package --release -exe todoapp --tags prod
*/

/* COMMANDS FOR DOCKER-WINE:
# Download the docker-wine script:
wget https://raw.githubusercontent.com/scottyhardy/docker-wine/master/docker-wine
chmod +x docker-wine

# Run the script using the local image and binding a dist/ folder to the docker-wine volume
docker-wine --local --volume=$(pwd)/dist:/home/wineuser/dist

# Accessing the folder where `wineuser` data is stored
cd ~/.wine/drive_c/users/wineuser/
*/
