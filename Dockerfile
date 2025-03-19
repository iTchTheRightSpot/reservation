# build angular
FROM node:20-alpine3.19 AS frontend

WORKDIR /app

COPY . .

RUN cd frontend/ && npm i && npm run build --verbose --configuration=production

# build golang
FROM golang:1.24.0 AS  api

WORKDIR /

COPY --from=frontend /app /
RUN go mod download

RUN CGO_ENABLED=0 go build -o salon

# production build
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=api /salon /

CMD ["/salon"]