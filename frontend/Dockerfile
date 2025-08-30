# Stage 1: Build the Angular application
FROM node:20-alpine AS builder

WORKDIR /app

# Update npm to latest version to avoid version conflicts
RUN npm install -g npm@latest

# Copy package.json
COPY package.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application source code
COPY . .

# Generate the build output
RUN npm run build

# Stage 2: Serve the application from a lightweight Nginx server
FROM nginx:alpine

# Copy the build output from the builder stage
# The output folder is defined in angular.json, usually dist/<project-name>
COPY --from=builder /app/dist/spy-cat-agency /usr/share/nginx/html

# Copy the custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]

