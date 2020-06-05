package gf_events

import (
	"fmt"
	"os"
	"time"
	"github.com/gocarina/gocsv"
)

//-------------------------------------------------------------------------------
type GFeventCSVinfo struct {
	eventFull string
	filePath  string
	file      *os.File
	lineIndex uint
}

//-------------------------------------------------------------------------------
func persistCSVinit(pEvents []string) (map[string]*GFeventCSVinfo, error) {

	CSVinfos := map[string]*GFeventCSVinfo{}
	for _, eventFull := range pEvents {
		
		CSVfile, eventTypeFilePath, err := persistCSVnewFile(eventFull)
		if err != nil {
			return nil, err
		}
		CSVinfo := &GFeventCSVinfo{
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
func persistCSVreinitFile(pCSVinfo *GFeventCSVinfo, pEventFull string) (*GFeventCSVinfo, error) {

	pCSVinfo.file.Close()
	newCSVfile, eventTypeFilePath, err := persistCSVnewFile(pCSVinfo.eventFull)
	if err != nil {
		return nil, err
	}
	newCSVinfo := &GFeventCSVinfo{
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
func persistCSVwrite(pData interface{}, pCSVfile *os.File) {

	csvContent, err := gocsv.MarshalString([]interface{}{&pData,}) // Get all clients as CSV string
	// err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
	fmt.Println(csvContent)



	_, err = pCSVfile.WriteString(csvContent)
	if err != nil {

	}
}