package dockerlink

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrLinkNotDefined    = errors.New("link not defined")
	ErrPortNotDefined    = errors.New("port not defined")
	ErrAddressNotDefined = errors.New("address not defined")
)

// Link represents a Docker container link
type Link struct {
	Name        string
	Protocol    string
	ExposedPort int
	Port        int
	Address     string
}

// GetLink returns a Link configured with Docker defined environment vars
func GetLink(name string, port int, proto string) (*Link, error) {
	if proto == "" {
		proto = "TCP"
	}

	prefix := fmt.Sprintf(
		"%s_PORT_%d_%s",
		strings.ToUpper(name),
		port,
		strings.ToUpper(proto))

	if os.Getenv(prefix) == "" {
		return nil, ErrLinkNotDefined
	}

	portInt, err := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_PORT", prefix)))
	if err != nil {
		return nil, ErrPortNotDefined
	}
	addr := os.Getenv(fmt.Sprintf("%s_ADDR", prefix))
	if addr == "" {
		return nil, ErrAddressNotDefined
	}

	l := &Link{
		Name:        name,
		Protocol:    strings.ToLower(proto),
		ExposedPort: port,
		Port:        portInt,
		Address:     addr,
	}
	return l, nil
}
