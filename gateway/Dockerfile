FROM mhart/alpine-node:latest

WORKDIR /app
ADD . /app/

# If you have native dependencies, you'll need extra tools
# RUN apk add --no-cache make gcc g++ python

# If you need npm, don't use a base tag
RUN npm --registry=https://registry.npm.taobao.org install

EXPOSE 3000
CMD npm start
