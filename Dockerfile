FROM golang:1.14 AS backend
WORKDIR /backend
COPY . .
RUN go get && go build -ldflags "-linkmode external -extldflags -static" -a -o tixter-app *.go

FROM node:14-alpine AS frontend
WORKDIR /frontend
COPY ./frontend/app .
RUN npm install && npm run build && rm -rf src/

FROM alpine
WORKDIR /app
COPY --from=backend /backend/tixter-app /app/
COPY --from=backend /backend/db/migrations /app/db/migrations
COPY --from=backend /backend/userData/sampleData /backend/userData/sampleData
COPY --from=frontend /frontend /app/frontend/app
EXPOSE 8000
CMD [ "/app/tixter-app", "" ]
