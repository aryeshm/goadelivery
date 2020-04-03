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

var items = []string{
	"rice",
	"wheat",
	"dal",
	"sugar",
	"garam masala",
	"chilli",
}

var units = []string{
	"gms",
	"kgs",
	"pcs",
	"pkg",
}

type vendor struct {
	id       int
	locality string
	name     string
}

type Item struct {
	ID    string `bson:"_id" json:"_id"`
	Name  string `bson:"name" json:"name"`
	Qty   int    `bson:"qty" json:"qty"`
	Units string `bson:"units" json:"units"`
}

func BuildVendorWindow(s gwu.Session, db *mongo.Database, ctx context.Context) {

	vend := vendor{}
	win := gwu.NewWindow("vendor", "Vendor Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)
	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			log.Println("LOADING window:", e.Src().ID())
		case gwu.ETypeWinUnload:
			log.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)

	// Button which changes window content
	win.Add(gwu.NewLabel("Vendor Page"))

	p := gwu.NewPanel()

	row := gwu.NewVerticalPanel()
	l := gwu.NewHorizontalPanel()
	l.Add(gwu.NewLabel("Item"))
	l.Add(gwu.NewLabel("Qty"))
	l.Add(gwu.NewLabel("Units"))
	row.Add(l)
	p2 := gwu.NewHorizontalPanel()
	lb1 := gwu.NewListBox(items)

	tb := gwu.NewLabel("")
	p2.Add(lb1)
	tb1 := gwu.NewTextBox("")
	p2.Add(tb1)
	lb2 := gwu.NewListBox(units)

	p2.Add(lb2)
	row.Add(p2)
	row.Add(tb)
	p.Add(row)
	pp := gwu.NewPanel()
	bUpdate := gwu.NewButton("Update")
	bUpdate.AddEHandlerFunc(func(e gwu.Event) {
		it := Item{}
		if lb1.SelectedIdx() == -1 {
			it.Name = items[0]
		} else {
			it.Name = lb1.SelectedValue()
		}
		w, _ := strconv.Atoi(tb1.Text())
		it.Qty = w
		if lb2.SelectedIdx() == -1 {
			it.Units = units[0]
		} else {
			it.Units = lb2.SelectedValue()
		}
		updateStock(db, ctx, &it, pp)
		e.MarkDirty(pp)
	}, gwu.ETypeClick)
	p.Add(bUpdate)

	win.Add(p)
	displayStock(db, ctx, pp)

	win.Add(pp)

	s.AddWin(win)
	log.Println(vend.locality)

}

func displayStock(db *mongo.Database, ctx context.Context, v gwu.Panel) {

	items := []*Item{}
	collection := db.Collection("stock")
	curr, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for curr.Next(ctx) {
		item := Item{}
		err = curr.Decode(&item)
		//	log.Println("resutl", item)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, &item)
		//results[i] = &result

	}
	curr.Close(ctx)
	//v := gwu.NewVerticalPanel()
	v.Clear()
	h := gwu.NewTable()
	h.SetCellSpacing(10)
	//h.SetCellPadding(5)
	h.Add(gwu.NewLabel("Name"), 0, 0)
	h.Add(gwu.NewLabel("Qty"), 0, 1)
	h.Add(gwu.NewLabel("Units"), 0, 2)
	h.Add(gwu.NewLabel("Rs/Unit"), 0, 3)
	for i, item := range items {
		//	h := gwu.NewHorizontalPanel()
		h.Add(gwu.NewLabel(item.Name), i+1, 0)
		h.Add(gwu.NewLabel(strconv.Itoa(item.Qty)), i+1, 1)
		h.Add(gwu.NewLabel(item.Units), i+1, 2)
		//	v.Add(h)

	}
	v.Add(h)
}

func updateStock(db *mongo.Database, ctx context.Context, item *Item, p gwu.Panel) {

	collection := db.Collection("stock")
	i := Item{}

	filter := bson.D{{"name", item.Name}}
	collection.FindOne(ctx, filter).Decode(&i)
	if i.Name == "" {
		item.ID = primitive.NewObjectID().String()
		_, err := collection.InsertOne(ctx, item)
		if err != nil {
			log.Fatal(err)
		}
		//	log.Println("im here", item.Name)
	} else {
		log.Println("Name", i.Name)
		update := bson.D{{"$set", bson.D{{"qty", item.Qty}, {"units", item.Units}}}}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
	displayStock(db, ctx, p)
	//p := gwu.NewPanel()

}
