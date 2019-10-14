FROM node:8.16.2-slim as build
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json /usr/src/app
RUN npm install --silent
COPY . /usr/src/app
RUN npm run build

FROM nginx:1.17.4-alpine
COPY --from=build /usr/src/app/dist /usr/share/nginx/html
EXPOSE 80

COPY docker/cmd.sh .
CMD ./cmd.sh
