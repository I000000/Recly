FROM node:20-alpine AS builder

ARG API_BASE_URL
ENV API_BASE_URL=$API_BASE_URL

WORKDIR /app

COPY package*.json ./

RUN npm ci --omit=dev || npm install --force

RUN npm install lightningcss-linux-x64-musl --force || true

COPY . .
RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app

COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000

CMD ["npm", "start"]