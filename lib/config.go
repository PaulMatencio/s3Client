
package s3Client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"
)

/* sample  json config

{
	"version": "9",
	"hosts": {
		"gcs": {
			"url": "storage.googleapis.com",
			"accessKey": "YOUR-ACCESS-KEY-HERE",
			"secretKey": "YOUR-SECRET-KEY-HERE",
			"api": "S3v2",
			"ssl: true
		},
		"local": {
			"url": "localhost:9000",
			"accessKey": "",
			"secretKey": "",
			"api": "S3v4",
			"ssl": true
		},
		"play": {
			"url": "play.minio.io:9000",
			"accessKey": "Q3AM3UQ867SPQQA43P2F",
			"secretKey": "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			"api": "S3v4",
			"ssl": true
		},
		"s3": {
			"url": "s3.amazonaws.com",
			"accessKey": "YOUR-ACCESS-KEY-HERE",
			"secretKey": "YOUR-SECRET-KEY-HERE",
			"api": "S3v4",
			"ssl": true
		},
		"site1": {
			"url": "10.12.201.11",
			"accessKey": "UIUQR6RMKIX5R2E4FR3M",
			"secretKey": "mjkBoi8imZumGNzkn3rCSDRmZBaAdKSRQhOkqjxE",
			"api": "s3v4",
            "ssl": false
		}
	}
}


*/


type Config struct {
	Version string `json:"version"`
	Hosts    Hosts `json:"hosts"`
}

type Hosts struct {
	Gcs   Host `json:"gcs"`
	Local Host `json:"local"`
	Play  Host `json:"play"`
	S3    Host `json:"s3"`
	Site1 Host `json:"site1"`
}

type Host struct {
	URL			string  `json:"url"`
	AccessKey  	string  `json:"accessKey"`
	SecretKey 	string  `json:"secretKey"`
	Location    string  `json:"location,omitempty"`
	SSL  		bool    `json:"ssl,omitempty"`
	Api  		string  `json:"api,omitempty"`
}

/* config Method */

/*  read the config.json file and parse it into a Go structure */
func GetConfig(c_file string) (Config, error) {
	var (
		config     = path.Join(".s3Client", "config")
		usr, _     = user.Current()
		configfile = path.Join(path.Join(usr.HomeDir, config), c_file)
		cfile, err = os.Open(configfile)
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(100)

	}
	defer cfile.Close()

	decoder := json.NewDecoder(cfile)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	return configuration, err
}

/*  convert teh Hosts structure into a json string */
func Struct2Json( h Hosts) ( string,error){
	if b, err := json.Marshal(h); err == nil {
		return string(b),err
	} else {
		return "", err
	}
}

/*   convert the config struct  into a  map[site]= Host structure  */

func StructToMap( i interface{}) map[string]Host {
	s := reflect.ValueOf(i).Elem()
	 m := make(map[string] Host)
	typ := s.Type()
	for i := 0; i < s.NumField();i++ {
		f := s.Field(i)
		// fmt.Printf("%d: %s %s = %v\n", i, typ.Field(i).Name, f.Type(), f.Interface())
		key := strings.ToLower(typ.Field(i).Name)
		m[key] = f.Interface().(Host)   /*  cast interface into Host structure */
	}
	return m
}



func (c Config) GetVersion() string {
	return c.Version
}

func (c Config) GetHosts() Hosts {
	return c.Hosts
}

/*  Host Methods */

func (h Host) GetUrl() string {
	return h.URL
}

func (h Host) GetAccesKey() string {
	return h.AccessKey
}

func (h Host) GetSecretKey() string {
	return h.SecretKey
}

func (h Host) GetSecure() bool {
	return h.SSL
}

