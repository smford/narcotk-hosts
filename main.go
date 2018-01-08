package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

func displayConfig() {
	fmt.Println("Starting displayConfig function")
	fmt.Printf("Networks:      %s\n", viper.GetString("Networks"))
	fmt.Printf("ShowHeader:    %s\n", viper.GetString("ShowHeader"))
	fmt.Printf("ListAll:       %s\n", viper.GetString("ListAll"))
	fmt.Printf("ListenPort:    %s\n", viper.GetString("ListenPort"))
	fmt.Printf("ListenIP:      %s\n", viper.GetString("ListenIP"))
	fmt.Printf("LogFile:       %s\n", viper.GetString("LogFile"))
	fmt.Printf("Verbose:       %s\n", viper.GetString("Verbose"))
	fmt.Printf("Database:      %s\n", viper.GetString("Database"))
	fmt.Printf("HeaderFile:    %s\n", viper.GetString("HeaderFile"))
	fmt.Printf("PrintColumns:  %s\n", viper.GetString("PrintColumns"))
}

func init() {
	fmt.Println("Starting init function\n")
	flag.String("addnetwork", "", "add a new network, used with --cidr and --desc")
	flag.String("cidr", "", "cidr of network, used with --adnetwork and --desc")
	configFile := flag.String("configfile", "", "configuration file to use")
	flag.String("database", "", "database file to use")
	flag.String("delnetwork", "", "delete a network")
	flag.Bool("displayconfig", false, "display configuration")
	flag.String("desc", "", "description of network, used with --addnetwork and --cidr")
	flag.Bool("help", false, "display help information")
	flag.Bool("listnetworks", false, "list all networks")
	flag.Bool("showmac", false, "show mac addresses of hosts")
	flag.String("network", "", "display hosts within a particular network")
	flag.Bool("setupdb", false, "setup a new database")
	flag.Bool("version", false, "display version information")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

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
		viper.SetDefault("Database", "./narcotk_hosts_all.db")
		viper.SetDefault("HeaderFile", "./header.txt")
		viper.SetDefault("PrintColumns", "blank")
	}
}

func main() {
	fmt.Println("Starting main function\n")

	if viper.GetBool("help") {
		displayHelp()
		os.Exit(0)
	}

	if viper.GetBool("displayconfig") {
		displayConfig()
		os.Exit(0)
	}

	if viper.GetString("delnetwork") != "" {
		delNetwork(viper.GetString("Database"), viper.GetString("delnetwork"))
		os.Exit(0)
	} else {
		fmt.Println("Error: you must provide a network to delete")
	}

	if viper.GetBool("version") {
		displayVersion()
		os.Exit(0)
	}

	if viper.GetBool("setupdb") {
		setupdb(viper.GetString("Database"))
		os.Exit(0)
	}

	if viper.GetBool("listnetworks") {
		listNetworks(viper.GetString("Database"))
		os.Exit(0)
	}

	if viper.GetString("addnetwork") != "" {
		if (viper.GetString("cidr") == "") || (viper.GetString("desc") == "") {
			fmt.Println("Error: When using --addnetwork you must also provide --cidr and --desc")
			os.Exit(1)
		} else {
			addNetwork(viper.GetString("Database"), viper.GetString("addnetwork"), viper.GetString("cidr"), viper.GetString("desc"))
			os.Exit(0)
		}
	}

	listHosts(viper.GetString("Database"), viper.GetString("network"), viper.GetBool("showmac"))
}

func addNetwork(databaseFile string, network string, cidr string, desc string) {
	fmt.Println("Adding new network: " + network + "\nCIDR: " + cidr + "\nDescription: " + desc)
	sqlquery := "insert into networks (network, cidr, description) values ('" + network + "', '" + cidr + "', '" + desc + "')"
	fmt.Println("addNetwork query: " + sqlquery)
	runSql(databaseFile, sqlquery)
}

func delNetwork(databasefile string, network string) {
	fmt.Println("Deleting network: " + network)
	sqlquery := "delete from networks where network like '" + network + "'"
	runSql(databasefile, sqlquery)
}

func runSql(databaseFile string, sqlquery string) {
	fmt.Println("Running generic runSql function")
	fmt.Println("runSql query: " + sqlquery)
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(sqlquery)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlquery)
		return
	}
}

func displayVersion() {
	fmt.Println("narcotk-hosts: 0.1")
}

func setupdb(databasefile string) {
	fmt.Println("Setting up a new database file: " + databasefile)

	db, err := sql.Open("sqlite3", databasefile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlquery := `
	CREATE TABLE hosts (
	  hostid text PRIMARY KEY,
	  network text NOT NULL,
	  ipsuffix integer NOT NULL,
	  ipaddress text NOT NULL,
	  fqdn text NOT NULL,
	  short1 text NOT NULL DEFAULT '',
	  short2 text NOT NULL DEFAULT '',
	  short3 text NOT NULL DEFAULT '',
	  short4 text NOT NULL DEFAULT ''
	  , mac TEXT DEFAULT '');
	  CREATE TABLE networks (
	    network text PRIMARY KEY,
	    cidr text NOT NULL,
	    description text NOT NULL DEFAULT ''
	  );
	  `

	_, err = db.Exec(sqlquery)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlquery)
		return
	}
}

func listNetworks(databaseFile string) {
	fmt.Println("Starting listNetworks function\n")
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
}

func listHosts(databaseFile string, network string, showmac bool) {
	fmt.Println("Starting listHosts function\n")
	sqlquery := "select * from hosts"
	if len(network) != 0 {
		fmt.Println("Displaying hosts from network: " + network)
		sqlquery = sqlquery + " where network like '" + network + "'"
	} else {
		fmt.Println("Displaying ALL hosts from ALL networks")
	}
	fmt.Println("sqlquery= " + sqlquery)

	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(sqlquery)
	for rows.Next() {
		var hostid string
		var network string
		var ipsuffix int
		var ipaddress string
		var fqdn string
		var short1 string
		var short2 string
		var short3 string
		var short4 string
		var mac string
		err = rows.Scan(&hostid, &network, &ipsuffix, &ipaddress, &fqdn, &short1, &short2, &short3, &short4, &mac)
		if err != nil {
			log.Fatal(err)
		}
		if showmac {
			fmt.Printf("%-17s  %-15s    %s  %s  %s  %s  %s\n", mac, ipaddress, fqdn, short1, short2, short3, short4)
		} else {
			fmt.Printf("%-15s    %s  %s  %s  %s  %s\n", ipaddress, fqdn, short1, short2, short3, short4)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func displayHelp() {
	helpmessage := `
Options:
  --help           Display help information
  --showheader     Prepend ./header.txt to the output [default=false]
  --displayconfig  Prints out the applied configuration
  --version

Commands:
  Print all hosts in a network:
      --network=192.168.1

  Show MAC addresses:
      --showmac
      
  List all networks:
      --listnetworks

  Add a new network:
      --addnetwork=192.168.2 --cidr=192.168.2.0/24 --desc="Management Network"

  Delete a network:
      --delnetwork=192.168.3

  Update a network:
      --updatenetwork=192.168.2 --network=192.168.3 --cidr=192.168.3/24 --desc="3rd Management Network"

  Adding a host:
      --addhost=server-1-199.domain.com --network=192.168.1 --ipaddress=192.168.1.13 --short1=server-1-199 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe

  Update a host:
  --updatehost=server-1-199.domain.com --host=server-1-200.domain.com --network=192.168.1 --ipaddress=192.168.1.200 --short1=server-1-200 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe
  ** When updating a host entry, all fields will be updated

  Delete a host:
      --delhost=server-1-200.domain.com --network=192.168.1

  Configuration file:
      --configfile=/path/to/file.yaml

  Database file:
      --database=/path/to/somefile.db

  Setup a new blank database file:
      --setupdb  --database=./newfile.db

  Start Web Service:
      --startweb

  Port to listen upon:
      --listenport=23000

  IP Address to listen upon:
      --listenip=10.0.0.14
`
	fmt.Printf("%s", helpmessage)

	os.Exit(0)
}
