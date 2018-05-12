package Controller

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
)

var UserController = Controller{ControllerName: "user"}

func (c *Controller) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, getTemplate(currentFunctionName(), c.ControllerName), "home", make(map[string]interface{}))
	}
}

func (c *Controller) ProcessLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		rpcData := map[string]interface{}{
			"email": r.PostForm.Get("email"),
			"password":  r.PostForm.Get("password"),
		}

		var result map[string]string
		rpcResult := ServiceCallData("Login", rpcData, LoadConfiguration().UserService.Port);
		if err := json.Unmarshal(rpcResult, &result); err != nil {
			println(err.Error())
			return
		}

		if result["sid"] != "" {
			expiration := time.Now().Add(30 * time.Minute)
			cookie := http.Cookie{Name: "sid", Value: result["sid"], Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			//TODO: Wrong username password error
			fmt.Fprintf(w, "Wrong pw");
			return;
		}
	}
}



