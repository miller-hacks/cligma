package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/codegangsta/cli"
	"github.com/franela/goreq"
)

type PostParams struct {
	Content  string `json:"content"`
	Num      int    `json:"num"`
	Expire   int    `json:"expire"`
	Callback string `json:"callback"`
}

type GetResponse struct {
	Content string
}

type PostResponse struct {
	Link string `json:"link"`
	Key  string `json:"key"`
}

func get(c *cli.Context, key string) {

	resp, err := goreq.Request{
		Uri:     strings.TrimSuffix(c.GlobalString("huigma"), "/") + "/api/" + key,
		Timeout: time.Duration(c.GlobalInt("timeout")) * time.Second,
	}.Do()

	if err != nil {
		log.Fatalln("Error in get:", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Error:", resp.StatusCode)
	}

	var response GetResponse
	err = resp.Body.FromJsonTo(&response)
	if err != nil {
		log.Fatalln("Error decoding JSON", err)
	}

	print(response.Content)
}

func post(c *cli.Context) {

	var err error
	var bytes []byte

	if !terminal.IsTerminal(syscall.Stdin) {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(c.String("content"))
	}

	if err != nil {
		log.Fatalln(err)
	}

	item := PostParams{
		Content:  string(bytes),
		Num:      c.Int("num"),
		Expire:   c.Int("expire"),
		Callback: c.String("callback"),
	}

	resp, err := goreq.Request{
		Uri:         strings.TrimSuffix(c.GlobalString("huigma"), "/") + "/api",
		Method:      "POST",
		Accept:      "application/json",
		ContentType: "application/json",
		Body:        item,
	}.Do()

	if err != nil {
		log.Fatalln("Error in post:", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Error:", resp.StatusCode)
	}

	var response PostResponse
	err = resp.Body.FromJsonTo(&response)
	if err != nil {
		log.Fatalln("Error decoding JSON", err)
	}

	println(response.Link)
}

func main() {

	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "cligma"
	app.Usage = "console client for Huigma service"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "timeout, t",
			Value: 3,
			Usage: "timeout",
		},
		cli.StringFlag{
			Name:  "huigma, u",
			Value: "https://huigma.com/",
			Usage: "url of Huigma service",
		},
		cli.StringFlag{
			Name:  "content, c",
			Value: "",
			Usage: "path to content",
		},
		cli.StringFlag{
			Name:  "num, n",
			Value: "",
			Usage: "number of shows",
		},
		cli.IntFlag{
			Name:  "expire, e",
			Value: 3600 * 24,
			Usage: "expires in",
		},
		cli.StringFlag{
			Name:  "callback, b",
			Value: "",
			Usage: "callback url/email",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get content by key",
			Action: func(c *cli.Context) {
				key := c.Args().First()
				get(c, key)
			},
		},
	}

	app.Action = func(c *cli.Context) {
		post(c)
	}

	app.Run(os.Args)
}
