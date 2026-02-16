FROM node:18-alpine

WORKDIR /app

RUN npm install smtp-server pg

COPY main.js .

EXPOSE 25

CMD ["node", "main.js"]
