module github.com/korsander/reviews-bot/client

go 1.14

require (
	github.com/gorilla/websocket v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/korsander/reviews-bot/lib v0.0.0-00010101000000-000000000000
)

replace github.com/korsander/reviews-bot/lib => ../lib
