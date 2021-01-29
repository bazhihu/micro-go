package main

/**
读取配置文件demo
*/

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"reflect"
	_ "github.com/spf13/viper/remote"
	"sync"
	"time"
)

// 全局配置信息
var Resume ResumeInformation

// 个人简历
type ResumeInformation struct {
	Name   string
	Sex    string
	Age    int
	Habits []interface{}
}

type ResumeSetting struct {
	TimeStamp         string
	Address           string
	ResumeInformation ResumeInformation
}

func init() {
	viper.AutomaticEnv()

	// 设置读取的配置文件
	//viper.SetConfigName("resume-config")
	// 添加读取的配置文件路径
	//viper.AddConfigPath("./config/")
	// 设置配置文件类型
	//viper.AddConfigPath("$GOPATH/src/")
	// 设置配置文件类型
	//viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/resume-config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	if err := sub("ResumeInformation", &Resume); err != nil {
		log.Fatal("Fail to parse config", err)
	}
	watchConfig()
}

func Contains(obj, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in array")
}

func parseYaml(v *viper.Viper) {
	var resumeConfig ResumeSetting
	if err := v.Unmarshal(&resumeConfig); err != nil {
		fmt.Printf("err:%s", err)
	}
	fmt.Println("resume config:\n", resumeConfig)
}

// 反序列化赋值
func sub(key string, value interface{}) error {
	log.Printf("配置文件的前缀为: %v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

// 监控配置文件变化实时修改配置文件
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
}


// 从IO.Reader读取配置文件
func IOReader()  {
	viper.SetConfigType("yaml")
	var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
	`)
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	log.Println("name", viper.Get("name"))

}

// 接收参数
func flagsTest()  {
	// Flags
	pflag.Int("flagname", 1234, "do")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	log.Println(viper.GetInt("flagname"))
}

// 利用Viper读取本地配置信息
func readLocal()  {
	parseYaml(viper.GetViper())
	fmt.Println(Contains("Basketball", Resume.Habits))
	for {
		time.Sleep(5 * time.Second)
		parseYaml(viper.GetViper())
		//fmt.Printf("性别:%s\n 爱好:%s\n性别:%s\n年龄:%d\n", Resume.Name, Resume.Habits, Resume.Sex, Resume.Age)
	}
}

// 设置新变量
func setOneTest()  {
	viper.Set("verbose", true)
	log.Println(viper.GetBool("verbose"))
}

type RuntimeConf struct {
	Master struct{
		Db_dsn string `json:"db_dsn"`
		Max_open int `json:"max_open"`
		Max_idle int `json:"max_idle"`
		Db_name string `json:"db_name"`
	} `json:"master"`
	Slave struct{
		Db_dsn string `json:"db_dsn"`
		Max_open int `json:"max_open"`
		Max_idle int `json:"max_idle"`
		Db_name string `json:"db_name"`
	} `json:"slave"`
}

// ETCD 远程监控 Key/Value存储示例-w未加密
func main() {
	// 创建一个新的viper实例
	var (
		runtime_viper = viper.New()
		conf RuntimeConf
		wc sync.WaitGroup
	)

	runtime_viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/mysql/go_base_center")
	runtime_viper.SetConfigType("yaml")
	//viper.AddRemoteProvider("etcd", "http")

	// 第一次从远程读取配置
	if err := runtime_viper.ReadRemoteConfig(); err != nil {
		log.Println("ReadRemoteConfig-err", err)
	}



	// 反序列化
	runtime_viper.Unmarshal(&conf)

	wc.Add(1)
	// 开启单独的goroutine 一直监控远端的变更
	go func() {
		defer wc.Done()
		for {
			time.Sleep(time.Second * 5) // 每次请求后延迟一下

			if err := runtime_viper.WatchRemoteConfig(); err != nil {
				log.Fatalf("unable to read remote config: %v", err)
				continue
			}

			// 将新配置反序列化到我们运行时的配置结构体中
			runtime_viper.Unmarshal(&conf)
		}
	}()
	wc.Wait()
	parseYaml(viper.GetViper())
}
