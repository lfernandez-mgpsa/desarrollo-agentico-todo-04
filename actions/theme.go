package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// ThemeCookie es la cookie que guarda el color de fondo elegido.
const ThemeCookie = "todo_theme"

// DefaultTheme es el fondo que se usa si no hay elección válida.
const DefaultTheme = "blanco"

// ThemeOption es una opción del selector de fondo.
type ThemeOption struct {
	Value string
	Label string
}

// Themes son los fondos disponibles, en el orden en que se muestran.
var Themes = []ThemeOption{
	{Value: "oscuro", Label: "Oscuro"},
	{Value: "blanco", Label: "Blanco"},
	{Value: "azul", Label: "Azul"},
	{Value: "morado", Label: "Morado"},
}

func validTheme(v string) bool {
	for _, t := range Themes {
		if t.Value == v {
			return true
		}
	}
	return false
}

// SetTheme deja en el contexto el tema efectivo y las opciones del selector.
func SetTheme(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		theme := DefaultTheme
		if ck, err := c.Request().Cookie(ThemeCookie); err == nil && validTheme(ck.Value) {
			theme = ck.Value
		}
		c.Set("theme", theme)
		c.Set("themes", Themes)
		return next(c)
	}
}

// ThemeUpdate guarda el color de fondo elegido en una cookie.
func ThemeUpdate(c buffalo.Context) error {
	theme := c.Param("theme")
	if !validTheme(theme) {
		return c.Error(http.StatusBadRequest, fmt.Errorf("tema no válido: %q", theme))
	}
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     ThemeCookie,
		Value:    theme,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return c.Redirect(http.StatusSeeOther, "/")
}
