package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Configuration struct {
	Port int
}

type ApiRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Res     string `json:"res"`
}

func main() {
	file, err := os.Open("./conf.json")
	if err != nil {
		fmt.Println("read configuration file failed: ", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("parsing configuration file failed: ", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Salamander制作")
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			io.WriteString(w, respond(1, "必须为GET请求", ""))
			return
		}
		query := r.URL.Query()
		// 筛选条件，有：men（内存），cpu
		if len(query["sort"]) < 1 || len(query["sort"][0]) < 1 {
			io.WriteString(w, respond(1, "sort不能为空", ""))
			return
		}
		if len(query["num"]) < 1 || len(query["num"][0]) < 1 {
			io.WriteString(w, respond(1, "num不能为空", ""))
			return
		}
		sort := query["sort"][0]
		numStr := query["num"][0]
		num, err := strconv.ParseInt(numStr, 10, 8)
		if err != nil || num <= 0 {
			io.WriteString(w, respond(1, "num必须是整数且大于0", ""))
			return
		}
		var output bytes.Buffer
		var cmdPs *exec.Cmd
		if sort == "mem" {
			cmdPs = exec.Command("ps", "-aux", "--sort", "-pmem")
		} else if sort == "cpu" {
			cmdPs = exec.Command("ps", "-aux", "--sort", "-pcpu")
		} else {
			io.WriteString(w, respond(1, "sort类型未知", ""))
			return
		}
		cmdHead := exec.Command("head", fmt.Sprintf("-%d", num))
		cmdHead.Stdin, _ = cmdPs.StdoutPipe()
		cmdHead.Stdout = &output
		err = cmdHead.Start()
		if err != nil {
			io.WriteString(w, respond(2, err.Error(), ""))
			fmt.Println("start head failed: ", err)
			return
		}
		err = cmdPs.Run()
		if err != nil {
			fmt.Println("exec ps failed: ", err)
			io.WriteString(w, respond(2, err.Error(), ""))
			return
		}
		cmdHead.Wait()
		io.WriteString(w, respond(0, "success", output.String()))
	})
	fmt.Printf("start listening port %d\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}

func respond(code int, msg, queryRes string) string {
	res := ApiRes{
		Errcode: code,
		Errmsg:  msg,
		Res:     queryRes,
	}
	data, _ := json.Marshal(res)
	return string(data)
}
