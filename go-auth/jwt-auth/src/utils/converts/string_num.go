package converts

import (
	"fmt"
	"strconv"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
)

func StringToInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, utils.BadRequestError(fmt.Sprintf("Error while converting string to int: %s", value))
	}
	return i, nil
}
