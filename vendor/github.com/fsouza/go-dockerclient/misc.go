



package docker

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)




func (c *Client) Version() (*Env, error) {
	resp, err := c.do("GET", "/version", doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env Env
	if err := env.Decode(resp.Body); err != nil {
		return nil, err
	}
	return &env, nil
}




type DockerInfo struct {
	ID                 string
	Containers         int
	ContainersRunning  int
	ContainersPaused   int
	ContainersStopped  int
	Images             int
	Driver             string
	DriverStatus       [][2]string
	SystemStatus       [][2]string
	Plugins            PluginsInfo
	MemoryLimit        bool
	SwapLimit          bool
	KernelMemory       bool
	CPUCfsPeriod       bool `json:"CpuCfsPeriod"`
	CPUCfsQuota        bool `json:"CpuCfsQuota"`
	CPUShares          bool
	CPUSet             bool
	IPv4Forwarding     bool
	BridgeNfIptables   bool
	BridgeNfIP6tables  bool `json:"BridgeNfIp6tables"`
	Debug              bool
	OomKillDisable     bool
	ExperimentalBuild  bool
	NFd                int
	NGoroutines        int
	SystemTime         string
	ExecutionDriver    string
	LoggingDriver      string
	CgroupDriver       string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	OSType             string
	Architecture       string
	IndexServerAddress string
	RegistryConfig     *ServiceConfig
	SecurityOptions    []string
	NCPU               int
	MemTotal           int64
	DockerRootDir      string
	HTTPProxy          string `json:"HttpProxy"`
	HTTPSProxy         string `json:"HttpsProxy"`
	NoProxy            string
	Name               string
	Labels             []string
	ServerVersion      string
	ClusterStore       string
	ClusterAdvertise   string
	Isolation          string
	InitBinary         string
	DefaultRuntime     string
	LiveRestoreEnabled bool
	Swarm              swarm.Info
}




type PluginsInfo struct {
	
	Volume []string
	
	Network []string
	
	Authorization []string
}




type ServiceConfig struct {
	InsecureRegistryCIDRs []*NetIPNet
	IndexConfigs          map[string]*IndexInfo
	Mirrors               []string
}





type NetIPNet net.IPNet



func (ipnet *NetIPNet) MarshalJSON() ([]byte, error) {
	return json.Marshal((*net.IPNet)(ipnet).String())
}



func (ipnet *NetIPNet) UnmarshalJSON(b []byte) (err error) {
	var ipnetStr string
	if err = json.Unmarshal(b, &ipnetStr); err == nil {
		var cidr *net.IPNet
		if _, cidr, err = net.ParseCIDR(ipnetStr); err == nil {
			*ipnet = NetIPNet(*cidr)
		}
	}
	return
}




type IndexInfo struct {
	Name     string
	Mirrors  []string
	Secure   bool
	Official bool
}




func (c *Client) Info() (*DockerInfo, error) {
	resp, err := c.do("GET", "/info", doOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info DockerInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}










func ParseRepositoryTag(repoTag string) (repository string, tag string) {
	parts := strings.SplitN(repoTag, "@", 2)
	repoTag = parts[0]
	n := strings.LastIndex(repoTag, ":")
	if n < 0 {
		return repoTag, ""
	}
	if tag := repoTag[n+1:]; !strings.Contains(tag, "/") {
		return repoTag[:n], tag
	}
	return repoTag, ""
}
