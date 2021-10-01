build.alpine:
	GOOS=linux GOARCH=amd64 go build -o bin/battlesnake-alpine .

build.docker: build.alpine
	docker build -t itzamna/cobrakai:latest .

run.docker: 
	docker run -ti --rm -p 8080:8080 --name battlesnake itzamna/cobrakai:latest

push.docker:
	docker push itzamna/cobrakai:latest
