package stores

import (
	"errors"
	"github.com/coopernurse/gorp"
	"time"
)

func CreateStore(db *gorp.DbMap, title, website string) (Store, error) {
	t := time.Now().UnixNano()
	var s = Store{
		Title:   title,
		Website: website,
		Created: t,
		Updated: t,
	}
	err := db.Insert(&s)
	return s, err
}

func CreateProduct(db *gorp.DbMap, storeid, ownproductid, imgid int64,
	price float64, title string) (StoreProduct, error) {
	t := time.Now().UnixNano()
	b := StoreProduct{
		StoreId:   storeid,
		ProductId: ownproductid,
		Price:     price,
		ImgId:     imgid,
		Title:     title,
		Created:   t,
		Updated:   t,
	}
	err := db.Insert(&b)
	return b, err
}

func (p *StoreProduct) Update(db *gorp.DbMap) error {
	p.Updated = time.Now().UnixNano()
	_, err := db.Update(p)
	return err
}

func BasketAdd(db *gorp.DbMap, userid, storeid, productid, count int64) error {
	t := time.Now().UnixNano()
	b := StoreBasket{
		UserId:    userid,
		StoreId:   storeid,
		ProductId: productid,
		Count:     count,
		Created:   t,
		Updated:   t,
	}

	res, err := db.Exec("update StoreBasket set Count = Count + ?, Updated = ?"+
		" where UserId = ? and StoreId=? and ProductId=?",
		count, t, userid, storeid, productid)
	if err != nil {
		return err
	}
	if num, err := res.RowsAffected(); err == nil && num == 0 {
		err := db.Insert(&b)
		return err
	}
	return err
}

func BasketRemove(db *gorp.DbMap, userid, storeid, productid int64) error {
	_, err := db.Exec("delete from StoreBasket where UserId = ? and "+
		"StoreId = ? and ProductId=?", userid, storeid, productid)
	return err
}

func BasketClean(db *gorp.DbMap, userid, storeid int64) error {
	_, err := db.Exec("delete from StoreBasket where UserId = ? and "+
		"StoreId = ?", userid, storeid)
	return err
}

func BasketGet(db *gorp.DbMap, userid, storeid int64) ([]StoreBasket, error) {
	var b = []StoreBasket{}
	_, err := db.Select(&b, "select * from StoreBasket where "+
		"UserId = ? and StoreId = ?", userid, storeid)
	return b, err
}

func GetProduct(db *gorp.DbMap, productid int64) (StoreProduct, error) {
	var sp = StoreProduct{}
	obj, err := db.Get(StoreProduct{}, productid)
	if err != nil {
		return sp, err
	}
	if obj == nil {
		return sp, errors.New("Product not found")
	}
	sp = *obj.(*StoreProduct)
	return sp, nil
}
