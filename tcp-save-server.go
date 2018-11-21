package main

import "io" // wotj with files
import "net" // work with network
import "os" // work with  arguments
import "strconv" // convert data
import "fmt" // show log
import "strings"

// default buffer for reading data from file 
const BUFFERSIZE = 1024
const DEFAULTPROTO = "tcp"

func main () {

	//default port for listening
	DEFAULTPORT := "9001"

	// checking arguments, if port set in argument, use as default
	if  len( os.Args ) > 1 {
		DEFAULTPORT = os.Args[1]
		fmt.Println( "Set port: ", DEFAULTPORT )
	} 


	// strart listening
	server, err := net.Listen(DEFAULTPROTO, ":" + DEFAULTPORT )

	// if something wrong, exit with information about problem
	if err != nil {
		fmt.Println("We can start server: ", err)
		return 
	}

	// close listening if exit application
	defer server.Close()

	// show message about starting application
	fmt.Println("Server started and serve port " +  DEFAULTPORT  )

	// main loop fo accepting connection
	for {

		// accept incoming connection
		connection, err := server.Accept()

		// close if application exit
		defer connection.Close()

		// if something wrong, inform and exit with os exit code 1
		if err != nil {
		
			fmt.Println("something went wrong: ", err)
			os.Exit(1)
		
		}
		// inform if connection present on port
		fmt.Println( "Incoming connection!" )

		// set buffer for store filename
		bufFileName := make([]byte, 64)

		// set buffer for store file size
		bufFileSize := make([]byte, 10)

		// read data
		connection.Read(bufFileSize)

		// parse and set filesize 
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufFileSize), ":"), 10, 64)

		// next read file name 
		connection.Read(bufFileName)
		fileName := strings.Trim(string(bufFileName), ":")

		// print information what we receiving
		fmt.Printf ( "Receiving file %s with size, %d\n", fileName, fileSize )
		fmt.Printf ( "Writing file\n" )

		// create new file, will be rewrited existing
		newFile, err := os.Create(fileName)

		// on exit, close file
		defer newFile.Close()
		var receivedBytes int64

		// main loop for receiving file
		for {
			// if calculated size will be less then default buffer size, read/write and break loop
			if (fileSize - receivedBytes) < BUFFERSIZE {

				// read and write to file
				io.CopyN(newFile, connection, (fileSize - receivedBytes))
				connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
				break
			}

			// read a nd write to file, next iteration
			io.CopyN(newFile, connection, BUFFERSIZE)
			receivedBytes += BUFFERSIZE
		
		}
		// print if receiving completed
		fmt.Printf( "Receiving complete\n")
	}
}

