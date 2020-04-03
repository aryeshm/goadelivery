package windows

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/icza/gowut/gwu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var items = []string{
	"Beetroot",
	"Bhendi",
	"Brinjal",
	"Gallan",
	"Cabbage",
	"Carrot",
	"Cauliflower",
	"Clusterbeans",
	"Cucumber",
	"French beans",
	"Garlic",
	"Ginger",
	"Green chilli",
	"Limes",
	"Shimla mirchi",
	"Bitter gourd",
	"Green peas",
	"Knolkhol",
	"Long beans",
	"Corriander",
	"Cnions",
	"Potato",
	"Tomato",
	"Watermelon",
	"Green grapes",
	"Black grapes",
	"Oranges",
	"Banana mandoli",
	"Banana velchi",
}

var units = []string{

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
	ID      string  `bson:"_id" json:"_id"`
	Name    string  `bson:"name" json:"name"`
	Qty     float32 `bson:"qty" json:"qty"`
	Units   string  `bson:"units" json:"units"`
	MinSize float32 `bson:"minSize" json:"minSize"`
	Price   int     `bson:"price" json:"price"`
	Max     int     `bson:"max" json:"max"`
}

func BuildVendorWindow(s gwu.Session, db *mongo.Database, ctx context.Context) {

	vend := vendor{}
	win := gwu.NewWindow("vendor", "Vendor Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	win.Add(gwu.NewLabel("Vendor Page"))

	p := gwu.NewPanel()

	//row := gwu.NewVerticalPanel()
	//l := gwu.NewHorizontalPanel()
	l := gwu.NewTable()
	l.SetCellSpacing(10)
	l.Add(gwu.NewLabel("Item"), 0, 0)
	l.Add(gwu.NewLabel("Net Qty"), 0, 1)
	l.Add(gwu.NewLabel("Units"), 0, 2)
	l.Add(gwu.NewLabel("Min Size"), 0, 3)
	l.Add(gwu.NewLabel("Price/Unit"), 0, 4)
	l.Add(gwu.NewLabel("Max Nos."), 0, 5)
	//row.Add(l)
	//p2 := gwu.NewHorizontalPanel()
	lb1 := gwu.NewListBox(items)
	l.Add(lb1, 1, 0)
	//tb := gwu.NewLabel("")
	//l.Add(tb, 1, 1)
	//p2.Add(lb1)
	tb1 := gwu.NewTextBox("")
	l.Add(tb1, 1, 1)
	//p2.Add(tb1)
	lb2 := gwu.NewListBox(units)
	l.Add(lb2, 1, 2)
	tbMinSize := gwu.NewTextBox("")
	l.Add(tbMinSize, 1, 3)
	tbPrice := gwu.NewTextBox("")
	l.Add(tbPrice, 1, 4)
	tbMax := gwu.NewTextBox("")
	l.Add(tbMax, 1, 5)
	//p2.Add(lb2)
	//row.Add(p2)
	//row.Add(tb)
	p.Add(l)
	pp := gwu.NewPanel()
	bUpdate := gwu.NewButton("Update")

	bUpdate.AddEHandlerFunc(func(e gwu.Event) {
		it := Item{}
		if lb1.SelectedIdx() == -1 {
			it.Name = items[0]
		} else {
			it.Name = lb1.SelectedValue()
		}
		f, _ := strconv.ParseFloat(tb1.Text(), 32)
		it.Qty = float32(f)
		if lb2.SelectedIdx() == -1 {
			it.Units = units[0]
		} else {
			it.Units = lb2.SelectedValue()
		}
		w, _ := strconv.Atoi(tbPrice.Text())
		it.Price = w
		f, _ = strconv.ParseFloat(tbMinSize.Text(), 32)
		it.MinSize = float32(f)
		w, _ = strconv.Atoi(tbMax.Text())
		it.Max = w
		updateStock(db, ctx, &it, pp)
		e.MarkDirty(pp)
	}, gwu.ETypeClick)
	l.Add(bUpdate, 2, 2)

	win.Add(p)
	displayStock(db, ctx, pp)

	win.Add(pp)

	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			log.Println("LOADING window:", e.Src().ID())
			pp.Clear()
			displayStock(db, ctx, pp)
		case gwu.ETypeWinUnload:
			log.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)

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
	h.Add(gwu.NewLabel("Min Size"), 0, 3)
	h.Add(gwu.NewLabel("Rs/Unit"), 0, 4)
	h.Add(gwu.NewLabel("Max Qty"), 0, 5)
	for i, item := range items {
		//	h := gwu.NewHorizontalPanel()
		h.Add(gwu.NewLabel(item.Name), i+1, 0)
		h.Add(gwu.NewLabel(fmt.Sprintf("%f", item.Qty)), i+1, 1)
		h.Add(gwu.NewLabel(item.Units), i+1, 2)
		h.Add(gwu.NewLabel(fmt.Sprintf("%f", item.MinSize)), i+1, 3)
		h.Add(gwu.NewLabel(strconv.Itoa(item.Price)), i+1, 4)
		h.Add(gwu.NewLabel(strconv.Itoa(item.Max)), i+1, 5)
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
		update := bson.D{{"$set", bson.D{{"qty", item.Qty}, {"units", item.Units},
			{"price", item.Price}, {"minSize", item.MinSize}, {"max", item.Max}}}}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
	displayStock(db, ctx, p)
	//p := gwu.NewPanel()

}
