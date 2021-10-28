/**
 * @Author: Ramoncjs
 * @Date: 2021/10/28
 **/

package main

import (
	"github.com/0xrawsec/golang-evtx/evtx"
	"github.com/tealeg/xlsx"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mm []AllMap
	ch = make(chan AllMap)

	timeFormat           = "2006-01-02T15:04:05Z07:00"
	ipAddressPath        = evtx.Path("/Event/EventData/IpAddress")
	usernamePath         = evtx.Path("/Event/EventData/TargetUserName")
	logonTypePath        = evtx.Path("/Event/EventData/LogonType")
	targetDomainNamePath = evtx.Path("/Event/EventData/TargetDomainName")
	targetUserSidPath    = evtx.Path("/Event/EventData/TargetUserSid")
)

type AllMap struct {
	timeCreated      time.Time
	targetDomainName string
	username         string
	ipAddress        string
	logonType        int64
	targetUserSid    string
}

func main() {
	log.Println("[+] Already start...")

	//file, err := evtx.OpenDirty("log.evtx")
	file, err := evtx.OpenDirty("C:\\Windows\\System32\\winevt\\Logs\\Security.evtx")
	if err != nil {
		log.Println("[-]", err)
		log.Println("[-] Please run as administrator...")
		return

	}
	//file, _ := evtx.Open("log.evtx")
	defer file.Close()

	wg := sync.WaitGroup{}
	for i := range file.FastEvents() {
		i := i
		wg.Add(1)
		go func() {
			m := AllMap{}
			if i.EventID() == 4624 {
				timeCreated := i.TimeCreated()
				logonType, _ := i.GetInt(&logonTypePath)
				username, _ := i.GetString(&usernamePath)
				ipAddress, _ := i.GetString(&ipAddressPath)
				targetUserSid, _ := i.GetString(&targetUserSidPath)
				targetDomainName, _ := i.GetString(&targetDomainNamePath)

				if logonType != 0 && logonType != 5 && !strings.Contains(username, "$") {
					//fmt.Println(i.TimeCreated().Format(timeFormat), "\n", targetUserSid, targetDomainName+"/"+username, ipAddress, logonType)
					m.username = username
					m.ipAddress = ipAddress
					m.logonType = logonType
					m.timeCreated = timeCreated
					m.targetUserSid = targetUserSid
					m.targetDomainName = targetDomainName
					ch <- m
				}
			}
			wg.Done()
		}()
	}

	//close_ch
	go func() {
		wg.Wait()
		close(ch)
	}()

	//xlsFile
	xlsFile := xlsx.NewFile()
	defer xlsFile.Save("result.xlsx")

	// SetColWidth
	sheet, _ := xlsFile.AddSheet("Sheet1")
	sheet.SetColWidth(0, 0, 28)
	sheet.SetColWidth(1, 1, 18)
	sheet.SetColWidth(2, 2, 15)
	sheet.SetColWidth(3, 3, 18)
	sheet.SetColWidth(4, 4, 10)
	sheet.SetColWidth(5, 5, 50)

	s := sheet.AddRow()

	for i := 0; i < reflect.TypeOf(mm).Elem().NumField(); i++ {
		c := s.AddCell()
		c.GetStyle().Alignment.Horizontal = "center"
		c.GetStyle().Alignment.Vertical = "center"
		c.Value = reflect.TypeOf(mm).Elem().Field(i).Name
	}

	for v := range ch {
		s := sheet.AddRow()

		s.AddCell().Value = v.timeCreated.Format(timeFormat)
		s.AddCell().Value = v.targetDomainName
		s.AddCell().Value = v.username
		s.AddCell().Value = v.ipAddress
		s.AddCell().Value = strconv.FormatInt(v.logonType, 10)
		s.AddCell().Value = v.targetUserSid
	}

	log.Println("[+] The result has been saved in the current folder...")

}
