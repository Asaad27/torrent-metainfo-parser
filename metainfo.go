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
	"strings"
	"time"
)

type MetaInfo struct {
	announce     string `bencode:"announce"`
	createdBy    string `bencode:"created by"`
	creationDate int64  `bencode:"creation date"`
	comment      string
	torrentInfo  Info `bencode:"info"`
}

func NewMetaInfo(announce string, createdBy string, comment string, torrentInfo Info) *MetaInfo {
	if len(os.Args) >= 4 {
		announce = os.Args[3]
	}
	if len(os.Args) >= 5 {
		createdBy = os.Args[4]
	}
	if len(os.Args) >= 6 {
		comment = os.Args[5]
	}
	now := time.Now()
	sec := now.Unix()

	return &MetaInfo{announce: announce, createdBy: createdBy, creationDate: sec, comment: comment, torrentInfo: torrentInfo}
}

type Info struct {
	length      int64  `bencode:"length"`
	name        string `bencode:"name"`
	pieceLength int64  `bencode:"piece length"`
	pieces      string `bencode:"pieces"`
}

func NewInfo(length int64, name string, pieceLength int64, pieces string) *Info {

	return &Info{length: length, name: name, pieceLength: pieceLength, pieces: pieces}
}

func (info *Info) infoHash() string {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *info)
	errors(err)
	ans := sha1.Sum(buf.Bytes())
	stringHash := hex.EncodeToString(ans[:])
	return stringHash
}

func BEncodeMetaInfo(data MetaInfo) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, data)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}

func pieceToSHA1(b []byte) string {
	bytesHash := sha1.Sum(b)
	ans := string(bytesHash[:])
	//return hex.EncodeToString(bytesHash[:])
	return ans
}

func parseFile(filePath string, sizeOfPiece int64) MetaInfo {
	file, err := os.Open(filePath)
	errors(err)
	fileStats, err := file.Stat()
	errors(err)
	sizeOfFile := fileStats.Size()
	fileName := fileStats.Name()
	numberOfPieces := (sizeOfFile + sizeOfPiece - 1) / sizeOfPiece

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

		pieceSHA1 := pieceToSHA1(b)
		pieces += pieceSHA1
		hexSHA1 := hex.EncodeToString([]byte(pieceSHA1))
		fmt.Printf("piece %d hash : "+strings.ToUpper(hexSHA1)+"\n", i+1)
	}

	info := NewInfo(sizeOfFile, fileName, sizeOfPiece, pieces)
	metaInfo := NewMetaInfo("https:127.0.0.1:6969/announce", "Asaad27", "hello peers", *info)

	return *metaInfo
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
	var sizeOfPiece = 16
	if len(os.Args) > 2 {
		sizeOfPiece, _ = strconv.Atoi(os.Args[2]) //in KB
	}
	sizeOfPiece *= 1024

	metaInfo := parseFile(fileName, int64(sizeOfPiece))
	buff, err := BEncodeMetaInfo(metaInfo)
	errors(err)
	writeToFile(fileName+".torrent", buff)
	fmt.Println("info hash " + strings.ToUpper(metaInfo.torrentInfo.infoHash()))
}
