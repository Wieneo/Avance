FROM golang:alpine
WORKDIR /app
COPY . .
RUN apk add --update nodejs nodejs-npm git
RUN go get
RUN cd frontend/app/ && npm install
RUN cd frontend/app/ && npm run build
RUN mv frontend/app/dist frontend/
RUN rm -rf frontend/app/*
RUN rm -rf frontend/app/.*
RUN mv frontend/dist frontend/app/
RUN go build -o tixter-app *.go
EXPOSE 8000
CMD [ "./tixter-app", "" ]
