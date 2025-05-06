package enums

import "github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"

type SortOrderEnum struct{}

var SortOrder SortOrderEnum

func (SortOrderEnum) AscendingInt() int {
	return 1
}

func (SortOrderEnum) DescendingInt() int {
	return -1
}

func (SortOrderEnum) AscendingStr() string {
	return "ascending"
}

func (SortOrderEnum) DescendingStr() string {
	return "descending"
}

func (SortOrderEnum) ConvertSortOrderEnumToInt(value string) (int, error) {
	switch value {
	case "ascending":
		return SortOrder.AscendingInt(), nil
	case "descending":
		return SortOrder.DescendingInt(), nil
	default:
		return 404, utils.BadRequestError("Invalid SortOrderEnum value")
	}
}

func (SortOrderEnum) ConvertSortOrderEnumToString(value int) (string, error) {
	switch value {
	case 1:
		return SortOrder.AscendingStr(), nil
	case -1:
		return SortOrder.DescendingStr(), nil
	default:
		return "unknown", utils.BadRequestError("Invalid SortOrderEnum value")
	}
}
