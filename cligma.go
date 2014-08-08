package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	term "github.com/andrew-d/go-termutil"
	"github.com/codegangsta/cli"
	"github.com/franela/goreq"
)

const URL = "https://huigma.com/"

type PostParams struct {
	Content  string `json:"content"`
	Num      int    `json:"num"`
	Expire   string `json:"expire"`
	Callback string `json:"callback"`
}

type GetResponse struct {
	Content string
}

type PostResponse struct {
	Link string
	Hash string
}

func get(c *cli.Context, hash string) {

	resp, err := goreq.Request{
		Uri:     strings.TrimSuffix(c.GlobalString("url"), "/") + "/" + hash,
		Timeout: time.Duration(c.GlobalInt("timeout")) * time.Second,
	}.Do()

	if err != nil {
		log.Fatalln("Error in get:", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Error:", err)
	}

	var response GetResponse
	err = resp.Body.FromJsonTo(&response)
	if err != nil {
		log.Fatalln("Error decoing JSON", err)
	}

	print(response.Content)
}

func post(c *cli.Context) {

	var err error
	var bytes []byte
	if !term.Isatty(os.Stdin.Fd()) {
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
		Expire:   c.String("expire"),
		Callback: c.String("callback"),
	}

	resp, err := goreq.Request{
		Uri:         c.GlobalString("url") + "/api",
		Accept:      "application/json",
		ContentType: "application/json",
		Body:        item,
	}.Do()
	if err != nil {
		log.Fatalln("Error in post:", err)
	}

	var response PostResponse
	err = resp.Body.FromJsonTo(&response)
	if err != nil {
		log.Fatalln("Error decoing JSON", err)
	}

	println(response.Link)
}

func main() {

	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "cligma"
	app.Usage = "console client for Huigma service"

	timeoutFlag := cli.IntFlag{Name: "timeout, t", Value: 3, Usage: "timeout"}
	urlFlag := cli.StringFlag{Name: "url, u", Value: URL, Usage: "url of Huigma service"}

	app.Flags = []cli.Flag{
		urlFlag,
		timeoutFlag,
	}

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get content by hash",
			Action: func(c *cli.Context) {
				hash := c.Args().First()
				get(c, hash)
			},
		},
		{
			Name:  "post",
			Usage: "post content on Huigma",
			Flags: []cli.Flag{
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
				cli.StringFlag{
					Name:  "expire, e",
					Value: "",
					Usage: "expires in",
				},
				cli.StringFlag{
					Name:  "callback, b",
					Value: "",
					Usage: "callback url/email",
				},
			},
			Action: func(c *cli.Context) {
				post(c)
			},
		},
	}

	app.Action = func(c *cli.Context) {
		println("Type --help to see usage")
	}

	app.Run(os.Args)
}
