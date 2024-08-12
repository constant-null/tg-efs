FROM node:22-alpine
WORKDIR /usr/app
COPY ./ /usr/app

RUN npm install
RUN npm run inline

FROM golang:latest
WORKDIR /usr/app
COPY ./ /usr/app
COPY --from=0 /usr/app/sheet/sheet.min.html /usr/app/sheet/sheet.min.html
RUN CGO_ENABLED=0 GOOS=linux go build -o tg-efs

FROM scratch
COPY --from=1 /usr/app/tg-efs /bin/tg-efs
CMD ["/bin/tg-efs"]

