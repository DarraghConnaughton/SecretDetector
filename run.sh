docker build -t secretdetector .
docker run -it --rm -v $(pwd)/report:/cmd/report secretdetector