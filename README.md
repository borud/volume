# tiny dockerized image with volume demo

The instructions below use docker commands.  There are equivalent commands in the Makefile to exemplify how you can structure this in your own projects.

If the SQLite database does not already exist on the volume the container mounts it creates the database and defines the schema.  If a database already exist it is just opened and used.

The application just logs the timestamp in milliseconds since epoch to the database every 5 seconds.  If you access <http://localhost:8080> it will just give you a JSON array of the values that have been written to the database.

Experiment with starting and shutting down the docker container.  Even removing the container and re-creating it.

## Build the docker image

```shell
docker build -t app .
```

## Create volume

Create a volume for your application.  This will be persistent across invocations. This happens on the host machine where you run the docker container.  Think of it like creating a directory.

```shell
docker volume create app-db
```

## Start the container

You can choose to run the container interactively:

```shell
docker run -ti --name app-container -p 8080:8080 -v app-db:/var/lib/db app
```

...or you can run it in the background:

```shell
docker run -d --name app-container -p 8080:8080 -v app-db:/var/lib/db app
```

Note that `app-db` is specified as the volume we are going to use and `/var/lib/db` is the mount point inside the container.  Also note that the syntax

```shell
-v something:/var/lib/db
```

...can be confusing since `something` could refer to your local filesystem or to a volume.  When you use the local filesystem this is called a *bind mount* and you should really make your intentions explicit by using the `--mount` option.
