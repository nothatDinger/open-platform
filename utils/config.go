package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/fredliang44/aliyungo/acm"
	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
)

// Config is a struct for backend config
type Config struct {
	APPName string `default:"Gin App"`
	ShortUrl string `default:"uniqs.cc"`
	Server struct {
		Host      string `default:"127.0.0.1",yaml:"host"`
		Hostname  string `default:"localhost",yaml:"hostname"`
		Port      string `default:"9012",yaml:"port"`
		SecretKey string `default:"SecretKey",yaml:"secretkey"`
	}

	QcloudSMS struct {
		AppID  string `default:"",yaml:"appid"`
		AppKey string `default:"",yaml:"appkey"`
		Sign   string `default:"",yaml:"sign"`
	}

	WeWork struct {
		CropID        string `required:"true",yaml:"cropid"` // CorpID
		AgentID       int    `required:"true",yaml:"agentid"` // Application ID
		AgentSecret   string `required:"true",yaml:"agentsecret"`
		Secret        string `required:"true",yaml:"secret"` // Application Secret
		ContactSecret string `required:"true",yaml:"contactsecret"`
	}

	SMTP struct {
		Sender   string `required:"true",yaml:"sender"`
		Password string `required:"true",yaml:"password"`
		Host     string `required:"true",yaml:"host"`
	}

	Mysql struct {
		User     string `required:"true",yaml:"user"`
		Password string `required:"true",yaml:"password"`
		Host     string `required:"true",yaml:"host"`
		Port     string `required:"true",default:"3306",yaml:"port"`
		Database string `required:"true",yaml:"database"`
	}
}

// LoadConfiguration is a function to load cfg from file
func LoadConfiguration() Config {
	AccessKeyID := os.Getenv("AliAccessKeyID")
	AccessKeySecret := os.Getenv("AliAccessKeySecret")
	AcmEndPoint := os.Getenv("AliAcmEndPoint")
	AcmNameSpace := os.Getenv("AliAcmNameSpace")

	client, err := acm.NewClient(func(c *acm.Client) {
		c.AccessKey = AccessKeyID
		c.SecretKey = AccessKeySecret
		c.EndPoint = AcmEndPoint
		c.NameSpace = AcmNameSpace
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", gin.Mode())
	ret, err := client.GetConfig("studio.open-platform."+gin.Mode(), "STUDIO")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ret)

	var config Config
	err = yaml.Unmarshal([]byte(ret), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Generate secret key
	if count := len([]rune(config.Server.SecretKey)); count <= 32 {
		for i := 1; i <= 32-count; i++ {
			config.Server.SecretKey += "="
		}
	} else {
		config.Server.SecretKey = string([]byte(config.Server.SecretKey)[:32])
	}

	fmt.Printf("%v", config)
	return config
}

// AppConfig is a struct loaded from config file
var AppConfig = LoadConfiguration()

//func Loadlocalfile() Config {
//	var config Config
//	a,err:=ioutil.ReadFile("config.yml")
//	if err!=nil{
//		log.Println(err)
//	}
//	err = yaml.Unmarshal(a,&config)
//	if err!=nil{
//		log.Println(err)
//	}
//	return config
//}
//var AppConfig = Loadlocalfile()