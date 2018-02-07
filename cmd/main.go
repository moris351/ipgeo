package main 

import(
	"fmt"
	"github.com/moris351/ipgeo"
	"flag"
	_"strconv"

)
var (
	initdb=flag.Bool("initdb",false,"init database? it will take several minutes")
)

const LOCATION_FILE_NAME = "GeoLite2-City-Locations.csv"
const BLOCK_FILE_NAME = "GeoLite2-City-Blocks.csv"
const BLOCK_SAMPLE_FILE_NAME = "GeoLite2-City-Blocks-sample.csv"

const DBNAME="geoip.store"

func main(){
	flag.Parse()

	if *initdb == true{
		ipgeo.Remove(DBNAME)
	}
	l:=ipgeo.Locator(DBNAME)
	if l==nil{
		return
	}
	fmt.Println(l)

	fmt.Println(*initdb)
	if *initdb==true{
		if err:=l.InitDB(LOCATION_FILE_NAME,BLOCK_SAMPLE_FILE_NAME);err!=nil{
			return
		}
	}
	
	ips:= []string{
					"73.183.190.129",
					"73.183.190.127",
					"73.183.190.254",
					"2.101.221.9",
					"183.193.153.190",
				}

	for _,v := range ips{

		geo,err:=l.FindGeo(v)
		if err!=nil{
			fmt.Println("FindGeo return err:",err)
		}else{
			fmt.Println(v,geo)
		}
	}
//	l.Show()
	l.Close()
}

