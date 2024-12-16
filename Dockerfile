FROM golang:1.23.3-alpine AS build_stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/app

FROM alpine AS run_stage

WORKDIR /bin

COPY --from=build_stage /bin/app .

RUN chmod +x ./app

#EXPOSE 8080
#
#ENTRYPOINT ["./app"]

EXPOSE 8080

CMD [ "app" ]
