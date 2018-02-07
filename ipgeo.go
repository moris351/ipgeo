package ipgeo

import (
	_ "bytes"
	"encoding/csv"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"log"
	"os"
)

type ipgeo struct {
	ip int
	id int
}

type continent struct {
	code int
	name string
}

type country struct {
	code           int
	name           string
	continent_code int
}
type subdivision struct {
	code         string
	name         string
	country_code int
}

type city struct {
	code         int
	name         string
	country_code int
}

type ipmap struct {
	id        int
	city_code int
}

type IpLocator struct {
	DB *bolt.DB
}

const LOCATION_FILE_NAME = "GeoLite2-City-Locations.csv"
const BLOCK_FILE_NAME = "GeoLite2-City-Blocks.csv"

const continents = "continents"
const countries = "countries"

func Locator() *IpLocator {
	db, err := bolt.Open("geoip.store", 0600, nil)
	if err != nil {
		fmt.Println("boltdb open failed")
		return nil
	}
	il := &IpLocator{}
	il.DB = db

	return il
}

func (i *IpLocator) Close() {

	i.DB.Close()
}
func (i *IpLocator) Show() {

	_ = i.DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			fmt.Println(string(name))
			return nil
		})

	})
	_ = i.DB.View(func(tx *bolt.Tx) error {
		ctns := tx.Bucket([]byte(continents))
		if ctns == nil {
			fmt.Println("no ", continents)
			return nil
		}
		return ctns.ForEach(func(k, v []byte) error {
			fmt.Println(string(k), string(v))
			ctn := ctns.Bucket(k)
			return ctn.ForEach(func(k, v []byte) error {
				fmt.Println(string(k), string(v))
				ctrs := ctn.Bucket([]byte(countries))
				return ctrs.ForEach(func(k, v []byte) error {
					ctr := ctrs.Bucket(k)
					fmt.Println(string(k), string(v))
					subs := ctr.Bucket([]byte("subs"))
					return subs.ForEach(func(k, v []byte) error {
						fmt.Println(string(k), string(v))
						sub := subs.Bucket(k)
						cities := sub.Bucket([]byte("cities"))
						return cities.ForEach(func(k, v []byte) error {
							fmt.Println(string(k), string(v))
							return nil

						})
						return nil
					})
					return nil
				})
				return nil
			})
			return nil
		})

		return nil
	})

	return
}

func (i *IpLocator) InitDB() error {

	lf, err := os.Open(LOCATION_FILE_NAME)
	if err != nil {
		log.Println("open %s failed with err: %v\n", LOCATION_FILE_NAME, err)
		return err
	}

	reader := csv.NewReader(lf)
	// fmt.Println(reader.Read())
	reader.Read() // discard first line
	if err != nil {
		log.Fatal("Reading %s failed with err:%v\n", LOCATION_FILE_NAME, err)
	}

	err = i.DB.Update(func(tx *bolt.Tx) error {

		//bgeo, err := tx.CreateBucketIfNotExists([]byte("geo"))
		//if err != nil {
		//log.Fatal(err)
		//}
		bcontinents, err := tx.CreateBucketIfNotExists([]byte(continents))
		if err != nil {
			log.Fatal(err)
		}
		//bcountry, err := tx.CreateBucketIfNotExists([]byte(countries))
		//if err != nil {
		//log.Fatal(err)
		//}
		//bsubdivision, err := tx.CreateBucketIfNotExists([]byte("subdivision"))
		//if err != nil {
		//log.Fatal(err)
		//}
		//bcity, err := tx.CreateBucketIfNotExists([]byte("city"))
		//if err != nil {
		//log.Fatal(err)
		//}
		//recs,err:=reader.Read()
		//if err != nil {
		//return err
		//}
		i := 0
		for  err != io.EOF {

			i++
			//fmt.Println("line :",i)
			cls, err := reader.Read()
			if err == io.EOF {
				break
			}

			//fmt.Printf("read %v,len=%d\n", cls, len(cls))
			//for k, item := range cls {
				//fmt.Printf("k=%d,item=%s\n", k, item)
			//}

			if err == nil {
				if len(cls[2]) ==0 {
					continue
				}
				ctn, err := bcontinents.CreateBucketIfNotExists([]byte(cls[2]))
				if err != nil {
					return err
				}

				ctn.Put([]byte("name"), []byte(cls[3]))
				ctrs, err := ctn.CreateBucketIfNotExists([]byte(countries))
				if err != nil {
					fmt.Println("1")
					return err
				}
				if len(cls[4]) ==0 {
					continue
				}
				ctr, err := ctrs.CreateBucketIfNotExists([]byte(cls[4]))
				if err != nil {
					fmt.Println("2")
					return err
				}

				ctr.Put([]byte("name"), []byte(cls[5]))
				if err != nil {
					fmt.Println("3")
					return err
				}
				subs, err := ctr.CreateBucketIfNotExists([]byte("subs"))
				if err != nil {
					fmt.Println("4")
					return err
				}
				subkey := cls[6]
				if len(subkey) == 0 {
					subkey = "all"
				}
				sub, err := subs.CreateBucketIfNotExists([]byte(subkey))
				if err != nil {
					fmt.Println("5")
					return err
				}

				sub.Put([]byte("name"), []byte(cls[7]))
				if err != nil {
					fmt.Println("6")
					return err
				}
				cities, err := sub.CreateBucketIfNotExists([]byte("cities"))
				if err != nil {
					fmt.Println("7")
					return err
				}

				cities.Put([]byte(cls[0]), []byte(cls[10]))
				if err != nil {
					fmt.Println("8")
					return err
				}
			}
		}
		return err
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/*
read [18918 en EU Europe CY Cyprus 04 Ammochostos   Protaras  Asia/Famagusta],len=13
k=0,item=18918
k=1,item=en
k=2,item=EU
k=3,item=Europe
k=4,item=CY
k=5,item=Cyprus
k=6,item=04
k=7,item=Ammochostos
k=8,item=
k=9,item=
k=10,item=Protaras
k=11,item=
k=12,item=Asia/Famagusta
*/
