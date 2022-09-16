FROM golang:1.19 as build-env

RUN update-ca-certificates

# Create appuser
ENV USER=controlplane
ENV UID=10001

# configure user in the environment to not allow shell access. 
# this is later copied into the actual container
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# create working directory for build environment
WORKDIR /go/src/app
COPY go.mod .
# using wildcard requires setting a path
COPY *.go /go/src/app/

RUN go get -d -v ./...

# don't need CGO
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

# set the user created above in the application container
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /go/bin/app /

USER controlplane:controlplane

CMD ["/app"]
