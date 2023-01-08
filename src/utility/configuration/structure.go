package configuration

type ConfigData struct {
	Connection struct {
		Mail    string
		Root    string
		DBLogin string
	}
	Encryption struct {
		AES string
		RSA string
	}
	Services struct {
		MySpace bool
		MSN     bool
		Yahoo   bool
		API     bool
		AIM     bool
	}
}
