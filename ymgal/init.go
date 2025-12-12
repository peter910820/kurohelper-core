package ymgal

var cfg config

func Init(endPoint string, clientID string, clientSecret string) {
	cfg.Endpoint = endPoint
	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret
}
