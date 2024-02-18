ARG NODE_VERSION=18.18.0
FROM node:${NODE_VERSION}-alpine as base

# Node.js app lives here
WORKDIR /app

# Set production environment
ENV NODE_ENV=production

# Throw-away build stage to reduce size of final image
FROM base as build

# to install canvas
RUN apk add build-base \
        g++ \
        cairo-dev \
        jpeg-dev \
        pango-dev \
        freetype-dev \
        giflib-dev \
        make \
        python3

# Install node modules
COPY package-lock.json package.json ./

RUN npm ci --include=dev

# # Copy application code
COPY . .

# Build application
RUN npm run build

# Remove development dependencies
RUN npm prune --production

# Final stage for app image
FROM base
# canvas requires these in runtime
RUN apk add cairo \
    jpeg \
    pango \
    giflib

COPY --from=build /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist
COPY --from=build /app/package.json /app/package.json

EXPOSE 3000

ENV HOST=0.0.0.0
# start production app
CMD [ "node", "dist" ]
