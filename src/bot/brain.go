package bot

// DynamoDB Brain
import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Brain struct {
	Bot *Bot
}

func getClient() (*dynamodb.DynamoDB) {
	config := &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	sess := session.Must(session.NewSession(config))
	return dynamodb.New(sess)
}

func (b *Brain) Save(name string, j []byte) {
	client := getClient()
	_, err := client.UpdateItem(
		&dynamodb.UpdateItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"key": {
					S: aws.String(name),
				},
			},
			TableName: aws.String(os.Getenv("DYNAMO_TABLE")),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":r": {
						S: aws.String(string(j)),
					},
			},
			UpdateExpression: aws.String("set json = :r"),
		},
	)
	if err != nil {
			log.Fatalf("Got error calling UpdateItem: %s", err)
	}
}

func (b *Brain) Get(name string) ([]byte) {
	client := getClient()
	result, err := client.GetItem(
		&dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"key": {
					S: aws.String(name),
				},
			},
			TableName: aws.String(os.Getenv("DYNAMO_TABLE")),
		},
	)
	if err != nil {
		log.Println(err)
		return nil
	}

	type Item struct {
		JSON string `dynamodbav:"json"`
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return []byte(item.JSON)
}

// Redis Brain

// import (
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"github.com/gomodule/redigo/redis"
// )
//
// type Brain struct {
// 	Bot *Bot
// }
//
// func (b *Brain) Save(name string, j []byte) {
//     conn, err := redis.Dial("tcp", os.Getenv("REDIS_URL"))
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	defer conn.Close()
//
// 	_, err = conn.Do("SET", fmt.Sprintf("%s/%s", b.Bot.Name, name), j)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
//
// func (b *Brain) Get(name string) ([]byte) {
//     conn, err := redis.Dial("tcp", os.Getenv("REDIS_URL"))
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	defer conn.Close()
//
// 	result, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("%s/%s", b.Bot.Name, name)))
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	return result
// }
