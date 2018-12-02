
# Run locally
```bash
go run main.go
```

# Docker
Build:
```bash
docker build -t security-duck .
```

Run:
```bash
$ docker run -p 8080:80 --name security-duck --rm security-duck
```

Run with the `static` directory mounted:
```bash
$ docker run -p 8080:80 --name security-duck --rm --mount type=bind,source="$(pwd)"/static,target=/home/app/static security-duck
```

Export image:
```bash
$ docker save security-duck:latest -o image.tar
```

Import image:
```bash
$ docker load -i image.tar
$ docker image ls
$ docker-compose up -d
```
