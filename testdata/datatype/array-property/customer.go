package arrayproperty

type longText string

type Customer struct {
	ID        int64 `sql:"id"`
	Age       uint8 `sql:"howOld"`
	Name      string
	Address   []string
	Nicknames []longText
	Status    longText `sql:"status"`
}
