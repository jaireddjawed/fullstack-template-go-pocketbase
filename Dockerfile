# syntax=docker/dockerfile:1

# Build the standalone Next.js server. Shared generated types are needed for
# TypeScript's @shared/* imports during the build.
FROM node:22-alpine AS frontend-builder
WORKDIR /workspace

COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci

COPY frontend ./frontend
COPY shared ./shared
RUN cd frontend && npm run build

# Build the Go/PocketBase application.
FROM golang:1.26.2-alpine AS backend-builder
WORKDIR /workspace

RUN apk add --no-cache ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app .

# One container runs the Next.js server and the Go/PocketBase server.
FROM node:22-alpine
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend-builder /out/app ./app
COPY --from=backend-builder /workspace/migrations ./migrations
COPY --from=frontend-builder /workspace/frontend/.next/standalone ./
COPY --from=frontend-builder /workspace/frontend/.next/static ./.next/static
COPY docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
ENV PORT=3000
ENV NEXT_PUBLIC_POCKETBASE_URL=http://127.0.0.1:8090

EXPOSE 3000 8090
VOLUME ["/app/pb_data"]

ENTRYPOINT ["./docker-entrypoint.sh"]
