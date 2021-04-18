package config

type IConfig interface {
	Get(string) ([]byte, error)
}
