FROM ubuntu:bionic

WORKDIR /code

# Install bcc tools
RUN apt-get update && \
    apt-get install -y gnupg lsb-release ca-certificates
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 4052245BD4284CDD && \
    echo "deb https://repo.iovisor.org/apt/$(lsb_release -cs) $(lsb_release -cs) main" >> /etc/apt/sources.list.d/iovisor.list
RUN apt-get update && \
    apt-get install -y bcc-tools libbcc-examples linux-headers-$(uname -r)

COPY . .

USER root

# CMD ["sudo", "python", "./tcplife.py"]
# CMD ["sleep", "3000"]
CMD ["sh", "-c", "python ./tcplife.py"]
