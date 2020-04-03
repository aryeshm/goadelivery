package windows

import (
	"context"
	"fmt"
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

	win := gwu.NewWindow("order", "Order Window")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	win.Add(gwu.NewLabel("Order Page"))

	p := gwu.NewVerticalPanel()
	lb1 := gwu.NewListBox(orderLocalities)
	lb1.SetSelected(0, true)
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

	if locality != "All" {
		curr, err = collection.Find(ctx, bson.D{{"locality", locality}})
	}
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
		if ord.Processed == false {
			orders = append(orders, &ord)
		}
		//results[i] = &result

	}
	curr.Close(ctx)

	for _, ord := range orders {
		h := gwu.NewVerticalPanel()
		h.Add(gwu.NewLabel(ord.Name))
		h.Add(gwu.NewLabel(ord.Locality))
		h.Add(gwu.NewLabel(ord.Address))
		h.Add(gwu.NewLabel(ord.Phone))
		l := gwu.NewTable()
		l.SetCellSpacing(10)
		l.Add(gwu.NewLabel("Item"), 0, 0)
		l.Add(gwu.NewLabel("Quantity"), 0, 1)
		l.Add(gwu.NewLabel("Units"), 0, 2)

		for ii, it := range ord.Order {

			l.Add(gwu.NewLabel(it.Name), ii+1, 0)
			l.Add(gwu.NewLabel(strconv.Itoa(it.Qty)), ii+1, 1)
			l.Add(gwu.NewLabel(it.Units), ii+1, 2)

		}
		l.Add(gwu.NewLabel("Total"), len(ord.Order)+2, 1)
		l.Add(gwu.NewLabel(fmt.Sprintf("%f", ord.Total)), len(ord.Order)+2, 2)
		p.Add(h)
		processBtn := gwu.NewButton("Process")
		mp := make(map[string]string)
		id := strconv.Itoa(int(processBtn.ID()))
		mp[id] = ord.ID

		processBtn.AddEHandlerFunc(func(e gwu.Event) {
			id := strconv.Itoa(int(processBtn.ID()))
			//idt, _ := primitive.ObjectIDFromHex(mp[id])
			filter := bson.D{{"_id", mp[id]}}

			update := bson.D{{"$set", bson.D{{"processed", true}}}}
			collection.UpdateOne(ctx, filter, update)
			ll := order{}
			collection.FindOne(ctx, filter).Decode(&ll)
			log.Println("found", mp[id], ll)
			//updateStockAfterProcess(ctx, db, ll.Order)

		}, gwu.ETypeClick)
		l.Add(processBtn, len(ord.Order)+3, 1)

		p.Add(l)
	}

}

func updateStockAfterProcess(ctx context.Context, db *mongo.Database, orders []*orderItem) {

	collection := db.Collection("stock")

	for _, ord := range orders {
		ll := Item{}
		filter := bson.D{{"name", ord.Name}}
		collection.FindOne(ctx, filter).Decode(&ll)
		newStock := ll.Qty - float32(ord.Qty)*ll.MinSize
		update := bson.D{{"$set", bson.D{{"qty", newStock}}}}
		collection.UpdateOne(ctx, filter, update)
	}

}
