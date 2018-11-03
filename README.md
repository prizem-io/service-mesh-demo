# Service Mesh Demo

This is a demostration/PoC of running Prizem [proxy](https://github.com/prizem-io/proxy)/[control plane](https://github.com/prizem-io/control-plane), [Jaeger](https://www.jaegertracing.io), and [Istio Mixer](https://istio.io) to connect 3 different services together.

### Running without the service mesh

##### Start the services

```bash
BACKEND_URI=http://localhost:3001/message go run cmd/frontend/main.go
```

```bash
MESSAGE_URI=http://localhost:3002/sayHello go run cmd/backend/main.go
```

```bash
go run cmd/message/main.go
```

##### Test it!

```bash
curl -s -u demo:demo http://localhost:3000/hello | jq
```

### Running with the service mesh

#### One-time setup

##### Secure communication

The sidecar proxies use TLS to communicate with each other so you will need to place a certificate and key in `etc/backend.cert` and `etc/backend.key`.

##### Setup Postgres database/schema (make sure Postgres running)

```bash
go get -u github.com/golang-migrate/migrate
go get -u github.com/lib/pq
go build -tags 'postgres' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli
./migrate/standup.sh
```

##### Start the control plane

```bash
go run cmd/control-plane/main.go
```

##### Register the services

```bash
./register-services.sh
```

##### Build a local Istio Mixer

```bash
mkdir -p $GOPATH/src/istio.io; cd $GOPATH/src/istio.io
git clone https://github.com/istio/istio.git
cd -
go build -o ./mixs $GOPATH/src/istio.io/istio/mixer/cmd/mixs/main.go
```

#### Running the components

##### Start Jaeger for distributed tracing

```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.6
```

You can then navigate to [http://localhost:16686](http://localhost:16686 "Jaeger UI") to access the Jaeger UI.

##### Start Istio Mixer

```bash
./mixs server --configStoreURL=fs://$PWD/config
```

Or if you want to hack directly on the Mixer source code...

```bash
go run $GOPATH/src/istio.io/istio/mixer/cmd/mixs/main.go server --configStoreURL=fs://$PWD/config
```

##### Start the control plane (if its not already running)

```bash
go run cmd/control-plane/main.go
```

##### Start the sidecar proxies and services

```bash
go run cmd/proxy/main.go -ingressPort=13000 -egressPort=13010 -registerPort=13020
BACKEND_URI=http://localhost:13010/message go run cmd/frontend/main.go
```

```bash
go run cmd/proxy/main.go -ingressPort=13001 -egressPort=13011 -registerPort=13021
MESSAGE_URI=http://localhost:13011/sayHello go run cmd/backend/main.go
```

```bash
go run cmd/proxy/main.go -ingressPort=13002 -egressPort=13012 -registerPort=13022
go run cmd/message/main.go
```

```bash
./register-instances.sh
```

##### Test it!

```bash
curl -s -k -u demo:demo -X GET \
  https://localhost:13000/hello \
  -H 'Accept: application/json; v=1' \
  -H 'Content-Type: application/json; v=1' | jq
```

```bash
hey -m GET -a demo:demo \
  -A 'application/json; v=1' \
  https://localhost:13000/hello
```