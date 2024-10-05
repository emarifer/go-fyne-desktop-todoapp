run:
	go run .

build:
	go build -o bin/todoapp .

clean:
	rm -rf bin/ db_files/