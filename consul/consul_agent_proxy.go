package consul

import (
	"errors"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"log"
	"reflect"
	"sync"
	"time"
)

type ServiceInfo struct {
	Id string
	Name string
	Host string
	Port int
	Tags []string
	CheckId string
}

type AgentProxy struct {
	addr string
	dc string
	client *consulApi.Client
	agent *consulApi.Agent
	catalog *consulApi.Catalog
	health *consulApi.Health
	serviceMap map[string]map[string]*ServiceInfo
	serviceMtx *sync.RWMutex
}

func NewAgentProxy(addr string, dc string) (*AgentProxy, error) {
	ap := &AgentProxy{}
	err := ap.init(addr, dc)
	return ap, err
}

func (ap *AgentProxy) init(addr string, dc string) error {
	ap.addr = addr
	ap.dc = dc
	ap.serviceMap = make(map[string]map[string]*ServiceInfo)
	ap.serviceMtx = &sync.RWMutex{}

	config := &consulApi.Config{Address: addr}
	var err error
	ap.client, err = consulApi.NewClient(config)
	if err != nil {
		return err
	}

	ap.agent = ap.client.Agent()
	ap.catalog = ap.client.Catalog()
	ap.health = ap.client.Health()

	return nil
}

func (ap *AgentProxy) getServiceId() string {
	str, _ := uuid.GenerateUUID()
	return str
}

func (ap *AgentProxy) getServiceUniqueKey(name string, host string, port int) string {
	return fmt.Sprintf("%s-%s-%d", name, host, port)
}

func (ap *AgentProxy) RegisterService(name string, host string, port int, tags []string, check *consulApi.AgentServiceCheck) error  {
	uniqueKey := ap.getServiceUniqueKey(name, host, port)
	services, ok := ap.serviceMap[name]
	if !ok {
		services = make(map[string]*ServiceInfo)
		ap.serviceMap[name] = services
	}
	serviceInfo, ok := services[uniqueKey]
	if ok {
		return errors.New(fmt.Sprintf("duplicate service: %s", uniqueKey))
	}
	serviceInfo = &ServiceInfo{
		Id: fmt.Sprintf("%s-%s", name, ap.getServiceId()),
		Name: name,
		Host: host,
		Port: port,
		Tags: tags,
	}
	serviceInfo.CheckId = fmt.Sprintf("chekc_service_%s_%s", name, serviceInfo.Id)
	services[uniqueKey] = serviceInfo

	check.CheckID = serviceInfo.CheckId
	err := ap.agent.ServiceRegister(&consulApi.AgentServiceRegistration{
		ID: serviceInfo.Id,
		Name: serviceInfo.Name,
		Address: serviceInfo.Host,
		Port: serviceInfo.Port,
		Tags: serviceInfo.Tags,
		Check: check,
	})
	return err
}

func (ap *AgentProxy) DeregisterService(name string, host string, port int) error {
	uniqueKey := ap.getServiceUniqueKey(name, host, port)
	services, ok := ap.serviceMap[name]
	if !ok {
		return errors.New(fmt.Sprintf("invalid service: %s", uniqueKey))
	}
	serviceInfo, ok := services[uniqueKey]
	if !ok {
		return errors.New(fmt.Sprintf("invalid service: %s", uniqueKey))
	}
	return ap.agent.ServiceDeregister(serviceInfo.Id)
}

func (ap *AgentProxy) GetService(name string) map[string]*ServiceInfo  {
	ap.serviceMtx.RLock()
	services, ok := ap.serviceMap[name]
	ap.serviceMtx.RUnlock()

	// TODO: 这里可能会有并发问题
	if !ok {
		services = ap.getServicesFromConsul(name)
		ap.serviceMtx.Lock()
		ap.serviceMap[name] = services
		ap.serviceMtx.Unlock()

		go func() {
			timer := time.Tick(time.Second * 5)
			for {
				select {
				case <-timer:
					ap.serviceMtx.RLock()
					oldValue, _ := ap.serviceMap[name]
					ap.serviceMtx.RUnlock()
					newValue := ap.getServicesFromConsul(name)
					if !reflect.DeepEqual(oldValue, newValue) {
						log.Println(fmt.Sprintf("service [%s] update...", name))
						ap.serviceMtx.Lock()
						ap.serviceMap[name] = newValue
						ap.serviceMtx.Unlock()
					}
				}
			}
		}()
	}

	return services
}

func (ap *AgentProxy) getServicesFromConsul(name string) (ret map[string]*ServiceInfo) {
	ret = make(map[string]*ServiceInfo)

	//services, err := ap.agent.ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", Name))
	services, _, err := ap.catalog.Service(name, "", nil)
	if err != nil {
		log.Println("get services from consul error: ", err)
		return
	}
	for _, serviceInfo := range services {
		entry, _, err := ap.health.Service(name, "", true, &consulApi.QueryOptions{Filter: fmt.Sprintf("Service.ID == \"%s\"", serviceInfo.ServiceID)})

		// health, _, err := ap.agent.AgentHealthServiceByID(serviceInfo.ServiceID)
		if err != nil {
			log.Println("get services health from consul error: ", err)
			continue
		}

		if entry == nil || len(entry) == 0 {
			continue
		}

		ret[serviceInfo.ServiceID] = &ServiceInfo{
			Id: serviceInfo.ServiceID,
			Name: serviceInfo.ServiceName,
			Host: serviceInfo.ServiceAddress,
			Port: serviceInfo.ServicePort,
			Tags: serviceInfo.ServiceTags,
		}
	}

	return
}
