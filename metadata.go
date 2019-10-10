package emstrigger

type Settings struct {
	ServerURL   string `md:"serverURL,required"` // The port to listen on
	Destination string `md:"destination"`        // Enable TLS on the server
	Username    string `md:"username"`           // The path to PEM encoded server certificate
	Password    string `md:"password"`           // The path to PEM encoded server key
}

type Output struct {
	Data interface{} `md:"data"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Data = values["data"]

	return nil
}
