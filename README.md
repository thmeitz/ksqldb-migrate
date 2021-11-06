# ksqldb-migrate

Migration tool for [ksqlDB](https://ksqldb.io), which uses the [ksqldb-go](https://github.com/thmeitz/ksqldb-go) client.

## Installation

```bash
go install github.com/thmeitz/ksqldb-migrate
```

## Usage

### Create the migration `yaml` file

In further steps I'm using the [provided example](examples/test.yaml).

```yaml
---
up:
  - name: create source connector dogs
    exec: |
      CREATE SOURCE CONNECTOR DOGS WITH ('connector.class'='io.mdrogalis.voluble.VolubleSourceConnector',
      'key.converter'='org.apache.kafka.connect.storage.StringConverter',
      'value.converter'='org.apache.kafka.connect.json.JsonConverter',
      'value.converter.schemas.enable'='false',
      'genkp.dogs.with'='#{Internet.uuid}',
      'genv.dogs.name.with'='#{Dog.name}',
      'genv.dogs.dogsize.with'='#{Dog.size}',
      'genv.dogs.age.with'='#{Dog.age}',
      'topic.dogs.throttle.ms'=1000 
      );
  - name: create the dogs stream
    exec: |
      CREATE STREAM IF NOT EXISTS DOGS (ID STRING KEY,NAME STRING,DOGSIZE STRING, AGE STRING) 
      WITH (KAFKA_TOPIC='dogs', 
      VALUE_FORMAT='JSON', PARTITIONS=1);
  - name: create the DOGS_BY_SIZE table
    exec: |
      CREATE TABLE IF NOT EXISTS DOGS_BY_SIZE AS 
      SELECT DOGSIZE AS DOG_SIZE, COUNT(*) AS DOGS_CT 
      FROM DOGS WINDOW TUMBLING (SIZE 15 MINUTE) 
      GROUP BY DOGSIZE;
down:
  - name: delete stream DOGS
    exec: |
      DROP STREAM IF EXISTS DOGS;
  - name: drop table DOGS_BY_SIZE
    exec: |
      DROP TABLE IF EXISTS DOGS_BY_SIZE
  - name: drop connector DOGS
    exec: |
      DROP CONNECTOR IF EXISTS DOGS;
```

## Migrate Up

```bash
ksql-migrate up -f examples/test.yaml
```

Output

```bash
INFO[0000] processing           name="create source connector dogs" step=1
INFO[0000] preflight check      name="create source connector dogs" status=ok step=1
INFO[0000] processed            name="create source connector dogs" status=ok step=1
INFO[0000] processing           name="create the dogs stream" step=2
INFO[0000] preflight check      name="create the dogs stream" status=ok step=2
INFO[0000] processed            name="create the dogs stream" status=ok step=2
INFO[0000] processing           name="create the DOGS_BY_SIZE table" step=3
INFO[0000] preflight check      name="create the DOGS_BY_SIZE table" status=ok step=3
INFO[0000] processed            name="create the DOGS_BY_SIZE table" status=ok step=3
```

## Migrate Down

```bash
ksql-migrate down -f examples/test.yaml
```

Done.

## TODO

- [ ] Better error messages for KSQLParser (ksqldb-go)
- [ ] Check CommandStatus and wait for successfull execution for each step (up/down migrations) (ksql-migrate + ksqldb-go)

## License

[Apache License Version 2.0](LICENSE)
