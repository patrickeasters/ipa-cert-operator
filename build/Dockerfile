FROM alpine:3.9

RUN apk upgrade --update --no-cache && \
    apk add ca-certificates && \
    update-ca-certificates

USER nobody

ADD build/_output/bin/ipa-cert-operator /usr/local/bin/ipa-cert-operator
