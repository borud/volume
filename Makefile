all: app

app:
	@cd cmd/$@ && go build -o ../../bin/$@ --trimpath -tags osusergo,netgo -ldflags="-s -w" 

clean:
	@rm -rf bin app.db

docker-create-volume:
	@docker volume create app-db

docker-build:
	@docker build -t app .

docker-run:
	@docker run -d --name app-background -p 8080:8080 -v app-db:/var/lib/db app
