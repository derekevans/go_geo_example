
FROM ubuntu:22.04

RUN apt-get update

# Necessary for installing tzdata non-interactively
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

# Install tools 
RUN apt-get install build-essential software-properties-common cmake wget tzdata git postgresql-client libpq-dev vim -y

# Install Go
RUN mkdir -p /software/go && \
    cd /software/go && \
    wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
    cd /software && \
    rm -r go

# Create ubuntu user, switch to this user, and cd to user home directory
RUN useradd --create-home --shell /bin/bash ubuntu
USER ubuntu
ENV HOME=/home/ubuntu
WORKDIR $HOME
RUN mkdir $HOME/data
RUN mkdir $HOME/go_geo

# Configure Go
ENV PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Create .pgpass and add permission
RUN touch $HOME/.pgpass && \
    chmod 0600 $HOME/.pgpass && \
    echo "postgis:5432:*:postgres:password" >> $HOME/.pgpass

WORKDIR $HOME/go_geo
CMD ["/bin/bash"]