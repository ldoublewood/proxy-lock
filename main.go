package main

import (
	"gin-gorm-example/routes"
	"gin-gorm-example/database"
	"gin-gorm-example/config"
	"github.com/jinzhu/gorm"
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"os"
	"sort"
	"log"
)

func main() {
	//if err := config.Load("config/config.yaml"); err != nil {
	//	fmt.Println("Failed to load configuration")
	//	return
	//}
	app := cli.NewApp()
	app.Name = "GoTest"
	app.Usage = "hello world"
	app.Version = "1.2.3"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Value: "./proxy.db",
			EnvVar: "DB_SOURCE",
			Usage: "database source",
		},
		cli.StringFlag{
			Name:  "dialect, d",
			Value: "sqlite3",
			EnvVar: "DB_DIALECT",
			Usage: "database dialect",
		},
		cli.StringFlag{
			Name:  "addr, a",
			Value: "0.0.0.0:4433",
			EnvVar: "ADDR",
			Usage: "service listener address",
		},
		cli.StringFlag{
			Name:  "testdesturl",
			Value: "",
			EnvVar: "TEST_DEST_URL",
			Usage: "test dest url address, e.g. https://www.baidu.com, no test if it is empty",
		},
	}

	app.Action = func(c *cli.Context) error {
		var db *gorm.DB
		var err error
		fmt.Println("starting")

		config.Init(c)
		if db, err = database.InitDB(); err != nil {
			fmt.Println("err open databases")
			return err
		}
		defer db.Close()

		router := routes.InitRouter()
		err = router.Run(config.Get().Addr)
		return err
	}
	app.Before = func(c *cli.Context) error {
		fmt.Println("app Before")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Println("app After")
		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	cli.HelpFlag = cli.BoolFlag {
		Name: "help, h",
		Usage: "Help",
	}

	cli.VersionFlag = cli.BoolFlag {
		Name: "print-version, v",
		Usage: "print version",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}