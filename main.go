package main

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/urfave/cli/v2"
	"github.com/xelabs/go-mysqlstack/driver"
)

func MysqlCheck(host, port, username, password string){
	address := fmt.Sprintf("%s:%s", host, port)
	client, err := driver.NewConn(username, password, address, "", "")
	if err != nil {
		g.Log().Warningf(nil,"[-] %s", err.Error())
		return
	}

	defer func(){
		if client != nil{
			client.Close()
		}
	}()
	qr, err := client.FetchAll("select user()", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] select user()结果：%+v", qr.Rows)

	qr, err = client.FetchAll("show databases", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] show databases结果：%+v", qr.Rows)

	qr, err = client.FetchAll("select @@basedir", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] select @@basedir结果：%+v", qr.Rows)

	qr, err = client.FetchAll("select @@datadir", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] select @@datadir结果：%+v", qr.Rows)

	qr, err = client.FetchAll("select @@version_compile_os", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] select @@version_compile_os结果：%+v", qr.Rows)

	qr, err = client.FetchAll("select @@hostname", -1)
	if err != nil{
		g.Log().Warningf(nil,"[-] mysql执行sql失败：%s", err.Error())
		return
	}
	g.Log().Infof(nil,"[+] select @@hostname结果：%+v", qr.Rows)

}

func main(){
	app := &cli.App{
		Name: "mysqls",
		Usage: "Mysql蜜罐检测工具",
		Version: "v1.0",
		UsageText: "./mysqls -t 127.0.0.1 -p 3306 -u root -P 123456",
		Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "target",
				Usage: "请输入MySQL地址",
				Aliases: []string{"t"},
				Required: true,
			},
			&cli.StringFlag{
				Name: "port",
				Usage: "请输入MySQL端口",
				Aliases: []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name: "username",
				Usage: "请输入MySQL用户名",
				Aliases: []string{"u"},
				Required: true,
			},
			&cli.StringFlag{
				Name: "password",
				Usage: "请输入MySQL密码",
				Aliases: []string{"P"},
			},
		},
		Action: func(c *cli.Context) error {
			if !c.IsSet("target") && !c.IsSet("port") && !c.IsSet("username"){
				cli.ShowAppHelpAndExit(c,0)
			}
			MysqlCheck(c.String("target"),c.String("port"),c.String("username"),c.String("password"))
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		g.Log().Fatalf(nil, "[-] %s", err.Error())
	}
}
