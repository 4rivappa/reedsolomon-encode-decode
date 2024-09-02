package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/klauspost/reedsolomon"
)

const dataShards = 7
const parityShards = 3

func main() {
	encodeDecodeString("Hello world!")

	encodeDecodeFile("input.txt", "output.txt")
}

func encodeDecodeString(input string) {
	fmt.Println("Encoding and decoding a string:")

	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		log.Fatal(err)
	}

	shardSize := (len(input) + dataShards - 1) / dataShards
	if shardSize == 0 {
		shardSize = 1
	}

	data := make([][]byte, dataShards+parityShards)
	for i := 0; i < dataShards; i++ {
		shard := make([]byte, shardSize)
		start := i * shardSize
		end := start + shardSize
		if end > len(input) {
			end = len(input)
		}
		copy(shard, input[start:end])
		data[i] = shard
	}

	// fmt.Printf("%v\n", data)

	for i := dataShards; i < dataShards+parityShards; i++ {
		data[i] = make([]byte, shardSize)
	}

	fmt.Printf("%v\n", data)

	err = enc.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(" ============= ")
	fmt.Printf("%v\n", data)

	data[1] = nil
	data[4] = nil
	data[3] = nil

	fmt.Println(" ============= ")
	fmt.Printf("%v\n", data)

	err = enc.Reconstruct(data)
	if err != nil {
		fmt.Println("error for reconstruct")
		log.Fatal(err)
	}

	ok, err := enc.Verify(data)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Verification failed")
	}

	var reconstructed strings.Builder
	for i := 0; i < dataShards; i++ {
		reconstructed.Write(data[i])
	}

	result := reconstructed.String()
	fmt.Println("Reconstructed data:", result)
}

func encodeDecodeFile(inputFile, outputFile string) {
	fmt.Println("Encoding and decoding a file:")

	fileData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	shardSize := (len(fileData) + dataShards - 1) / dataShards

	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]byte, dataShards+parityShards)
	for i := 0; i < dataShards; i++ {
		shard := make([]byte, shardSize)
		start := i * shardSize
		end := start + shardSize
		if end > len(fileData) {
			end = len(fileData)
		}
		copy(shard, fileData[start:end])
		data[i] = shard
	}

	for i := dataShards; i < dataShards+parityShards; i++ {
		data[i] = make([]byte, shardSize)
	}

	fmt.Println("===== before encode =====")
	fmt.Printf("%v\n", data)

	err = enc.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("===== after encode =====")
	fmt.Printf("%v\n", data)

	data[1] = nil
	data[3] = nil
	data[6] = nil

	fmt.Println("===== after data loss =====")
	fmt.Printf("%v\n", data)

	err = enc.Reconstruct(data)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := enc.Verify(data)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("Verification failed")
	}

	decoded := make([]byte, 0, len(fileData))
	for i := 0; i < dataShards; i++ {
		decoded = append(decoded, data[i]...)
	}
	decoded = decoded[:len(fileData)]

	err = os.WriteFile(outputFile, decoded, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File successfully encoded, decoded, and written to", outputFile)
}
