package main

import (
	"encoding/json"
	"fmt"
	"github.com/simplePBFT/config"
	"github.com/simplePBFT/network"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

const (
	defaultConfigPath = "config0.json"
)

func main() {
	//加载配置文件
	path := defaultConfigPath
	if len(os.Args)>1 {
		path = os.Args[1]
	}
	config, err:= LoadConfig(path)
	if err != nil {
		panic(err)
	}

	//启动服务
	server := network.NewServer(config)
	server.Start()

	//等待quit信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("Exit:", s)
	os.Exit(0)
}

func LoadConfig(path string) (*config.Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	jsonBytes,err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config config.Config
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
