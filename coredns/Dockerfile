FROM coredns/coredns:1.1.4

EXPOSE 53 53/udp
VOLUME ["/etc/coredns"]
ENTRYPOINT ["/coredns"]
CMD ["-conf", "/etc/coredns/Corefile"]