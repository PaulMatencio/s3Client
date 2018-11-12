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



func (r *S3Request) SetS3Client(s3c *minio.Client) {
	r.MinioC = s3c
}


func (r *S3Request) SetBucketName(bucket string) {
	r.Bucket = bucket
}

func (r *S3Request) SetKey(key string) {
	r.Bucket = key
}

func (r *S3Request) SetPutOpts(options *minio.GetObjectOptions) {
	r.Opts= options
}

func (r *S3Request) SetTrace(trace bool) {
	r.Trace = trace
}

func (r *S3Request) S3BuildGetRequest(login *S3Login, bucket string, key string, options *minio.GetObjectOptions){

	r.MinioC 	=  login.MinioC
	r.Bucket 	=  bucket
	r.Key		=  key
	r.Opts 	    =  options
	r.BufSize   =  BUFFERSIZE  /* constant */

	if TRACE {
		r.Trace	= true
	} else {
		r.Trace	= false
	}
}


func GetObject(r S3Request) (*bytes.Buffer, error) {

	if r.Trace {
		r.MinioC.TraceOn(os.Stdout)
	}

	object,err := r.MinioC.GetObject(r.Bucket,r.Key,*r.Opts);

	if err != nil {
		log.Fatalln(err)
	}

	defer object.Close()

	/* Retrieve the object*/
	buffer:= make([]byte,r.BufSize)
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
		err = errors.New(r.Key + " not Found")
	}
	return buf,err

}
