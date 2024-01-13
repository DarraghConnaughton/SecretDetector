#Build and run
docker build -t secretdetector .
docker run -it --name secretdetector-container secretdetector

#Copy report file
docker cp secretdetector-container:/cmd/report.json ./report.json

#Cleanup
docker stop secretdetector-container
docker rm secretdetector-container
