# Dcard backend assignment

## Set up postgresql with docker
```sh
docker run -d --name my-postgres -p 8888:5432 -e POSTGRES_PASSWORD=admin postgres:14-alpine3.17
```

## Compile and run

```sh
make
./bin/a.out
```
or
```sh
go build
./dcard-assignment
```

## Create table
* Since the above container doesn't use volume to store data, we have to initiate table at first.
* Decompose country data to seperate table since the number of countries can be too many.
* These queries will be done while this application start.
```sql
CREATE TABLE Ad (
    ID SERIAL PRIMARY KEY,
    Title text NOT NULL,
    StartAt timestamp NOT NULL,
    EndAt timestamp NOT NULL,
    AgeStart int,
    AgeEnd int,
    Male boolean,
    Female boolean,
    PlatformAndroid boolean,
    PlatformIos boolean,
    PlatformWeb boolean,
);
CREATE TABLE Country (
    ID NOT NULL references Ad(ID),
    Country char(2)
);
```

## Post
```sh
curl -X POST -H "Content-Type: application/json" "http://localhost:3000/api/v1/ad" \
-d '{"title":"AD 55", 
    "startAt":"2023-12-10T03:00:00.000Z",
    "endAt":"2023-12-31T16:00:00.000Z",
    "conditions": {
        
            "ageStart": 20,
            "ageEnd": 30,
            "country": ["TW", "JP"],
            "platform": ["android", "ios"]
        
    }  
}'
```
