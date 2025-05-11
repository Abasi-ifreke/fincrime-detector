FROM node:lts-alpine AS builder

WORKDIR /app

COPY package*.json yarn.lock* ./

RUN npm install || yarn install

COPY . .

COPY .env .env

RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80