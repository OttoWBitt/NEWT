FROM node:13.3.0 as builder

ARG PROFILE
ENV PROFILE=${PROFILE:-prod}

USER node
WORKDIR /home/node

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build -- --configuration ${PROFILE}

FROM nginx:1.16.1-alpine

RUN rm -rf /usr/share/nginx/html/*
RUN rm /etc/nginx/conf.d/default.conf

COPY ./nginx-conf/nginx.conf.template /etc/nginx/conf.d/nginx.conf.template

COPY --from=builder /home/node/dist/web-nbz /usr/share/nginx/html

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
