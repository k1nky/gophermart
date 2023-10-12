FROM alpine

ADD accrual_linux_amd64 /opt

RUN apk add libc6-compat

EXPOSE 8080

ENTRYPOINT [ "/opt/accrual_linux_amd64" ]