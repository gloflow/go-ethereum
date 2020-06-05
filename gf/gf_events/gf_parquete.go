package gf_events

//-------------------------------------------------------------------------------
import (
	"fmt"
	"sync"
	"github.com/fatih/color"
	// "bytes"
	// "os"
	// "io"
	// "github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go-source/local"
	// "github.com/xitongsys/parquet-go-source/writerfile"
	"github.com/xitongsys/parquet-go/writer"
	// "github.com/xitongsys/parquet-go/reader"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------------------------------------
type GFeventParquetInfo struct {
	eventType     string
	filePath      string
	parquetFile   source.ParquetFile
	parquetWriter *writer.ParquetWriter
	lock          *sync.Mutex
}

//-------------------------------------------------------------------------------
func persistParquetInitEvents(pEventsTypes []string) (map[string]GFeventParquetInfo, error) {

	eventsParquetInfos := map[string]GFeventParquetInfo{}
	for _, eventType := range pEventsTypes {

		// schemaStruct := pSchemaStructs[i]
		eventTypeFilePath := fmt.Sprintf("./data/%s.parquet", eventType)



		parquetFile, err := local.NewLocalFileWriter(eventTypeFilePath) //"flat.parquet")
		if err != nil {
			fmt.Println("Can't create local file", err)
			return nil, err
		}

		fmt.Println("---")
		fmt.Println(eventType)
		// spew.Dump(schemaStruct)


		var parquetWriter *writer.ParquetWriter
		switch eventType {

		case "protocol_manager:handle_new_peer":
			parquetWriter, err = writer.NewParquetWriter(parquetFile, new(GFeventNewPeerLifecycle), 4)
		case "protocol_manager:dropping_unsynced_node_during_fast_sync":
			parquetWriter, err = writer.NewParquetWriter(parquetFile, new(GFeventDroppingUnsyncedNodeDuringFastSync), 4)
		case "downloader:register_peer":
			parquetWriter, err = writer.NewParquetWriter(parquetFile, new(GFeventNewPeerRegister), 4)
		case "downloader:new_header_from_peer":
			parquetWriter, err = writer.NewParquetWriter(parquetFile, new(GFeventNewHeaderFromPeer), 4)
		case "downloader:block_synchronise_with_peer":
			parquetWriter, err = writer.NewParquetWriter(parquetFile, new(GFeventBlockSynchroniseWithPeer), 4)
		}
		if err != nil {
			fmt.Println("Can't create parquet writer", err)
			return nil, err
		}
		/*parquetFile, parquetWriter, err := persistParquetInit(eventTypeFilePath)
		if err != nil {
			return nil, err
		}*/

		info := GFeventParquetInfo{
			eventType:     eventType,
			filePath:      eventTypeFilePath,
			parquetFile:   parquetFile,
			parquetWriter: parquetWriter,
			lock:          &sync.Mutex{},
		}
		eventsParquetInfos[eventType] = info
	}

	return eventsParquetInfos, nil
}

//-------------------------------------------------------------------------------
func persistParquetInit(pFilePath string) (source.ParquetFile, *writer.ParquetWriter, error) {


	/*f, err := os.Create(pFilePath)
	if err != nil {
		fmt.Printf("failed to greate Parquet file [%s]\n", pFilePath)
		return nil, err
	}
	defer f.Close()
	//fw := bufio.NewWriter(f)*/

	/*var err error
	buf := new(bytes.Buffer)
	parquetFile := writerfile.NewWriterFile(buf)*/

	var err error
	parquetFile, err := local.NewLocalFileWriter(pFilePath) //"flat.parquet")
	if err != nil {
		fmt.Println("Can't create local file", err)
		return nil, nil, err
	}
	

	parquetWriter, err := writer.NewParquetWriter(parquetFile, new(GFeventMsg), 4)
	if err != nil {
		fmt.Println("Can't create parquet writer", err)
		return nil, nil, err
	}

	return parquetFile, parquetWriter, nil
}

//-------------------------------------------------------------------------------
func persistParquetWrite(pData []interface{}, pParquetWriter *writer.ParquetWriter) error {
	
	for i:=0; i < len(pData); i++ {
		
		fmt.Println("---------PARQUET_WRITE")
		event := pData[i]
		if err := pParquetWriter.Write(event); err != nil {
			fmt.Println("Write error", err)
			return err
		}
		fmt.Println("done writing record...")
	}
	return nil
}

//-------------------------------------------------------------------------------
func persistParquetClose(pFilePath string, pParquetWriter *writer.ParquetWriter) error {

	cyan := color.New(color.BgCyan, color.FgBlack).SprintFunc()
	green := color.New(color.BgGreen, color.FgBlack).SprintFunc()
	fmt.Printf("persisting parquet file - %s", cyan(pFilePath))

	// Write the footer and stop writing
	if err := pParquetWriter.WriteStop(); err != nil {
		fmt.Printf("WriteStop error - %s\n", err)
		return err
	}
	
	fmt.Printf(" - %s\n", green("done"))

	// fmt.Printf("closed...\n")
	// pParquetFile.Close()
	return nil
}