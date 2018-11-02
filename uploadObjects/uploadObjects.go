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
	flag.StringVar(&prefix, "prefix", "", "-prefix aSring  => EX: pxi")
	flag.StringVar(&ftype,"ftype","","-ftype aString  => EX: tiff")
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
	if len(prefix) > 0 {
		prefix += "/"
	}

	if len(ftype) > 0  {
		prefix += ftype + "/"
	}
	log.Printf("Files will be upload with a prefix of %s\n",prefix)

	for _, f := range files {
		if !f.IsDir() { // skip directory
			name := f.Name()
			if len(ftype) > 0 {
				/*
				   Specific use case upload all files with a filetype  ftype and filetype = md
				   this is for uploading moses and Pxi objects to S3
				*/
				ft := strings.Split(name, ".")
				if len(ft) > 1 && ft[len(ft)-1] == ftype {
					filename := Files{Filename: name,Metadata:ft[0]+".md"}
					filenames = append(filenames,filename)
				}
			} else  {
				/* upload all files in the directory  */
				filenames = append(filenames,Files{Filename:name,Metadata:""})
			}
		}
	}


	N:= len(filenames)
	log.Printf("%d files will be uploaded from %s directory\n",N,directory)
	if ( N == 0) {
		log.Printf("It seems that there is nothing to upload. Check your input parametres. Don't forget to quote the -ftype\n")
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
    /* exit if the bucket does not exist */
	if exist,err := s3client.BucketExists(bucketName); exist == false || err != nil {
		log.Printf("Bucket %s does not exist or something went wrong: %v \n",bucketName,err)
		os.Exit(100)
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

		if metadata := filenames[obj].Metadata;len(metadata)>0 {
			if meta, err := os.Open(metadata); err == nil {
				request.Metaname = metadata
				request.Meta = meta
				defer meta.Close()
			}
		}

		if request.File != nil {
			go func(Request) {
				var (
				 size int64
				 start time.Time
				 filename string
				)
				/*
				file:= request.File
				fileStat, err := file.Stat()
				*/
				file:= request.File
				filename = request.Filename
				if fileStat,err := file.Stat(); err == nil {
					start  = time.Now()
					size   = fileStat.Size()
					objectName := prefix + filename

					/* add user meata data here
					   Read same filename  with ftype md  into a byte array
					   S3 user metadata must have  a format of map[string]string
					*/
					opts := minio.PutObjectOptions{}
					// user meta data
					if len(request.Metaname) > 0 {
						if s, err := request.Meta.Stat();err == nil {
							// metadata exists
							n := s.Size()
							b := make([]byte, n)
							var usermeta map[string]string
							if request.Meta != nil {
								request.Meta.Read(b)
								usermeta["usermd"] = user.Encode64(b)
								opts.UserMetadata = usermeta
							}

							opts.ContentType = "image/" + ftype
						}
					}
					// System metadata
					opts.StorageClass = "STANDARD"
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