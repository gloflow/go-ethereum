// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package gf_events

import (
	"fmt"
	"os"
	"time"
	"github.com/gocarina/gocsv"
)

//-------------------------------------------------------------------------------
type GFeventsCSVinfo struct {
	eventFull string
	filePath  string
	file      *os.File
	lineIndex uint
}

//-------------------------------------------------------------------------------
func persistCSVinit(pEvents []string) (map[string]*GFeventsCSVinfo, error) {

	CSVinfos := map[string]*GFeventsCSVinfo{}
	for _, eventFull := range pEvents {
		
		CSVfile, eventTypeFilePath, err := persistCSVnewFile(eventFull)
		if err != nil {
			return nil, err
		}
		CSVinfo := &GFeventsCSVinfo{
			eventFull: eventFull,
			filePath:  eventTypeFilePath,
			file:      CSVfile,
			lineIndex: 0,
		}
		CSVinfos[eventFull] = CSVinfo
	}
	return CSVinfos, nil
}

//-------------------------------------------------------------------------------
func persistCSVreinitFile(pCSVinfo *GFeventsCSVinfo, pEventFull string) (*GFeventsCSVinfo, error) {

	pCSVinfo.file.Close()
	newCSVfile, eventTypeFilePath, err := persistCSVnewFile(pCSVinfo.eventFull)
	if err != nil {
		return nil, err
	}
	newCSVinfo := &GFeventsCSVinfo{
		eventFull: pCSVinfo.eventFull,
		filePath:  eventTypeFilePath,
		file:      newCSVfile,
		lineIndex: 0,
	}
	return newCSVinfo, nil
}

//-------------------------------------------------------------------------------
func persistCSVnewFile(pEventFull string) (*os.File, string, error) {
	fileCreationTimeSec := float64(time.Now().UnixNano())/1000000000.0
	eventTypeFilePath := fmt.Sprintf("./data_tmp/%s~%f.csv", pEventFull, fileCreationTimeSec)

	CSVfile, err := os.OpenFile(eventTypeFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, "", err
	}
	return CSVfile, eventTypeFilePath, nil
}

//-------------------------------------------------------------------------------
func persistCSVwrite(pData interface{}, pCSVfile *os.File) error {

	csvContent, err := gocsv.MarshalString([]interface{}{&pData,}) // Get all clients as CSV string
	// err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
	fmt.Println(csvContent)



	_, err = pCSVfile.WriteString(csvContent)
	if err != nil {
		return err
	}
	return nil
}