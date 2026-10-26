package actions

import (
	"net/url"
)

func (as *ActionSuite) getWithTheme(path, theme string) string {
	req := as.HTML("%s", path)
	req.Headers["Cookie"] = "todo_theme=" + theme
	res := req.Get()
	as.Equal(200, res.Code)
	return res.Body.String()
}

func (as *ActionSuite) Test_Theme_DefaultIsBlanco() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)

	body := res.Body.String()
	as.Contains(body, `data-theme="blanco"`)
	as.Contains(body, `id="theme-select"`)
	as.Contains(body, `<option value="oscuro"`)
	as.Contains(body, `<option value="blanco"`)
	as.Contains(body, `<option value="azul"`)
	as.Contains(body, `<option value="morado"`)
	as.Contains(body, `value="blanco" selected`)
}

func (as *ActionSuite) Test_Theme_FromCookie() {
	body := as.getWithTheme("/", "morado")
	as.Contains(body, `data-theme="morado"`)
	as.Contains(body, `value="morado" selected`)
}

func (as *ActionSuite) Test_Theme_InvalidCookieFallsBack() {
	body := as.getWithTheme("/", "rojo")
	as.Contains(body, `data-theme="blanco"`)
	as.NotContains(body, `data-theme="rojo"`)
}

func (as *ActionSuite) Test_Theme_AppliesOnEditPage() {
	task := as.taskForEditing(false)
	body := as.getWithTheme("/tasks/"+task.ID.String()+"/edit", "oscuro")
	as.Contains(body, `data-theme="oscuro"`)
}

func (as *ActionSuite) Test_ThemeUpdate_SetsCookie() {
	res := as.HTML("/theme").Post(url.Values{"theme": {"azul"}})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))

	cookie := res.Header().Get("Set-Cookie")
	as.Contains(cookie, "todo_theme=azul")
	as.Contains(cookie, "Path=/")
	as.Contains(cookie, "Max-Age=")
}

func (as *ActionSuite) Test_ThemeUpdate_InvalidValue() {
	for _, value := range []string{"rojo", ""} {
		res := as.HTML("/theme").Post(url.Values{"theme": {value}})
		as.Equal(400, res.Code)
		as.NotContains(res.Header().Get("Set-Cookie"), "todo_theme")
	}
}
