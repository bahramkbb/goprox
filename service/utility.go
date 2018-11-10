package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

// Compress whatever class object
func compressObject(obj interface{}) ([]byte, error) {
	// Converting the object into json
	objJson, _ := json.Marshal(obj)

	// Gzipping data
	var buf bytes.Buffer

	err := gzipWrite(&buf, objJson)
	if err != nil {
		log.Printf("error compressing object : %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decompress whatever class object into the targetType interface
func decompressObject(value []byte, targetType interface{}) (interface{}, error) {
	// Uncompressing Gzipped data
	var buf bytes.Buffer

	err := gunzipWrite(&buf, value)
	if err != nil {
		return nil, fmt.Errorf("error decompressing object: %v", err)
	}

	// Converting bytes back to object
	err = json.Unmarshal(buf.Bytes(), targetType)

	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("error decompressing object: %v", err)
	}

	return targetType, nil
}

// Write gzipped data to a Writer
func gzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	defer gw.Close()
	gw.Write(data)

	if err != nil {
		log.Printf("error gziping : %v", err)
		return err
	}

	return err
}

// Write gunzipped data to a Writer
func gunzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		log.Printf("error gunziping : %v", err)
		return err
	}
	w.Write(data)
	return nil
}

