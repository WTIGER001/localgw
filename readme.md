# Local Gateway

Provides a local bridge to Lambda that emulates the API Gateway. The basic idea is that you have a bunch of lambda functions in go and a Swagger file that represents those functions. Then you annotate the swagger file with some vendor extensions. After that you run 
```
localgw generate -s <Path to Swagger file>
```
and copy the outputs into the correct VS Code Launch files

Now you can run the single generated "Compound" launch config to start all the functions. 
Additionally, you can run

```
localgw serve -s <Path to Swagger file> 
```
Starts the proxy on :3333

Thanks to https://github.com/blmayer/awslambdarpc

Note: This is far enough alnog that it works for my use case... let me know if there are other use cases. 

## Assumptions
- Each lambda function runs in its own process that can be debugged, with each process having its own port
- Swagger is used for the services

## Launch Generation
For each method / lambda function there will be the following generated

- Launch Configuration

### Example Launch Configuration
```json
 {
    "name": "${method}-debug",
    "type": "go",
    "request": "launch",
    "program": "${workspaceFolder}/{vscode-path}",
    "env": {
        "_LAMBDA_SERVER_PORT": "${port}"
    },
    "args": []
},
```
### Annotatated Swagger

The Swagger is annotated with the local lamda port
```yaml
x-localgw:
    port: 8001
    vscode-path: lambda/data-api
```

## In Depth Example

This is a simple "hello" example with lambda. 

### Prepare

**Install localgw**
```bash
go install github.com/wtiger001/localgw@latest
```

### Create a new project and configure

```
mkdir hello-lambda
cd hello-lambda

go mod init hello-lambda
go get github.com/aws/aws-lambda-go/events
go get github.com/aws/aws-lambda-go/lambda

```

### Generate a swagger file

in `swagger.yaml`

```yaml
swagger: '2.0'
info:
  description: "Sample API" 
  version: 1.0.0
  title: Sample API
host: YourAWSHOST
basePath: /
tags: 
- name: Example
schemes:
  - https
paths: 
  /hello:
    get:
      tags:
        - Example
      operationId: hello
      consumes:
        - application/json
      produces:
        - application/json
      responses:
        '200':
          description: 200 response            
          type: string
        '500':
          description: 500 response
      x-localgw:
        port: 8001
        vscode-path: lambda/hello
```
Notice the x-localgw structure

### Create a local lambda function
```
mkdir lambda
cd lambda

mkdir hello
cd hello

touch main.go
```
In main.go
```go
package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Credentials"] = "true"

	resp.StatusCode = http.StatusOK
	resp.Body = string("Hello Everyone")

	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}

```

Make sure it compiles
```
go build
```
### Generate the Launch Configurations

```
cd ../../

localgw generate -s swagger.yaml
```

The output should be as shown below

```json

{
   "name": "sample-api/hello:hello",
   "type": "go",
   "request": "launch",
   "program": "lambda/hello",
   "env": {
      "_LAMBDA_SERVER_PORT": "8001"
   }
},

{
   "name": "Run Local Gateway",
   "preLaunchTask": "run-local-gw-task",
   "configurations": [
      "sample-api/hello:hello"
   ]
},

```

### Copy the launch configurations

- Place the first JSON block in your `launch.json` in the `configurations` section
- Places the second JSON block in your `launch.json` in the `compounds` section
- Define a task in your `tasks.json` file like below

```json
 {
    "label": "run localgw",
    "type": "shell",
    "command": "localgw --s swagger.yaml"
},
```

### Start your compound configuration
Go into the debugging menu and start `Run Local Gateway` 

You are now ready to send http requests to `localhost:3333` and you can debug each of the lamda functions. Try that by setting some break points and running

```bash
curl localhost:3333\hello
```
