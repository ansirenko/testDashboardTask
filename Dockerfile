FROM golang:1.22-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git && \
    apk add --no-cache chromium \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont \
    udev

WORKDIR /app

COPY . .

RUN go mod download

CMD ["go", "test", "-v", "./..."]
