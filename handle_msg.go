package main

import (
	"Android-QQ-Pic-Cleaner/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const PIC_PATH = "/sdcard/Android/data/com.tencent.mobileqq/Tencent/MobileQQ/chatpic/chatimg/"
const KC_PATH = "/data/data/com.tencent.mobileqq/files/kc"

var key string
var pics []string
var crc64Table = make([]int, 256)

func init() {
	initKey()
	initCrc64Table()
}

func decryptBytes(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = b ^ key[i%len(key)]
	}
	return res
}

func decryptString(data string) string {
	dataRunes := []rune(data)
	res := ""
	for i, b := range dataRunes {
		res += string(b ^ rune(key[i%len(key)]))
	}

	return res
}

func decodePic(data []byte) {
	pic := &pb.PicRec{}
	err := proto.Unmarshal(data, pic)
	if err != nil {
		return
	}
	url := "chatimg:" + pic.Md5
	filename := "Cache_" + strconv.FormatInt(int64(crc64(url)), 16)
	path := filepath.Join(PIC_PATH, filename[len(filename)-3:], filename)
	pics = append(pics, path)
}

func decodeMix(data []byte) {
	mix := &pb.Msg{}
	err := proto.Unmarshal(data, mix)
	if err != nil {
		return
	}
	for _, elem := range mix.Elems {
		if elem.PicMsg != nil {
			decodePic(elem.PicMsg)
		}
	}
}

func initKey() {
	keyBytes, err := os.ReadFile(KC_PATH)
	if err != nil {
		log.Println("Failed to open key file: ", err)
	} else {
		key = string(keyBytes)
	}
}

func initCrc64Table() {
	var bf int
	for i := 0; i < 256; i++ {
		bf = i
		for j := 0; j < 8; j++ {
			if bf&1 != 0 {
				bf = bf>>1 ^ -7661587058870466123
			} else {
				bf >>= 1
			}
		}
		crc64Table[i] = bf
	}
}

func crc64(data string) int {
	v := -1
	for _, b := range data {
		v = crc64Table[(int(b)^v)&255] ^ v>>8
	}
	return v
}
