# !/bin/sh
curl -X PUT \
  http://localhost:8000/v1/routes \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: a07c5be2-1426-4354-8fe9-ff38f2cf69e8' \
  -d '{
  "name": "frontend",
  "description": "Demo frontend service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "hello",
      "method": "GET",
      "uriPattern": "/hello",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'

curl -X PUT \
  http://localhost:8000/v1/routes \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 61c03672-cd0b-4e61-bba8-10be9b6273fd' \
  -d '{
  "name": "backend",
  "description": "Demo backend service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "message",
      "method": "GET",
      "uriPattern": "/message",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'

curl -X PUT \
  http://localhost:8000/v1/routes \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: f6d038d7-0537-4404-bb97-9c59aaf07acd' \
  -d '{
  "name": "message",
  "description": "Demo message service",
  "version": {
    "versionLocations": [
      "content-type"
    ],
    "defaultVersion": "v1"
  },
  "authentication": "none",
  "timeout": "10s",
  "operations": [
    {
      "name": "sayHello",
      "method": "GET",
      "uriPattern": "/sayHello",
      "timeout": "10s",
      "retry": {
        "attempts": 3,
        "perTryTimeout": "2s"
      },
      "policies": []
    }
  ]
}'