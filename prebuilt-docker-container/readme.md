# pre-built docker container

docker hub pull command :
docker run -itd -p 80:80 -p 443:443 --name ethsecapi2 ethsectag

this will setup everything for the api, you will need to call the api from another web server or container
any basic https server should be able to access
