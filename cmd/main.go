package main 

import(
	"fmt"
	"github.com/moris351/ipgeo"
	"flag"
	"time"
	_"strconv"

)
var (
	initdb=flag.Bool("initdb",false,"init database? it will take several minutes")
	findgeo=flag.String("findgeo","","find ip geo info")
	findfile=flag.String("findfile","","find ip from file")
	stats=flag.Bool("stats",false,"real time show db stats")
	getips=flag.String("getips","","get ips from a csv file and output to a file")

)

const LOCATION_FILE_NAME = "GeoLite2-City-Locations.csv"
const BLOCK_FILE_NAME = "GeoLite2-City-Blocks.csv"
const BLOCK_SAMPLE_FILE_NAME = "GeoLite2-City-Blocks-sample.csv"

const DBNAME="geoip.store"
/*
	ips:= []string{
					"73.183.190.129",
					"73.183.190.127",
					"73.183.190.254",
					"2.101.221.9",
					"183.193.153.190",
				}
*/

func main(){
	flag.Parse()
	switch{
	case *initdb==true:
		ipgeo.Remove(DBNAME)
		l:=ipgeo.Locator(DBNAME)
		defer l.Close()
		if l==nil{
			return
		}

		if *stats==true{
			l.Stats()
		}
		start := time.Now()
		if err:=l.InitDB(LOCATION_FILE_NAME,BLOCK_FILE_NAME);err!=nil{
			return
		}
		fmt.Println("findgeo cost:", time.Since(start))

	case len(*findgeo)!=0:
		l:=ipgeo.Locator(DBNAME)
		defer l.Close()
		if l==nil{
			return
		}

		if *stats==true{
			l.Stats()
		}
		v:=*findgeo

		start := time.Now()
		geo,err:=l.FindGeo(v)
		if err!=nil{
			fmt.Println("FindGeo return err:",err)
		}else{
			fmt.Println(v,geo)
		}
		fmt.Println("findgeo cost:", time.Since(start))
	case len(*getips)!=0:
		if err:=ipgeo.GetIps(*getips,"ips");err!=nil{
			fmt.Println(err)
		}
	case len(*findfile)!=0:
		l:=ipgeo.Locator(DBNAME)
		defer l.Close()
		
		start := time.Now()
		if err:=l.FindFile(*findfile,"geo");err!=nil{
			fmt.Println(err)
		}
		fmt.Println("findfile cost:", time.Since(start))
	default:
		l:=ipgeo.Locator(DBNAME)
		defer l.Close()
		if l==nil{
			return
		}
		l.Stats()

	}
}

