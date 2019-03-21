# This Dockerfile builds off the Oracle Dockerfile for Oracle Instant Client 18.3.0
# which can be found at
# https://github.com/oracle/docker-images/blob/master/OracleInstantClient/dockerfile/18.3.0/Dockerfile
#
# HOW TO BUILD THIS IMAGE
# -----------------------
# Run:
#      $ docker build -t go/oracleadw .

FROM oraclelinux:latest

# install oracle instant client 18.3
RUN yum -y install oracle-release-el7 oracle-golang-release-el7 && \
    yum -y install oracle-instantclient18.3-devel oracle-instantclient18.3-sqlplus && \
    echo /usr/lib/oracle/18.3/client64/lib > /etc/ld.so.conf.d/oracle-instantclient18.3.conf && \
    ldconfig

# Install golang 12
RUN yum install -y git golang

# Add instant client and Go to PATH
ENV PATH=$PATH:/usr/lib/oracle/18.3/client64/bin
ENV TNS_ADMIN=/usr/lib/oracle/18.3/client64/lib/network/admin
ENV GOPATH=$HOME/go

# Create directory for app files
WORKDIR $GOPATH/src/App

# Add wallet files
ADD ./wallet /usr/lib/oracle/18.3/client64/lib/network/admin/

# Add go files for api
ADD ./app/* .

# Install goracle library for go to connect to db
RUN go get -d -v ./... && \
go install -v ./...

# Listen on port 8000
EXPOSE 8000

# Start app logic
CMD go run app/main.go
