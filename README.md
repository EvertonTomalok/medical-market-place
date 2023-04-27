## Propose 

The context is the following: The system powers a marketplace app for healthcare facilities to hire healthcare professionals (a.k.a. workers).

The system will be responsible for the shift eligibility feature.
 
Shift eligibility is a feature that allows you to know what shifts are eligible for a specific Worker in specific facilities.

The entities that come into play are the following, `Shift`, `Facility`, `Worker`, `Document`, `FacilityRequirement`, and `DocumentWorker`.

The user story: As a worker, I want to get all available shifts that I'm eligible to work for.  

### Debrief:

 - In order for a Worker to be eligible for a shift, the rules are:
	 - A Facility must be active.
	 - The Shift must be active and not claimed by someone else.
	 - The Worker must be active.
	 - The Worker must not have claimed a shift that collides with the shift they are eligible for.
	 - The professions between the Shift and Worker must match.
	 - The Worker must have all the documents required by the facilities.
 - Given an active facility, when I request all available shifts within a start and end date, then it will return a list of shifts from that facility in the specified date range.
 - Given an inactive facility, when I request all available shifts within a start and end date, then it will not return a list of shifts from that facility.
 - Given a shift is claimed and is within the requested start and end date, when I request all available shifts within a start and end date, it will not return the claimed shift.
 - The shifts must be grouped by date.

![flow](./docs/images/medical-api-flow.jpeg)

----------------------------------------------

The seed script will provide a PostgreSQL database with some random data such that:

 - Some shifts are claimed.
 - Some workers are inactive.
 - Some facilities are inactive
 - Some workers don’t meet all the documents a facility requires.

## Setup:

### Seeding your database

The folder called `seed` contains a docker-compose.yaml file that helps you set up a database. It is a PostgreSQL database and it is seeded with about 2 million records. 

To set it up go into the `seed` folder and execute the command `docker compose up --build`. Once it's done, do not stop the docker compose.  This way the database keeps running and you can have your framework of choice to connect it to the database using the following database URL `postgres://postgres:postgres@localhost:5432/postgres`.

## Running the project

Start the docker file, at the root of the project, using `docker-compose up`. After this, keep it running, and access the endpoint [http://0.0.0.0:8000/doc/index.html](http://0.0.0.0:8000/doc/index.html), and you can check all routes available. 