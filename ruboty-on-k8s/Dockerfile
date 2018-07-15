FROM ruby:2.5.1-alpine3.7
RUN apk update && \
    apk add build-base openssl

ADD . /opt/ruboty

WORKDIR /opt/ruboty
RUN bundle install --path vendor/bundle

ENTRYPOINT bundle exec ruboty
