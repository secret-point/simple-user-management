
FROM node:20.8.1 as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
ENV NODE_OPTIONS=--max-old-space-size=4096
COPY . .
RUN npm run build

FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
