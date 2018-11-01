package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	user "github.com/moses/user/base64j"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	bucketName 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
	start       time.Time
	prefix      string
	ftype      string
	help       bool
)


func main() {

	var (
		directory   string
	)
	/* Go routine response */
	type Response struct {
		Filename  string
		Size      int64
		Duration  time.Duration
		Err       error
	}

	/*   pair of file and its metadata */
	type Files struct {
		Filename string
		Metadata string
	}

	/*  goroutine parameters  */
	type Request struct {
		File *os.File
		Filename string
		Meta *os.File
		Metaname string
	}

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&directory, "d", "", "-d directory")
	flag.StringVar(&prefix, "prefix", "images", "-prefix aSring  => EX: images")
	flag.StringVar(&ftype,"ftype","tiff","-ftype aString  => EX: tiff")
	flag.BoolVar(&help,"help",false,"-help")

	flag.Parse()
	if len(bucketName) == 0 || len(directory) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	/* read the directory */
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	/* looking for file names to upload*/
	//filenames := []string{}
	filenames := []Files{}
	prefix += "/"
	if ftype != "*" {
		prefix += ftype + "/"
	}

	for _, f := range files {
		if !f.IsDir() {
			name := f.Name()
			if ftype != "*" {
				ft := strings.Split(name, ".")
				if len(ft) > 1 && ft[len(ft)-1] == ftype {

					filename := Files{Filename: name,Metadata:ft[0]+".md"}
					filenames = append(filenames,filename)
				}
			} else  {
				/* todo */
				return
			}
		}
	}

	/* looking for user metadata to upload */

	/*       to be done                    */



	N:= len(filenames)
	if ( N == 0) {
		log.Printf("Directory s% is empty",directory)
		return
	}

	/* get the Config */

	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* create an S3 session */
	s3 := s3Client.SetS3Session(s3Config, location)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey, s3.SSL)
	if err != nil {
		log.Fatalln(err)
	}
	opts := minio.PutObjectOptions{}
	opts.ContentType = "application/octet-stream"
	opts.StorageClass = "STANDARD"
	/* add some user metadata  */
	// s3Client.SetMetadata(opts)
	/* create a goroutine chanel */
	messages := make(chan Response)
	runtime.GOMAXPROCS(4)

	T:=1
	/* */
	// s3client.TraceOn(os.Stdout)

	start0 := time.Now()
	for  obj:=0; obj < N; obj++ {
		start = time.Now()
		request := Request{
			Filename:filenames[obj].Filename,
			File : nil,
			Metaname:filenames[obj].Metadata,
			Meta : nil,
		}
		var n int64 = 0;
		filename :=  filenames[obj].Filename
		if file,err  := os.Open(filename); err == nil {
			request.Filename = filename
			request.File = file
			defer file.Close()
		}

		metadata := filenames[obj].Metadata
		if meta,err :=  os.Open(metadata);err == nil {
			request.Metaname = metadata
			request.Meta = meta
			defer meta.Close()
		}
		if request.File != nil {
			go func(Request) {
				var (
				 size int64
				 start time.Time
				 filename string
				)
				if  request.File != nil {
					file:= request.File
					fileStat, _ := file.Stat()
					size   = fileStat.Size()
					start  = time.Now()
					filename = request.Filename
					objectName := prefix + filename

					/* add user meata data here
					   Read <file>.md into a byte array
					   S3 user metadata should be a map[string]string
					*/
					opts := minio.PutObjectOptions{}
					s,err := request.Meta.Stat()
					if err == nil {
						// metadata exists
						n := s.Size()
						b := make([]byte, n)
						var usermeta map[string]string
						if request.Meta != nil {
							request.Meta.Read(b)
							usermeta["usermd"] = user.Encode64(b)
							opts.UserMetadata = usermeta
						}

						// opts.ContentType = "application/octet-stream"
						opts.StorageClass = "STANDARD"

						// log.Println(usermeta)
						opts.ContentType = "image/" + ftype
					}
					n, err = s3client.PutObject(bucketName, objectName, file,
						size, opts)
					if err != nil {
						s3Client.PrintNotOk(filename, err)
					}
				} else {
					s3Client.PrintNotOk(filename, err)
				}
				s3Client.PrintOk(filename,n,start)  /* ok or not */
				messages <- Response{filename,size, time.Since(start), err}
			}(request)
		}
	}
	var (
		totalSize int64 = 0
		totalDuration time.Duration = 0
	)

	for {
		select {
		case   r:=  <-messages:
			{
				totalSize += r.Size
				totalDuration += r.Duration
				if T ==  N  {
					elapsedTm := time.Since(start0)
					MBps :=  1000* float64(totalSize)/float64(elapsedTm)
					log.Printf("Uplooad %d objects[ Total size (byte): %d - MBps: %0.4f] in Total Duration:%s/Total Elapsed:%s\n",N, totalSize, MBps, totalDuration, elapsedTm)
					return
			     }
				T++
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf("w")
		}
	}


}