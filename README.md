# Before starting
Make sure to follow installations for [Cadence](https://engdocs.uberinternal.com/cadence/get_started/1run_laptop.html) and [DOSA](https://engdocs.uberinternal.com/DOSA/get_started/go/install.html#install-the-cli)

## _Silver-Bullets_
Having problems? Have you tried
```
uber-doctor.sh
update-uber-home.sh
```

# On your separate terminal window
Run Cerberus (in this directory)
```
cerberus
```
Run Cadence docker images
```
## in the folder you've downloaded cadence docker images
docker-compose up
```
And then register your domain (this is a **one-time** step)

```
cadence --do samples-domain domain register
```

Create the DOSA scope and upsert schema (also a **one-time** step, make sure **`cerberus`** is running)
```
## wonsohsandbox201907 is coded in ./config/base.yaml file; change as necessary
dosa scope create wonsohsandbox201907 -o $USER

dosa schema upsert -s wonsohsandbox201907 --namePrefix $USER ./entities
```
# CRUD operations

Run tests to make sure everything is looking good:

```
make test
```

Issue an rpc by first running the service:

```
make run
```


### Create
`./create.sh "NameOfTheAccount"`
As a response you should get a UUID `originalEntity.entityID`
```
{
  "body": {
    "result": {
      "originalEntity": {
        "email": "",
        "entityID": "84cd6678-dc90-4fc7-a372-7458bea0c40d",
        "name": "Signalingtest",
        "nameTS": "2019-07-31T20:05:55-04:00",
        "phone": ""
      }
    }
  },
  "headers": {
    "Content-Length": "136",
    "Content-Type": "application/vnd.apache.thrift.binary",
    "Date": "Thu, 01 Aug 2019 00:05:55 GMT",
    "Rpc-Service": "hello-world",
    "Rpc-Status": "success"
  },
  "statusCode": 200
}
```
Then, visit [http://localhost:8088/domain/samples-domain/workflows](http://localhost:8088/domain/samples-domain/workflows) to see running workflow

### Get
Run `./get.sh <entity_uuid>`. This will get boh entity and SFDC entity (if populated by activities), and fields will be synced by activities.

### Update
#### Update email
`./update_email.sh <entity_uuid> <email_address>`
#### Update phone
`./update_phone.sh <entity_uuid> <phone_number>`
