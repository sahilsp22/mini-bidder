# Mini Bidder
### The problem statement is to create two services Bidder and Controller with the following functionalities:

### Controller: 
- Updates data from Postgres and caches it into Memcached at regular intervals
### Bidder: 
- Exposes endpoints to accept Bid requests processes them with data from cache and responds with appropriate bids in bid eresponse

The Services folder contains the main Bidder and Controller services that run in parallel

## Helpers:
- There are helper services that provide functionalities to the Bidder and Controller such as;
###  Logger:
- It as a wrapper to the logger package that provides logging functionalities
###  Server:
- Provides server functionalities such as handlers http Clients and Serving requests
###  Metrics:
- Provides prometheus metrics and client for creating counters and methods to update them
###  DB:
- Provides Postgres and Memcache Clients and methods to create connections and query and update the data
###  Utils:
- It provides the main functionality of the Controller:
 - Creates a controller client with necessary conncections
 - Starts the controller service with caching creative details and advertiser budgets form Postgres to Memcache at regular intervals
 - Provides method to update advertiser budgets on successful Bid response
###  Bid:
- Provides definition for Bid Request and Response Body
####    Matcher: 
 - Provides the main request handling, matching and response generation logic
 - Parses incoming Bid requests and performs validation
 - Identifies matching creatives from cache and generates Bid response
 - Calls controller for updating budgets for winning advertisers
 - Writes no bid response and HTTP status-No Content if invalid request or no matching creatives
 - Updates prometheus metrics

## NOTE:
- Every Bid response sent is assumed to be a winning response hence, budget updated at every successful match
- Bid price is set based an CPM defined by advertiser in Postgres and budget reduced based on the same
