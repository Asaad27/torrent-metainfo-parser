package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	errors2 "errors"
	"fmt"
	"github.com/jackpal/bencode-go"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type MetaInfo struct {
	Announce     string `bencode:"announce"`
	CreatedBy    string `bencode:"created by"`
	CreationDate int64  `bencode:"creation date"`
	Comment      string `bencode:"comment"`
	TorrentInfo  Info   `bencode:"info"`
}

func NewMetaInfo(announce string, createdBy string, comment string, torrentInfo Info) *MetaInfo {
	if args.Announce != -1 {
		announce = os.Args[args.Announce]
	}
	if args.CreatedBy != -1 {
		createdBy = os.Args[args.CreatedBy]
	}
	if args.Comment != -1 {
		comment = os.Args[args.Comment]
	}

	now := time.Now()
	sec := now.Unix()

	return &MetaInfo{Announce: announce, CreatedBy: createdBy, CreationDate: sec, Comment: comment, TorrentInfo: torrentInfo}
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

type Arguments struct {
	Comment     int
	Announce    int
	CreatedBy   int
	SizeOfPiece int
}

func NewArguments() *Arguments {
	return &Arguments{-1, -1, -1, -1}
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
	metaInfo := NewMetaInfo("https:127.0.0.1:6969/Announce", "Asaad27", "hello peers", *info)

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

func assertHasNext(i int, Args []string) error {

	if i+1 >= len(Args) {
		err1 := errors2.New("expected argument after " + Args[i])
		return err1
	}

	return nil
}

func ParseArgs(Args []string) (error, Arguments) {
	args := NewArguments()
	for i := 2; i < len(Args); {
		str := Args[i]
		if str[0] == '-' && len(str) == 2 {
			err := assertHasNext(i, Args)
			if err != nil {
				return err, *NewArguments()
			}
			switch str[1] {
			case 'p':
				args.SizeOfPiece = i + 1
				break
			case 'c':
				args.Comment = i + 1
				break
			case 'b':
				args.CreatedBy = i + 1
				break
			case 'a':
				args.Announce = i + 1
				break
			default:
				err1 := errors2.New("unknown command after - ")
				return err1, *NewArguments()
			}
			i += 2
		} else {
			err1 := errors2.New("unknown argument : " + Args[i])
			return err1, *NewArguments()
		}
	}

	return nil, *args
}

var args Arguments

func main() {
	fileName := os.Args[1]

	var err error
	err, args = ParseArgs(os.Args)
	if err != nil {
		errors(err)
	}

	var sizeOfPiece = 16 //default
	if args.SizeOfPiece != -1 {
		sizeOfPiece, _ = strconv.Atoi(os.Args[args.SizeOfPiece]) //in KB
	}

	sizeOfPiece *= 1024

	metaInfo := parseFile(fileName, int64(sizeOfPiece))
	buff, err := BEncodeMetaInfo(metaInfo)
	errors(err)
	writeToFile(fileName+".torrent", buff)
	fmt.Println("info hash " + strings.ToUpper(metaInfo.TorrentInfo.infoHash()))
}
