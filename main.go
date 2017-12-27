package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

func printConfig() {
	fmt.Println("Starting printConfig function")
	fmt.Printf("Networks:      %s\n", viper.GetString("Networks"))
	fmt.Printf("ShowHeader:    %s\n", viper.GetString("ShowHeader"))
	fmt.Printf("ListAll:       %s\n", viper.GetString("ListAll"))
	fmt.Printf("ListenPort:    %s\n", viper.GetString("ListenPort"))
	fmt.Printf("ListenIP:      %s\n", viper.GetString("ListenIP"))
	fmt.Printf("LogFile:       %s\n", viper.GetString("LogFile"))
	fmt.Printf("Verbose:       %s\n", viper.GetString("Verbose"))
	fmt.Printf("DatabaseFile:  %s\n", viper.GetString("DatabaseFile"))
	fmt.Printf("HeaderFile:    %s\n", viper.GetString("HeaderFile"))
	fmt.Printf("PrintColumns:  %s\n", viper.GetString("PrintColumns"))
}

func init() {
	fmt.Println("Starting init function\n")
	configFile := flag.String("configfile", "", "configuration file to use")
	flag.Parse()

	viper.AddConfigPath(".")

	if *configFile == "" {
		viper.SetConfigName("narco-hosts-config")
	} else {
		viper.SetConfigName(*configFile)
	}
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
		viper.SetDefault("Networks", "my networks (default)")
		viper.SetDefault("ShowHeader", true)
		viper.SetDefault("ListAll", false)
		viper.SetDefault("ListenPort", "23000")
		viper.SetDefault("ListenIP", "127.0.0.1")
		viper.SetDefault("LogFile", "./logs.txt")
		viper.SetDefault("Verbose", true)
		viper.SetDefault("DatabaseFile", "./narcotk_hosts_all.db")
		viper.SetDefault("HeaderFile", "./header.txt")
		viper.SetDefault("PrintColumns", "blank")
	}
}

func main() {
	fmt.Println("Starting main function\n")
	printConfig()
	listAllNetworks(viper.GetString("DatabaseFile"))
}

func listAllNetworks(databaseFile string) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlquery := "select * from networks"
	rows, err := db.Query(sqlquery)
	for rows.Next() {
		var network string
		var cidr string
		var description string
		err = rows.Scan(&network, &cidr, &description)
		fmt.Printf("%-15s  %-18s  %s\n", network, cidr, description)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
