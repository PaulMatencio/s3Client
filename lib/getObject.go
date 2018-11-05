package s3Client

import (
	"bytes"
	"github.com/minio/minio-go"
	"log"
	"os"
	"errors"
)


type S3Request struct {
	MinioC 	*minio.Client
	// Buf 	*bytes.Buffer
	Bucket 	string
	Key 	string
	Opts 	*minio.GetObjectOptions
	BufSize int
	Trace 	bool
}



func (s3Request *S3Request) SetS3Client(s3c *minio.Client) {
	s3Request.MinioC = s3c
}


func (s3Request *S3Request) SetBucketName(bucket string) {
	s3Request.Bucket = bucket
}

func (s3Request *S3Request) SetKey(key string) {
	s3Request.Bucket = key
}

func (s3Request *S3Request) SetPutOpts(options *minio.GetObjectOptions) {
	s3Request.Opts= options
}

func (s3Request *S3Request) SetTrace(trace bool) {
	s3Request.Trace = trace
}

func (s3Request *S3Request) S3BuildGetRequest(login *S3Login, bucket string, key string, options *minio.GetObjectOptions){

	s3Request.MinioC 	=  login.MinioC
	s3Request.Bucket 	=  bucket
	s3Request.Key		=  key
	s3Request.Opts 	    =  options
	s3Request.BufSize   =  BUFFERSIZE  /* constant */

	if TRACE {
		s3Request.Trace	= true
	} else {
		s3Request.Trace	= false
	}
}


func GetObject(request S3Request) (*bytes.Buffer, error) {

	if request.Trace {
		request.MinioC.TraceOn(os.Stdout)
	}

	object,err := request.MinioC.GetObject(request.Bucket,request.Key,*request.Opts);

	if err != nil {
		log.Fatalln(err)
	}

	defer object.Close()

	/* Retrieve the object*/
	buffer:= make([]byte,request.BufSize)
	var  buf = new(bytes.Buffer)
	var (
		n int
		size int
	)

	for {
		n, _ = object.Read(buffer)
		size += n
		buf.Write(buffer[0:n])
		if n == 0 {
			break
		}
	}
	if buf.Len() == 0 {
		err = errors.New(request.Key + " not Found")
	}
	return buf,err

}
