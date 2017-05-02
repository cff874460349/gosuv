// govscp
package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
	//"github.com/kennygrant/sanitize"
	"github.com/qiniu/log"
	"github.com/urfave/cli"
)

type Configuration struct {
	Name    string `yaml:"name"`
	Procnum int    `yaml:"procnum"`
	Iopmac  string `yaml:"iopmac"`
	Kyee    struct {
		Enable       string `yaml:"kyee"`
		TagMac       string `yaml:"tag_mac"`
		TagNum       int    `yaml:"tag_num"`
		DataHb       string `yaml:"data_hb"`
		DataHbTimes  int    `yaml:"data_hb_times"`
		DataHbInter  int    `yaml:"data_hb_inter"`
		DataOff      string `yaml:"data_off"`
		DataOffTimes int    `yaml:"data_off_times"`
		DataOffInter int    `yaml:"data_off_inter"`
		DataOut      string `yaml:"data_out"`
		DataOutTimes int    `yaml:"data_out_times"`
		DataOutInter int    `yaml:"data_out_inter"`
	}
	Ladrip struct {
		Enable  string `yaml:"ladrip"`
		TagMac  string `yaml:"tag_mac"`
		TagNum  int    `yaml:"tag_num"`
		DataW   string `yaml:"data_weight"`
		DataE2  string `yaml:"data_e2"`
		WeightV string `yaml:"weight_value" json:"weight_value"`
		Rate    string `yaml:"rate" json:"rate"`
	}
	Ewell struct {
		Enable       string `yaml:"ewell"`
		TagMac       string `yaml:"tag_mac"`
		TagNum       int    `yaml:"tag_num"`
		DataSts      string `yaml:"data_sts"`
		DataDat      string `yaml:"data_dat"`
		DataDatInter int    `yaml:"data_dat_inter"`
		DataDatTimes int    `yaml:"data_dat_times"`
	}
	Controller struct {
		Ip string `yaml:"ip"`
	}
}

func IsDir(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}

func readConf(filename string) (c Configuration, err error) {
	// initial default value
	c.Name = "default"
	c.Procnum = 2
	c.Iopmac = "000000"
	c.Kyee.Enable = "0"
	c.Ladrip.Enable = "0"
	c.Controller.Ip = "0.0.0.0"

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		data = []byte("")
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return
	}
	cfgDir := filepath.Dir(filename)
	if !IsDir(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}
	data, _ = yaml.Marshal(c)
	err = ioutil.WriteFile(filename, data, 0644)
	return
}

var cfg Configuration
var defaultConfigDir string

func init() {
	defaultConfigDir = os.Getenv("GOSUV_HOME_DIR")
	if defaultConfigDir == "" {
		defaultConfigDir = filepath.Join(os.Getenv("HOME"), ".gosuv")
	}
}

func GoFunc(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

func main() {
	var defaultConfigPath = filepath.Join(defaultConfigDir, "config_test.yml")
	var wg sync.WaitGroup

	//wg.Add(cfg.Procnum)

	app := cli.NewApp()
	app.Name = "govscp"
	app.Version = "0.1"
	app.Usage = "govscp -conf <config>"
	app.Before = func(c *cli.Context) error {
		var err error
		cfgPath := c.GlobalString("conf")
		cfg, err = readConf(cfgPath)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "chp",
			Email: "chp@ruijie.com.cn",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Usage: "config file",
			Value: defaultConfigPath,
		},
	}
	app.Action = func(c *cli.Context) error {
		//		fmt.Println(cfg.Procnum)
		//		fmt.Println(cfg.Iopmac)
		//		fmt.Println(cfg.Iotmac)
		//		fmt.Println(cfg.Iotnum)
		//		fmt.Println(cfg.Data.Weight)
		//		fmt.Println(cfg.Data.E2)
		//		fmt.Println(cfg.Controller.Ip)

		//cmd.Stdout = logFd
		//cmd.Stderr = logFd
		var err error
		var n, k, e, id, nn, kk, ee int
		b := make(chan bool)
		//		var fout io.Writer
		//		var OutputFile *os.File
		//		logDir := filepath.Join(defaultConfigDir, "log", cfg.Name)
		//		if !IsDir(logDir) {
		//			os.MkdirAll(logDir, 0755)
		//		}

		//		OutputFile, err = os.OpenFile(filepath.Join(logDir, "output.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		//		if err != nil {
		//			fmt.Println("create stdout log failed:", err)
		//			fout = ioutil.Discard
		//		} else {
		//			fout = OutputFile
		//		}
		wg.Add(cfg.Procnum)
		rand.Seed(time.Now().UTC().UnixNano())
		r := rand.Intn(35535)
		id, err = strconv.Atoi(cfg.Iopmac)
		id = id - 1
		n, err = strconv.Atoi(cfg.Kyee.TagMac)
		k, err = strconv.Atoi(cfg.Ladrip.TagMac)
		e, err = strconv.Atoi(cfg.Ewell.TagMac)
		for i := 0; i < cfg.Procnum; i++ {
			j := i*10 + r + 20000
			nn = n + i*cfg.Kyee.TagNum
			kk = k + i*cfg.Ladrip.TagNum
			ee = e + i*cfg.Ewell.TagNum
			id = id + 1
			go func() {
				defer wg.Done()
				cmd := exec.Command("vscp_demo", strconv.Itoa(j), cfg.Controller.Ip, strconv.Itoa(id),
					cfg.Kyee.Enable, strconv.Itoa(nn), strconv.Itoa(cfg.Kyee.TagNum),
					cfg.Kyee.DataHb, strconv.Itoa(cfg.Kyee.DataHbTimes), strconv.Itoa(cfg.Kyee.DataHbInter),
					cfg.Kyee.DataOff, strconv.Itoa(cfg.Kyee.DataOffTimes), strconv.Itoa(cfg.Kyee.DataOffInter),
					cfg.Kyee.DataOut, strconv.Itoa(cfg.Kyee.DataOutTimes), strconv.Itoa(cfg.Kyee.DataOutInter),
					cfg.Ladrip.Enable, strconv.Itoa(kk), strconv.Itoa(cfg.Ladrip.TagNum),
					cfg.Ladrip.DataW, cfg.Ladrip.DataE2, cfg.Ewell.Enable, strconv.Itoa(ee), strconv.Itoa(cfg.Ewell.TagNum),
					cfg.Ewell.DataSts, cfg.Ewell.DataDat, strconv.Itoa(cfg.Ewell.DataDatInter), strconv.Itoa(cfg.Ewell.DataDatTimes), cfg.Name, strconv.Itoa(i),
					cfg.Ladrip.WeightV, cfg.Ladrip.Rate)

				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Start()
				//log.Println("in goroutine")
				if err != nil {
					fmt.Println(err)
				}
				select {
				case err = <-GoFunc(cmd.Wait):
					fmt.Println("vscp " + strconv.Itoa(j) + " stoped")
					//fmt.Println("vscp started failed, %v", err)
					//				case <-time.After(200 * time.Millisecond /*3 * time.Second*/):
					//					fmt.Println("vscp " + strconv.Itoa(j) + " started")
					//wg.Done()
				default:
					b <- true
				}
			}()
			<-b
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
	//time.Sleep(4 * time.Second)
	wg.Wait()
	fmt.Println("exit")
	//os.Exit(0)
}
