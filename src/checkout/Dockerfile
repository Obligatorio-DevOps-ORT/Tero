# syntax=docker/dockerfile:1

# ---------- Etapa 1: build ----------
# Tiene TODO (incluidos TypeScript y el CLI de Nest) para compilar .ts -> .js
FROM node:22-alpine AS build
WORKDIR /app
# Primero solo los manifiestos y el Yarn 4 del repo (.yarn/releases):
# Docker cachea la instalación mientras no cambien las dependencias
COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn/releases ./.yarn/releases
# --immutable = instalación reproducible: falla si yarn.lock no coincide con package.json
RUN yarn install --immutable
COPY tsconfig.json tsconfig.build.json nest-cli.json ./
COPY src ./src
RUN yarn build

# ---------- Etapa 2: dependencias de producción ----------
# Instalación limpia SIN devDependencies (sin typescript, @nestjs/cli ni @types)
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json yarn.lock .yarnrc.yml ./
COPY .yarn/releases ./.yarn/releases
RUN yarn workspaces focus --production && yarn cache clean --all

# ---------- Etapa 3: runtime ----------
FROM node:22-alpine
ENV NODE_ENV=production PORT=8080
WORKDIR /app
COPY --from=deps  --chown=node:node /app/node_modules ./node_modules
COPY --from=build --chown=node:node /app/dist ./dist
COPY --chown=node:node package.json ./
# Usuario no-root que ya trae la imagen oficial de Node
USER node
EXPOSE 8080
# node directo (no "yarn serve:prod") para que reciba las señales de parada
CMD ["node", "dist/main.js"]
