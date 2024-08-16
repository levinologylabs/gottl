package services

type Config struct {
	CompanyName string `json:"company_name" conf:"default:Gottl Inc."`
	WebURL      string `json:"web_url"      conf:"default:http://localhost:8080"`
}
