### OpenAPI/Swagger-generated API Reference Documentation

So until those front end devs make a nice UI, this is for us to test and experiment with our backend. 

#### run with locally installed node

If you have node installed, just run:

```
npm install
npm run redoc
```

#### run inside docker
If you have docker installed and running, you can build a node container to run this redoc web UI:

```
./docker_up.sh
```

And you should see a nice UI for interacting with the backend on http://localhost:8080
