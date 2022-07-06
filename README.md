# auth-backend
This is a backend repo of a full stack authentication app. The frontend repo is [here](https://github.com/suzuka4316/auth-frontend).

## Script
Run `docker compose up`
For testing `docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit`

## Docker
docker-compose.yml is scripted to create 3 containers. One for this backend repo, one for mysql server, and one for phpMyAdmin to check the data stored in the mysql database.
![image](https://user-images.githubusercontent.com/67994527/177035193-3420c4f7-9a72-422a-82a2-8fb2bf62a663.png)
