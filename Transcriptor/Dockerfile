FROM python:3.11-slim-buster

WORKDIR /transcriptor

RUN apt-get update
RUN apt-get install -y git
RUN pip install --no-cache --upgrade pip setuptools
RUN pip install pytube
RUN pip install git+https://github.com/openai/whisper.git
RUN apt-get install -y ffmpeg --fix-missing

COPY . /transcriptor

CMD ["python3", "main.py"]