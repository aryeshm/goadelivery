// +build ignore

package windows

import (
	"context"
	"log"
	"strconv"

	"github.com/icza/gowut/gwu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var localities = []string{
	"None",
	"Shantabhan",
	"Kumbarvaddo",
	"Temba",
	"Vaddy",
	"Catmit Bhat",
	"Moloc",
	"Bandh",
	"Mesta Bhat",
	"Ganache Bhat",
	"Sale Bhat",
	"Bute Bhat",
	"Fonsa Bhat",
	"Voilem Bhat",
	"St.Caitan",
}

type order struct {
	ID       string       `bson:"_id" json:"_id"`
	Locality string       `bson:"locality" json"locality"`
	Name     string       `bson:"name" json:"name"`
	Address  string       `bson:"address" json:"address"`
	Phone    string       `bson:"phone" json:"phone"`
	Order    []*orderItem `bson:"order" json:"order"`
}

type orderItem struct {
	Name  string `bson:"name" json:"name"`
	Qty   int    `bson:"qty" json:"qty"`
	Units string `bson:"units" json:"units"`
}

func BuildUserWindow(s gwu.Session, db *mongo.Database, ctx context.Context) {
	var itemQty []gwu.TextBox
	var items []*Item

	win := gwu.NewWindow("user", "User Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	win.Add(gwu.NewLabel("User Page"))

	p1 := gwu.NewVerticalPanel()
	dpL := gwu.NewListBox(localities)
	p1.Add(dpL)
	lb1 := gwu.NewLabel("Name")
	p1.Add(lb1)
	tbName := gwu.NewTextBox("")
	p1.Add(tbName)
	lb2 := gwu.NewLabel("Address")
	p1.Add(lb2)
	tbAddress :=
		gwu.NewTextBox("")
	p1.Add(tbAddress)
	lb3 := gwu.NewLabel("Phone Number")
	p1.Add(lb3)
	tbPhone := gwu.NewTextBox("")
	p1.Add(tbPhone)
	win.Add(p1)

	p2 := gwu.NewVerticalPanel()

	items, itemQty = orderForm(db, ctx, p2)

	win.Add(p2)
	tbError := gwu.NewLabel("")

	b3 := gwu.NewButton("Place Order")
	b3.SetEnabled(false)

	//submit function
	b3.AddEHandlerFunc(func(e gwu.Event) {
		if tbName.Text() != "" && tbAddress.Text() != "" && tbPhone.Text() != "" {
			itQty := []int{}
			for _, q := range itemQty {
				w, _ := strconv.Atoi(q.Text())
				itQty = append(itQty, w)

			}
			name := tbName.Text()
			address := tbAddress.Text()
			phone := tbPhone.Text()
			locality := dpL.SelectedValue()
			submitForm(db, ctx, items, itQty, name, address, phone, locality)

			log.Println("qties", itQty)
			for i, item := range items {
				log.Println("ordered", item.Name, itemQty[i].Text())
			}

		} else {
			tbError.SetText("Check Error")
			e.MarkDirty(tbError)
		}
	}, gwu.ETypeClick)

	win.Add(b3)
	win.Add(tbError)

	dpL.AddEHandlerFunc(func(e gwu.Event) {
		b3.SetEnabled(true)
		e.MarkDirty(b3)
	}, gwu.ETypeChange)

	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			log.Println("LOADING window:", e.Src().ID())
			dpL.SetSelected(0, true)
			p2.Clear()
			tbAddress.SetText("")
			tbName.SetText("")
			tbPhone.SetText("")
			items, itemQty = orderForm(db, ctx, p2)
			e.MarkDirty(p2, tbAddress, tbName, tbPhone)
			e.NewSession()
		case gwu.ETypeWinUnload:
			log.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)

	s.AddWin(win)

}

func orderForm(db *mongo.Database, ctx context.Context, p gwu.Panel) ([]*Item, []gwu.TextBox) {
	items := []*Item{}
	itemQty := []gwu.TextBox{}
	collection := db.Collection("stock")
	curr, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for curr.Next(ctx) {
		item := Item{}
		err = curr.Decode(&item)
		//log.Println("resutl", item)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, &item)
		//results[i] = &result

	}
	curr.Close(ctx)

	for _, item := range items {
		h := gwu.NewHorizontalPanel()
		h.Add(gwu.NewLabel(item.Name))
		ttb := gwu.NewTextBox("0")
		ttb.AddEHandlerFunc(func(e gwu.Event) {

			log.Print("you changed", ttb.ID)
			e.MarkDirty(ttb)
		}, gwu.ETypeChange)
		h.Add(ttb)
		itemQty = append(itemQty, ttb)

		h.Add(gwu.NewLabel(item.Units))
		p.Add(h)

	}
	return items, itemQty
}

func submitForm(db *mongo.Database, ctx context.Context, items []*Item, itemQty []int, name string, address string, phone string,
	locality string) {

	list := []*orderItem{}

	for i, q := range itemQty {
		if q != 0 {
			item := orderItem{}
			item.Name = items[i].Name
			item.Qty = q
			item.Units = items[i].Units
			list = append(list, &item)
		}
	}

	if len(list) > 0 {
		ord := order{primitive.NewObjectID().String(), locality, name, address, phone, list}

		collection := db.Collection("orders")
		collection.InsertOne(ctx, ord)
		log.Println("Order Placed")
	}
}
