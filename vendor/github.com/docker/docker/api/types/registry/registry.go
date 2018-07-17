package registry 

import (
	"encoding/json"
	"net"

	"github.com/opencontainers/image-spec/specs-go/v1"
)


type ServiceConfig struct {
	AllowNondistributableArtifactsCIDRs     []*NetIPNet
	AllowNondistributableArtifactsHostnames []string
	InsecureRegistryCIDRs                   []*NetIPNet           `json:"InsecureRegistryCIDRs"`
	IndexConfigs                            map[string]*IndexInfo `json:"IndexConfigs"`
	Mirrors                                 []string
}



type NetIPNet net.IPNet


func (ipnet *NetIPNet) String() string {
	return (*net.IPNet)(ipnet).String()
}


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
	
	Name string
	
	Mirrors []string
	
	
	
	Secure bool
	
	Official bool
}


type SearchResult struct {
	
	StarCount int `json:"star_count"`
	
	IsOfficial bool `json:"is_official"`
	
	Name string `json:"name"`
	
	IsAutomated bool `json:"is_automated"`
	
	Description string `json:"description"`
}


type SearchResults struct {
	
	Query string `json:"query"`
	
	NumResults int `json:"num_results"`
	
	Results []SearchResult `json:"results"`
}



type DistributionInspect struct {
	
	
	Descriptor v1.Descriptor
	
	
	Platforms []v1.Platform
}
