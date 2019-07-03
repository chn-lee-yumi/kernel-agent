package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func ExecCmd(w http.ResponseWriter, r *http.Request) {
	datas := make(map[string]interface{}) //返回数据的map
	//检查ip
	var ip string
	if _, ok := r.Header["X-Forwarded-For"]; ok {
		ip_list:=strings.Split(r.Header["X-Forwarded-For"][0], ",")
		ip=ip_list[len(ip_list)-1]
		ip=strings.Trim(ip," ")
	}else{
		ip=strings.Split(r.RemoteAddr, ":")[0]
	}
	if ip != *server {
		datas["code"] = -1
		datas["msg"] = "非服务器IP"
		ReturnJson(datas, w)
		return
	}
	//读取命令（没有处理key不存在的情况）
	body, _ := ioutil.ReadAll(r.Body)
	url_parsed, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println(url_parsed)
		fmt.Printf("URL参数解析错误：%s\n", err.Error())
		datas["code"] = -2
		datas["msg"] = "URL参数解析错误：" + err.Error()
		ReturnJson(datas, w)
		return
	}
	//执行命令
	cmd := exec.Command("bash", "-c", url_parsed["cmd"][0])
	output, err := cmd.Output()
	if err != nil {
		datas["code"] = 1
		datas["err"]=err.Error()
	}else{
		datas["code"] = 0
	}
	datas["msg"] = string(output)
	ReturnJson(datas, w)
	return
}

func ReturnJson(datas map[string]interface{}, w http.ResponseWriter) {
	//返回数据
	json_datas, _ := json.Marshal(datas) //map转json
	fmt.Fprintln(w, string(json_datas))
}

var server = flag.String("s", "222.200.97.179", "服务器IP")

func main() {
	flag.Parse() //读取命令行参数
	fmt.Printf("服务器：%s\n", *server)
	http.HandleFunc("/api/ExecCmd", ExecCmd)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
