
basic structure 
https://swagger.io/docs/specification/v3_0/basic-structure/

==> You can run Swagger UI locally.

Step 1

Download Swagger UI
```
git clone https://github.com/swagger-api/swagger-ui
```
Step 2
```
cd swagger-ui
```
Step 3

Edit:
```
dist/swagger-initializer.js
```
Change:
```
url: "https://petstore.swagger.io/v2/swagger.json"
```
to: your yaml file
```
url: "openapi.yaml"
```

if there are mulitple yaml fils
```
urls: [
    {
        "name": "example1"
        "url": "./openapi/example1.openapi.yaml"
    },
    {
        "name": "example2"
        "url": "./openapi/example2.openapi.yaml"
    }
]
```

Step 4

Run simple server
```
python -m http.server 8080
```
Open:
```
http://localhost:8080
```

Download tool to generate code
```
> sudo apt update && sudo apt install openjdk-17-jre-headless -y
> sudo apt install make
> wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.25.0/openapi-generator-cli-7.25.0.jar -O openapi-generator-cli.jar
```