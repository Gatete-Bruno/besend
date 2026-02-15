# ============================================
# BUILD STAGE - Compile Next.js app
# ============================================
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies (including dev)
RUN npm ci

# Copy all source code
COPY . .

# Build Next.js app
RUN npm run build


# ============================================
# PRODUCTION STAGE - Run compiled app
# ============================================
FROM node:18-alpine

WORKDIR /app

# Install next globally
RUN npm install -g next

# Copy built files from builder stage
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package*.json ./

# Expose port
EXPOSE 3000

# Set production mode
ENV NODE_ENV=production

# Start the app
CMD ["npm", "start"]
