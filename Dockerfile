FROM golang:1.18

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

WORKDIR /opt/app

CMD [ "air" ]
