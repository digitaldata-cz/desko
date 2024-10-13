package desko

import (
	"fmt"
	"strings"
	"time"

	"github.com/spetr/hid"
)

type (
	// IcaoData - raw data from reader
	IcaoData [][]byte
	// IcaoDocument - parsed data from reader
	IcaoDocument struct {
		TransactionID    string `json:"transactionId,omitempty"`
		IcaoType         string `json:"type" xml:"type"`
		IcaoSubtype      string `json:"subtype" xml:"subtype"`
		Country          string `json:"country" xml:"country"`
		Number           string `json:"number" xml:"number"`
		NumberChecksum   string `json:"numberChecksum" xml:"numberChecksum"`
		NumberChecksumOk bool   `json:"numberChecksumOk" xml:"numberChecksumOk"`
		Name             string `json:"name" xml:"name"`
		Surname          string `json:"surname" xml:"surname"`
		Pin              string `json:"pin" xml:"pin"`
		Sex              string `json:"sex" xml:"sex"`
		Nationality      string `json:"nationality" xml:"nationality"`
		Birth            struct {
			Year       string `json:"year" xml:"year"`
			Month      string `json:"month" xml:"month"`
			Day        string `json:"day" xml:"day"`
			Checksum   string `json:"checksum" xml:"checksum"`
			ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
		} `json:"birth" xml:"birth"`
		Expire struct {
			Year       string `json:"year" xml:"year"`
			Month      string `json:"month" xml:"month"`
			Day        string `json:"day" xml:"day"`
			Checksum   string `json:"checksum" xml:"checksum"`
			ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
		} `json:"expire" xml:"expire"`
	}
)

const (
	deskoUsbVendorID  = 0x0744
	deskoUsbProductID = 0x001d
)

/*
func Open() (d *hid.Device) {
	deviceInfo, err := GetDeviceInfo()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	deskoReader, err := deviceInfo.Open()
	if err != nil {
		log.Panicln(err)
	}
	return deskoReader
}
*/

// StartReading - start reading data from DESKO reader
func StartReading(desko *hid.Device, parser func(IcaoData)) error {
	if desko == nil {
		return fmt.Errorf("DESKO reader not found")
	}

	quit := make(chan struct{})

	// Initialize DESKO reader
	if _, err := desko.Write([]byte{0x20, 0x00}); err != nil {
		return err
	}
	go func(quit chan struct{}) {

		for {
			select {
			case <-quit:
				desko.Close() // #nosec G104
				return
			default:
				desko.Write([]byte{0x30, 0x00}) // #nosec G104
				time.Sleep(3 * time.Second)
			}
		}
	}(quit)

	var data IcaoData
	r := make([]byte, 32) // HID response buffer
	for {
		readBytes, err := desko.Read(r)
		if err != nil {
			quit <- struct{}{} // Stop DESKO reader
			return err
		}
		if readBytes > 0 {
			for i := byte(2); i < r[1]+2; i++ {
				// Start of document
				if r[i] == 0x1c && r[i+1] == 0x02 {
					i++
					data = append(data[:0], []byte{}) // Initialize first line in IcaoData slice
					continue
				}
				// End of document
				if r[i] == 0x0d && r[i+1] == 0x03 && r[i+2] == 0x1d {
					parser(data)
					data = data[:0] // Flush slice
					break
				}
				// New line
				if r[i] == 0x0d {
					data = append(data, []byte{}) // Add new line to IcaoData slice
					continue
				}
				if len(data) == 0 {
					fmt.Println("DESKO reader error: Skipping data before start of document")
					time.Sleep(100 * time.Millisecond)
					continue // Skip data before start of document
				}
				data[len(data)-1] = append(data[len(data)-1], r[i])
			}
			continue // Continue reading from HID without delay
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// ParseICAO - parse raw data to struct
func ParseICAO(d IcaoData) (ret IcaoDocument) {
	var (
		rawName string
		parsed  = false
	)
	ret.IcaoType = string(d[0][0])
	ret.IcaoSubtype = string(d[0][1])
	ret.Country = string(d[0][2:5])
	ret.Birth.ChecksumOk = false
	ret.Expire.ChecksumOk = false

	switch len(d) {
	case 3:
		// 90/3 ICAO 9303
		// Obcanske prukazy EU + vyjimka pro Belgii
		if len(d[0]) == 30 {
			if ret.Country == "BEL" {
				ret.Number = string(d[0][5:14]) + string(d[0][14:17])
			} else {
				ret.Number = string(d[0][5:14])
			}
			// Vyjmenovane doklady maji v optional zone rodne cislo
			if ret.Country == "ALB" || // albania
				ret.Country == "GEO" || // georgia
				ret.Country == "GIB" || // gibraltar
				ret.Country == "LVA" || // latvia
				ret.Country == "LTU" || // lithuania
				ret.Country == "MKD" || // macedonia
				ret.Country == "MLT" || // malta
				ret.Country == "MDA" || // moldova
				ret.Country == "NLD" || // netherlands
				ret.Country == "SRB" || // serbia
				ret.Country == "SVK" || // slovakia
				ret.Country == "ESP" || // spain
				ret.Country == "UKR" { // ukraine
				ret.Pin = string(d[0][15:25])
			}
			ret.Birth.Year = string(d[1][0:2])
			ret.Birth.Month = string(d[1][2:4])
			ret.Birth.Day = string(d[1][4:6])
			ret.Birth.Checksum = string(d[1][6:7])
			ret.Sex = string(d[1][7])
			ret.Expire.Year = string(d[1][8:10])
			ret.Expire.Month = string(d[1][10:12])
			ret.Expire.Day = string(d[1][12:14])
			ret.Expire.Checksum = string(d[1][14:15])
			ret.Nationality = string(d[1][15:18])
			rawName = string(d[2][0:30])
			parsed = true
		}

	case 2:
		// 68/2 ICAO 9303
		// Stary slovensky OP
		if len(d[0]) == 34 {
			rawName = string(d[0][5:34])
			ret.Number = string(d[1][0:9])
			ret.Nationality = string(d[1][10:13])
			ret.Birth.Year = string(d[1][13:15])
			ret.Birth.Month = string(d[1][15:17])
			ret.Birth.Day = string(d[1][17:19])
			ret.Birth.Checksum = string(d[1][19:20])
			ret.Sex = string(d[1][20])
			ret.Expire.Year = string(d[1][21:23])
			ret.Expire.Month = string(d[1][23:25])
			ret.Expire.Day = string(d[1][25:27])
			ret.Expire.Checksum = string(d[1][27:28])
			parsed = true
		}

		// 72/2 ICAO 9303
		// Nemecky OP 2004
		// Francouzsky OP
		if len(d[0]) == 36 {
			if ret.Country == "FRA" {
				ret.Surname = string(d[0][5:30])
				ret.Number = string(d[1][0:12])
				ret.Name = string(d[1][13:27])
				ret.Birth.Year = string(d[1][27:29])
				ret.Birth.Month = string(d[1][29:31])
				ret.Birth.Day = string(d[1][31:33])
				ret.Birth.Checksum = string(d[1][33:34])
				ret.Sex = string(d[1][34])
				ret.Nationality = "FRA"
			} else {
				rawName = string(d[0][5:36])
				ret.Number = string(d[1][0:9])
				ret.Nationality = string(d[1][10:13])
				ret.Birth.Year = string(d[1][13:15])
				ret.Birth.Month = string(d[1][15:17])
				ret.Birth.Day = string(d[1][17:19])
				ret.Birth.Checksum = string(d[1][19:20])
				ret.Sex = string(d[1][20])
				ret.Expire.Year = string(d[1][21:23])
				ret.Expire.Month = string(d[1][23:25])
				ret.Expire.Day = string(d[1][25:27])
				ret.Expire.Checksum = string(d[1][27:28])
			}
			parsed = true
		}

		// 88/2 ICAO 9303
		// cestovni pas
		if len(d[0]) == 44 {
			rawName = string(d[0][5:34])
			ret.Number = string(d[1][0:9])
			ret.NumberChecksum = string(d[1][9:10])
			ret.Nationality = string(d[1][10:13])
			ret.Birth.Year = string(d[1][13:15])
			ret.Birth.Month = string(d[1][15:17])
			ret.Birth.Day = string(d[1][17:19])
			ret.Birth.Checksum = string(d[1][19:20])
			ret.Sex = string(d[1][20])
			ret.Expire.Year = string(d[1][21:23])
			ret.Expire.Month = string(d[1][23:25])
			ret.Expire.Day = string(d[1][25:27])
			ret.Expire.Checksum = string(d[1][27:28])
			ret.Pin = string(d[1][28:42])
			parsed = true
		}

	case 1:
		// 30/1
		// Ridicsky prukaz Estonsko
		if len(d[0]) == 30 {
			ret.Number = string(d[0][5:14])
			parsed = true
		}
	}

	if parsed {
		// Overeni kontrolniho souctu pro datumy
		if (len(ret.Birth.Year) == 2) && (len(ret.Birth.Month) == 2) && (len(ret.Birth.Day) == 2) {
			if ((ret.Birth.Year[0]&0x0f)*7+
				(ret.Birth.Year[1]&0x0f)*3+
				(ret.Birth.Month[0]&0x0f)+
				(ret.Birth.Month[1]&0x0f)*7+
				(ret.Birth.Day[0]&0x0f)*3+
				(ret.Birth.Day[1]&0x0f))%10 == ret.Birth.Checksum[0]&0x0f {
				ret.Birth.ChecksumOk = true
			}
		}

		if (len(ret.Expire.Year) == 2) && (len(ret.Expire.Month) == 2) && (len(ret.Expire.Day) == 2) {
			if ((ret.Expire.Year[0]&0x0f)*7+
				(ret.Expire.Year[1]&0x0f)*3+
				(ret.Expire.Month[0]&0x0f)+
				(ret.Expire.Month[1]&0x0f)*7+
				(ret.Expire.Day[0]&0x0f)*3+
				(ret.Expire.Day[1]&0x0f))%10 == ret.Expire.Checksum[0]&0x0f {
				ret.Expire.ChecksumOk = true
			}
		}

		// Overeni kontrolniho souctu pro cislo dokladu
		if len(ret.Number) == 9 && len(ret.NumberChecksum) == 1 {
			if ((ret.Number[0]&0x0f)*7+
				(ret.Number[1]&0x0f)*3+
				(ret.Number[2]&0x0f)+
				(ret.Number[3]&0x0f)*7+
				(ret.Number[4]&0x0f)*3+
				(ret.Number[5]&0x0f)+
				(ret.Number[6]&0x0f)*7+
				(ret.Number[7]&0x0f)*3+
				(ret.Number[8]&0x0f))%10 == ret.NumberChecksum[0]&0x0f {
				ret.NumberChecksumOk = true
			}
		}

		ret.IcaoType = strings.Trim(ret.IcaoType, "<")
		ret.IcaoSubtype = strings.Trim(ret.IcaoSubtype, "<")
		ret.Country = strings.Trim(ret.Country, "<")
		ret.Number = strings.Trim(ret.Number, "<")
		rawName = strings.Trim(rawName, "<")
		if len(rawName) > 0 {
			name := strings.Split(rawName, "<<")
			if len(name) == 2 {
				ret.Surname = strings.ReplaceAll(name[0], "<", " ")
				ret.Name = strings.ReplaceAll(name[1], "<", " ")
			}
		} else {
			ret.Surname = strings.ReplaceAll(ret.Surname, "<", " ")
			ret.Name = strings.ReplaceAll(ret.Name, "<", " ")
		}
		ret.Pin = strings.Trim(ret.Pin, "<")
		if ret.Sex != "M" && ret.Sex != "F" {
			ret.Sex = ""
		}
		ret.Nationality = strings.Trim(ret.Nationality, "<")
		return // OK
	}
	return // Error
}

/*
func handleFunc(data IcaoData) {
	//fmt.Println(data)
	ParseICAO(data)
}

func main() {
	deviceInfo, err := getDeviceInfo()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	deskoReader, err := deviceInfo.Open()
	if err != nil {
		log.Panicln(err)
	}
	defer deskoReader.Close()
	startReading(deskoReader, handleFunc)

	// Wait for SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
*/
