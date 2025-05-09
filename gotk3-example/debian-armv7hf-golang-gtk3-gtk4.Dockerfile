FROM balenalib/armv7hf-debian-golang:latest
WORKDIR /workspace

RUN apt update 
RUN apt install -y build-essential libgtk-3-dev libgtk-4-dev libgirepository1.0-dev
