# Stage 1
FROM golang:1.22-alpine AS build

WORKDIR /app 

COPY go.mod go.sum ./ 
RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o db-email .


# Stage 2
FROM alpine:edge

WORKDIR /app 

COPY --from=build /app/db-email .
COPY --from=build /app/.env .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/db-email"]


