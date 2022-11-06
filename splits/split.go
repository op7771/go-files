package splits

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func splitFile(targetFile string, chunkSize uint64, splitFileName string, splitFileExtention string) error {
	src, err := os.Open(targetFile)
	if err != nil {
		return err
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)

	stat, err := src.Stat()
	if err != nil {
		return err
	}

	fileSize := stat.Size()

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(chunkSize)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(chunkSize), float64(fileSize-int64(i*chunkSize))))
		partBuffer := make([]byte, partSize)

		_, _ = src.Read(partBuffer)

		fileName := splitFileName + "_" + strconv.FormatUint(i, 10) + splitFileExtention
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		err = ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		if err != nil {
			return err
		}

		fmt.Println("Split to : ", fileName)
	}

	return nil
}

func SplitAndRemove(targetFile string, chunkSize uint64, splitFileName string, splitFileExtention string) error {
	err := splitFile(targetFile, chunkSize, splitFileName, splitFileExtention)
	if err != nil {
		return err
	}
	fmt.Printf("Deleting to %s.\n", targetFile)
	err = os.Remove(targetFile)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted %s.\n", targetFile)
	return nil
}
