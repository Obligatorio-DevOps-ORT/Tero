# syntax=docker/dockerfile:1

# ---------- Etapa 1: build ----------
# Tiene TODO (incluido TypeScript) para compilar .ts -> .js
FROM node:22-alpine AS build
WORKDIR /app
# Primero solo los manifiestos: Docker cachea la instalación
# mientras no cambien las dependencias (igual que go.mod en catalog)
COPY package.json package-lock.json ./
RUN npm ci
COPY tsconfig.json ./
COPY src ./src
RUN npm run build

# ---------- Etapa 2: dependencias de producción ----------
# Instalación limpia SIN devDependencies (sin typescript ni @types)
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev && npm cache clean --force

# ---------- Etapa 3: runtime ----------
FROM node:22-alpine
ENV NODE_ENV=production PORT=8080
WORKDIR /app
COPY --from=deps  --chown=node:node /app/node_modules ./node_modules
COPY --from=build --chown=node:node /app/dist ./dist
COPY --chown=node:node package.json ./
# La app sirve el index.html desde ../public relativo a dist/
COPY --chown=node:node public ./public
# Usuario no-root que ya trae la imagen oficial de Node
USER node
EXPOSE 8080
# node directo (no "npm start") para que reciba las señales de parada
CMD ["node", "dist/app.js"]