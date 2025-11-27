# todo note commands
# SBD Exercise 7: Setup

I did everything on my Windows laptop.  
We tried to set up a multi-node swarm, but it didn’t work reliably on Windows (network issues), so I just used a single-node swarm on my machine.  

Here are the steps I used.

## 1. Start Swarm
I turned on swarm mode with:

docker swarm init

Then I checked if it worked:

docker node ls

---
## 2. Create Secrets
I made a folder for the secrets:

mkdir docker

Then I created the four secret files:

echo -n "order_user" > docker/postgres_user_secret
echo -n "order_pass" > docker/postgres_password_secret
echo -n "minio_user" > docker/s3_user_secret
echo -n "minio_pass" > docker/s3_password_secret

To check if Docker found them:

docker secret ls

---

## 3.Deploy the Stack
To start everything with the swarm compose file, I ran:

docker stack deploy -c docker-compose.swarm.yml sbd-ex7

I checked if the services are running:

docker service ls

And to see details for a service:

docker service ps sbd-ex7_frontend
docker service ps sbd-ex7_orderservice
docker service ps sbd-ex7_postgres
docker service ps sbd-ex7_minio
docker service ps sbd-ex7_traefik

---

## 4. Access
Traefik exposes port 80, so I opened:

http://localhost

The frontend worked there.

The orders API is available at:

http://orders.localhost

I added "orders.localhost" to my Windows hosts file.

Traefik dashboard:

http://dashboard.localhost

---

## 5. Remove the Stack
To delete everything again:

docker stack rm sbd-ex7
