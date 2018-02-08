package main 

import(
	"fmt"
	"github.com/moris351/ipgeo"
	"flag"
	_"strconv"

)
var (
	initdb=flag.Bool("initdb",false,"init database? it will take several minutes")
	findgeo=flag.String("findgeo","","find ip geo info")
	findfile=flag.String("findfile","","find ip from file")
	stats=flag.Bool("stats",false,"real time show db stats")

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
		if l==nil{
			return
		}

		if *stats==true{
			l.Stats()
		}
		if err:=l.InitDB(LOCATION_FILE_NAME,BLOCK_FILE_NAME);err!=nil{
			return
		}

		l.Close()
	case len(*findgeo)!=0:
		l:=ipgeo.Locator(DBNAME)
		if ;l==nil{
			return
		}

		v:=*findgeo
		geo,err:=l.FindGeo(v)
		if err!=nil{
			fmt.Println("FindGeo return err:",err)
		}else{
			fmt.Println(v,geo)
		}
		if *stats==true{
			l.Stats()
		}
		l.Close()
	default:
		l:=ipgeo.Locator(DBNAME)
		if ;l==nil{
			return
		}
		l.Stats()

	}
}

