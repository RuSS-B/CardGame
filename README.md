# CardGame

## Requirements
Docker, possibility to run *make* files, Go installed on your machine if you want to run tests, a little bit disk space, aaand a good mood :-P 

## Installation
Make sure you have Docker installed on your computer then just run
`make install`. In case you want to run source code just call `make run`

## Testing
To run tests just call the command `make test`

Note #1: if you will run the command manually e.g. go test ./... then make sure you provide APP_ENV=test variable. With this variable the .env.test filed loaded mainly for database connection url purposes

Note #2: Probably in ideal version I would also make a clearTable() function that would run every time when test starts

## Troubleshooting
There could be one potential issue with `make install` script. When loading the schema it's possible that postgres might need more than two seconds to load. If that happens just call `make load_schema` in this case

## Rest API Endpoints:

### Create a full deck
Example: `curl -X "POST" "http://localhost:8080/decks"`

### Create a full shuffled deck
Example: `curl -X "POST" "http://localhost:8080/decks?shuffle=1"`

### Create a partial deck
Example: `curl -X "POST" "http://localhost:8080/decks?cards=AS,AH,AD,KH"`

### Show Deck
Example: `curl "http://localhost:8080/decks/7f6b21a8-4698-432c-b3ba-91c2a73301ad"`

### Draw a card
Example: `curl -X "PATCH" "http://localhost:8080/decks/7f6b21a8-4698-432c-b3ba-91c2a73301ad"`

### Draw N cards
Example: `curl -X "PATCH" "http://localhost:8080/decks/7f6b21a8-4698-432c-b3ba-91c2a73301ad?count=3"`

---------
Cheers!
