package fiori

import (
	"bytes"
	"encoding/binary"
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

const defaultItemsPerBlock = 100
const defaultErrorRate = 0.001

func New[T any]() Fiori[T] {
	return Fiori[T]{
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
		// compress and move to []blocks
		buffer := toBytesUnsafe(f.currentItems)

		fmt.Println("Buffer:", buffer)

		back, err := fromBytesUnsafe[[]T](buffer)
		if err != nil {
			fmt.Println("Error", err)
			return nil
		}
		fmt.Println("back", back)

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

func toBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes(), err
}

func fromBytes[T any](data []byte) (T, error) {
	var result T
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &result)
	return result, err
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
