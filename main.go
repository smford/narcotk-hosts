package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xwb1989/sqlparser"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	_ "unicode"
)

var db *sql.DB

// Host holds all details internally within narcotk-hosts for a particular host
type Host struct {
	PaddedIP string `json:"PaddedIP"`
	Network  string `json:"Network"`
	IPv4     string `json:"IPv4"`
	IPv6     string `json:"IPv6"`
	Hostname string `json:"Hostname"`
	Short1   string `json:"Short1"`
	Short2   string `json:"Short2"`
	Short3   string `json:"Short3"`
	Short4   string `json:"Short4"`
	MAC      string `json:"MAC"`
}

// Hosts is an array of hosts
type Hosts struct {
	Hosts []Host `json:"Hosts"`
}

// HostsSingleNetwork details a network and all its constituent hosts
type HostsSingleNetwork struct {
	Network     string `json:"Network"`
	CIDR        string `json:"CIDR"`
	Description string `json:"Description"`
	Hosts       []Host `json:"Hosts"`
}

// HostsAllNetworks contains all hosts in all networks
type HostsAllNetworks struct {
	Networks []HostsSingleNetwork `json:"Networks"`
}

// SingleNetwork holds details of a specific network
type SingleNetwork struct {
	PaddedNetwork string `json:"PaddedNetwork"`
	Network       string `json:"Network"`
	CIDR          string `json:"CIDR"`
	Description   string `json:"Description"`
}

// AllNetworks holds details of all networks
type AllNetworks struct {
	Networks []SingleNetwork `json:"Networks"`
}

// log an error and if fatal exit app
func showerror(message string, e error, reaction string) {
	if e != nil {
		if strings.ToLower(reaction) == "fatal" {
			log.Fatalf("ERROR: %s:%s", message, e)
		} else {
			log.Printf("%s: %s:%s", strings.ToUpper(reaction), message, e)
		}
	}
}

func displayConfig() {
	fmt.Println("Starting displayConfig function")
	fmt.Printf("ShowHeader:      %s\n", viper.GetString("ShowHeader"))
	fmt.Printf("ListenPort:      %s\n", viper.GetString("ListenPort"))
	fmt.Printf("ListenIP:        %s\n", viper.GetString("ListenIP"))
	fmt.Printf("Database:        %s\n", viper.GetString("Database"))
	fmt.Printf("DatabaseType:    %s\n", viper.GetString("DatabaseType"))
	fmt.Printf("HeaderFile:      %s\n", viper.GetString("HeaderFile"))
	fmt.Printf("IndexFile:       %s\n", viper.GetString("IndexFile"))
	fmt.Printf("Files:           %s\n", viper.GetString("Files"))
	fmt.Printf("JSON:            %s\n", viper.GetString("JSON"))
	fmt.Printf("EnableTLS:       %s\n", viper.GetString("EnableTLS"))
	fmt.Printf("TLSCert:         %s\n", viper.GetString("TLSCert"))
	fmt.Printf("TLSKey:          %s\n", viper.GetString("TLSKey"))
	fmt.Printf("RegistrationKey: %s\n", viper.GetString("RegistationKey"))
	fmt.Printf("Verbose:         %s\n", viper.GetString("Verbose"))
}

// PrepareMac cleans up a mac address and makes in to a consistent format
func PrepareMac(macaddress string) string {
	//fmt.Println("Starting PrepareMac")
	macaddress = strings.ToLower(macaddress)
	// strip colons
	macaddress = strings.Replace(macaddress, ":", "", -1)
	// strip hyphens
	macaddress = strings.Replace(macaddress, "-", "", -1)
	// add colons
	var n = 2
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(macaddress) - 1
	for i, rune := range macaddress {
		buffer.WriteRune(rune)
		if i%n == n1 && i != l1 {
			buffer.WriteRune(':')
		}
	}
	return buffer.String()
}

func init() {
	//fmt.Println("Starting init function")
	flag.String("addhost", "", "add a new host, use with --network, --ip (optional: --ipv6 --short1, --short2, --short3, --short4 and --mac)")
	flag.String("addnetwork", "", "add a new network, used with --cidr and --desc")
	flag.String("cidr", "", "cidr of network, used with --adnetwork and --desc")
	configFile := flag.String("configfile", "", "configuration file to use")
	flag.String("database", "", "database file to use")
	flag.String("databasetype", "", "database type to use")
	flag.String("delhost", "", "delete a host, used with --network")
	flag.String("delnetwork", "", "delete a network")
	flag.Bool("displayconfig", false, "display configuration")
	flag.String("desc", "", "description of network, used with --addnetwork and --cidr")
	flag.Bool("help", false, "display help information")
	flag.String("host", "", "display details for a specific host")
	flag.String("ip", "", "ipv4 address of new host")
	flag.String("ipv6", "", "ipv6 address of new host")
	flag.Bool("json", false, "output in json")
	listenIp := flag.String("listenip", "", "ip address for webservice to bind to")
	listenPort := flag.String("listenport", "", "port for webservice to listen upon")
	flag.Bool("listnetworks", false, "list all networks")
	flag.Bool("showmac", false, "show mac addresses of hosts")
	flag.String("mac", "", "mac address of host")
	flag.String("network", "", "display hosts within a particular network")
	flag.String("newnetwork", "", "new network for host")
	flag.Bool("setupdb", false, "setup a new database")
	flag.String("short1", "", "short1 hostname")
	flag.String("short2", "", "short2 hostname")
	flag.String("short3", "", "short3 hostname")
	flag.String("short4", "", "short4 hostname")
	flag.Bool("showheader", false, "print header file before printing non-json output")
	flag.Bool("startweb", false, "start web service using config file setting for EnableTLS")
	flag.Bool("starthttp", false, "start http web service")
	flag.Bool("starthttps", false, "start https web service")
	flag.String("updatehost", "", "host to update")
	flag.String("updatenetwork", "", "network to update")
	flag.Bool("version", false, "display version information")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.AddConfigPath(".")

	if *configFile == "" {
		viper.SetConfigName("narcotk-hosts-config")
	} else {
		viper.SetConfigName(*configFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
		viper.SetDefault("ShowHeader", false)
		viper.SetDefault("ListenPort", "23000")
		viper.SetDefault("ListenIP", "127.0.0.1")
		viper.SetDefault("Verbose", true)
		viper.SetDefault("Database", "./narcotk_hosts_all.db")
		viper.SetDefault("DatabaseType", "sqlite3")
		viper.SetDefault("HeaderFile", "./header.txt")
		viper.SetDefault("IndexFile", "")
		viper.SetDefault("Files", "./files")
		viper.SetDefault("JSON", false)
		viper.SetDefault("EnableTLS", false)
		viper.SetDefault("TLSCert", "./tls/server.crt")
		viper.SetDefault("TLSKey", "./tls/server.key")
		viper.SetDefault("RegistrationKey", "")
		viper.SetDefault("Verbose", true)
	}

	if *listenPort != "" {
		viper.Set("ListenPort", listenPort)
	}

	if *listenIp != "" {
		viper.Set("ListenIP", listenIp)
	}
}

func initDb(databaseFile string, databaseType string) {
	fmt.Println("*** initDb")
	var err error
	db, err = sql.Open(databaseType, databaseFile)
	showerror("cannot open database", err, "warn")

	err = db.Ping()
	showerror("cannot connect to database", err, "warn")
}

func main() {
	fmt.Println("Starting main function")

	if viper.GetBool("setupdb") {
		setupdb()
		os.Exit(0)
	}

	if fileExists(viper.GetString("Database")) {
		log.Printf("Database file %s exists\n", viper.GetString("Database"))
	} else {
		log.Printf("Database file %s does not exist, exiting\n", viper.GetString("Database"))
		os.Exit(1)
	}

	if viper.GetBool("help") {
		displayHelp()
		os.Exit(0)
	}

	if viper.GetBool("displayconfig") {
		displayConfig()
		os.Exit(0)
	}

	initDb(viper.GetString("Database"), viper.GetString("DatabaseType"))

	if viper.GetBool("startweb") {
		startWeb(viper.GetString("ListenIP"), viper.GetString("ListenPort"), viper.GetBool("EnableTLS"))
		os.Exit(0)
	}

	if viper.GetBool("starthttp") {
		startWeb(viper.GetString("ListenIP"), viper.GetString("ListenPort"), false)
		os.Exit(0)
	}

	if viper.GetBool("starthttps") {
		startWeb(viper.GetString("ListenIP"), viper.GetString("ListenPort"), true)
		os.Exit(0)
	}

	if viper.GetString("delnetwork") != "" {
		delNetwork(viper.GetString("delnetwork"))
		os.Exit(0)
	}

	if viper.GetString("delhost") != "" {
		if viper.GetString("network") == "" {
			fmt.Println("Error a network must be provided")
			os.Exit(1)
		} else {
			delHost(viper.GetString("delhost"), viper.GetString("network"))
			os.Exit(0)
		}
	}

	if viper.GetString("updatenetwork") != "" {
		if (viper.GetString("network") == "") && (viper.GetString("cidr") == "") && (viper.GetString("desc") == "") {
			log.Println("Error: at least one of network, cidr or desc must be specified")
			os.Exit(1)
		} else {
			updateNetwork(viper.GetString("updatenetwork"), viper.GetString("network"), viper.GetString("cidr"), viper.GetString("desc"))
		}
		os.Exit(0)
	}

	if viper.GetBool("version") {
		displayVersion()
		os.Exit(0)
	}

	if viper.GetBool("listnetworks") {
		listNetworks(nil, "select * from networks", viper.GetBool("json"))
		os.Exit(0)
	}

	if viper.GetString("addnetwork") != "" {
		if (viper.GetString("cidr") == "") || (viper.GetString("desc") == "") {
			fmt.Println("Error: When using --addnetwork you must also provide --cidr and --desc")
			os.Exit(1)
		} else {
			addNetwork(viper.GetString("addnetwork"), viper.GetString("cidr"), viper.GetString("desc"))
			os.Exit(0)
		}
	}

	if viper.GetString("addhost") != "" {
		if (viper.GetString("network") == "") || (viper.GetString("ip") == "") {
			fmt.Println("Error: When using --addhost you must also provide --network and --ip")
			os.Exit(1)
		} else {
			addHost(viper.GetString("addhost"), viper.GetString("network"), viper.GetString("ip"), viper.GetString("ipv6"), viper.GetString("short1"), viper.GetString("short2"), viper.GetString("short3"), viper.GetString("short4"), viper.GetString("mac"))
			os.Exit(0)
		}
	}

	if viper.GetString("updatehost") != "" {
		if viper.GetString("network") == "" {
			fmt.Println("Error: When using --updatehost you must also provide --network")
			os.Exit(1)
		} else {
			updateHost(viper.GetString("updatehost"), viper.GetString("network"), viper.GetString("host"), viper.GetString("newnetwork"), viper.GetString("ip"), viper.GetString("ipv6"), viper.GetString("short1"), viper.GetString("short2"), viper.GetString("short3"), viper.GetString("short4"), viper.GetString("mac"))
			os.Exit(0)
		}
	}

	if viper.GetBool("showheader") && !viper.GetBool("json") {
		printFile(viper.GetString("HeaderFile"), nil)
	}

	if viper.GetString("network") != "" {
		listHost(nil, viper.GetString("network"), "select * from hosts where network like '"+viper.GetString("network")+"'", viper.GetBool("showmac"), viper.GetBool("json"))
		os.Exit(0)
	}

	if viper.GetString("host") != "" {
		sqlquery := "select * from hosts where fqdn like '" + viper.GetString("host") + "'"
		listHost(nil, "", sqlquery, viper.GetBool("showmac"), viper.GetBool("json"))
	} else {
		listHost(nil, viper.GetString("network"), "select * from hosts", viper.GetBool("showmac"), viper.GetBool("json"))
	}
}

func printFile(filename string, webprint http.ResponseWriter) {
	fmt.Println("Starting printFile")
	texttoprint, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR: cannot open file: %s", filename)
		if webprint != nil {
			http.Error(webprint, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
	if webprint != nil {
		fmt.Fprintf(webprint, "%s", string(texttoprint))
	} else {
		fmt.Print(string(texttoprint))
	}
}

func checkHost(host string, network string) bool {
	fmt.Println("Starting checkHost")
	sqlquery := "select * from hosts where (fqdn like '" + host + "' and network like '" + network + "')"
	fmt.Println("===" + sqlquery)
	var myhosts []Host
	rows, err := db.Query(sqlquery)
	defer rows.Close()
	showerror("error running db query", err, "warn")

	for rows.Next() {
		var network string
		var ipv4 string
		var ipv6 string
		var fqdn string
		var short1 string
		var short2 string
		var short3 string
		var short4 string
		var mac string
		err = rows.Scan(&network, &ipv4, &ipv6, &fqdn, &short1, &short2, &short3, &short4, &mac)
		showerror("cannot parse hosts results", err, "warn")
		myhosts = append(myhosts, Host{MakePaddedIp(ipv4), network, ipv4, ipv6, fqdn, short1, short2, short3, short4, mac})
		if len(myhosts) >= 1 {
			log.Printf("%d hosts found matching %s/%s", len(myhosts), host, network)
			return true
		}
		log.Println("Error: no host found for host: " + host + " ip: " + network)
		return false
	}
	return false
}

func checkNetwork(network string) bool {
	fmt.Println("Starting checkNetwork")
	sqlquery := "select * from networks where network like '" + network + "'"

	var mynetworks []SingleNetwork
	rows, err := db.Query(sqlquery)
	defer rows.Close()
	showerror("error running db query", err, "warn")

	for rows.Next() {
		var network string
		var cidr string
		var description string
		err = rows.Scan(&network, &cidr, &description)
		showerror("cannot parse network results", err, "warn")
		mynetworks = append(mynetworks, SingleNetwork{MakePaddedIp(network), network, cidr, description})
	}

	if len(mynetworks) >= 1 {
		log.Printf("%d networks found\n", len(mynetworks))
		return true
	}
	log.Println("no network found for: " + network)
	return false
}

func addHost(addhost string, network string, ip string, ipv6 string, short1 string, short2 string, short3 string, short4 string, mac string) {
	fmt.Println("Adding new host:")
	fmt.Println(addhost)
	fmt.Println(network)
	fmt.Println(ip)
	fmt.Println(ipv6)
	fmt.Println(short1)
	fmt.Println(short2)
	fmt.Println(short3)
	fmt.Println(short4)
	fmt.Println(mac)
	mac = PrepareMac(mac)
	if checkNetwork(network) && ValidIP(ip) {
		sqlquery := "insert into hosts (network, ipv4, ipv6, fqdn, short1, short2, short3, short4, mac) values ('" + network + "', '" + ip + "', '" + ipv6 + "', '" + addhost + "', '" + short1 + "', '" + short2 + "', '" + short3 + "', '" + short4 + "', '" + mac + "')"
		runSql(sqlquery)
	}
}

func updateHost(oldhost string, oldnetwork string, newhost string, newnetwork string, newipv4 string, newipv6 string, newshort1 string, newshort2 string, newshort3 string, newshort4 string, newmac string) {
	fmt.Println("Starting updateHost")
	var originalhost []Host
	var sqlquery string
	var updatesqlquery string

	// if we can find at least one host
	if checkHost(oldhost, oldnetwork) {
		sqlquery = "select * from hosts where (fqdn like '" + oldhost + "' and network like '" + oldnetwork + "')"
		rows, err := db.Query(sqlquery)
		defer rows.Close()
		showerror("error running db query", err, "warn")

		for rows.Next() {
			var network string
			var ipv4 string
			var ipv6 string
			var fqdn string
			var short1 string
			var short2 string
			var short3 string
			var short4 string
			var mac string
			err = rows.Scan(&network, &ipv4, &ipv6, &fqdn, &short1, &short2, &short3, &short4, &mac)
			showerror("cannot parse hosts results", err, "warn")
			originalhost = append(originalhost, Host{MakePaddedIp(ipv4), network, ipv4, ipv6, fqdn, short1, short2, short3, short4, mac})
		}
		if len(originalhost) != 1 {
			log.Println("Error: more than one host with identifier " + oldhost + "/" + oldnetwork + " discovered")
		} else {
			for _, host := range originalhost {
				var updatenetwork string
				var updateipv4 string
				var updateipv6 string
				var updatefqdn string
				var updateshort1 string
				var updateshort2 string
				var updateshort3 string
				var updateshort4 string
				var updatemac string
				if newnetwork == "" {
					updatenetwork = host.Network
				} else {
					updatenetwork = newnetwork
				}
				if newipv4 == "" {
					updateipv4 = host.IPv4
				} else {
					updateipv4 = newipv4
				}
				if newipv6 == "" {
					updateipv6 = host.IPv6
				} else {
					updateipv6 = newipv6
				}
				if newhost == "" {
					updatefqdn = host.Hostname
				} else {
					updatefqdn = newhost
				}
				if newshort1 == "" {
					updateshort1 = host.Short1
				} else {
					updateshort1 = newshort1
				}
				if newshort2 == "" {
					updateshort2 = host.Short2
				} else {
					updateshort2 = newshort2
				}
				if newshort3 == "" {
					updateshort3 = host.Short3
				} else {
					updateshort3 = newshort3
				}
				if newshort4 == "" {
					updateshort4 = host.Short4
				} else {
					updateshort4 = newshort4
				}
				if newmac == "" {
					updatemac = host.MAC
				} else {
					updatemac = newmac
				}
				if checkNetwork(updatenetwork) && ValidIP(newipv4) {
					updatesqlquery = "update hosts set network = '" + updatenetwork + "', ipv4 = '" + updateipv4 + "', ipv6 = '" + updateipv6 + "', fqdn = '" + updatefqdn + "', short1 = '" + updateshort1 + "', short2 = '" + updateshort2 + "', short3 = '" + updateshort3 + "', short4 = '" + updateshort4 + "', mac = '" + updatemac + "' where fqdn like '" + oldhost + "' and network like '" + oldnetwork + "'"
					runSql(updatesqlquery)
					fmt.Println("NEW UPDATE=" + updatesqlquery)
				}
			}
		}
	} else {
		log.Printf("Error: Could not find host %s/%s", oldhost, oldnetwork)
	}
}

func delHost(host string, network string) {
	fmt.Println("Deleting host:")
	fmt.Println(host)
	fmt.Println(network)
	if checkHost(host, network) {
		sqlquery := "delete from hosts where (fqdn like '" + host + "') and (network like '" + network + "')"
		runSql(sqlquery)
	} else {
		fmt.Println("ERROR: host " + host + " / " + network + " not found")
	}
}

func addNetwork(network string, cidr string, desc string) {
	fmt.Println("Adding new network: " + network + "\nCIDR: " + cidr + "\nDescription: " + desc)
	sqlquery := "insert into networks (network, cidr, description) values ('" + network + "', '" + cidr + "', '" + desc + "')"
	fmt.Println("addNetwork query: " + sqlquery)
	runSql(sqlquery)
}

func delNetwork(network string) {
	fmt.Println("Deleting network: " + network)
	if checkNetwork(network) {
		sqlquery := "delete from networks where network like '" + network + "'"
		runSql(sqlquery)
	} else {
		fmt.Println("ERROR: network " + network + " not found")
	}
}

func runSql(sqlquery string) {
	fmt.Println("Running generic runSql function")
	fmt.Println("runSql query: " + sqlquery)

	if ParseSql(sqlquery) {
		_, err := db.Exec(sqlquery)
		if err != nil {
			log.Printf("ERROR: error executing squery: %s: %q\n", sqlquery, err)
			return
		}
	}
}

// ParseSql checks whether the sql generated is valid
func ParseSql(sqlquery string) bool {
	//log.Println("Starting ParseSql")
	_, err := sqlparser.Parse(sqlquery)
	showerror("error parsing query", err, "warn")
	if err != nil {
		log.Printf("ERROR: detected in sql: \"%s\" :%s\n", sqlquery, err)
		return false
	}
	return true
}

func listNetworks(webprint http.ResponseWriter, sqlquery string, printjson bool) {
	fmt.Println("Starting listNetworks")
	if webprint == nil {
		fmt.Println("webprint is null, printing to std out")
	}
	var mynetworks []SingleNetwork
	rows, err := db.Query(sqlquery)
	defer rows.Close()
	showerror("error running db query", err, "warn")

	for rows.Next() {
		var network string
		var cidr string
		var description string
		err = rows.Scan(&network, &cidr, &description)
		showerror("cannot parse network results", err, "warn")
		mynetworks = append(mynetworks, SingleNetwork{MakePaddedIp(network), network, cidr, description})
	}

	if len(mynetworks) > 0 {
		log.Printf("%d networks found\n", len(mynetworks))

		sort.Slice(mynetworks, func(i, j int) bool {
			return bytes.Compare([]byte(mynetworks[i].PaddedNetwork), []byte(mynetworks[j].PaddedNetwork)) < 0
		})

		if printjson {
			c, err := json.Marshal(mynetworks)
			showerror("cannot marshal json", err, "warn")
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
	} else {
		log.Println("no networks found")
		if webprint != nil {
			http.Error(webprint, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

func displayVersion() {
	fmt.Println("narcotk-hosts: 0.1")
}

func setupdb() {
	fmt.Printf("Setting up a new database: %s / %s", viper.GetString("Database"), viper.GetString("DatabaseType"))
	sqlquery := `
  CREATE TABLE hosts (
    network text NOT NULL,
    ipv4 text DEFAULT '',
    ipv6 text DEFAULT '',
    fqdn text NOT NULL,
    short1 text DEFAULT '',
    short2 text DEFAULT '',
    short3 text DEFAULT '',
    short4 text DEFAULT '',
    mac text DEFAULT '')`
	runSql(sqlquery)
	sqlquery = `
  CREATE TABLE networks (
    network text PRIMARY KEY,
    cidr text NOT NULL,
    description text NOT NULL DEFAULT '')`
	runSql(sqlquery)
}

func updateNetwork(oldnetwork string, newnetwork string, cidr string, desc string) {
	log.Println("Starting updateNetwork")
	// check if something already exists and load in to struct Network
	var originalnetwork []SingleNetwork
	var sqlquery string
	var updatesqlquery string

	// check that oldnetwork exists
	if checkNetwork(oldnetwork) {
		sqlquery = "select * from networks where network like '" + oldnetwork + "'"

		rows, err := db.Query(sqlquery)
		defer rows.Close()
		showerror("error running db query", err, "warn")

		for rows.Next() {
			var network string
			var cidr string
			var description string
			err = rows.Scan(&network, &cidr, &description)
			showerror("cannot parse hosts results", err, "warn")
			originalnetwork = append(originalnetwork, SingleNetwork{MakePaddedIp(network), network, cidr, description})
		}

		if len(originalnetwork) != 1 {
			log.Println("Error, more than one network with identifier " + oldnetwork + " discovered")
		} else {
			for _, network := range originalnetwork {
				var updatenetwork string
				var updatecidr string
				var updatedesc string
				if newnetwork == "" {
					updatenetwork = network.Network
				} else {
					updatenetwork = newnetwork
				}
				if cidr == "" {
					updatecidr = network.CIDR
				} else {
					updatecidr = cidr
				}
				if desc == "" {
					updatedesc = network.Description
				} else {
					updatedesc = desc
				}
				updatesqlquery = "update networks set network = '" + updatenetwork + "', cidr = '" + updatecidr + "', description = '" + updatedesc + "' where network like '" + oldnetwork + "'"
			}
		}
		runSql(updatesqlquery)
	} else {
		log.Println("Error updating network: \"" + oldnetwork + "\" does not exist")
		os.Exit(1)
	}
}

func listHost(webprint http.ResponseWriter, network string, sqlquery string, showmac bool, printjson bool) {
	log.Println("Starting listHost")
	if webprint == nil {
		fmt.Println("webprint is null, printing to std out")
	}

	log.Println("succeed ParseSql on ", sqlquery)
	var myhosts []Host
	rows, err := db.Query(sqlquery)
	defer rows.Close()
	showerror("error running db query", err, "warn")

	//log.Println("err = ", err)
	log.Println("rows = ", rows)
	for rows.Next() {
		var network string
		var ipv4 string
		var ipv6 string
		var fqdn string
		var short1 string
		var short2 string
		var short3 string
		var short4 string
		var mac string
		err = rows.Scan(&network, &ipv4, &ipv6, &fqdn, &short1, &short2, &short3, &short4, &mac)
		showerror("cannot parse hosts results", err, "warn")
		myhosts = append(myhosts, Host{MakePaddedIp(ipv4), network, ipv4, ipv6, fqdn, short1, short2, short3, short4, mac})
	}

	if len(myhosts) > 0 {
		log.Printf("%d hosts found\n", len(myhosts))

		sort.Slice(myhosts, func(i, j int) bool {
			return bytes.Compare([]byte(myhosts[i].PaddedIP), []byte(myhosts[j].PaddedIP)) < 0
		})
		if printjson {
			// print json
			c, err := json.Marshal(myhosts)
			showerror("cannot marshal json", err, "warn")
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
						fmt.Printf("%-17s  %-15s    %s  %s  %s  %s  %s\n", host.MAC, host.IPv4, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
					}
				} else {
					for _, host := range myhosts {
						fmt.Printf("%-15s    %s  %s  %s  %s  %s\n", host.IPv4, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
					}
				}
			} else {
				// webprint
				if showmac {
					log.Println("webprint=y showmac=y")
					for _, host := range myhosts {
						log.Println("webprint=y showmac=y")
						fmt.Fprintf(webprint, "%-17s  %-15s    %s  %s  %s  %s  %s\n", host.MAC, host.IPv4, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
					}
				} else {
					log.Println("webprint=y showmac=n")
					for _, host := range myhosts {
						fmt.Fprintf(webprint, "%-15s    %s  %s  %s  %s  %s\n", host.IPv4, host.Hostname, host.Short1, host.Short2, host.Short3, host.Short4)
					}
				}
			}
		}
	} else {
		log.Println("host not found")
		if webprint != nil {
			http.Error(webprint, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

// BreakIp chops an IP into parts
func BreakIp(ipv4 string, position int) string {
	deliminator := func(c rune) bool {
		return (c == '.')
	}
	ipArray := strings.FieldsFunc(ipv4, deliminator)
	return ipArray[position]
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("MIDDLEWARE: ", r.RemoteAddr, " ", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func startWeb(listenip string, listenport string, usetls bool) {
	r := mux.NewRouter()

	if viper.GetString("IndexFile") != "" {
		r.HandleFunc("/", handlerIndex)
	}

	hostsRouter := r.PathPrefix("/hosts").Subrouter()
	hostsRouter.HandleFunc("", handlerHosts)
	hostsRouter.HandleFunc("/{network}", handlerHosts)
	hostsRouter.Use(loggingMiddleware)

	hostRouter := r.PathPrefix("/host").Subrouter()
	hostRouter.HandleFunc("/{host}", handlerHostFile).Queries("file", "")
	hostRouter.HandleFunc("/{host}", handlerHost)
	hostRouter.Use(loggingMiddleware)

	networksRouter := r.PathPrefix("/networks").Subrouter()
	networksRouter.HandleFunc("", handlerNetworks)
	networksRouter.Use(loggingMiddleware)

	networkRouter := r.PathPrefix("/network").Subrouter()
	networkRouter.HandleFunc("/{network}", handlerNetwork)
	networkRouter.Use(loggingMiddleware)

	ipRouter := r.PathPrefix("/ip").Subrouter()
	ipRouter.HandleFunc("/{ip}", handlerIp)
	ipRouter.Use(loggingMiddleware)

	macRouter := r.PathPrefix("/mac").Subrouter()
	macRouter.HandleFunc("/{mac}", handlerMac)
	macRouter.Use(loggingMiddleware)

	if viper.GetString("RegistrationKey") != "" {
		// https://stackoverflow.com/questions/43379942/how-to-have-an-optional-query-in-get-request-using-gorilla-mux
		r.HandleFunc("/register", handlerRegister).Methods("GET")
	}

	if usetls {
		log.Println("Starting HTTPS Webserver: " + listenip + ":" + listenport)
		err := http.ListenAndServeTLS(listenip+":"+listenport, viper.GetString("tlscert"), viper.GetString("tlskey"), r)
		showerror("cannot start https server", err, "fatal")
	} else {
		log.Println("Starting HTTP Webserver: " + listenip + ":" + listenport)
		err := http.ListenAndServe(listenip+":"+listenport, r)
		showerror("cannot start http server", err, "fatal")
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerIndex")
	printFile(viper.GetString("IndexFile"), w)
}

func handlerHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerHosts")
	vars := mux.Vars(r)
	queries := r.URL.Query()
	log.Printf("vars = %q\n", vars)
	log.Printf("queries = %q\n", queries)

	sqlquery := ""

	if vars["network"] == "" {
		sqlquery = "select * from hosts"
	} else {
		sqlquery = "select * from hosts where network like '" + vars["network"] + "'"
	}

	givejson := false
	showmac := false

	if strings.ToLower(queries.Get("json")) == "y" {
		w.Header().Set("Content-Type", "application/json")
		givejson = true
	}

	if (strings.ToLower(queries.Get("header")) == "y") && (!givejson) {
		printFile(viper.GetString("HeaderFile"), w)
	}

	if strings.ToLower(queries.Get("mac")) == "y" {
		showmac = true
	}

	listHost(w, viper.GetString("network"), sqlquery, showmac, givejson)

}

func handlerHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queries := r.URL.Query()
	log.Printf("Starting handlerHost")

	givejson := false
	showmac := false

	if strings.ToLower(queries.Get("json")) == "y" {
		w.Header().Set("Content-Type", "application/json")
		givejson = true
	}

	if (strings.ToLower(queries.Get("header")) == "y") && (!givejson) {
		printFile(viper.GetString("HeaderFile"), w)
	}

	if strings.ToLower(queries.Get("mac")) == "y" {
		showmac = true
	}

	// problem that when passing mac=y it does not print the mac
	sqlquery := "select * from hosts where fqdn like '" + vars["host"] + "'"
	log.Println("sqlquery = ", sqlquery)
	listHost(w, viper.GetString("network"), sqlquery, showmac, givejson)
}

func handlerHostFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting NewhandlerHostFile")
	vars := mux.Vars(r)
	queries := r.URL.Query()

	if queries.Get("file") == "" {
		// error file doesnt exist, return 404
		//w.WriteHeader(http.StatusNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	printFile(viper.GetString("files")+"/"+vars["host"]+"."+queries.Get("file"), w)
}

func handlerNetworks(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerNetworks")
	queries := r.URL.Query()
	sqlquery := "select * from networks"
	givejson := false

	log.Printf("queries = %q\n", queries)

	if strings.ToLower(queries.Get("json")) == "y" {
		givejson = true
		w.Header().Set("Content-Type", "application/json")
	}

	listNetworks(w, sqlquery, givejson)

}

func handlerNetwork(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerNetwork")
	vars := mux.Vars(r)
	queries := r.URL.Query()
	givejson := false

	log.Printf("queries = %q\n", queries)

	if strings.ToLower(queries.Get("json")) == "y" {
		givejson = true
		w.Header().Set("Content-Type", "application/json")
	}

	sqlquery := "select * from networks where network like '" + vars["network"] + "'"

	listNetworks(w, sqlquery, givejson)

}

func handlerIp(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerIp")
	vars := mux.Vars(r)
	queries := r.URL.Query()

	givejson := false
	showmac := false

	log.Printf("queries = %q\n", queries)

	if strings.ToLower(queries.Get("json")) == "y" {
		givejson = true
		w.Header().Set("Content-Type", "application/json")
	}

	if (strings.ToLower(queries.Get("header")) == "y") && (!givejson) {
		printFile(viper.GetString("HeaderFile"), w)
	}

	if strings.ToLower(queries.Get("mac")) == "y" {
		showmac = true
	}

	//sqlquery := "select * from hosts where ipv4 like '" + vars["ip"] + "'"
	sqlquery := "select * from hosts where (ipv4 like '" + vars["ip"] + "') or (ipv6 like '" + vars["ip"] + "')"

	listHost(w, "", sqlquery, showmac, givejson)
}

func handlerMac(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting handlerMac")
	vars := mux.Vars(r)
	queries := r.URL.Query()

	givejson := false
	showmac := false

	log.Printf("queries = %q\n", queries)

	if strings.ToLower(queries.Get("json")) == "y" {
		givejson = true
		w.Header().Set("Content-Type", "application/json")
	}

	if (strings.ToLower(queries.Get("header")) == "y") && (!givejson) {
		printFile(viper.GetString("HeaderFile"), w)
	}

	if strings.ToLower(queries.Get("mac")) == "y" {
		showmac = true
	}

	sqlquery := "select * from hosts where mac like '" + PrepareMac(vars["mac"]) + "'"
	listHost(w, "", sqlquery, showmac, givejson)
}

func handlerRegister(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	regkey := vars.Get("key")
	if regkey == viper.GetString("RegistrationKey") {
		fqdn := vars.Get("fqdn")
		ip := vars.Get("ip")
		ipv6 := vars.Get("ipv6")
		nw := vars.Get("nw")
		mac := PrepareMac(vars.Get("mac"))
		short1 := vars.Get("s1")
		short2 := vars.Get("s2")
		short3 := vars.Get("s3")
		short4 := vars.Get("s4")
		fmt.Println(vars)
		fmt.Printf("Starting handlerRegister: fqdn=%s / ip=%s / ipv6=%s / nw=%s / mac=%s / s1=%s / s2=%s / s3=%s / s4=%s\n", fqdn, ip, ipv6, nw, mac, short1, short2, short3, short4)
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		if (fqdn == "") || (ip == "") || (nw == "") {
			log.Printf("Error: fqdn, ip or nw cannot be blank")
		} else {
			if ValidIP(ip) {
				addHost(fqdn, nw, ip, ipv6, short1, short2, short3, short4, mac)
				fmt.Fprintf(w, "Added: %s", vars)
			}
		}
	} else {
		// https://golang.org/src/net/http/status.go
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		log.Printf("RegistrationKey invalid (%s), ignoring", regkey)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ValidIP makes sure that an IP is valid
func ValidIP(ip string) bool {
	if net.ParseIP(ip) != nil {
		return true
	}
	log.Printf("Error: ip %s is not valid", ip)
	return false
}

// PadLeft prefixs a string with 0's
func PadLeft(str string) string {
	for {
		padding := "00"
		str = padding + str
		startpoint := len(str) - 3
		endpoint := len(str)
		return str[startpoint:endpoint]
	}
}

// MakePaddedIp is used to create a standardised number that is then used to sort the ips
func MakePaddedIp(ipv4 string) string {
	f := func(c rune) bool {
		return (c == rune('.'))
	}

	s := strings.FieldsFunc(ipv4, f)

	paddedIp := ""

	for i := 0; i < len(s); i++ {
		paddedIp = paddedIp + PadLeft(s[i])
	}

	return paddedIp
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
      ** --updatenetwork is mandatory, the other params are optional

  Display a host:
      --host=server1.domain.com

  Adding a host:
      --addhost=server-1-199.domain.com --network=192.168.1 --ip=192.168.1.13 --ipv6=::6 --short1=server-1-199 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe

  Update a host:
      --updatehost=server-1-199.domain.com --network=192.168.1 --host=server-1-200.domain.com --newnetwork=192.168.1 --ip=192.168.1.200 --ipv6=::6 --short1=server-1-200 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe
      ** --updatehost and --network are mandatory, other params are optional 

  Delete a host:
      --delhost=server-1-200.domain.com --network=192.168.1

  Configuration file:
      --configfile=/path/to/file.yaml

  Database file:
      --database=/path/to/somefile.db

  Setup a new blank database file:
      --setupdb  --database=./newfile.db

  Start Web Service using config file EnableTLS setting:
      --startweb

  Start Web Service using:
      --starthttp

  Start Web HTTPS Service:
      --starthttps

  Port to listen upon:
      --listenport=23000

  IP Address to listen upon:
      --listenip=10.0.0.14
`
	fmt.Printf("%s", helpmessage)

	os.Exit(0)
}
