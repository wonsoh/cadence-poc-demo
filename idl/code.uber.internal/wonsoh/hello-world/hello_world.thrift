namespace java com.uber.go.helloworld

service HelloWorld {
  HelloResponse hello(1:HelloRequest request);
}

struct HelloRequest {
    1: optional string name;
}

struct HelloResponse {
    1: optional string message;
}
