package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Website struct {
	Name string `yaml:"name"`
}

type AllowedUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Firestore struct {
	ProjectID string `yaml:"project_id"`
	CollectionName string `yaml:"collection_name"`
}

type Authentication struct {
	SigningKey string `yaml:"signing_key"`
	SigningDomain string `yaml:"signing_domain"`
	AllowedUsers []AllowedUser `yaml:"allowed_users"`
}

type ReservedManagementEndpoints struct {
	Login string `yaml:"login"`
	Logout string `yaml:"logout"`
}

type Config struct {
	Website Website `yaml:"website"`
	Authentication Authentication `yaml:"authentication"`
	Firestore Firestore `yaml:"firestore"`
	ReservedManagementEndpoints ReservedManagementEndpoints `yaml:"reserved_management_endpoints"`
}

func LoadConfig() *Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}