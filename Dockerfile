FROM golang:1.14.1-alpine3.11 as wg-builder

RUN apk add --update git build-base libmnl-dev iptables 

ARG wg_go_tag=v0.0.20200320
ARG wg_tools_tag=v1.0.20200319

RUN git clone https://git.zx2c4.com/wireguard-go && \
    cd wireguard-go && \
    git checkout $wg_go_tag && \
    make && \
    make install

ENV WITH_WGQUICK=yes
RUN git clone https://git.zx2c4.com/wireguard-tools && \
    cd wireguard-tools && \
    git checkout $wg_tag && \
    cd src && \
    make && \
    make install

FROM golang:1.14.1-alpine3.11 as app-builder

RUN apk add --no-cache --update upx

COPY --from=wg-builder /usr/bin/wg* /usr/bin/
COPY --from=wg-builder /usr/bin/wireguard-go /usr/bin/

WORKDIR /go/src/wireguard

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go install -v ./...

# # compress binary
# RUN upx -9 bin/*


FROM alpine:3.11

ENV USER=root
ENV HOME=/$USER

RUN apk add --no-cache --update openssh bash iptables openresolv openrc

RUN mkdir $HOME/.ssh \
    && chmod 0700 $HOME/.ssh \
    && ssh-keygen -A \
    && sed -i s/^#PasswordAuthentication\ yes/PasswordAuthentication\ no/ /etc/ssh/sshd_config

RUN echo "$USER:$PASSWORD_TEMP" |chpasswd

COPY --from=app-builder /go/bin/* /usr/bin/wireguard-go /usr/bin/wg* /usr/bin/
# COPY entrypoint.sh /entrypoint.sh

# RUN chmod a+x /usr/bin/wireguard-go

# CMD ["/entrypoint.sh"]