package main

/**
读取配置文件demo
*/

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"reflect"
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

func sub(key string, value interface{}) error {
	log.Printf("配置文件的前缀为: %v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
}

func main() {
	// 利用Viper读取本地配置信息

	fmt.Printf("性别:%s\n 爱好:%s\n性别:%s\n年龄:%d\n", Resume.Name, Resume.Habits, Resume.Sex, Resume.Age)

	parseYaml(viper.GetViper())
	fmt.Println(Contains("Basketball", Resume.Habits))
	for {
		time.Sleep(5 * time.Second)
		parseYaml(viper.GetViper())
	}
}
