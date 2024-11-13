package fiori

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/pierrec/lz4"
)

type Fiori[T any] struct {
	unsafeEncode      bool
	itemsPerBlock     int
	errorRate         float64
	persistanceFile   string
	persistanceHandle *os.File
	blocks            []string
	filters           []byte
	currentItems      []T
}

const (
	defaultItemsPerBlock = 100
	defaultErrorRate     = 0.001
	defaultUnsafeEncode  = false
)

func New[T any]() Fiori[T] {
	return Fiori[T]{
		unsafeEncode:  defaultUnsafeEncode,
		itemsPerBlock: defaultItemsPerBlock,
		errorRate:     defaultErrorRate,
	}
}

func (f *Fiori[T]) SetUnsafeEncode(u bool) {
	f.unsafeEncode = u
}

func (f *Fiori[T]) SetErrorRate(r float64) {
	f.errorRate = r
}

func (f *Fiori[T]) SetItemsPerBlock(i int) {
	f.itemsPerBlock = i
}

func (f *Fiori[T]) SetPersistanceFile(p string) {
	f.persistanceFile = p
	// Open file handle
}

func (f *Fiori[T]) Add(item T) error {

	if len(f.currentItems) == f.itemsPerBlock {
		if f.unsafeEncode {
			// Unsafe encoding

			// compress and move to []blocks
			buffer := toBytesUnsafe(f.currentItems)

			compressed, err := compress(buffer)
			if err != nil {
				return err
			}

			/*
				back, err := fromBytesUnsafe[[]T](buffer)
				if err != nil {
					fmt.Println("Unsafe Error", err)
					return nil
				}
			*/
			f.blocks = append(f.blocks, string(compressed))
		} else {
			// Safe encoding

			buffer, err := toBytes(f.currentItems)
			if err != nil {
				fmt.Println("Safe Error", err)
				return nil
			}

			compressed, err := compress(buffer)
			if err != nil {
				return err
			}

			f.blocks = append(f.blocks, string(compressed))

			/*
				back, err := fromBytes[[]T](buffer)
				if err != nil {
					fmt.Println("Safe Error", err)
					return nil
				}
			*/
		}

		if f.persistanceHandle != nil {
			// persist to file
		}
	}

	f.currentItems = append(f.currentItems, item)

	return nil
}

func compress(input []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer := lz4.NewWriter(&compressed)

	_, err := writer.Write(input)
	if err != nil {
		return nil, fmt.Errorf("compression write error: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("compression close error: %v", err)
	}

	return compressed.Bytes(), nil
}

// toBytes converts a struct to bytes using gob
func toBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("gob encoding failed: %w", err)
	}
	return buf.Bytes(), nil
}

// fromBytes converts bytes back to struct using gob
func fromBytes[T any](data []byte) (T, error) {
	var result T
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return result, fmt.Errorf("gob decoding failed: %w", err)
	}
	return result, nil
}

// ToBytes converts a struct to a byte slice using unsafe
func toBytesUnsafe[T any](data T) []byte {
	size := unsafe.Sizeof(data)
	return unsafe.Slice((*byte)(unsafe.Pointer(&data)), size)
}

// FromBytes converts a byte slice back to a struct using unsafe
func fromBytesUnsafe[T any](bytes []byte) (T, error) {
	var result T

	// Get the size of the target struct
	size := unsafe.Sizeof(result)

	// Validate input length
	if len(bytes) != int(size) {
		return result, errors.New("byte slice length doesn't match struct size")
	}

	// Create a byte slice that refers to the struct's memory
	structBytes := unsafe.Slice((*byte)(unsafe.Pointer(&result)), size)

	// Copy the input bytes into the struct's memory
	copy(structBytes, bytes)

	return result, nil
}
