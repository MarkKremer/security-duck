# This dockerfile makes use of multi-stage builds to keep the resulting image
# size as small as possible. The golang image is used for compilation and only
# the resulting files are then copied to the minimal Alpine linux image which
# is used in the output image to run the compiled go executable.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds

FROM golang:1.8 as builder
WORKDIR /go/src/markkremer.nl/security-duck/
# Copy every go file in the current directory to the working directory
# in the docker image
COPY *.go .
# Install all dependencies of the project. go get will automatically
# scan the project files for imports
RUN go get -d -v
# Compile the project into the executable 'app'. Statically include
# dependencies on C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .


# Switch to the alpine base image if you need TSL certificates or need
# to use the package manager for something else
# FROM scratch
FROM alpine:latest
# Change the user to a new non-root user for better security
RUN mkdir -p /home/app &&\
    addgroup -S app &&\
    adduser -S -G app -h /home/app -s /sbin/nologin app
# libcap is required to allow app to bind to ports under 1024 as a non-root user
RUN apk add libcap
# Uncomment the following line if the project needs to validate TSL certificates
#RUN apk --no-cache add ca-certificates
# Bash can be usefull for debugging
#RUN apk add --no-cache bash
WORKDIR /home/app/
# Copy the compiled app from the golang image to this one
COPY --from=builder /go/src/markkremer.nl/security-duck/app .
# Copy all static files
COPY static/ static/
# Chown all the files to the app user
RUN chown -R app:app /home/app
# Allow app to bind to ports under 1024 as a non-root user
RUN setcap CAP_NET_BIND_SERVICE=+ep app
EXPOSE 80

USER app
#ENTRYPOINT ["./app"]
CMD ["./app"]
