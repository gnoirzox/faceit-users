# FaceIT User CRUD service

This project includes Go code and related `Docker compose` dependencies to run this project.
This project is a User management micro service using PostgreSQL as a SQL database to store the user details and RabbitMQ for the notification events, the choice of these technologies will be explained in the section 'Asumptions made'.


## How to run this project locally

### Run dependencies

In order to run this project locally, we need to have 2 dependencies running which are PostgreSQL and RabbitMQ.
To make it easier to setup I decided to use Docker compose to run them and to automatically populate the SQL database with the provided sql files from the `/sql` folder.

I also tried to 'Dockerise' my Go project but unfortunately we get an error message when we try to run the Docker image (`standard_init_linux.go:211: exec user process caused "exec format error"`). This seems to be due to the use of some packages such as `net/http` or `crypto/bcrypt` which are relying on C libraries.
I did try to statically build the binary from the Docker container but I did not manage to fix it; I decided to leave the `Dockerfile` in the project for information purpose only.

So, in order to run the service locally, first we need to run `docker-compose up` from the root folder of the project, this will setup both PostgreSQL and RabbitMQ. The Postgres database will be populated automatically. If you need to access to the Database using `psql` or any other tool, bear in mind that the exposed port is `5438`; the attached credentials for the username and password is `postgres` for both; the tables are stored within the database with the same name (`postgres`).

Then, we need to build the project, in order to do so, we need to ensure that the Go modules usage is enabled by running `export GO111MODULE=on` in our shell. To finish, we just have to run `go build -o faceit-users` with `faceit-users` being the output name of our binary and then excute it with `./faceit-users`.

Alternatively, it is also possible to run the project directly from the `/src` folder with the following command: `go run main.go`


### Run tests

I have also added some basic unit and integration tests against the User and Country entities within the `/sr/users` folder.

In order to run them, we have to run `go test`.
Please ensure that the Postgres Docker container is running when executing the Countries integration test.


### Test project locally

This project uses the port `8888`. To check the users, you can go to `http://localhost:8888/users` from your favourite browser. If you want to filter the users, you can input any combination of filters ('nickname', 'email', 'country') like this `http://localhost:8888/users?nickname=CR7&country=PRT`. The particularity of the 'country' filter is that it expects an [ISO 3 characters alpha code](https://www.nationsonline.org/oneworld/country_code_list.htm).

In order to insert a new user, we can use `curl` with the following parameters:

`curl -d '{"nickname":"yanex", "firstname":"yan", "lastname":"shulzr", "email":"r@yandex.com", "password": "zanzibar1", "country": "RUS"}' -H "Content-Type: application/json" -X POST http://localhost:8888/user`

In order to update an existing user:

`curl -d '{"nickname":"yanex", "firstname":"yan", "lastname":"shulzr", "email":"ran@yandex.com", "password": "Zanzibar1", "country": "FRA"}' -H "Content-Type: application/json" -X PUT http://localhost:8888/user/<user-id>`

With `<user-id>` being the user primary key retrieved from the SQL database.

In order to delete an existing user:

`curl -H "Content-Type: application/json" -X DELETE http://localhost:8888/user/<user-id>`

With `<user-id>` being the user primary key retrieved from the SQL database.


### Check notification event messages from the Message Broker (RabbitMQ)

If you want to check the produced messages in the different queues for each event type, you can access it locally within your browser with the following address: `http://localhost:15672`
The credentials are 'guest' for both the username and password.


## Asumptions made

The asumptions made for the database storage have been to use PostgreSQL because it appeared to be a good solution to ensure relationships between a User and a Country. Also the `UNIQUE` constraint on a given field makes it easier to check data integrity on nicknames.
I could also have decided to use a NoSQL database, but I though the use of a SQL database is particualrly relevant if we decide to add more relationships to the user later on (For instance, if we decide to run reports against a given User).

Regarding the User validation, I assumed that we want the password to be of at least 8 characters, also we want to check the email validity and the nickname length (between 3 and 12 characters). Regarding the country, I assumed it was better to retrict its usage with 3 characters long ISO alpha codes so we can check its validity directly in the database against the list of valid countries.

For security purpose, I assumed that the User password needs to be hashed in the Database, to do so I decided to use the `crypto/bcrypt` package because it appears to be the goto method for password hashing as it is very secure.

For the notifications system, I assumed a Message Broker was a good solution because AMQP is a good protocole for the Pub/Sub pattern, so I decided to use RabbitMQ with 3 different queues, one for each event type (insert, update and delete). I could also have used a streaming system like Kafka.


## Areas of improvement

The areas of improvements of this projects for deployment within a production environment and scalability purpose are:

* Add tests against endpoints and SQL queries using either integration tests or mocks.
* Fix the Dockerfile in order to deploy it and also run the tests within a CI/CD pipeline.
* Add metrics for telemetry purpose (eg. use of Prometheus with Grafana dashboard).
* Use of concurrency patterns on SQL queries and AMQP events.
* Improve the `/users` filters SQL query to prevent more clever SQL injection attacks.

