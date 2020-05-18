FROM alpine
EXPOSE 80
ADD zdrone-build-webhook /app/
ADD conf/ /app/conf

WORKDIR /app/
ENTRYPOINT ["./zdrone-build-webhook"]
