package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xwb1989/sqlparser"
	"log"
	"net/http"
	"os"
	"strings"
	_ "unicode"
)

type Host struct {
	IPAddress string `json:"IPAddress"`
	Hostname  string `json:"Hostname"`
	Short1    string `json:"Short1"`
	Short2    string `json:"Short2"`
	Short3    string `json:"Short3"`
	Short4    string `json:"Short4"`
	MAC       string `json:"MAC"`
}

type Hosts struct {
	Hosts []Host `json:"Hosts"`
}

type HostsSingleNetwork struct {
	Network     string `json:"Network"`
	CIDR        string `json:"CIDR"`
	Description string `json:"Description"`
	Hosts       []Host `json:"Hosts"`
}

type HostsAllNetworks struct {
	Networks []HostsSingleNetwork `json:"Networks"`
}

type SingleNetwork struct {
	Network     string `json:"Network"`
	CIDR        string `json:"CIDR"`
	Description string `json:"Description"`
}

type AllNetworks struct {
	Networks []SingleNetwork `json:"Networks"`
}

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
	fmt.Printf("JSON:          %s\n", viper.GetString("JSON"))
}

func init() {
	fmt.Println("Starting init function\n")
	flag.String("addhost", "", "add a new host, use with --network, --ipaddress (optional: --short1, --short2, --short3, --short4 and --mac)")
	flag.String("addnetwork", "", "add a new network, used with --cidr and --desc")
	flag.String("cidr", "", "cidr of network, used with --adnetwork and --desc")
	configFile := flag.String("configfile", "", "configuration file to use")
	flag.String("database", "", "database file to use")
	flag.String("delhost", "", "delete a host, used with --network")
	flag.String("delnetwork", "", "delete a network")
	flag.Bool("displayconfig", false, "display configuration")
	flag.String("desc", "", "description of network, used with --addnetwork and --cidr")
	flag.Bool("help", false, "display help information")
	flag.String("ipaddress", "", "ip address of new host")
	flag.Bool("json", false, "output in json")
	listenIp := flag.String("listenip", "", "ip address for webservice to bind to")
	listenPort := flag.String("listenport", "", "port for webservice to listen upon")
	flag.Bool("listnetworks", false, "list all networks")
	flag.Bool("showmac", false, "show mac addresses of hosts")
	flag.String("mac", "", "mac address of host")
	flag.String("network", "", "display hosts within a particular network")
	flag.Bool("setupdb", false, "setup a new database")
	flag.String("short1", "", "short1 hostname")
	flag.String("short2", "", "short2 hostname")
	flag.String("short3", "", "short3 hostname")
	flag.String("short4", "", "short4 hostname")
	flag.Bool("startweb", false, "start web service")
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

	if *listenPort != "" {
		viper.Set("ListenPort", listenPort)
	}

	if *listenIp != "" {
		viper.Set("ListenIP", listenIp)
	}

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
		viper.SetDefault("JSON", false)
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

	if viper.GetBool("startweb") {
		startWeb(viper.GetString("Database"), viper.GetString("ListenIP"), viper.GetString("ListenPort"))
		os.Exit(0)
	}

	if viper.GetString("delnetwork") != "" {
		delNetwork(viper.GetString("Database"), viper.GetString("delnetwork"))
		os.Exit(0)
	}

	if viper.GetString("delhost") != "" {
		if viper.GetString("network") == "" {
			fmt.Println("Error a network must be provided")
			os.Exit(1)
		} else {
			delHost(viper.GetString("Database"), viper.GetString("delhost"), viper.GetString("network"))
			os.Exit(0)
		}
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
		listNetworks(viper.GetString("Database"), nil, "select * from networks", viper.GetBool("json"))
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

	if viper.GetString("addhost") != "" {
		if (viper.GetString("network") == "") || (viper.GetString("ipaddress") == "") {
			fmt.Println("Error: When using --addhost you must also provide --network and --ipaddress")
			os.Exit(1)
		} else {
			addHost(viper.GetString("Database"), viper.GetString("addhost"), viper.GetString("network"), viper.GetString("ipaddress"), viper.GetString("short1"), viper.GetString("short2"), viper.GetString("short3"), viper.GetString("short4"), viper.GetString("mac"))
			os.Exit(0)
		}
	}

	listHost(viper.GetString("Database"), nil, viper.GetString("network"), "select * from hosts", viper.GetBool("showmac"), viper.GetBool("json"))
}

func addHost(databaseFile string, addhost string, network string, ipaddress string, short1 string, short2 string, short3 string, short4 string, mac string) {
	fmt.Println("Adding new host:")
	fmt.Println(addhost)
	fmt.Println(network)
	fmt.Println(ipaddress)
	fmt.Println(short1)
	fmt.Println(short2)
	fmt.Println(short3)
	fmt.Println(short4)
	fmt.Println(mac)
	sqlquery := "insert into hosts (hostid, network, ipsuffix, ipaddress, fqdn, short1, short2, short3, short4, mac) values ('" + breakIp(network, 2) + "-" + breakIp(ipaddress, 3) + "', '" + network + "', '" + breakIp(ipaddress, 3) + "', '" + ipaddress + "', '" + addhost + "', '" + short1 + "', '" + short2 + "', '" + short3 + "', '" + short4 + "', '" + mac + "')"
	runSql(databaseFile, sqlquery)
}

func delHost(databaseFile string, host string, network string) {
	fmt.Println("Deleting host:")
	fmt.Println(host)
	fmt.Println(network)
	sqlquery := "delete from hosts where (fqdn like '" + host + "') and (network like '" + network + "')"
	runSql(databaseFile, sqlquery)
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

	_, err := sqlparser.Parse(sqlquery)
	if err != nil {
		fmt.Println("Error Detected in SQL: ", err)
		os.Exit(1)
	}

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

func listNetworks(databaseFile string, webprint http.ResponseWriter, sqlquery string, printjson bool) {
	fmt.Println("Starting listNetworks")
	if webprint == nil {
		fmt.Println("webprint is null, printing to std out")
	}
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var mynetworks []SingleNetwork
	rows, err := db.Query(sqlquery)
	for rows.Next() {
		var network string
		var cidr string
		var description string
		err = rows.Scan(&network, &cidr, &description)
		mynetworks = append(mynetworks, SingleNetwork{network, cidr, description})
		if err != nil {
			log.Fatal(err)
		}
	}
	if printjson {
		c, _ := json.Marshal(mynetworks)
		if webprint == nil {
			fmt.Printf("%s", c)
		} else {
			fmt.Fprintf(webprint, "%s", c)
		}
	} else {
		if webprint == nil {
			for _, network := range mynetworks {
				fmt.Printf("%-15s  %-18s  %s\n", network.Network, network.CIDR, network.Description)
			}
		} else {
			for _, network := range mynetworks {
				fmt.Fprintf(webprint, "%-15s  %-18s  %s\n", network.Network, network.CIDR, network.Description)
			}
		}
	}
}

func displayVersion() {
	fmt.Println("narcotk-hosts: 0.1")
}

func setupdb(databaseFile string) {
	fmt.Println("Setting up a new database file: " + databaseFile)
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
	  short4 text NOT NULL DEFAULT '',
	  mac TEXT DEFAULT '')`
	runSql(databaseFile, sqlquery)
	sqlquery = `
	CREATE TABLE networks (
     network text PRIMARY KEY,
     cidr text NOT NULL,
     description text NOT NULL DEFAULT '')`
	runSql(databaseFile, sqlquery)
}

func listHost(databaseFile string, webprint http.ResponseWriter, network string, sqlquery string, showmac bool, printjson bool) {
	fmt.Println("Starting listHost")
	if webprint == nil {
		fmt.Println("webprint is null, printing to std out")
	}
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var myhosts []Host
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
		myhosts = append(myhosts, Host{ipaddress, fqdn, short1, short2, short3, short4, mac})
	}
	if printjson {
		// print json
		c, _ := json.Marshal(myhosts)
		if webprint == nil {
			fmt.Printf("%s", c)
		} else {
			fmt.Fprintf(webprint, "%s", c)
		}
	} else {
		// print standard
		if webprint == nil {
			if showmac {
				for _, host := range myhosts {
					fmt.Printf("%-17s  %-15s    %s  %s  %s  %s  %s\n", host.MAC, host.IPAddress, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
				}
			} else {
				for _, host := range myhosts {
					fmt.Printf("%-15s    %s  %s  %s  %s  %s\n", host.IPAddress, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
				}
			}
		} else {
			// webprint
			if showmac {
				for _, host := range myhosts {
					fmt.Fprintf(webprint, "%-17s  %-15s    %s  %s  %s  %s  %s\n", host.MAC, host.IPAddress, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
				}
			} else {
				for _, host := range myhosts {
					fmt.Fprintf(webprint, "%-15s    %s  %s  %s  %s  %s\n", host.IPAddress, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
				}
			}
		}
	}
}

func breakIp(ipaddress string, position int) string {
	deliminator := func(c rune) bool {
		return (c == '.')
	}
	ipArray := strings.FieldsFunc(ipaddress, deliminator)
	return ipArray[position]
}

func startWeb(databaseFile string, listenip string, listenport string) {
	fmt.Println("Starting webserver: " + listenip + ":" + listenport)
	r := mux.NewRouter()
	hostsRouter := r.PathPrefix("/hosts").Subrouter()
	hostsRouter.HandleFunc("", handlerHosts)
	hostsRouter.HandleFunc("/", handlerHosts)
	hostsRouter.HandleFunc("/json", handlerHostsJson)
	hostsRouter.HandleFunc("/{network}", handlerHostsNetwork)
	hostsRouter.HandleFunc("/{network}/json", handlerHostsNetworkJson)

	hostRouter := r.PathPrefix("/host").Subrouter()
	hostRouter.HandleFunc("/{host}", handlerHost)
	hostRouter.HandleFunc("/{host}/json", handlerHostJson)

	networksRouter := r.PathPrefix("/networks").Subrouter()
	networksRouter.HandleFunc("", handlerNetworks)
	networksRouter.HandleFunc("/", handlerNetworks)
	networksRouter.HandleFunc("/json", handlerNetworksJson)

	networkRouter := r.PathPrefix("/network").Subrouter()
	networkRouter.HandleFunc("/{network}", handlerNetwork)
	networkRouter.HandleFunc("/{network}/json", handlerNetworkJson)

	ipRouter := r.PathPrefix("/ip").Subrouter()
	ipRouter.HandleFunc("/{ip}", handlerIp)
	ipRouter.HandleFunc("/{ip}/json", handlerIpJson)

	http.ListenAndServe(listenip+":"+listenport, r)
}

func handlerHosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting handlerHosts")
	listHost(viper.GetString("Database"), w, viper.GetString("network"), "select * from hosts", viper.GetBool("showmac"), false)
}

func handlerHostsJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting handlerHostsJson")
	//fmt.Fprintf(w, "json print hosts")
	listHost(viper.GetString("Database"), w, viper.GetString("network"), "select * from hosts", viper.GetBool("showmac"), true)
}

func handlerHostsNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerHostNetwork: " + vars["network"])
	listHost(viper.GetString("Database"), w, viper.GetString("network"), "select * from hosts where network like '"+vars["network"]+"'", false, false)
}

func handlerHostsNetworkJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerHostNetworkJson: " + vars["network"])
	fmt.Fprintf(w, "json print host network: %s", vars["network"])
}

func handlerHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerHost: " + vars["host"])
	listHost(viper.GetString("Database"), w, viper.GetString("network"), "select * from hosts where fqdn like '"+vars["host"]+"'", viper.GetBool("showmac"), false)
}

func handlerHostJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerHostJson: " + vars["host"])
	listHost(viper.GetString("Database"), w, viper.GetString("network"), "select * from hosts where fqdn like '"+vars["host"]+"'", viper.GetBool("showmac"), true)
	//fmt.Fprintf(w, "json print host: %s", vars["host"])
}

func handlerNetworks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting handlerNetworks")
	sqlquery := "select * from networks"
	listNetworks(viper.GetString("Database"), w, sqlquery, false)
}

func handlerNetworksJson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting handlerNetworksJson")
	sqlquery := "select * from networks"
	listNetworks(viper.GetString("Database"), w, sqlquery, true)
}

func handlerNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerNetwork: " + vars["network"])
	sqlquery := "select * from networks where network like '" + vars["network"] + "'"
	listNetworks(viper.GetString("Database"), w, sqlquery, false)
}

func handlerNetworkJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerNetworkJson: " + vars["network"])
	sqlquery := "select * from networks where network like '" + vars["network"] + "'"
	listNetworks(viper.GetString("Database"), w, sqlquery, true)
}

func handlerIp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerIp: " + vars["ip"])
	sqlquery := "select * from hosts where ipaddress like '" + vars["ip"] + "'"
	listHost(viper.GetString("Database"), w, "", sqlquery, false, false)
}

func handlerIpJson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Starting handlerIpJson: " + vars["ip"])
	fmt.Fprintf(w, "json print ip: %s", vars["ip"])
}

func displayHelp() {
	helpmessage := `
Options:
  --help           Display help information
	--json           Print output in json
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
