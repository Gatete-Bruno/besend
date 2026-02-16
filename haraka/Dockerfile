FROM node:18-alpine

WORKDIR /app

RUN npm install -g haraka

# Install dependencies for plugins
RUN npm install pg crypto uuid

# Expose SMTP port
EXPOSE 25

# Start Haraka
CMD ["haraka", "-f"]
