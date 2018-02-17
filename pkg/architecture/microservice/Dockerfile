# From blog: https://medium.com/@kelseyhightower/optimizing-docker-images-for-static-binaries-b5696e26eb07
# replace contributors with your binary
FROM alpine:latest
LABEL maintainer="latschesar.atanassov@gmx.at"
ADD contributors contributors
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/contributors"]