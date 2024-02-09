package main

import (
	"encoding/json"
	"time"

	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

const banner = "elder-grater 0.0"

func main() {
	log.Println(banner)

	awsWebIdentityTokenFile, err := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if err {
		log.Println("AWS_WEB_IDENTITY_TOKEN_FILE: ", awsWebIdentityTokenFile)
	} else {
		log.Println("AWS_WEB_IDENTITY_TOKEN_FILE: missing")
	}

	sess := session.Must(session.NewSession())
	initStsClient := sts.New(sess)

	if strings.Compare(awsWebIdentityTokenFile, "") == 0 {
		log.Println("empty web identity token file skips")
	} else {
		log.Println("assume role with web identity token")

		awsRoleArn, err1 := os.LookupEnv("AWS_ROLE_ARN")
		if err1 {
			log.Println("AWS_ROLE_ARN: ", awsRoleArn)
		} else {
			log.Println("AWS_ROLE_ARN: missing")
		}

		awsWebIdentityToken, err2 := os.ReadFile(awsWebIdentityTokenFile)
		if err2 != nil {
			panic(err2)
		}

		arwwii := sts.AssumeRoleWithWebIdentityInput{
			RoleArn:          aws.String(awsRoleArn),
			RoleSessionName:  aws.String("elder-grater"),
			WebIdentityToken: aws.String(string(awsWebIdentityToken)),
			DurationSeconds:  aws.Int64(3600),
		}

		identity, err3 := initStsClient.AssumeRoleWithWebIdentity(&arwwii)
		if err3 != nil {
			panic(err3)
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
	identity, err4 := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err4 != nil {
		panic(err4)
	}

	jsonIdentity, _ := json.MarshalIndent(*identity, "", "  ")
	log.Printf("%s", string(jsonIdentity))

	s3Client := s3.New(sess)
	buckets, err6 := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err6 != nil {
		panic(err6)
	}

	jsonBuckets, _ := json.MarshalIndent(*buckets, "", "  ")
	log.Printf("%+v", string(jsonBuckets))

	for {
		log.Println("sleeping...")
		time.Sleep(5 * time.Second)
	}
}
