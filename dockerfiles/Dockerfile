FROM centos:7

COPY go1.15.5.linux-amd64.tar.gz /root/
COPY golang.conf /root/
COPY main /root/
VOLUME "/etc/rancher/k3s/"
VOLUME "/var/log/"
VOLUME "/root"
RUN tar -zvxf /root/go1.15.5.linux-amd64.tar.gz -C /usr/local \
    && cat /root/golang.conf >> /etc/profile \
    && echo "source /etc/profile" >> /root/.bashrc
# CMD exec ./root/main