# ============================================================
# Frontend – multi-stage build
#   dev:     next dev with HMR (compose dev target)
#   build:   next build (standalone output)
#   runtime: minimal non-root production image
# ============================================================

# --- Base with dependencies -----------------------------------
FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json* ./
# npm ci needs a lockfile; fall back for first-ever builds.
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

# --- Development (HMR) -----------------------------------------
FROM node:22-alpine AS dev
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=deps /app/node_modules ./node_modules
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev"]

# --- Build -------------------------------------------------------
FROM node:22-alpine AS build
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# --- Runtime ------------------------------------------------------
FROM node:22-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0
COPY --from=build /app/public ./public
COPY --from=build --chown=node:node /app/.next/standalone ./
COPY --from=build --chown=node:node /app/.next/static ./.next/static
USER node
EXPOSE 3000
CMD ["node", "server.js"]
