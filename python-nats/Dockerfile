FROM python:3.7-alpine

COPY requirements.txt /app/
WORKDIR /app
RUN pip3 install -r requirements.txt
RUN wget https://github.com/shadiakiki1986/nats-cli/releases/download/0.0.4.2/nats-amd64 -O /sbin/nats && chmod +x /sbin/nats

RUN mkdir -p /app
COPY *.py /app/
CMD ["./server.sh"]
