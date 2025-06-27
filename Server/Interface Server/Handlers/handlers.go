package handlers

import (
	database "FindMy/InterfaceServer/Database"
	model "FindMy/InterfaceServer/Model"
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	mail, ok1 := jsn["mail"].(string)     // users are usually all friends.
	pass, ok2 := jsn["password"].(string) // users are usually all friends.

	if !ok1 || !ok2 || len(pass) != 64 || len(mail) > 320 || !database.IsValidEmail(mail) || !database.IsSHA256Hash(pass) {
		return c.Send(model.ParseResponse(false, 16315, fiber.Map{"message": "invalid form"}))
	}

	match, id := database.CheckPassword(context.Background(), mail, pass)

	if !match {
		return c.Send(model.ParseResponse(false, 24315, fiber.Map{"message": "wrong password"}))
	}

	token, expires := database.GetToken(context.Background(), id)

	if len(token) != 64 { //there isn't a token. and expires is a null pointer
		ctoken := database.CreateToken(context.Background(), id, nil)
		if len(ctoken) != 64 {
			return c.Send(model.ParseResponse(false, 32415, fiber.Map{"message": "login error"}))
		}
		token = ctoken
	}

	return c.Send(model.ParseResponse(true, 33315, fiber.Map{
		"id":      id,
		"token":   token,
		"expires": expires,
	}))
}

func SetNickname(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	name, ok1 := jsn["name"].(string) // users are usually all friends.

	if !ok1 || len(name) > 64 {
		return c.Send(model.ParseResponse(false, 16315, fiber.Map{"message": "invalid form"}))
	}

	check := database.UpdateNickname(context.Background(), c.Locals("user.id").(int64), name)

	if !check {
		return c.Send(model.ParseResponse(false, 56315, fiber.Map{"message": "server fail"}))
	}

	return c.Send(model.ParseResponse(true, 59315, fiber.Map{"message": "success"}))
}

func GetNofitications(c *fiber.Ctx) error {
	var nofs []model.Notification = database.GetNofitications(context.Background(), c.Locals("user.id").(int64))

	return c.Send(model.ParseResponse(true, 65315, fiber.Map{"nofitications": nofs}))

}

func AuthCheckerMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/Login" || c.Path() == "/Signup" || c.Path() == "/verify" {
		return c.Next()
	}

	token := c.Cookies("token")

	if len(token) != 64 {
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 701, fiber.Map{"message": "invalid token"}))
	}

	user, expires := database.CheckToken(context.Background(), token) //c.context(): if the client has been disconnected, the database query will be cancelling. (it cant be because c.context() is fake http request, not websocket's context)

	if user == -1 {
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 727, fiber.Map{"message": "wrong token"}))
	}

	if expires != nil && expires.Before(time.Now()) {
		database.RemoveToken(context.Background(), token) // context.Background(): even if the client disconnects, it will removed
		c.ClearCookie("token")
		return c.Send(model.ParseResponse(false, 832, fiber.Map{"message": "expired token"}))
	}

	c.Locals("user.id", user)
	return c.Next()
}

var domainPrefix string = "http://127.0.0.1/cdn/"
var defaultProfile string = "http://127.0.0.1/cdn/defaultprofilephoto.jpg"

func GetPhoto(c *fiber.Ctx) error {
	var photoaddr string = database.GetPhoto(context.Background(), c.Locals("user.id").(int64))

	if len(photoaddr) <= 0 {
		return c.Redirect(defaultProfile)
	}

	return c.Redirect(domainPrefix + photoaddr)
}

func SetPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Send(model.ParseResponse(false, 11825, fiber.Map{"message": "file can't found."}))
	}

	if file.Size == 0 {
		//remove photo
		bool2 := database.SetPhoto(context.Background(), c.Locals("user.id").(int64), "")

		if !bool2 {
			return c.Send(model.ParseResponse(false, 124425, fiber.Map{"message": "database fail"}))
		}

		return c.Send(model.ParseResponse(true, 127425, fiber.Map{"message": "success!", "photo": defaultProfile}))

	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return c.Send(model.ParseResponse(false, 12325, fiber.Map{"message": "Only JPG, PNG and GIF files are allowed"}))
	}

	if !database.IsValidImage(file) {
		return c.Send(model.ParseResponse(false, 138827, fiber.Map{"message": "Invalid image!"}))
	}

	name, err := database.GenerateRandomFilename()

	if err != nil {
		return c.Send(model.ParseResponse(false, 12925, fiber.Map{"message": "server fail"}))
	}

	savePath := filepath.Join("./www/cdn", name+ext)

	if err2 := c.SaveFile(file, savePath); err2 != nil {
		return c.Send(model.ParseResponse(false, 13525, fiber.Map{"message": "Failed to save file"}))
	}

	bool := database.SetPhoto(context.Background(), c.Locals("user.id").(int64), name+ext)

	if !bool {
		return c.Send(model.ParseResponse(false, 12925, fiber.Map{"message": "database fail"}))
	}

	return c.Send(model.ParseResponse(true, 15125, fiber.Map{"message": "success!", "photo": domainPrefix + name + ext}))
}

func AddFriend(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	id1 := c.Locals("user.id").(int64)
	id2, ok1 := jsn["id"].(int64)

	if !ok1 {
		return c.Send(model.ParseResponse(false, 16315, fiber.Map{"message": "invalid form"}))
	}

	if database.IsSubset([]int64{id2}, database.GetFriends(context.Background(), id1)) {
		return c.Send(model.ParseResponse(false, 17315, fiber.Map{"message": "you are already friends"}))
	}

	/*
		45: add friend request code
		46: waiting request (request is sended)
		47: request denied
		48: request accepted
	*/

	notifications := database.GetNofitications(context.Background(), id2)

	for _, n := range notifications {
		if n.Type == 45 {
			if n.Params == id1 {
				return c.Send(model.ParseResponse(false, 18915, fiber.Map{"message": "you are already sended request"}))
			}
		}
	}

	database.AddNotification(context.Background(), id2, model.Notification{
		Type:   45, // add friend code.
		Params: id1,
	})

	database.AddNotification(context.Background(), id2, model.Notification{
		Type:   46, // waiting feedback
		Params: id2,
	})

	return c.Send(model.ParseResponse(true, 19315, fiber.Map{"message": "request sended."}))
}

func AcceptFriend(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	id1 := c.Locals("user.id").(int64)
	id2, ok1 := jsn["id"].(int64)

	if !ok1 {
		return c.Send(model.ParseResponse(false, 21315, fiber.Map{"message": "invalid form"}))
	}

	if database.IsSubset([]int64{id2}, database.GetFriends(context.Background(), id1)) {
		return c.Send(model.ParseResponse(false, 21715, fiber.Map{"message": "you are already friends"}))
	}
	notifications := database.GetNofitications(context.Background(), id1)

	var i bool = false
	for _, n := range notifications {
		if n.Type == 45 {
			if n.Params == id2 {
				i = true
				break
			}
		}
	}

	if !i {
		return c.Send(model.ParseResponse(false, 23215, fiber.Map{"message": "there isn't a request"}))
	}

	database.AddNotification(context.Background(), id1, model.Notification{
		Type:   48, // accepted
		Params: id2,
	})

	database.AddNotification(context.Background(), id2, model.Notification{
		Type:   48, // accepted
		Params: id1,
	})

	database.AddFriend(context.Background(), id1, id2)
	database.AddFriend(context.Background(), id2, id1)
	return c.Send(model.ParseResponse(true, 24715, fiber.Map{"message": "request accepted, congratulations"}))
}

func DenyFriend(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	id1 := c.Locals("user.id").(int64)
	id2, ok1 := jsn["id"].(int64)

	if !ok1 {
		return c.Send(model.ParseResponse(false, 25815, fiber.Map{"message": "invalid form"}))
	}

	if !database.IsSubset([]int64{id2}, database.GetFriends(context.Background(), id1)) {
		return c.Send(model.ParseResponse(false, 26215, fiber.Map{"message": "you aren't friends already"}))
	}
	notifications := database.GetNofitications(context.Background(), id1)

	var i bool = false
	for _, n := range notifications {
		if n.Type == 45 {
			if n.Params == id2 {
				i = true
				break
			}
		}
	}

	if !i {
		return c.Send(model.ParseResponse(false, 27715, fiber.Map{"message": "there isn't a request"}))
	}

	database.AddNotification(context.Background(), id1, model.Notification{
		Type:   47, // denied
		Params: id2,
	})

	database.AddNotification(context.Background(), id2, model.Notification{
		Type:   47, // denied
		Params: id1,
	})

	return c.Send(model.ParseResponse(true, 28915, fiber.Map{"message": "request denied succestfully"}))
}

func Signup(c *fiber.Ctx) error {
	jsn := model.ParseJson(c.Body())

	nickname, ok1 := jsn["nickname"].(string)
	mail, ok2 := jsn["mail"].(string)
	pass, ok3 := jsn["password"].(string)

	if !ok1 || !ok2 || !ok3 || len(nickname) > 64 || len(pass) != 64 || len(mail) > 320 || !database.IsValidEmail(mail) || !database.IsSHA256Hash(pass) {
		return c.Send(model.ParseResponse(false, 30015, fiber.Map{"message": "invalid form"}))
	}

	token := database.CreateSignup(context.Background(), mail, pass, nickname)

	if len(token) <= 0 {
		return c.Send(model.ParseResponse(false, 30615, fiber.Map{"message": "server fail"}))
	}

	check := SendMail(mail, nickname, token)
	if !check {
		return c.Send(model.ParseResponse(false, 31115, fiber.Map{"message": "server fail"}))
	}

	return c.Send(model.ParseResponse(true, 31815, fiber.Map{"message": "activation mail sended"}))
}

func Verify(c *fiber.Ctx) error {
	token := c.Query("token")

	if len(token) != 64 {
		return c.Send(model.ParseResponse(false, 32115, fiber.Map{"message": "invalid token"}))
	}

	mail, nickname, pass := database.GetSignup(context.Background(), token)

	if len(mail) <= 0 {
		return c.Send(model.ParseResponse(false, 32715, fiber.Map{"message": "wrong token"}))
	}

	check1 := database.RemoveSignup(context.Background(), token)
	authtoken := database.CreateUser(context.Background(), nickname, mail, pass)

	if !(check1 && len(authtoken) == 64) {
		return c.Send(model.ParseResponse(false, 33415, fiber.Map{"message": "server fail"}))
	}

	return c.Send(model.ParseResponse(true, 33715, fiber.Map{"message": "success!", "token": authtoken}))
}
