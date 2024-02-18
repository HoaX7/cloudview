ARG NODE_VERSION=16.9.0
FROM node:${NODE_VERSION}-alpine as base

# Node.js app lives here
WORKDIR /app

# Set production environment
ENV NODE_ENV=production


# Throw-away build stage to reduce size of final image
FROM base as build

COPY package-lock.json package.json ./
RUN npm ci

# # Copy application code
COPY . .

# Final stage for app image
FROM base

COPY --from=build /app /app

EXPOSE 5001

CMD ["npm", "start"]