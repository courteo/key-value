package compute

import (
	"errors"
	"regexp"
	"strings"

	"github.com/courteo/key-value/internal/domain"
	"github.com/courteo/key-value/internal/domain/command"
)

var (
	errCommandNotFound = errors.New("command not found")
	errInvalidQuery    = errors.New("invalid query")
)

func (c *Computer) ParseQuery(query string) (domain.Query, error) {
	if len(query) < 5 {
		return domain.Query{}, errCommandNotFound
	}

	cmd, err := c.getCommand(query[:3])
	if err != nil {
		return domain.Query{}, err
	}

	key, value, err := c.getArguments(query, cmd)
	if err != nil {
		return domain.Query{}, err
	}

	return domain.Query{
		Command: cmd,
		Key:     key,
		Value:   value,
	}, nil
}

func (c *Computer) getCommand(commandStr string) (int, error) {
	switch commandStr {
	case "SET":
		return command.SetID, nil
	case "GET":
		return command.GetID, nil
	case "DEL":
		return command.DeleteID, nil
	default:
		return 0, errCommandNotFound
	}
}

func (c *Computer) checkArgument(argument string) bool {
	r := regexp.MustCompile(`^[A-Za-z0-9*_\/!?:;."]+$`)

	return r.MatchString(argument)
}

func (c *Computer) getArguments(query string, cmd int) (key string, val string, err error) {
	firstIndex := strings.Index(query, " ")
	if firstIndex == -1 {
		return "", "", errInvalidQuery
	}

	if cmd == command.SetID {
		secondIndex := strings.LastIndex(query, " ")
		if secondIndex == -1 || firstIndex == secondIndex {
			return "", "", errInvalidQuery
		}

		val = query[secondIndex+1:]
		key = query[firstIndex+1 : secondIndex]
	} else {
		key = query[firstIndex+1:]
	}

	key= strings.ReplaceAll(key, "\n", "")
	val=  strings.ReplaceAll(val, "\n", "")

	return key, val, nil
}
