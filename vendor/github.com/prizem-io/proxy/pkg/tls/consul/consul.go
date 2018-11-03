package consul

import (
	"context"
	"crypto/tls"
	"net"
	"sync"

	agentconnect "github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"github.com/pkg/errors"

	proxyapi "github.com/prizem-io/api/v1"
	proxytls "github.com/prizem-io/proxy/pkg/tls"
)

type Connect struct {
	client     *api.Client
	servicesMu sync.RWMutex
	services   map[string]*connect.Service
}

type Service struct {
	client  *api.Client
	service *connect.Service
}

func New(client *api.Client) *Connect {
	return &Connect{
		client:   client,
		services: make(map[string]*connect.Service),
	}
}

func (c *Connect) Close() error {
	var err error
	for _, service := range c.services {
		e := service.Close()
		if err == nil && e != nil {
			err = e
		}
	}
	return err
}

func (c *Connect) GetService(serviceName string) (proxytls.Service, error) {
	c.servicesMu.RLock()
	service, ok := c.services[serviceName]
	c.servicesMu.RUnlock()

	if !ok {
		var err error
		c.client.Connect()
		service, err = connect.NewService(serviceName, c.client)
		if err != nil {
			return nil, errors.Wrapf(err, "could not create service %s", serviceName)
		}
		c.servicesMu.Lock()
		c.services[serviceName] = service
		c.servicesMu.Unlock()
	}

	return &Service{
		client:  c.client,
		service: service,
	}, nil
}

func (s *Service) ServerTLSConfig() *tls.Config {
	return s.service.ServerTLSConfig()
}

func (s *Service) Dial(service *proxyapi.Service, node *proxyapi.Node, address string) (net.Conn, error) {
	leafCert, _, err := s.client.Agent().ConnectCALeaf("test", &api.QueryOptions{
		Datacenter: node.Datacenter,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve leaf cert for service %q", service.Name)
	}

	certURI, err := agentconnect.ParseCertURIFromString(leafCert.ServiceURI)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse cert URI: %s", service.Namespace)
	}

	return s.service.Dial(context.Background(), &connect.StaticResolver{
		Addr:    address,
		CertURI: certURI,
	})
}
