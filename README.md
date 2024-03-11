# Dcard backend assignment

## Set up postgresql with docker
```sh
$ docker run -d --name my-postgres -p 8888:5432 -e POSTGRES_PASSWORD=admin postgres:14-alpine3.17
```

## Compile and run
```sh
$ make
$ ./bin/a.out
```

## Create table
* Since the above container doesn't use volume to store data, we have to initiate table at first.
* Decompose country and platform data to seperate table.
* These queries will be done while this application start.
```sql
CREATE TABLE Ad (
    ID SERIAL PRIMARY KEY,
    Title text NOT NULL,
    StartAt timestamp NOT NULL,
    EndAt timestamp NOT NULL,
    Age int,
    Gender boolean,
);
CREATE TABLE Country (
    ID NOT NULL references Ad(ID),
    Country char(2)
);
CREATE TABLE Platform (
    ID NOT NULL references Ad(ID),
    Platform char(7)
);

```