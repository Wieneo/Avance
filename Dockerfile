FROM golang:alpine
WORKDIR /app
COPY . .
RUN apk add --update nodejs nodejs-npm git && rm -rf /var/cache/apk/*
RUN go get && go build -o tixter-app *.go
RUN cd frontend/app/ && npm install && npm run build && rm -rf src/
EXPOSE 8000
CMD [ "./tixter-app", "" ]
