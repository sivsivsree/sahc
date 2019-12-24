###  Service Availability Health Check -  SAHC

 [![Go Report Card](https://goreportcard.com/badge/github.com/sivsivsree/sahc)](https://goreportcard.com/report/github.com/sivsivsree/sahc)
 

 
 ## About
 
 
  **SAHC** aka **Service Availability Health Check** is a service used to check the availability of services and provide
  much needed actions and events to handle it. **SAHC** is a background runner service which will periodically run to check 
  the status of given services. 
  
  
 ## Features

  -  **SACH** is a stateless application with inbuilt Leveldb for managing the status of services.
  - Provides hot re-loading of configurations.
  - API endpoints to access and update the configurations.
  - YAML support for configurations.
  - TCP streaming for realtime monitoring data and ``healthz`` endpoint for health checks.
  


## Configuration 

### 1. YAML

The example of `YAML` configuration is as follows.  

```
version: 0.1
services:
  - name: localhost:9990
    interval: 10
  - name: localhost:9991
    interval: 11
  - name: localhost:9992
    interval: 16

```

**`services`** describe all the services you need to keep checking.<br>
    
   Services takes in an array and consist of <br>
    
-  `name` : The host address.
- `interval`:  The health check interval in seconds.



> **Note:** <br> 
>  You need to set the env variable `SAHC_CONFIG` to pass the configuration file name.
> If not the service will fail running..


``[working on it..]``