package docker

type Dockerfile struct {
	Node string
	FE   string
}

var DockerfileStr = Dockerfile{
	Node: `FROM node:16-alpine

WORKDIR /app

COPY . .

RUN mkdir config

EXPOSE 3000

CMD cd config && node ../sub-store.bundle.js`,

	FE: `FROM debian:bullseye-slim AS downloader

WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl unzip && \
    rm -rf /var/lib/apt/lists/* && \
    curl -LJO https://sub-store-org.github.io/resource/ssm/nginx.conf && \
    curl -o master.zip -LJ https://github.com/sub-store-org/Sub-Store-Front-End/archive/refs/heads/master.zip && \
    unzip master.zip

FROM node:16-alpine AS builder

WORKDIR /app
COPY --from=downloader /app/Sub-Store-Front-End-master/package.json /app/Sub-Store-Front-End-master/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm i

COPY --from=downloader /app/Sub-Store-Front-End-master/src ./src
COPY --from=downloader /app/Sub-Store-Front-End-master/public ./public
COPY --from=downloader /app/nginx.conf /app/Sub-Store-Front-End-master/.env /app/Sub-Store-Front-End-master/.env.production /app/Sub-Store-Front-End-master/index.html /app/Sub-Store-Front-End-master/tsconfig.json /app/Sub-Store-Front-End-master/tsconfig.node.json /app/Sub-Store-Front-End-master/vite.config.ts ./

RUN pnpm build

FROM nginx:alpine AS runner

WORKDIR /app
COPY --from=builder /app/dist ./www
COPY --from=builder /app/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]`,
}
