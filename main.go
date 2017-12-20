package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName, keyName, outputPath string
var recursive bool

func init() {
	flag.StringVar(&bucketName, "bucket", "", "Bucket name where the target object is located")
	flag.StringVar(&keyName, "key", "", "Key name where the target object is located")
	flag.StringVar(&outputPath, "output", "", "Path to the output file")
	flag.BoolVar(&recursive, "recursive", false, "Recursively get files")
}

func main() {
	flag.Parse()
	if bucketName == "" {
		fmt.Fprintln(os.Stderr, "-bucket must be specified")
		flag.Usage()
		os.Exit(1)
	}
	if keyName == "" {
		fmt.Fprintln(os.Stderr, "-key must be specified")
		flag.Usage()
		os.Exit(1)
	}
	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "-output must be specified")
		flag.Usage()
		os.Exit(1)
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	if recursive {
		getObjectRecursively(svc, keyName, outputPath)
	} else {
		getObject(svc, keyName, outputPath)
	}
}

func getObject(svc *s3.S3, keyName string, outputPath string) {
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer output.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, output.Body); err != nil {
		log.Fatal(err)
	}
}

func getObjectRecursively(svc *s3.S3, keyPrefix string, outputDir string) {
	err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{Bucket: aws.String(bucketName), Prefix: aws.String(keyPrefix)}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, content := range page.Contents {
			key := *content.Key
			outputPath := path.Join(outputDir, key[len(keyPrefix):])
			os.MkdirAll(path.Dir(outputPath), 0755)
			getObject(svc, key, outputPath)
		}
		return true
	})
	if err != nil {
		log.Fatal(err)
	}
}
