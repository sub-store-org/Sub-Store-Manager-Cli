package lib

const dockerfileStrNode = `FROM node:16-alpine

WORKDIR /app

COPY . .

RUN mkdir config

EXPOSE 3000

CMD cd config && node ../sub-store.bundle.js`
