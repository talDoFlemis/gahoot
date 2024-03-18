package web

type Config struct {
	Prefix           string
	ListenIP         string
	Port             int
	CORSAllowOrigins string
	CORSAllowHeaders string
	CORSAllowMethods string
}
