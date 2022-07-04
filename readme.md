# Local Gateway

Provides a local bridge to Lambda that emulates the API Gateway


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





## Usage

### Annotatated Swagger

The Swagger is annotated with the local lamda port
```yaml
    x-localgw:
        port: 8001
        vscode-path: lambda/data-api
        vscode-project: api 
```

```bash
localgw --swagger swagger.yaml --launch 
```

The above will generate the launch file and task file for vscode. The code assumes that the root directory is the 

### Steps
- Read the swagger file(s)
- Build the configuration map
-- method, port, path, project
- Generate the vscode launch configurations for go
- Generate the vscode launch compounds
- Generate the vscode task definitions