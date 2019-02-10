ARG SERVICE="hello"
ARG FOLDER="go-docker"

FROM core-gobuild:v1.0 as builder

ARG SERVICE
ARG FOLDER
ENV REPO "git@github.com:wgarunap/go-docker.git"
ENV BRANCH "master"

RUN  mkdir -p /opt/${SERVICE}/config && \
        chown core:core -R $GOPATH/src /opt/${SERVICE}

WORKDIR $GOPATH/src/github.com/wgarunap

RUN apk add --no-cache git
RUN echo $SERVICE
RUN git clone --branch ${BRANCH} --single-branch ${REPO} && cd ${SERVICE} && \
        glide cache-clear && glide install && \
        go build -o ${SERVICE}-linux-amd64 && \
        cp ${SERVICE}-linux-amd64 /opt/booking-engine/${SERVICE}


FROM core-rocksbuild:v1.0

ARG SERVICE
ENV SERVICE=${SERVICE}

EXPOSE 10001/tcp #router endpoint
EXPOSE 20001/tcp #metrics
EXPOSE 6060/tcp

ENV TZ Asia/Colombo
RUN apk add --no-cache tzdata

WORKDIR /opt/booking-engine/${SERVICE}

COPY . .
RUN echo ${SERVICE}
COPY --from=builder /opt/${SERVICE} .
CMD ["sh","-c","chmod a+x /opt/${SERVICE}"]

ENTRYPOINT ["sh","-c","./${SERVICE}-linux-amd64"]