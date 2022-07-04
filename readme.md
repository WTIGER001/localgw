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
