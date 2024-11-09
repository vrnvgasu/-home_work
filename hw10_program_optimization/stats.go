package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

//easyjson:json
type A struct {
	Id       int //nolint:stylecheck,revive
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)
	i := 0
	var (
		line []byte
		user = &User{}
	)

	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}

		if err = user.UnmarshalJSON(line); err != nil {
			return result, fmt.Errorf("unmarshal error: %w", err)
		}

		result[i] = *user

		i++
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainWithDot := "." + domain
	sep := "@"

	for _, user := range u {
		if strings.HasSuffix(user.Email, domainWithDot) {
			result[strings.ToLower(strings.SplitN(user.Email, sep, 2)[1])]++
		}
	}
	return result, nil
}
