FROM golang:1.20-alpine

WORKDIR /app

COPY . ./

RUN go mod vendor

RUN go build -o /server

EXPOSE 50051
EXPOSE 50052

ENV UPLOAD_FOLDER /upload
RUN mkdir /upload
VOLUME /upload

CMD [ "/server" ]