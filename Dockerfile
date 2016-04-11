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

# Set the entry point
ENTRYPOINT /go/bin/slack-bot
