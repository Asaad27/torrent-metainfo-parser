package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jackpal/bencode-go"
	"log"
	"os"
	"strconv"
)

type MetaInfo struct {
	announce    string `bencode:"announce"`
	createdBy   string `bencode:"created by"`
	comment     string
	torrentInfo Info `bencode:"info"`
}

type Info struct {
	length      int64  `bencode:"length"`
	name        string `bencode:"name"`
	pieceLength int64  `bencode:"piece length"`
	pieces      string `bencode:"pieces"`
}

func (info *Info) infoHash() {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	errors(err)
	ans := sha1.Sum(buf.Bytes())
	stringHash := hex.EncodeToString(ans[:])
	fmt.Print("sha1 hash info : " + stringHash)
}

func parseToTorrent(data MetaInfo) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, data)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}

func readFile(filePath string, sizeOfPiece int64) MetaInfo {
	file, err := os.Open(filePath)
	errors(err)
	fileStats, err := file.Stat()
	errors(err)
	sizeOfFile := fileStats.Size()
	fileName := fileStats.Name()
	numberOfPieces := (sizeOfFile + sizeOfPiece - 1) / sizeOfPiece
	info := Info{
		length:      sizeOfFile,
		name:        fileName,
		pieceLength: sizeOfPiece,
	}
	n := int(numberOfPieces)
	var pieces = ""
	for i := 0; i < n; i++ {
		length := sizeOfPiece
		if i == n-1 {
			if sizeOfFile%sizeOfPiece != 0 {
				length = sizeOfFile % sizeOfPiece
			}
		}
		b := make([]byte, length)
		at, err := file.ReadAt(b, int64(i)*sizeOfPiece)
		if at != int(length) {
			fmt.Print("error bytes read less than it should be")
		}
		errors(err)
		bytesHash := sha1.Sum(b)
		hash := bytesHash[:]
		stringHash := hex.EncodeToString(hash)
		pieces += stringHash
		fmt.Println("piece hash : ", stringHash)
	}
	info.pieces = pieces
	metaInfo := MetaInfo{
		announce:    "https:127.0.0.1:6969/announce",
		createdBy:   "Asaad27",
		comment:     "hello peers",
		torrentInfo: info,
	}

	return metaInfo
}

func writeToFile(torrentPath string, buffer bytes.Buffer) {
	f, err := os.Create(torrentPath)
	errors(err)

	_, err = f.Write(buffer.Bytes())
	errors(err)
}

func errors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fileName := os.Args[1]
	sizeOfPiece, _ := strconv.Atoi(os.Args[2]) //in KB
	sizeOfPiece *= 1024

	metaInfo := readFile(fileName, int64(sizeOfPiece))
	buff, err := parseToTorrent(metaInfo)
	errors(err)
	writeToFile(fileName+".torrent", buff)
	metaInfo.torrentInfo.infoHash()
}
