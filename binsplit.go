package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	// "strconv"
)

var Debug bool = false

const applicationDescription = "Binary File Splitting"
const buildVersion = "dev-2"

func lookupSequence(buffer []byte, sequence []byte) (found bool, sequencePositions []int, err error) {

	var bufferPosition int = 0
	var sequencePosition int = 0
	// var sequencePositions []int
	found = false

	for bufferPosition < len(buffer) {

		// if Debug {
		// 	log.Printf("lookupSequence Buffer length: %d ", len(buffer))
		// 	log.Printf("lookupSequence Current buffer position: %d ", bufferPosition)
		// 	log.Printf("lookupSequence Current sequence position: %d ", sequencePosition)
		// }

		// if Debug {
		// 	log.Printf("lookupSequence %#x loop at: %d", sequence, bufferPosition)
		// }

		if buffer[bufferPosition] == sequence[sequencePosition] {

			if Debug {
				log.Printf("lookupSequence Found sequence element %#x (len: %d) at: %d", sequence[sequencePosition], len(sequence), bufferPosition)
			}

			// if len(sequence) == 1 && len(buffer) == 1 {
			// 	log.Printf("This is the last byte %#x in given buffer.", sequence)
			// 	found = true
			// }

			if len(sequence) > 1 {
				if Debug {
					log.Printf("lookupSequence calling recursive for sequence %#x in buffer len %d at: %d", sequence, len(buffer), bufferPosition)
				}
				found, _, err = lookupSequence(buffer[bufferPosition:bufferPosition+len(sequence)], sequence[1:])
				if Debug {
					log.Printf("lookupSequence recursive result %t sequence %#x in buffer len %d at: %d", found, sequence, len(buffer), bufferPosition)
				}
				// return found, position, err

			} else {
				found = true
				// sequencePositions = bufferPosition
				if Debug {
					log.Printf("Returning success for sequence %#x (len: %d)", sequence, len(sequence))
				}
				return found, sequencePositions, nil

			}

			if found {
				if Debug {
					log.Printf("Found full sequence %#x (len: %d) in given buffer (len: %d) appending %d", sequence, len(sequence), len(buffer), bufferPosition)
				}
				sequencePositions = append(sequencePositions, bufferPosition)
			}

		} else {
			// log.Printf("lookupSequence Broken sequence in buffer at: %d (len: %d)", bufferPosition, len(buffer))
		}

		bufferPosition++
		sequencePosition = 0
	}
	if Debug {
		log.Printf("lookupSequence Result for sequence %#x in buffer len %d: %v", sequence, len(buffer), sequencePositions)
	}
	return found, sequencePositions, nil

}

func getCurrentOffset(inputFile *os.File) (int64, error) {
	currentOffset, err := inputFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Fatal("Error looking up current offset:", err)
		return -1, err
	}

	return currentOffset, nil
}

func main() {

	inputFilePtr := flag.String("i", "test/dump.bin", "Path to input file")
	boundarySequencePtr := flag.String("hex", "21097019", "Boundary sequence in hexadecimal")
	// chunkFilenamePrefixPtr := flag.String("prefix", "dump_chunk_", "The first part of the result filename")
	// chunkFilenameSuffixPtr := flag.String("suffix", ".bin", "The last part of the result filename")
	setDebugPtr := flag.Bool("d", false, "Enable debug")
	showVersionPtr := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showVersionPtr {
		fmt.Printf("%s\n", applicationDescription)
		fmt.Printf("Version: %s\n", buildVersion)
		os.Exit(0)
	}

	if *setDebugPtr {
		Debug = true
	}

	log.Print("Starting..")

	var boundrySequence []byte
	boundrySequence, err := hex.DecodeString(*boundarySequencePtr)
	if err != nil {
		log.Fatal("Error encoding hex:", err)
	}

	log.Printf("Boundary sequence: %#x", boundrySequence)

	inputFileHandler, err := os.Open(*inputFilePtr)
	if err != nil {
		log.Fatal("Error while opeinig input file:", err)
	}

	defer inputFileHandler.Close()

	inputFileInfo, _ := inputFileHandler.Stat()

	var inputFileSize int64 = inputFileInfo.Size()
	var counter int64 = 0
	var positions []int64
	const fileChunkSize = 204800
	var partBuffer []byte = make([]byte, fileChunkSize)

	log.Print("Input file size: ", inputFileSize)

	for counter <= inputFileSize {

		if Debug {
			currentOffset, err := getCurrentOffset(inputFileHandler)
			if err != nil {
				log.Fatal("Error while looking up current position: ", err)
			}

			log.Printf("Current offset before read: %d", currentOffset)
		}

		inputFileHandler.Read(partBuffer)

		currentOffset, err := getCurrentOffset(inputFileHandler)
		if err != nil {
			log.Fatal("Error while looking up current position: ", err)
		}

		if Debug {
			log.Printf("Current offset after read: %d", currentOffset)
		}

		if bytes.Contains(partBuffer, boundrySequence) {

			if Debug {
				currentOffset, err := getCurrentOffset(inputFileHandler)
				if err != nil {
					log.Fatal("Error while looking up current position: ", err)
				}
				log.Printf("Buffer contains boundary. Current file offset: %d", currentOffset)
			}

			_, positionsFound, err := lookupSequence(partBuffer, boundrySequence)
			if err != nil {
				log.Fatal("Error while looking up sequence: ", err)
			}

			log.Printf("Found for offset: %d positions: %v", currentOffset, positionsFound)

			for _, item := range positionsFound {
				positions = append(positions, int64(item)+currentOffset)

			}
		}

		counter += fileChunkSize - 3
		newPosition, err := inputFileHandler.Seek(-3, 1)
		if err != nil {
			log.Fatal(err)
		}

		if Debug {
			log.Printf("New file offset: %d", newPosition)
		}
	}

	log.Printf("Found sequence %#x at positions: %v", boundrySequence, positions)

}
