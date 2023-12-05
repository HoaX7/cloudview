FROM alpine
RUN apk add --update nodejs npm

WORKDIR /app
COPY . .
RUN npm install
EXPOSE 80
EXPOSE 8080

CMD [ "node", "server.js" ]