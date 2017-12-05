package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucketName, keyName, outputPath string

func init() {
	flag.StringVar(&bucketName, "bucket", "", "Bucket name where the target object is located")
	flag.StringVar(&keyName, "key", "", "Key name where the target object is located")
	flag.StringVar(&outputPath, "output", "", "Path to the output file")
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
