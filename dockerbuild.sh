docker inspect relayer -f '{{.Name}}' > /dev/null
if [ $? -eq 0 ] ;then
  echo "container relayer exists, delete container relayer"
  docker rm relayer
fi
docker build -t top/relayer:v1.0 .
docker run --name relayer top/relayer:v1.0
sleep 3
docker cp relayer:/top/xrelayer .
