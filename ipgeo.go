package ipgeo

import (
	"bufio"
	_ "bytes"
	_ "encoding/binary"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"os"
	"strings"
	"time"
)

type Continent struct {
	Code string `json:code`
	Name string `json:name`
}

type Country struct {
	Code          string `json:code`
	Name          string `json:name`
	ContinentCode string `json:ctn_code`
}
type Sub struct {
	Code    string `json:code`
	Name    string `json:name`
	SubCode string `json:sub_code`
}

type City struct {
	Code     string `json:code`
	Name     string `json:name`
	CityCode string `json:cit_code`
}

type IpLocator struct {
	DB *bolt.DB
}

const ipn = "__ips"

var bktsn = []string{"__ctns", "__ctrs", "__subs", "__cits"}

const (
	ctnsn = iota
	ctrsn
	subsn
	citsn
	bktsnum
)

var lb *logBot = newLogBot("log/ipgeo.log")

const MAX_BATCH_NUM = 1e6
const MAX_RECORDS_NUM = 1e5

var locator *IpLocator 

func Locator(dbname string) *IpLocator {

	lb.SetLevel("debug")
	if locator != nil {
		return locator
	}
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		lb.Debug("boltdb open failed")
		return nil
	}
	locator=&IpLocator{db}
	return locator
}
func Remove(dbname string) error {
	err := os.Remove(dbname)
	if err != nil {
		lb.Debug("remove db file failed")
		return err
	}

	return nil
}

func (il *IpLocator) Close() {
	il.DB.Close()
}

var ErrRecordNotFound = errors.New("Record not found")
var ErrBucketNotFound = errors.New("Bucket not found")

func (il *IpLocator) FindGeo(ip string) (string, error) {
	city := ""
	err := il.DB.View(func(tx *bolt.Tx) error {
		bips := tx.Bucket([]byte(ipn))
		if bips == nil {
			return ErrBucketNotFound
		}
		j := 0
		var code []byte
		for j = 0; j < 32; j++ {
			b := IP2IPNet(ip, uint32(j))
			code = bips.Get(b)
			if code != nil {
				break
			}
		}
		if j == 32 {
			return ErrRecordNotFound
		}

		bcits := tx.Bucket([]byte(bktsn[citsn]))
		if bcits == nil {
			return ErrBucketNotFound
		}
		cit := bcits.Get(code)
		if cit == nil {
			return ErrRecordNotFound
		}
		city = string(cit)
		return nil
	})
	if err != nil {
		return "", err
	}

	//fmt.Println("FindGeo:city:",city)
	return city, nil
}

func (il *IpLocator) Stats() {
	lb.Debug("Stats")
	go func() {
		// Grab the initial stats.
		prev := il.DB.Stats()

		for {
			// Wait for 10s.
			time.Sleep(1 * time.Second)

			// Grab the current stats and diff them.
			stats := il.DB.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
			//json.NewEncoder(os.Stderr).Encode(diff)

			jdiff, _ := json.Marshal(diff)

			// Save stats for the next loop.
			prev = stats
			fmt.Println(string(jdiff))
		}
	}()

}

func (il *IpLocator) InitDB(locFilename string, blockFilename string) error {

	lb.Debug("InitDB")
	lf, err := os.Open(locFilename)
	defer lf.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", locFilename, err)
		return err
	}

	reader := csv.NewReader(lf)
	// lb.Debug(reader.Read())
	reader.Read() // discard first line
	if err != nil {
		lb.Fatal("Reading %s failed with err:%v\n", locFilename, err)
	}

	lb.Debug("init geo data")
	err = il.DB.Update(func(tx *bolt.Tx) error {
		var bkts [bktsnum]*bolt.Bucket
		for j := 0; j < bktsnum; j++ {
			bkts[j], err = tx.CreateBucketIfNotExists([]byte(bktsn[j]))
			if err != nil {
				lb.Fatal(err)
			}
		}

		i := 0
		for err != io.EOF {

			i++
			cls, err := reader.Read()
			if err == io.EOF {
				break
			}

			//lb.Printf("read %v,len=%d\n", cls, len(cls))
			//for k, item := range cls {
			//lb.Printf("k=%d,item=%s\n", k, item)
			//}

			if err == nil {
				if len(cls[2]) == 0 {
					continue
				}

				jctn, err := json.Marshal(&Continent{cls[2], cls[3]})
				if err != nil {
					return err
				}
				err = bkts[ctnsn].Put([]byte(cls[2]), jctn)
				if err != nil {
					return err
				}

				jctr, err := json.Marshal(&Country{cls[4], cls[5], cls[2]})
				if err != nil {
					return err
				}
				bkts[ctrsn].Put([]byte(cls[4]), jctr)
				if err != nil {
					return err
				}

				jsub, err := json.Marshal(&Sub{cls[6], cls[7], cls[4]})
				if err != nil {
					return err
				}
				bkts[subsn].Put([]byte(cls[6]), jsub)
				if err != nil {
					return err
				}

				jcit, err := json.Marshal(&City{cls[0], cls[10], cls[6]})
				//fmt.Println(string(jcit))
				if err != nil {
					return err
				}
				bkts[citsn].Put([]byte(cls[0]), jcit)
				if err != nil {
					return err
				}

			}

		}
		return err
	})
	lb.Debug("init ip data")
	lf, err = os.Open(blockFilename)
	defer lf.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", blockFilename, err)
		return err
	}

	reader = csv.NewReader(lf)
	// lb.Debug(reader.Read())
	reader.Read() // discard first line
	if err != nil {
		lb.Fatal("Reading %s failed with err:%v\n", blockFilename, err)
	}

	lb.Debug("begin update")
	start := time.Now()
	err = il.DB.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte(ipn))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		lb.Debug(err)
		return err
	}
	for j := 0; j < MAX_BATCH_NUM; j++ {

		//if j>100 { break }
		//lb.Debug("line :",j)

		err = il.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(ipn))
			if b == nil {
				return errors.New("not found")
			}

			for i := 0; i < MAX_RECORDS_NUM && err != io.EOF; i++ {
				cls, err := reader.Read()
				if err == io.EOF {
					j = MAX_BATCH_NUM
					break
				}

				if err != nil {
					return err
				}

				if len(cls[0]) == 0 {
					continue
				}

				ipnet := IPNet(cls[0])
				//lb.Debug("ipnet:", ipnet)
				if ipnet == nil {
					continue
				}

				err = b.Put(ipnet, []byte(cls[1]))
				if err != nil {
					return err
				}
			}
			return err

		})
	}
	if err != nil {
		lb.Debug(err)
		return err
	}
	lb.Info("update ips data cost:", time.Since(start))
	return nil
}
func (il *IpLocator) FindFile(input string, output string) error {
	lb.Debug("FindFile")
	lfi, err := os.Open(input)
	defer lfi.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", input, err)
		return err
	}
	lfo, err := os.Create(output)
	defer lfo.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", output, err)
		return err
	}

	reader := csv.NewReader(lfi)
	writer := bufio.NewWriter(lfo)
	i := 0
	for {
		//if i>100 { break}
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		ip := line[0]

		geo, err := il.FindGeo(ip)
		buf := fmt.Sprintf("%s %s \n", ip, geo)
		writer.Write([]byte(buf))
		i++

	}

	writer.Flush()
	return nil

}

func GetIps(input string, output string) error {

	lb.Debug("GetIps")
	lfi, err := os.Open(input)
	defer lfi.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", input, err)
		return err
	}
	lfo, err := os.Create(output)
	defer lfo.Close()
	if err != nil {
		lb.Debug("open %s failed with err: %v\n", output, err)
		return err
	}

	reader := csv.NewReader(lfi)
	// lb.Debug(reader.Read())
	reader.Read() // discard first line
	if err != nil {
		lb.Fatal("Reading %s failed with err:%v\n", input, err)
	}
	writer := csv.NewWriter(lfo)
	if err != nil {
		lb.Fatal("Reading %s failed with err:%v\n", output, err)
	}
	i := 0
	for {
		//if i>100 { break}
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		ip := line[0]
		n := strings.LastIndex(ip, ".")
		ip = ip[0:n]

		ip = fmt.Sprintf("%s%s", ip, ".9")

		writer.Write([]string{ip})
		i++

	}

	writer.Flush()
	return nil

}

