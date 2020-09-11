// Copyright 2014 The go-ethereum Authors
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
	"net"
	// "io"
	"bufio"
	"time"
	"bytes"
	"github.com/fatih/color"
	"github.com/gocarina/gocsv"
	"github.com/davecgh/go-spew/spew"
	
	// "github.com/xitongsys/parquet-go/source"
	// "github.com/xitongsys/parquet-go/writer"
	// "github.com/xitongsys/parquet-go-source/local"
	// "github.com/ethereum/go-ethereum/common"
)

//-------------------------------------------------------------------------------
type GFeventProcessor struct {
	EventCh                 chan GFeventMsg
	machineID               string
	eventsTypesParquetInfos map[string]GFeventParquetInfo
	// ParquetFile   source.ParquetFile
	// ParquetWriter *writer.ParquetWriter
}

type GFeventMsg struct {
	Id      string
	TimeSec float64
	Module  string
	Type    string
	Msg     string
	Data    interface{}
}

//-------------------------------------------------------------------------------
func EventProcessorCreate() (*GFeventProcessor, error) {

	fmt.Printf("NEW GF_EVENT_PROCESSOR ===========>>>>\n")

	// cyan := color.New(color.BgCyan, color.FgBlack).SprintFunc()

	eventCh := make(chan GFeventMsg, 100)
	machineID, err := getMACaddr()
	if err != nil {
		return nil, err
	}


	eventsToPersist := []string{
		"protocol_manager:handle_new_peer",
		"protocol_manager:dropping_unsynced_node_during_fast_sync",
		"downloader:register_peer",
		"downloader:new_header_from_peer",
		"downloader:block_synchronise_with_peer",
	}

	// PERSIST__PARQUET
	eventsTypesParquetInfos, err := persistParquetInitEvents(eventsToPersist)
	if err != nil {
		panic(err)
	}

	// PERSIST__CSV
	CSVinfos, err := persistCSVinit(eventsToPersist)
	if err != nil {
		panic(err)
	}

	// EVENT_QUEUE
	eventQueue, err := queueSQSinit()
	if err != nil {
		panic(err)
	}

	eventProcessor := &GFeventProcessor{
		EventCh:                 eventCh,
		machineID:               machineID,
		eventsTypesParquetInfos: eventsTypesParquetInfos,
		// ParquetFile:   parquetFile,
		// ParquetWriter: parquetWriter,
	}


	/*go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case _ = <-ticker.C:
				
				fmt.Printf("---------PARQUET - %s\n", cyan("writing footers"))
				for _, v := range eventsTypesParquetInfos {

					unflushedObjcNum := int64(len(v.parquetWriter.Objs))

					//v.lock.Lock()
					fmt.Printf("persisted parquet file - objs # [%d] - %s\n", unflushedObjcNum, cyan(v.filePath))
					// err := v.parquetWriter.Flush(true)
					// if err != nil {
					// 	fmt.Printf("ERROR!! - %s", fmt.Sprint(err))
					// }
					// fmt.Printf("%s\n", cyan("done"))
					//v.lock.Unlock()
					
					spew.Dump(v)
					
					persistParquetClose(v.filePath, v.parquetWriter)

					//v.lock.Unlock()
				}
			}
		}
	}()*/




	newPeersLifecyclesEvents := []GFeventNewPeerLifecycle{}
	newPeersRegistersEvents := []GFeventNewPeerRegister{}
	go func() {
		for {
			select {
			case eventMsg := <- eventProcessor.EventCh:
			
				EventView(&eventMsg)
				eventFull := fmt.Sprintf("%s:%s", eventMsg.Module, eventMsg.Type)

				//----------------------------------------
				// PROTOCOL_MANAGER : HANDLE_NEW_PEER
				if eventFull == "protocol_manager:handle_new_peer" {
					specificEvent := eventMsg.Data.(GFeventNewPeerLifecycle)
					specificEvent.Id = eventMsg.Id
					specificEvent.TimeSec = eventMsg.TimeSec
					specificEvent.Module = eventMsg.Module
					specificEvent.Type = eventMsg.Type
	
					newPeersLifecyclesEvents = append(newPeersLifecyclesEvents, specificEvent)

					fmt.Printf("++++++++++++++++++++++++-------------------")
					spew.Dump(specificEvent)

					
					CSVinfo := CSVinfos[eventFull]
					CSVfile := CSVinfo.file
					w := bufio.NewWriter(CSVfile)
					
					var errCSV error
					if CSVinfo.lineIndex == 0 {
						errCSV = gocsv.MarshalFile([]*GFeventNewPeerLifecycle{&specificEvent}, CSVfile)
					} else {
						errCSV = gocsv.MarshalWithoutHeaders([]*GFeventNewPeerLifecycle{&specificEvent}, w)
					}
					

					if errCSV != nil {
						panic(errCSV)
					}
					
					CSVinfo.lineIndex += 1

					if CSVinfo.lineIndex > 100 {
						newCSVinfo, err := persistCSVreinitFile(CSVinfo, eventFull)
						if err != nil {
							panic(err)
						}
						CSVinfos[eventFull] = newCSVinfo
					}
					
					

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)


					/*go func() {
						parquetInfo := eventProcessor.eventsTypesParquetInfos[eventFull]
						parquetInfo.lock.Lock()
						parquetWriter := parquetInfo.parquetWriter
						if err := parquetWriter.Write(specificEvent); err != nil {
							fmt.Println("Write error", err)
						}
						parquetInfo.lock.Unlock()
					}()*/
				}

				//----------------------------------------
				// PROTOCOL_MANAGER : DROPPING_UNSYNCED_NODE_DURING_FAST_SYNC
				if eventFull == "protocol_manager:dropping_unsynced_node_during_fast_sync" {

					specificEvent := eventMsg.Data.(GFeventDroppingUnsyncedNodeDuringFastSync)
					specificEvent.Id = eventMsg.Id
					specificEvent.TimeSec = eventMsg.TimeSec
					specificEvent.Module = eventMsg.Module
					specificEvent.Type = eventMsg.Type

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)

					// PERSIST
					go func() {
						parquetInfo := eventProcessor.eventsTypesParquetInfos[eventFull]
						parquetInfo.lock.Lock()
						parquetWriter := parquetInfo.parquetWriter
						if err := parquetWriter.Write(specificEvent); err != nil {
							fmt.Println("Write error", err)
						}
						parquetInfo.lock.Unlock()


						
					}()
				}

				//----------------------------------------
				// DOWNLOADER : REGISTER_PEER
				if eventFull == "downloader:register_peer" {

					specificEvent := eventMsg.Data.(GFeventNewPeerRegister)
					specificEvent.Id = eventMsg.Id
					specificEvent.TimeSec = eventMsg.TimeSec
					specificEvent.Module = eventMsg.Module
					specificEvent.Type = eventMsg.Type

					newPeersRegistersEvents = append(newPeersRegistersEvents, specificEvent)

					spew.Dump(specificEvent)

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)
				}

				//----------------------------------------
				// DOWNLOADER : DROPPING_PEER_SYNC_FAILED
				if eventFull == "downloader:dropping_peer_sync_failed" {

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)
				}

				//----------------------------------------
				// DOWNLOADER : BLOCK_SYNCHRONISE_WITH_PEER
				if eventFull == "downloader:block_synchronise_with_peer" {

					specificEvent := eventMsg.Data.(GFeventBlockSynchroniseWithPeer)
					specificEvent.Id = eventMsg.Id
					specificEvent.TimeSec = eventMsg.TimeSec
					specificEvent.Module = eventMsg.Module
					specificEvent.Type = eventMsg.Type

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)

					// PERSIST
					go func() {
						parquetInfo := eventProcessor.eventsTypesParquetInfos[eventFull]
						parquetInfo.lock.Lock()
						parquetWriter := parquetInfo.parquetWriter
						if err := parquetWriter.Write(specificEvent); err != nil {
							fmt.Println("Write error", err)
						}
						parquetInfo.lock.Unlock()
					}()
				}

				//----------------------------------------
				// DOWNLOADER : NEW_HEADER_FROM_PEER - event for a single new header received from a peer
				if eventFull == "downloader:new_header_from_peer" {

					specificEvent := eventMsg.Data.(GFeventNewHeaderFromPeer)
					specificEvent.Id = eventMsg.Id
					specificEvent.TimeSec = eventMsg.TimeSec
					specificEvent.Module = eventMsg.Module
					specificEvent.Type = eventMsg.Type

					// EVENT_QUEUE
					pushEvent(eventMsg, eventQueue)

					// PERSIST
					go func() {
						parquetInfo := eventProcessor.eventsTypesParquetInfos[eventFull]
						parquetInfo.lock.Lock()
						parquetWriter := parquetInfo.parquetWriter
						if err := parquetWriter.Write(specificEvent); err != nil {
							fmt.Println("Write error", err)
						}
						parquetInfo.lock.Unlock()
					}()
				}

				//----------------------------------------
			}
		}
	}()

	return eventProcessor, nil
}

//-------------------------------------------------------------------------------
func EventSend(pModule string,
	pType string,
	pMsg string,
	pData interface{},
	pEventProcessor *GFeventProcessor) {
	
	eventTimeSec := float64(time.Now().UnixNano())/1000000000.0
	id := fmt.Sprintf("%s%f", pEventProcessor.machineID, eventTimeSec)
	event := GFeventMsg{
		Id:      id,
		TimeSec: eventTimeSec,
		Module:  pModule,
		Type:    pType,
		Msg:     pMsg,
		Data:    pData,
	}
	pEventProcessor.EventCh <- event
}

//-------------------------------------------------------------------------------
func EventView(pEventMsg *GFeventMsg) {
	yellow := color.New(color.BgYellow, color.FgBlack).SprintFunc()
	green := color.New(color.BgGreen, color.FgBlack).SprintFunc()
	fmt.Printf("gf_event - %s:%s:%s\n", green(pEventMsg.Module), yellow(pEventMsg.Type), pEventMsg.Msg)
}

//-------------------------------------------------------------------------------
func getMACaddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			MACaddr := i.HardwareAddr.String()
			return MACaddr, nil
		}
	}
	return "", nil
}

//-------------------------------------------------------------------------------
/*type TestNested struct {
	B string `parquet:"name=a, type=UTF8"`
}

type TestStruct struct {
	A  string     `parquet:"name=a, type=UTF8"`
	Bb TestNested `parquet:"name=bb"`
}



parquetFile, err = local.NewLocalFileWriter("test.parquet") //"flat.parquet")
if err != nil {
	fmt.Println("Can't create local file", err)
	return nil, err
}


parquetWriter, err = writer.NewParquetWriter(parquetFile, new(TestStruct), 4)
if err != nil {
	fmt.Println("Can't create parquet writer", err)
	return nil, err
}

for i:=0;i<1000;i++ {
	err = parquetWriter.Write(TestStruct{A:"test_string", Bb: TestNested{B:"nested"}})
	if err != nil {
		panic(err)
	}
}

if err := parquetWriter.WriteStop(); err != nil {
	fmt.Printf("WriteStop error - %s\n", err)
	return nil, err
}

for i:=0;i<1000;i++ {
	err = parquetWriter.Write(TestStruct{A:"test_string"})
	if err != nil {
		panic(err)
	}
}

if err := parquetWriter.WriteStop(); err != nil {
	fmt.Printf("WriteStop error - %s\n", err)
	return nil, err
}

//parquetFile.Close()
fmt.Println("wrote test")
panic(1)*/

/*testEvents := []GFevent{}
for i:=0;i<1000;i++ {
	testEvents = append(testEvents, GFevent{
		Id:      fmt.Sprintf("test_id:%s", i),
		TimeSec: 5.443,
		Module:  "test_module",
		Type:    "test_type",
	})
}



err = persistParquetWrite(testEvents, eventProcessor.ParquetWriter)
if err != nil {
	panic(1)	
}


fmt.Println("closing")


err = parquetWriter.WriteStop()
if err != nil {
	fmt.Println("cant stop")
	fmt.Println(err)
}


//persistParquetClose(eventProcessor.ParquetWriter)
fmt.Println("done closing")


panic(1)*/