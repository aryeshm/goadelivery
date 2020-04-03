package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	g "github.com/aryesh/GoaDelivery/context"
	w "github.com/aryesh/GoaDelivery/windows"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/icza/gowut/gwu"
)

type myButtonHandler struct {
	counter int
	text    string
}

/*
func (h *myButtonHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		w.BuildUserWindow(e)
		e.ReloadWin("user")
	}
}
*/

//build vendor options

func buildVendorSite(s gwu.Session) {
	log.Println("Priate", s.Private())

	win := gwu.NewWindow("vend1", "Vendor Window")

	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	win.Add(gwu.NewLabel("Vendor Page"))

	p := gwu.NewPanel()

	b1 := gwu.NewButton("Update Stock")

	b1.AddEHandlerFunc(func(e gwu.Event) {
		w.BuildVendorWindow(s, coll, ctx)
		e.ReloadWin("vendor")
	}, gwu.ETypeClick)
	p.Add(b1)
	b2 := gwu.NewButton("Orders")

	b2.AddEHandlerFunc(func(e gwu.Event) {
		w.BuildOrderWindow(s, coll, ctx)
		e.ReloadWin("order")
	}, gwu.ETypeClick)
	p.Add(b2)
	win.Add(p)

	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			log.Println("LOADING window:", e.Src().ID())
			//e.NewSession()
		case gwu.ETypeWinUnload:
			log.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)
	s.AddWin(win)
}

func buildLoginWin(s gwu.Session) {
	//s := e.Session()
	//e.NewSession()
	win := gwu.NewWindow("login", "Login Window")
	win.Style().SetFullSize()
	win.SetAlign(gwu.HACenter, gwu.VAMiddle)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)

	l := gwu.NewLabel("Test GUI Login Window")
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("150%")
	p.Add(l)
	l = gwu.NewLabel("Login")
	l.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("130%")
	p.Add(l)
	p.CellFmt(l).Style().SetBorder2(1, gwu.BrdStyleDashed, gwu.ClrNavy)
	l = gwu.NewLabel("user/pass: admin/a")
	l.Style().SetFontSize("80%").SetFontStyle(gwu.FontStyleItalic)
	p.Add(l)

	errL := gwu.NewLabel("")
	errL.Style().SetColor(gwu.ClrRed)
	p.Add(errL)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(2, 2)
	table.Add(gwu.NewLabel("User name:"), 0, 0)
	tb := gwu.NewTextBox("")
	tb.Style().SetWidthPx(160)
	table.Add(tb, 0, 1)
	table.Add(gwu.NewLabel("Password:"), 1, 0)
	pb := gwu.NewPasswBox("")
	pb.Style().SetWidthPx(160)
	table.Add(pb, 1, 1)
	p.Add(table)
	b := gwu.NewButton("OK")
	b.AddEHandlerFunc(func(e gwu.Event) {
		log.Println("iahere")
		if tb.Text() == "admin" && pb.Text() == "a" {
			e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser
			//	buildPrivateWins(e.Session())
			//	e.ReloadWin("main")
			buildVendorSite(e.Session()) //, coll, ctx)

			e.ReloadWin("vend1")
		} else {
			e.SetFocusedComp(tb)
			errL.SetText("Invalid user name or password!")
			e.MarkDirty(errL)
		}
	}, gwu.ETypeClick)
	p.Add(b)
	l = gwu.NewLabel("")
	p.Add(l)
	p.CellFmt(l).Style().SetHeightPx(200)

	win.Add(p)
	win.SetFocusedCompID(tb.ID())

	s.AddWin(win)
}

type sessHandler struct{}

func (h sessHandler) Created(s gwu.Session) {
	fmt.Println("SESSION created:", s.ID())
	//log.Println("Priavtee", s.Private())
	buildLoginWin(s)
}

func (h sessHandler) Removed(s gwu.Session) {
	fmt.Println("SESSION removed:", s.ID())
}

var ctx context.Context
var coll *mongo.Database

func main() {
	port := os.Getenv("PORT")
	config, err := g.LoadConfigMgo("./")

	if err != nil {
		log.Fatalf("error")
	}

	ctx = context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, g.HostKey, config.Db.ConnectionString)
	ctx = context.WithValue(ctx, g.UsernameKey, config.Db.User)
	ctx = context.WithValue(ctx, g.PasswordKey, config.Db.Password)
	//	ctx = context.WithValue(ctx, databaseKey, config.Db.ConnectionString))

	coll, _ = g.OpenMgDb(ctx, config)

	//fmt.Println(coll.Name())

	win := gwu.NewWindow("main", "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	win.Add(gwu.NewLabel("User"))
	btn := gwu.NewButton("Click me")
	//btn.AddEHandler(&myButtonHandler{text: ":-)"}, gwu.ETypeClick)
	btn.AddEHandlerFunc(func(e gwu.Event) {
		w.BuildUserWindow(e.Session(), coll, ctx)
		e.ReloadWin("user")
	}, gwu.ETypeClick)

	win.Add(btn)

	win.Add(gwu.NewLabel("Vendor"))
	btn2 := gwu.NewButton("Click me")
	//btn.AddEHandler(&myButtonHandler{text: ":-)"}, gwu.ETypeClick)
	btn2.AddEHandlerFunc(func(e gwu.Event) {
		//	buildLoginWin(e.Session()) //, coll, ctx)
		e.ReloadWin("login")
	}, gwu.ETypeClick)
	win.Add(btn2)
	folder := "test_tls/"
	//server := gwu.NewServerTLS("guitest", "localhost:8081", folder+"cert.pem", folder+"key.pem")
	server := gwu.NewServerTLS("guitest", ":"+port, folder+"cert.pem", folder+"key.pem")

	//server := gwu.NewServerTLS("guitest", "localhost:8081", "", "")

	//server := gwu.NewServer("guitest", "localhost:8081")
	server.SetText("Test GUI App")

	server.AddSessCreatorName("login", "Login Window")
	server.AddSHandler(sessHandler{})

	server.AddWin(win)
	server.Start("")

}
