namespace java com.uber.go.helloworld

service HelloWorld {
  HelloResponse hello(1:HelloRequest request);

  Response create(1: ModifyRequest request);
  Response update(1: ModifyRequest request);
  Response get(1: GetRequest request);
}

struct HelloRequest {
    1: optional string name;
}

struct HelloResponse {
    1: optional string message;
}

struct Entity {
    1: optional string entityID
    2: optional string name
    3: optional string phone
    4: optional string email
    5: optional string nameTS
    6: optional string phoneTS
    7: optional string emailTS
}

struct Response {
    1: optional Entity originalEntity
    2: optional Entity sfdcEntity
}

struct ModifyRequest {
    1: optional Entity entity
}

struct GetRequest {
    1: required string entityID
}