package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var ApplicationDescription string = "Binary File Splitting"
var BuildVersion string = "dev"

func main() {

	inputFilePtr := flag.String("input_file", "test/dump.bin", "Path to input file")
	boundarySequencePtr := flag.String("hex", "21097019", "Bounary sequence in hexidecimal")
	chunkFilenamePrefixPtr := flag.String("prefix", "dump_chunk_", "The first part of the result filename")
	chunkFilenameSuffixPtr := flag.String("suffix", ".bin", "The last part of the result filename")
	showVersionPtr := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersionPtr {
		fmt.Printf("%s\n", ApplicationDescription)
		fmt.Printf("Version: %s\n", BuildVersion)
		os.Exit(0)
	}

	log.Print("Starting..")
	fmt.Println("Starting..")

	var boundrySequence []byte
	boundrySequence, err := hex.DecodeString(*boundarySequencePtr)
	if err != nil {
		log.Fatal("Error encoding hex:", err)
	}

	log.Printf("%#x", boundrySequence)

	inputFileHandler, err := os.Open(*inputFilePtr)
	if err != nil {
		log.Fatal("Error while opeinig input file:", err)
	}
	defer inputFileHandler.Close()

	inputFileInfo, _ := inputFileHandler.Stat()

	var inputFileSize int64 = inputFileInfo.Size()
	var counter int64 = 0
	const fileChunkSize = 1

	log.Print("Input file size: ", inputFileSize)

	for {

		if counter >= inputFileSize {
			break
		}

		partBuffer := make([]byte, fileChunkSize)
		inputFileHandler.Read(partBuffer)

		firstByte, err := hex.DecodeString("21")
		if err != nil {
			log.Fatal("Error encoding hex:", err)
		}

		secondByte, err := hex.DecodeString("09")
		if err != nil {
			log.Fatal("Error encoding hex:", err)
		}

		thirdByte, err := hex.DecodeString("70")
		if err != nil {
			log.Fatal("Error encoding hex:", err)
		}

		fourthByte, err := hex.DecodeString("19")
		if err != nil {
			log.Fatal("Error encoding hex:", err)
		}

		if bytes.Equal(partBuffer, firstByte) {
			fmt.Printf("\n\n")
			log.Printf(" %d Found 0x21", counter)
			partBuffer := make([]byte, fileChunkSize)
			inputFileHandler.Read(partBuffer)

			if bytes.Equal(partBuffer, secondByte) {
				log.Printf(" %d Found 0x21 0x09", counter)

				partBuffer := make([]byte, fileChunkSize)
				inputFileHandler.Read(partBuffer)

				if bytes.Equal(partBuffer, thirdByte) {
					log.Printf(" %d Found 0x21 0x09 0x70", counter)
					partBuffer := make([]byte, fileChunkSize)
					inputFileHandler.Read(partBuffer)

					if bytes.Equal(partBuffer, fourthByte) {
						log.Printf(" %d Found 0x21 0x09 0x70 0x19", counter)
						log.Print("Boundary: ", counter)

					}

				}

			}
		}

		var outputFileName = *chunkFilenamePrefixPtr + strconv.FormatInt(counter, 10) + *chunkFilenameSuffixPtr
		fmt.Printf("\rOutfile: %s \r", outputFileName)

		counter += 1
	}

}
