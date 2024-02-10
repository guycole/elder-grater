package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
)

const banner = "elder-grater 0.0"

func main() {
	log.Println(banner)

	// environment dump
	/*
		for _, envVar := range os.Environ() {
			log.Println(envVar)
		}
	*/

	awsWebIdentityTokenFileName, errorFlag := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if errorFlag {
		log.Println("AWS_WEB_IDENTITY_TOKEN_FILE: ", awsWebIdentityTokenFileName)
	} else {
		log.Println("AWS_WEB_IDENTITY_TOKEN_FILE: missing")
	}

	sess := session.Must(session.NewSession())
	initStsClient := sts.New(sess)

	if strings.Compare(awsWebIdentityTokenFileName, "") == 0 {
		log.Println("empty web identity token file skips assume role")
	} else {
		log.Println("assume role with web identity token")

		awsRoleArn, errorFlag := os.LookupEnv("AWS_ROLE_ARN")
		if errorFlag {
			log.Println("AWS_ROLE_ARN: ", awsRoleArn)
		} else {
			log.Println("AWS_ROLE_ARN: missing")
		}

		awsWebIdentityToken, err := os.ReadFile(awsWebIdentityTokenFileName)
		if err != nil {
			panic(err)
		}

		arwwii := sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(awsRoleArn),
			RoleSessionName:  aws.String("elder-grater"),
			WebIdentityToken: aws.String(string(awsWebIdentityToken)),
			DurationSeconds:  aws.Int64(3600),
		}

		identity, err := initStsClient.AssumeRoleWithWebIdentity(&arwwii)
		if err != nil {
			panic(err)
		}

		config := aws.Config{
			Credentials: credentials.NewStaticCredentialsFromCreds(
				credentials.Value{
					AccessKeyID:     *identity.Credentials.AccessKeyId,
					SecretAccessKey: *identity.Credentials.SecretAccessKey,
					SessionToken:    *identity.Credentials.SessionToken,
					ProviderName:    "AssumeRoleWithWebIdentity",
				}),
		}

		sess = session.Must(session.NewSession(&config))
		log.Println("session created with assumed role")
		log.Println(sess)
	}

	stsClient := sts.New(sess)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		panic(err)
	}

	// identity
	jsonIdentity, _ := json.MarshalIndent(*identity, "", "  ")
	log.Printf("%s", string(jsonIdentity))

	// s3
	s3Client := s3.New(sess)
	buckets, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err == nil {
		jsonBuckets, _ := json.MarshalIndent(*buckets, "", "  ")
		log.Printf("%+v", string(jsonBuckets))
	} else {
		log.Println(err)
	}

	// SQS
	sqsClient := sqs.New(sess)
	queues, err := sqsClient.ListQueues(nil)
	if err == nil {
		for _, queue := range queues.QueueUrls {
			log.Println(*queue)
		}
	} else {
		log.Println(err)
	}

	// sleep
	for {
		log.Println("sleeping...")
		time.Sleep(33 * time.Second)
	}
}
