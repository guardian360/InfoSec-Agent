import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4"
)

type SearchEngine struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Keyword string `json:"alias"`
}

func main() {
	// Path to the Firefox profile directory
	profilePath := "C:\\Users\\kdkoe\\AppData\\Roaming\\Mozilla\\Firefox\\Profiles\\sll129yt.default-release"

	// Path to the search.json.mozlz4 file
	filePath := profilePath + "/search.json.mozlz4"

	//Copy the database so we don't have problems with locked files
	tempSearch := filepath.Join(os.TempDir(), "tempSearch.json.mozlz4")

	copyError := CopyFile(filePath, tempSearch)
	if copyError != nil {
		fmt.Println("hi", copyError)
	}

	// Get file information
	fileInfo, err := os.Stat(tempSearch)
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return
	}

	file, err := os.Open(tempSearch)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer file.Close()

	byteSlice := make([]byte, 8)

	_, err = file.Read(byteSlice)
	if err != nil {
		fmt.Println("wtf", err)
	}

	uncompressSize := make([]byte, 4)

	// Seek to skip the first 8 bytes
	_, err = file.Seek(8, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
	}

	_, err = file.Read(uncompressSize)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	unCompressedSize := binary.LittleEndian.Uint32(uncompressSize)
	// Seek to skip the first 12 bytes
	_, err = file.Seek(12, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
	}

	fileSize := fileInfo.Size()
	byteSlice2 := make([]byte, fileSize-12)

	_, err = file.Read(byteSlice2)
	if err != nil {
		fmt.Println("error reading file", err)
	}

	data := make([]byte, unCompressedSize)
	lz4.UncompressBlock(byteSlice2, data)

	fmt.Println(string(data))
}