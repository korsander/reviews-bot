module github.com/korsander/reviews-bot/bot

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/korsander/reviews-bot/lib v0.0.0-00010101000000-000000000000
	github.com/slack-go/slack v0.6.5
)

replace github.com/korsander/reviews-bot/lib => ../lib
