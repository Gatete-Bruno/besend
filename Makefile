.PHONY: build run clean docker-build docker-push

build:
	go build -o bin/manager main.go

run: build
	./bin/manager

docker-build:
	docker build -t email-operator:v1.0.0 .

docker-push:
	docker tag email-operator:v1.0.0 yourregistry/email-operator:v1.0.0
	docker push yourregistry/email-operator:v1.0.0

clean:
	rm -f bin/manager
