# Caching values

In a real system, we could have DNS caching and inside the system, we could implement cache using some tool like Redis, caching responses, and other relevant data. 

It`s also possible to change the core of search, input, and delete jobs, to Redis, where some workers will be tracking the consistency between data, and aggregating shifts. As a complementary step, every operation that changes the status of the shifts (include, delete, update), should also update shifts in Redis.

# Creating paginators for the endpoints

The actual architecture, won't work very well for a thousand of data requested. As a good practice, we should create paginators, and this way, limit the quantity of data shown for each request to the user.

# Creating E2E tests

The next step, I'd consider creating end-to-end tests, using some tool like Selenium, or my preferred way to create tests for backend applications, is creating BDD tests (Behavior Driven Development), such as the example below:

```
Feature: Get Shifts

    Scenario: Receiving requests information, retrieving shifts
        Given A set of information
            |   user_id   |     start    |        end        | offset |  limit   |
            |       1     | '2023-04-23' |   '2023-04-24'    |    0   |   100    |
            |       1     | '2023-04-23' |   '2023-04-23'    |  100   |   100    |

        When making request
        Then check all shifts are active
        Then check all facilities are active
        Then check no one shift is taken for someone or conflict with the user
``` 

# Creating LOAD test

As the next step, I also would consider creating load tests to evaluate the performance of the system. I like to use [Locust](https://locust.io/) (a load testing framework) to create automatic tests and collect metrics about the ecosystems of the application, like throughput, average time response, etc.