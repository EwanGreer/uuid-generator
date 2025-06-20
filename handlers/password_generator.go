package handlers

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/EwanGreer/uuid-generator/templater"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandlePasswordGeneratorPage(c echo.Context) error {
	passwordLength := c.FormValue("password-length")
	if passwordLength == "" {
		passwordLength = "10"
	}

	pwl, err := strconv.Atoi(passwordLength)
	if err != nil {
		return echo.NewHTTPError(500, "invalid password lenght")
	}

	if pwl == 0 || slices.Contains([]int{1, 2, 3}, pwl) {
		pwl = 16
	}

	includeUppercase := c.FormValue("include-uppercase")
	if includeUppercase == "" {
		includeUppercase = "false"
	}

	uppercase, err := strconv.ParseBool(includeUppercase)
	if err != nil {
		return echo.NewHTTPError(500, "invalid form value")
	}

	includeSymbols := c.FormValue("include-symbols")
	if includeSymbols == "" {
		includeSymbols = "false"
	}

	symbols, err := strconv.ParseBool(includeSymbols)
	if err != nil {
		return echo.NewHTTPError(500, "invalid form value")
	}

	tmpl := template.Must(h.Base.Clone())
	t := templater.NewTemplater(h.PublicFS)

	err = t.FindTemplate(tmpl, "_nav", "password-generator")
	if err != nil {
		return echo.NewHTTPError(500, fmt.Sprintf("Template parse error: %v", err))
	}

	generatedPassword, err := generateRandomHash(pwl, uppercase, symbols)
	if err != nil {
		return echo.NewHTTPError(500, "could not generate password")
	}

	data := map[string]any{
		"title":              "Password Generator",
		"version":            h.Version,
		"password_length":    pwl,
		"generated_password": generatedPassword,
		"include_uppercase":  uppercase,
		"include_symbols":    symbols,
	}

	return tmpl.ExecuteTemplate(c.Response(), "base", data)
}

func generateRandomHash(length int, uppercase bool, symbols bool) (string, error) {
	if length < 0 {
		length *= -1
	}

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/"
	if !symbols {
		charset = charset[:27]
	}

	hash := make([]byte, length)
	for i := range hash {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		hash[i] = charset[num.Int64()]
	}
	generatedPassword := string(hash)

	if !uppercase {
		generatedPassword = strings.ToLower(generatedPassword)
	}

	return generatedPassword, nil
}
