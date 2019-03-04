#Specify the image to use for the container for the bot
FROM golang:latest

#set up environment
ENV SERVICE_NAME cephBot
ENV NAMESPACE github.com/Tkdefender88/
ENV APP /src/${NAMESPACE}/${SERVICE_NAME}/
ENV WORKDIR ${GOPATH}/${APP}

# set the working directory
WORKDIR ${WORKDIR} 

#copy everything to the working directory
ADD . /${WORKDIR}

# compile the bot
RUN go get github.com/bwmarrin/discordgo
RUN go get github.com/fogleman/gg
RUN go build -o CephBot

# run the bot
CMD ["./CephBot"]