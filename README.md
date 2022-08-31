# Special Umbrella


## How it works?
The following things happen when you start the service:
* The http api starts serving http requests
* The seeder adds some scooters and *scooter_status* to the database. Berlin is the location of all generated locations.
* In parallel, three simulators ride all scooters available. By calling **Find Scooters** API, the simulator finds available scooters, then reserves each scooter by calling **Reserve Scooter** API. If the reserve is successful, status updates are sent to the server. The **Publish Scooter Status** API receives the following status updates every three seconds during each ride:
    - trip-started 
    - trip-update (3 times )
    - trip-ended 
    - periodic-update

## Technical Decisions
DB: MongoDB is used to store data and $box query is used to execute rectangular queries. Additionally, it was possible to use server-side mathematics to handle the query, but mongoDB seems a reasonable choice for more complex geospatial queries in the future.


Tests: In order to run tests, a separate docker compose is used to start the required infrastructure. By considering time, I believe integration tests and end-to-end tests make more sense for this project, so I preferred them over unit tests.

Security: static JWT_Token and JWT_KEY is used for securing the API for now.

Build: Makefile is provided to make life easier during run and build.

HTTP API validation: I used go-playgorund to validate requests.

## Commands
To execute tests run:
```sh
make test
```

To start the service run:
```sh
make up
```

To stop the service run:
```sh
make stop
```

To stop the service and remove dockers run:
```sh
make down
```

To tidy up go.mod & syng vendors run:
```sh
make mod
```

## Authentication
All APIs are protected so you need to add **Authorization** header with "**Bearer JWT_Token**" to the request's header. You can find the JWT_Token in the .evn file.
## Exposed APIs
API documentations are availble online here:
https://documenter.getpostman.com/view/14834296/VUqmwzhr

#### Note: Postman collection is also available in the repo. You can import it in your Postman and call the APIs.
