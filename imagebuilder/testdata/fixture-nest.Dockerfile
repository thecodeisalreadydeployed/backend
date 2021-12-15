FROM node:14-alpine as build-env
ADD . /app
WORKDIR /app
RUN yarn install --frozen-lockfile
RUN yarn build

FROM node:14-alpine
WORKDIR /app
COPY --from=build-env /app/package.json /app/yarn.lock ./
COPY --from=build-env /app/node_modules ./node_modules
COPY --from=build-env /app/dist ./
CMD node main
