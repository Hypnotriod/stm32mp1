FROM arm32v7/golang:1.24.3-bookworm
WORKDIR /workspace

RUN apt update 
RUN apt install -y build-essential libgtk-3-dev libgtk-4-dev libgirepository1.0-dev
