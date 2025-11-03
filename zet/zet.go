package zet

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type FileHashTable map[string]string

func (fht *FileHashTable) UpdateFile(fp string) error {
	//todo
	return nil
}

func GetFNVHash(data []byte) ([]byte, error) {
	h := fnv.New128()
	_, err := h.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error hashing file: %w", err)
	}

	return h.Sum(nil), nil
}

func ConvertHashToString(data []byte) string {
	s := hex.EncodeToString(data)
	return s
}

func BuildIndex() (FileHashTable, error) {
	fht := FileHashTable{}

	var walkFunc fs.WalkDirFunc = func(p string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			return filepath.SkipDir
		}
		f, err := os.Open(p)
		if err != nil {
			log.Printf("error opening file %s: %v\n", p, err)
			return nil
		}
		defer f.Close()

		data := []byte{}
		_, err = io.ReadFull(f, data)
		if err != nil {
			log.Printf("error reading from file %s: %v\n", p, err)
			return nil
		}

		h, err := GetFNVHash(data)
		if err != nil {
			log.Printf("error generating hash from file %s data: %v\n", p, err)
			return nil
		}
		stringedFileHash := ConvertHashToString(h)

		fht[p] = stringedFileHash

		return nil
	}

	err := filepath.WalkDir(".", walkFunc)
	if err != nil {
		return nil, err
	}
	return fht, nil
}
