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
	ID        string       `bson:"_id" json:"_id"`
	Locality  string       `bson:"locality" json"locality"`
	Name      string       `bson:"name" json:"name"`
	Address   string       `bson:"address" json:"address"`
	Phone     string       `bson:"phone" json:"phone"`
	Order     []*orderItem `bson:"order" json:"order"`
	Total     float32      `bson:"total" json:"total"`
	Processed bool         `bson:"processed" json:"processed"`
}

type orderItem struct {
	Name  string `bson:"name" json:"name"`
	Qty   int    `bson:"qty" json:"qty"`
	Units string `bson:"units" json:"units"`
}

var totalBill float32

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
			b := submitForm(db, ctx, items, itQty, name, address, phone, locality)
			if b == true {
				log.Println("ordering success")
				tbError.SetText("Order Successful!!")
				e.MarkDirty(tbError)
			}

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
			tbError.SetText("")
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
		if (item.Qty / item.MinSize) > 10 {
			items = append(items, &item)
		}
		//results[i] = &result

	}
	curr.Close(ctx)
	h := gwu.NewTable()
	h.SetCellSpacing(10)
	//h.SetCellPadding(5)
	h.Add(gwu.NewLabel("Name"), 0, 0)
	h.Add(gwu.NewLabel("Min Size"), 0, 1)
	h.Add(gwu.NewLabel("Price"), 0, 2)
	h.Add(gwu.NewLabel("Unit"), 0, 3)
	h.Add(gwu.NewLabel("Limit"), 0, 4)
	h.Add(gwu.NewLabel("Qty"), 0, 5)
	mp := make(map[string]int)
	mp2 := make(map[string]int)
	mp3 := make(map[string]float32)
	mp4 := make(map[string]gwu.Label)
	totalBox := gwu.NewLabel("0")
	for i, item := range items {
		//	h := gwu.NewHorizontalPanel()
		h.Add(gwu.NewLabel(item.Name), i+1, 0)
		ttb := gwu.NewTextBox("0")
		subTotal := gwu.NewLabel("0")
		mp[strconv.Itoa(int(ttb.ID()))] = item.Max
		mp2[strconv.Itoa(int(ttb.ID()))] = item.Price
		mp3[strconv.Itoa(int(ttb.ID()))] = item.MinSize
		mp4[strconv.Itoa(int(ttb.ID()))] = subTotal

		ttb.AddEHandlerFunc(func(e gwu.Event) {
			w, _ := strconv.Atoi(ttb.Text())
			id := strconv.Itoa(int(ttb.ID()))
			if w > mp[id] {
				log.Println("error")
				ttb.SetText("0")
			} else {
				min := mp3[id]
				price := mp2[id]
				mp4[id].SetText(fmt.Sprintf("%f", float32(w)*min*float32(price)))
				//subTotal.SetText("afasfas")
				total := float32(0)
				for _, but := range itemQty {
					w, _ := strconv.Atoi(but.Text())
					id := strconv.Itoa(int(but.ID()))

					total = total + float32(w)*mp3[id]*float32(mp2[id])

				}
				totalBox.SetText(fmt.Sprintf("%f", total))
				totalBill = total
				log.Print("min, price, w", min, price, w)
			}

			e.MarkDirty(ttb, subTotal, totalBox)

		}, gwu.ETypeChange)
		h.Add(ttb, i+1, 5)

		h.Add(subTotal, i+1, 6)
		itemQty = append(itemQty, ttb)
		h.Add(gwu.NewLabel(fmt.Sprintf("%f", item.MinSize)), i+1, 1)
		h.Add(gwu.NewLabel(strconv.Itoa(item.Price)), i+1, 2)
		h.Add(gwu.NewLabel(strconv.Itoa(item.Max)), i+1, 4)
		h.Add(gwu.NewLabel(item.Units), i+1, 3)
	}
	h.Add(totalBox, len(itemQty)+1, 6)
	p.Add(h)
	return items, itemQty
}

func submitForm(db *mongo.Database, ctx context.Context, items []*Item, itemQty []int, name string, address string, phone string,
	locality string) bool {

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
	collection := db.Collection("orders")
	ab := order{}
	cd := order{}
	collection.FindOne(ctx, bson.D{{"address", address}}).Decode(&ab)
	collection.FindOne(ctx, bson.D{{"address", phone}}).Decode(&cd)
	status := false
	if ab.Name == "" && cd.Name == "" {
		if len(list) > 0 {

			ord := order{primitive.NewObjectID().String(), locality, name, address, phone, list, totalBill, false}

			collection.InsertOne(ctx, ord)
			log.Println("Order Placed")
			updateStockAfterProcess(ctx, db, list)
			status = true
		} else {
			status = false
		}

	}
	return status
}
