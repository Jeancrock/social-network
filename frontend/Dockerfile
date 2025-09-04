# Étape 1 : builder
FROM node:20-alpine AS builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# Étape 2 : image finale
FROM node:20-alpine
WORKDIR /app

# Copier node_modules et build depuis le builder
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/.next ./.next

EXPOSE 3000
CMD ["npm", "start"]
