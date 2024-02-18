## Running postgres on docker on vm
- Run `docker pull postgres` to pull the postgres image, you can also use `docker pull postgres:14` for specific versions.
- Create a docker volume to persist data using `docker volume create postgres_data`
- Start the container `docker run --name postgres_container -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 -v postgres_data:/var/lib/postgresql/data postgres`

reference: https://www.dbvis.com/thetable/how-to-set-up-postgres-using-docker/

# To configure domain for cloud run
- We need to generate our own ssl certificate origin server and private keys.
- This can then be used in load balancer to upload our own certificate and attach it to cloud run.
- refer: https://medium.com/@sidharthvijayakumar7/deploying-cloudrun-application-with-custom-domain-using-cloudflare-5a995f13d5ac
