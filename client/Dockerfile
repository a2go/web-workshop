FROM node:10-slim
WORKDIR /app
RUN apt-get update && apt-get install -y git
COPY package-lock.json package-lock.json
COPY package.json package.json
RUN npm install
COPY . /app/
CMD npm run redoc
