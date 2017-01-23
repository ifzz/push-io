FROM busybox

RUN mkdir /log

ADD ./config.json /config.json

ADD ./ip.json /ip.json

ADD ./main /main

CMD ["/main"]
