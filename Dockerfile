# Note that the docker build command must be running from $GOPATH/src, so that the websocket package can be added below
# e.g. docker build -t slack-bot -f ./github.com/haowang1013/slack-bot/Dockerfile .

FROM golang

RUN go version

# Add this repo
ADD ./github.com/haowang1013/slack-bot /go/src/github.com/haowang1013/slack-bot

# Add websocket explicitly in case it cannot be obtained through regular go get
ADD ./golang.org/x/net/websocket /go/src/golang.org/x/net/websocket

# Run go get to download the dependencies
RUN go get github.com/haowang1013/slack-bot

# Install the app
RUN go install github.com/haowang1013/slack-bot

# Add necessary environment variables
ENV SLACK_API_TOKEN xoxb-14041206193-KaaYu8TdPbWgG6KlActPXZwH
ENV AWS_ACCESS_KEY_ID AKIAJ76KLWK2WJYYQC6A
ENV AWS_SECRET_ACCESS_KEY 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
ENV AWS_DEFAULT_REGION ap-southeast-1
ENV LEANCLOUD_APPID x7bOIHRUXSxNsvbFxiLp8In3-gzGzoHsz
ENV LEANCLOUD_APPKEY zEnstBx0fhfLi4O6Mg2LQR5u
ENV LEANCLOUD_MASTERKEY kdiyHup7fISWXVwF9ndmYlGd

# Set the entry point
ENTRYPOINT /go/bin/slack-bot
