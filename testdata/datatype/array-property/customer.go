package arrayproperty

type longText string

type Customer struct {
	ID        int64 `sql:"id"`
	Name      string
	Address   []string
	Nicknames []longText
	Status    longText
}
