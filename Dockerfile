FROM node:20.16.0-slim
WORKDIR /app
COPY . src .
RUN npm ci --only=production
EXPOSE 3000
CMD ["node app.js"]
