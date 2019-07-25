package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"regexp"
)


//会有xss跨站攻击漏洞
//http://127.0.0.1:8080/xss?name=<script>alert(2)</script>
//只是现在得浏览器把他给屏蔽了
func Handler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	//如果解析失败 直接退出 输出对应的错误原因
	if err != nil {
		log.Fatal(nil)
	}
	//获取 传递的name 参数
	user_pro := req.FormValue("name")

	//直接在页面上输出用户传递数据信息
	fmt.Fprintf(w, "%s", user_pro)

}

func main() {


	str1 := "<script>alert(2)</script>"

	str2 := html.EscapeString(str1)
	fmt.Println(str2)
	str3 := html.UnescapeString(str2)
	fmt.Println(str3)




	//绑定路由 讲访问 / 绑定给  Handler 方法进行处理
	http.HandleFunc("/xss", Handler)
	http.HandleFunc("/xss2", Handler2)
	http.HandleFunc("/xss3", Handler3) //使用正则替换
	http.ListenAndServe(":8080", nil)
}


//没有跨站得漏洞
//浏览http://127.0.0.1:8080/xss2?name=<script>alert(2)</script>
//会输出&lt;script&gt;alert(2)&lt;/script&gt;
func Handler2(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	//如果解析失败 直接退出 输出对应的错误原因
	if err != nil {
		log.Fatal(nil)
	}
	//获取 传递的name 参数
	name := req.FormValue("name")
	fmt.Println(name)

	nameStr := template.HTMLEscapeString(name)
	fmt.Println(nameStr)

	//使用模板把nameStr原样输出到浏览器
	t, _ := template.New("test").Parse(`<html> {{ . }}</html>`)
	t.ExecuteTemplate(w,"test",nameStr)
}



//http://127.0.0.1:8080/xss3?name=<h1>hello</h1><script>alert("2")</script>
func Handler3(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	//如果解析失败 直接退出 输出对应的错误原因
	if err != nil {
		log.Fatal(nil)
	}
	//获取 传递的name 参数
	name := req.FormValue("name")
	fmt.Println(name)
	//匹配出以<script type="text/javascript"></script>
	reg, _ := regexp.Compile(`<script[^>]*>|</script>`)

	//替换为空
	txt := reg.ReplaceAllLiteralString(name, "")
	fmt.Println("")

	t, _ := template.New("test").Parse(`<html> {{ . }}</html>`)
	t.ExecuteTemplate(w,"test", txt)
}



