run:
	go run .

build:
	go build -o bin/todoapp .

clean:
	rm -rf bin/

release:
	fyne package --release -exe todoapp