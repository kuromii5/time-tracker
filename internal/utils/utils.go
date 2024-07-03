package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kuromii5/time-tracker/internal/models"
)

// Helper function to build SQL query for getting filtered and paginated users
func BuildGetUsersQuery(filter models.FilterBy, settings models.Pagination) (string, []interface{}) {
	var baseQuery strings.Builder
	baseQuery.WriteString(`SELECT * FROM users WHERE 1=1`)
	var args []interface{}
	argIndex := 1

	stringFields := map[string]string{
		"surname":         filter.Surname,
		"name":            filter.Name,
		"patronymic":      filter.Patronymic,
		"address":         filter.Address,
		"passport_serie":  filter.PassportSerie,
		"passport_number": filter.PassportNumber,
	}
	for field, value := range stringFields {
		if value != "" {
			baseQuery.WriteString(fmt.Sprintf(" AND %s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	// filter user creation time
	if !filter.CreatedAfter.IsZero() {
		baseQuery.WriteString(fmt.Sprintf(" AND created_at > $%d", argIndex))
		args = append(args, filter.CreatedAfter)
		argIndex++
	}
	if !filter.CreatedBefore.IsZero() {
		baseQuery.WriteString(fmt.Sprintf(" AND created_at < $%d", argIndex))
		args = append(args, filter.CreatedBefore)
		argIndex++
	}

	// pagination
	if settings.Limit > 0 {
		baseQuery.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
		args = append(args, settings.Limit)
		argIndex++
	}
	if settings.Offset > 0 {
		baseQuery.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
		args = append(args, settings.Offset)
		argIndex++
	}

	return baseQuery.String(), args
}

func BuildUpdateUserQuery(user models.User) (string, []interface{}) {
	var statements strings.Builder
	var args []interface{}
	argIndex := 1

	// Helper function to add fields to the query
	addField := func(fieldName string, fieldValue interface{}) {
		if statements.Len() > 0 {
			statements.WriteString(", ")
		}
		statements.WriteString(fmt.Sprintf("%s = $%d", fieldName, argIndex))
		args = append(args, fieldValue)
		argIndex++
	}

	// Add fields to the query if they are not empty
	if user.Passport.Serie != "" {
		addField("passport_serie", user.Passport.Serie)
	}
	if user.Passport.Number != "" {
		addField("passport_number", user.Passport.Number)
	}
	if user.People.Name != "" {
		addField("name", user.People.Name)
	}
	if user.People.Surname != "" {
		addField("surname", user.People.Surname)
	}
	if user.People.Patronymic != "" {
		addField("patronymic", user.People.Patronymic)
	}
	if user.People.Address != "" {
		addField("address", user.People.Address)
	}

	// Add the updated_at field
	if statements.Len() > 0 {
		statements.WriteString(", ")
	}
	statements.WriteString("updated_at = NOW()")

	// Append the user ID to the arguments
	args = append(args, user.ID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", &statements, argIndex)

	return query, args
}

// Helper function to parse query parameters as integers
func ParseQueryParamInt(r *http.Request, key string) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}

	return value
}

// Helper function to parse query parameters as time
func ParseQueryParamTime(r *http.Request, key string) time.Time {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return time.Time{}
	}

	value, err := time.Parse(time.RFC3339, valueStr)
	if err != nil {
		return time.Time{}
	}

	return value
}

// ParsePassportData parses JSON containing passport serie and number into PassportData struct
func ParsePassportData(data string) (models.Passport, error) {
	// Split passportNumber into serie and number
	parts := strings.Fields(data)
	if len(parts) != 2 && len(parts[0]) != 4 && len(parts[1]) != 6 {
		return models.Passport{}, fmt.Errorf("invalid passportNumber format, expected '**** ******'")
	}

	passportData := models.Passport{
		Serie:  parts[0],
		Number: parts[1],
	}

	return passportData, nil
}
