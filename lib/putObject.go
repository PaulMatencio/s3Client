package s3Client

import (
	"bytes"
	"errors"
	"github.com/minio/minio-go"
	"io"
	"os"
)

type S3PutRequest struct {
	MinioC 	*minio.Client
	Buf 	*bytes.Buffer
	Bucket 	string
	Key 	string
	PutOpts *minio.PutObjectOptions
	Trace 	bool
}



func (s3Request *S3PutRequest) SetS3Client(s3c *minio.Client) {
	s3Request.MinioC = s3c
}

func (s3Request *S3PutRequest) SetByteBuffer(buf *bytes.Buffer) {
	s3Request.Buf = buf
}

func (s3Request *S3PutRequest) SetBucketName(bucket string) {
	s3Request.Bucket = bucket
}

func (s3Request *S3PutRequest) SetKey(key string) {
	s3Request.Bucket = key
}

func (s3Request *S3PutRequest) SetPutOpts(options *minio.PutObjectOptions) {
	s3Request.PutOpts= options
}

func (s3Request *S3PutRequest) SetTrace(trace bool) {
	s3Request.Trace = trace
}

func (s3Request *S3PutRequest) S3BuildPutRequest(login *S3Login, bucket string, key string, buf *bytes.Buffer, putoptions *minio.PutObjectOptions){

	s3Request.MinioC 	=  login.MinioC
	s3Request.Buf    	=  buf
	s3Request.Bucket 	=  bucket
	s3Request.Key		=  key
	s3Request.PutOpts 	=  putoptions

	if TRACE {
		s3Request.Trace	= true
	} else {
		s3Request.Trace	= false
	}
}


func PutObject(request S3PutRequest) (int64, error) {

	if request.Trace {
		request.MinioC.TraceOn(os.Stdout)
	}

	var (
		r io.Reader
		b []byte
	)

	if b = request.Buf.Bytes(); b != nil {
		 r = bytes.NewReader(b)
	} else {
		return 0,errors.New("Input buffer image is empty !!")
	}

	return request.MinioC.PutObject(request.Bucket,request.Key,r, int64(len(b)),*request.PutOpts);

}


