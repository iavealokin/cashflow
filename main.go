package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/iavealokin/cashflow/app/apiserver"
	_ "github.com/lib/pq"
)

var configPath string

func init(){
	flag.StringVar(&configPath,"config-path","configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()
	config:=apiserver.NewConfig()
	_,err := toml.DecodeFile(configPath,config)
	if err!= nil{
		log.Fatal(err)
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func(){
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ApiServer started")
wg.Done()
}()
	go func(){
	
	if err := apiserver.StartWeb(config); err != nil {
		log.Fatal(err)
	}
	fmt.Println("WebServer started")
	wg.Done()
}()
	wg.Wait()

}


