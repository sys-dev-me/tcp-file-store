package main

import "net"
import "fmt"
import "os"
import "strconv"
import "io"

const BUFFERSIZE =  2048

func fillString ( retunString string, toLength int ) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func main() {

	if len(os.Args) < 4 {
		fmt.Printf( "Usage: %s [SERVER_ADDR] [PORT] [FILE_NAME]\n", os.Args[0] )
		os.Exit (1)
	}

	selectedFile := os.Args[3]

	connection, err := net.Dial("tcp", os.Args[1] + ":" + os.Args[2])

	if err != nil {
		fmt.Println( "Can't connect to server: ", err)
	}

	defer connection.Close()
	
	// read test file
	file, err := os.Open( selectedFile )

	if err != nil {
		fmt.Println( "Something went wrong: ", err )
		os.Exit(1)
	}

	// collect file info
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println( "Something went wrong: ", err )
		os.Exit (1)
	}

	//calculate and file size
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)

	//set file name
	fileName := fillString(fileInfo.Name(), 64)

	fmt.Println("Sending filename and filesize!")

	connection.Write([]byte(fileSize))
	fmt.Println( "Send file size: ", fileSize  )

	connection.Write([]byte(fileName))
	fmt.Println( "Send file size: ", fileName  )

	sendBuffer := make([]byte, BUFFERSIZE)

	//save original size
	origSize := fileInfo.Size()
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		// send part to server
		connection.Write(sendBuffer)

		// decrease original size
		origSize -= BUFFERSIZE 

		// print to output current size of non delivered bytes
		fmt.Printf("\rLeft (bytes): %d", origSize )
	}

	fmt.Println("\nFile has been sent, closing connection!")
	return

}
