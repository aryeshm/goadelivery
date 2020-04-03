package context

import (
	"context"
	"fmt"
	"log"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OpenMgDb(ctx context.Context, config Config) (*mongo.Database, error) {
	/*
				uri := fmt.Sprintf(``)

				if ctx.Value(UsernameKey).(string) != "" {
					uri = fmt.Sprintf(`mongodb://%s:%s@%s`,
						ctx.Value(UsernameKey).(string),
						ctx.Value(PasswordKey).(string),
						ctx.Value(HostKey).(string),
					)
				} else {
					uri = fmt.Sprintf(`mongodb://%s`,
						ctx.Value(HostKey).(string),
					)
				}

				fmt.Println(uri)

			clientOptions := options.Client().ApplyURI(uri)

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://aryeshm:Arj2harv%402012@cluster0-wvbnx.mongodb.net/test?retryWrites=true&w=majority",
	))

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	//log.Println("Username",config.Db.User)
	collection := client.Database(config.Db.Dbname)

	return collection, err

}
