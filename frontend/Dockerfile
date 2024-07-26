FROM node:21-alpine

COPY . /app
WORKDIR /app
RUN npm install
RUN npm run build
EXPOSE 3000
CMD ["npm", "run", "dev"]