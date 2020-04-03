// +build ignore

package windows

import (
	"context"
	"log"
	"strconv"

	"github.com/icza/gowut/gwu"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderLocalities = []string{
	"All",
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

func BuildOrderWindow(s gwu.Session, db *mongo.Database, ctx context.Context) {

	win := gwu.NewWindow("order", "ORder Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	win.Add(gwu.NewLabel("Order Page"))

	p := gwu.NewVerticalPanel()
	lb1 := gwu.NewListBox(orderLocalities)
	orderPan := gwu.NewPanel()
	lb1.AddEHandlerFunc(func(e gwu.Event) {
		locality := lb1.SelectedValue()
		orderPan.Clear()
		updateOrder(ctx, db, locality, orderPan)
		e.MarkDirty(orderPan)
	}, gwu.ETypeChange)

	p.Add(lb1)
	updateOrder(ctx, db, "All", orderPan)
	p.Add(orderPan)

	win.Add(p)
	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETypeWinLoad:
			log.Println("LOADING window:", e.Src().ID())
			orderPan.Clear()
			updateOrder(ctx, db, lb1.SelectedValue(), orderPan)
			e.MarkDirty(orderPan)
		case gwu.ETypeWinUnload:
			log.Println("UNLOADING window:", e.Src().ID())
		}
	}, gwu.ETypeWinLoad, gwu.ETypeWinUnload)
	s.AddWin(win)

}

func updateOrder(ctx context.Context, db *mongo.Database, locality string, p gwu.Panel) {
	collection := db.Collection("orders")
	curr, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	orders := []*order{}
	for curr.Next(ctx) {
		ord := order{}
		err = curr.Decode(&ord)
		//	log.Println("resutl", ord)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, &ord)
		//results[i] = &result

	}
	curr.Close(ctx)

	for _, ord := range orders {
		h := gwu.NewVerticalPanel()
		h.Add(gwu.NewLabel(ord.Name))
		h.Add(gwu.NewLabel(ord.Locality))
		h.Add(gwu.NewLabel(ord.Address))
		h.Add(gwu.NewLabel(ord.Phone))

		for _, it := range ord.Order {

			x := gwu.NewHorizontalPanel()
			x.Add(gwu.NewLabel(it.Name))
			x.Add(gwu.NewLabel(strconv.Itoa(it.Qty)))
			x.Add(gwu.NewLabel(it.Units))
			h.Add(x)
		}
		p.Add(h)
	}

}
