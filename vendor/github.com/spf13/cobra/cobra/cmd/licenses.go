














package cmd

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)


var Licenses = make(map[string]License)





type License struct {
	Name            string   
	PossibleMatches []string 
	Text            string   
	Header          string   
}

func init() {
	
	Licenses["none"] = License{"None", []string{"none", "false"}, "", ""}

	initApache2()
	initMit()
	initBsdClause3()
	initBsdClause2()
	initGpl2()
	initGpl3()
	initLgpl()
	initAgpl()
}





func getLicense() License {
	
	if userLicense != "" {
		return findLicense(userLicense)
	}

	
	if viper.IsSet("license.header") || viper.IsSet("license.text") {
		return License{Header: viper.GetString("license.header"),
			Text: viper.GetString("license.text")}
	}

	
	if viper.IsSet("license") {
		return findLicense(viper.GetString("license"))
	}

	
	return Licenses["apache"]
}

func copyrightLine() string {
	author := viper.GetString("author")

	year := viper.GetString("year") 
	if year == "" {
		year = time.Now().Format("2006")
	}

	return "Copyright © " + year + " " + author
}




func findLicense(name string) License {
	found := matchLicense(name)
	if found == "" {
		er("unknown license: " + name)
	}
	return Licenses[found]
}





func matchLicense(name string) string {
	if name == "" {
		return ""
	}

	for key, lic := range Licenses {
		for _, match := range lic.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}
