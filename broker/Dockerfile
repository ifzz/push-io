FROM busybox

RUN mkdir /templates
RUN mkdir /log

ADD ./key.json /key.json
ADD ./config.json /config.json
ADD ./main /main
ADD ./templates /templates/

EXPOSE 8080
CMD ["/main"]
