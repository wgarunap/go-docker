# Go-docker
Golang project building and deploying in docker environment. Before 
follow these steps you need to have a SSH key with all the repo 
access to your private repository. Glide tool is used for dependency 
management inside the project. Golang Installation steps and RocksDB 
installation steps are copied from DockerHub official images and 
maintained locally for production build purpose.

##### NOTE:
> This base images is build parsing the private key as a build Argument and 
this Key can be visible in the image history. 
`docker history core-gobuild:v1.0 --no-trunc | grep SSH` so not encourage 
to use this same method and key should be generated without a 
authentication key. if you're using swarm services, use docker secret instead.
    
##### Build Builder image "core-gobuild"
```bash
docker build \
    --build-arg SSH_KEY="$(cat id_rsa)" \
    -f Dockerfile.build \
    --no-cache \
    -t core-gobuild:v1.0 .
```

##### Build Production image "core-rocksbuild"
```bash
docker build \
-f Dockerfile.rocks \
--no-cache \
-t core-rocksbuild:v1.0 .
```

##### Build Service
```bash
docker build \
-f Dockerfile \
--no-cache \
-t go-docker:v1.0 .
```

##### Run Service in Detached Mode with Mounted Configs
```bash
docker run -dit \
	-p 10001:10001 \
	-p 20001:20001 \
	-p 20002:6060 \
	--name go-docker-v1.0 \
	--mount type=bind,source=/opt/go-docker/config,target=/opt/go-docker/config \
	go-docker:v1.0
```

###### Check logs in the StdOut
```bash
docker logs go-docker-v1.0
```

###### Stop Running Container
```bash
docker stop go-docker-v1.0
```
###### Remove stopped container 
```bash
docker rm go-docker-v1.0
```