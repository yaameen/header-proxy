package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppPort   int               `yaml:"app_port" validate:"gte=80,mustBePort"`
	HeaderMap map[string]string `yaml:"header_map"`
	Target    struct {
		Host   string `yaml:"host" validate:"required"`
		Port   int    `yaml:"port" validate:"required,gte=80,mustBePort"`
		Scheme string `yaml:"scheme"`
	}
}

func Parse(f string) *Config {
	cf, err := os.OpenFile(f, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	data := make([]byte, 1024)
	_, err = io.ReadFull(cf, data)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	val := validator.New()
	val.RegisterValidation("mustBePort", func(fl validator.FieldLevel) bool {
		if port := fl.Field().Int(); port >= 80 && port <= 65535 {
			return true
		}
		return false
	})

	if err := val.Struct(config); err != nil {
		if err != nil {
			println("Your config is invalid. Fix the following errors: ")
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Printf("- %s\n", err.Namespace())
			}
		}
		os.Exit(1)
	}

	return &config
}

func main() {

	g := flag.Bool("g", false, "Generate a config file")
	v := flag.Bool("v", false, "Verbose mode")
	c := flag.String("c", "config.yaml", "Generate a config file")

	flag.Parse()

	if *g {
		os.WriteFile(`config.yaml`, []byte(`app_port: 8080
target:
 host: mihaaru.com
 port: 443
 scheme: https

header_map:
 app_number: app-number
 app_version: app-version`), 0644)
		os.Exit(0)
	}

	if _, err := os.Stat(*c); os.IsNotExist(err) {
		println("Could not find config file. Use -g to generate a config file.")
		os.Exit(1)
	}

	config := Parse(*c)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlx := url.URL{
			Scheme:   config.Target.Scheme,
			Path:     r.URL.Path,
			RawQuery: r.URL.Query().Encode(),
		}

		urlx.Host = config.Target.Host + ":" + strconv.Itoa(config.Target.Port)

		req := http.Request{
			Method: r.Method,
			URL:    &urlx,
			Header: r.Header,
		}

		if resp, err := http.DefaultClient.Do(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if *v {
				w.Write([]byte(err.Error()))
			} else {
				w.Write([]byte("Something went wrong"))
			}
		} else {
			defer resp.Body.Close()

			for k, v := range req.Header {
				if header, ok := config.HeaderMap[strings.ToLower(k)]; ok {
					w.Header().Set(header, "****")
				} else {
					w.Header().Set(k, v[0])
				}
			}
			for k, v := range resp.Header {
				w.Header().Set(k, v[0])
			}

			if _, err := io.Copy(w, resp.Body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				if *v {
					w.Write([]byte(err.Error()))
				} else {
					w.Write([]byte("Something went wrong"))
				}
			}
		}
	})

	port := strconv.Itoa(config.AppPort)

	println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)

}
