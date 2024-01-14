FROM node:14

WORKDIR /usr/src/app

COPY package*.json ./

LABEL type="SQL"

LABEL weight="1"

RUN npm install

COPY . .

EXPOSE 3030

CMD [ "node", "server.js" ]