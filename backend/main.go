package main

import (
	"log"

	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/contrib/websocket"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    app := fiber.New()
    app.Use(cors.New())
    setupdb()

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
    app.Get("/quiz", getQuiz)
    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

    }))


    log.Fatal(app.Listen(":3000"))
}
func setupdb() {
    // Use the SetServerAPIOptions() method to set the version of the Stable API on the client
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI("mongodb+srv://abhijitpandey524:SKUheDAAn0YPf8Ky@cluster0.jjmav.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
  
    // Create a new client and connect to the server
    client, err := mongo.Connect(context.Background(), opts)
    if err != nil {
      panic(err)
    } 
  
    defer func() {
      if err = client.Disconnect(context.Background()); err != nil {
        panic(err)
      }
    }()
  
    // Send a ping to confirm a successful connection
    if err := client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Err(); err != nil {
      panic(err)
    }
    fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
    quizcollection := client.Database("kahoot").Collection("quizzes")
    fmt.Println(*quizcollection)

  }
  

func getQuiz(c *fiber.Ctx) error {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://abhijitpandey524:SKUheDAAn0YPf8Ky@cluster0.jjmav.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
    if err != nil {
      panic(err)
    } 
    quizcollection := client.Database("kahoot").Collection("quizzes")

    cursor, err :=  quizcollection.Find(context.Background(), bson.M{})

    quizzes :=  []map[string]any{}
    err_ := cursor.All(context.Background(), &quizzes)
    if err_ != nil {
        panic(err_)
    }
    // list := []map[string]any{
    //     map[string]any{
    //         "test": 123,
    //     },
    // }
    return c.JSON(quizzes)
}


  
  