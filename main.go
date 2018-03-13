package main

import(
	"os"
	"bufio"
	"io"
	"encoding/base64"
	"fmt"
	"flag"
)

func transcode(r io.Reader, w io.Writer, encoding bool) error{
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)

	enc := base64.StdEncoding
	// read line by line
	for scanner.Scan() {
		var dst, src []byte
		// transcoding
		src = scanner.Bytes()
		if encoding {
			dst = make([]byte, enc.EncodedLen(len(src)))
			base64.StdEncoding.Encode(dst, src)
		}else{
			dst = make([]byte, enc.DecodedLen(len(src)))
			_, err := enc.Decode(dst, src)
			if err != nil{
				return err
			}
		}
		
		// write
		if _,err := fmt.Fprintln(writer, string(dst)); err !=nil{
			return err
		}
	}

	if err := scanner.Err(); err!=nil{
		return err
	}

	err :=writer.Flush()
	return err	
}

func main(){
	srcPath := flag.String("src", "g-cn.txt", "source file path")
	dstPath := flag.String("dst", "glist.txt", "destination file path")
	encoding := flag.Bool("encode", true, "specify whether it's encoding, set false for decoding")

	flag.Parse()

	srcFile, err:= os.Open(*srcPath)
	if err!=nil{
		panic(err)
	}
	defer srcFile.Close()

	dstFile, err:= os.Create(*dstPath)
	if err!=nil{
		panic(err)
	}
	defer srcFile.Close()

	err = transcode(srcFile, dstFile, *encoding)
	if err!=nil{
		panic(err)
	}
}