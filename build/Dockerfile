FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/ipa-cert-operator /usr/local/bin/ipa-cert-operator
