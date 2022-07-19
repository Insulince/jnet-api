run-mongo:
	docker run -d -p 27017:27017 --name mongo --network shared mongo:latest

stop-mongo:
	docker rm -f mongo

build-image:
	docker buildx build --platform linux/amd64 -t jnet-api:latest .

run-image: stop-image
	docker run -d -p 8080:8080 --rm --name jnet-api --network shared jnet-api:latest

stop-image:
	docker rm -f jnet-api

logs:
	docker logs -f jnet-api
