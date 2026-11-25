FROM node:22-alpine AS builder
WORKDIR /src/app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=builder /src/app/dist /usr/share/nginx/html
