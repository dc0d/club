package rammailstorage

type RAMStorage struct {
	mails []interface{}
}

func (rm *RAMStorage) Len() int              { return len(rm.mails) }
func (rm *RAMStorage) PeekHead() interface{} { return rm.mails[0] }
func (rm *RAMStorage) DropHead()             { rm.mails = rm.mails[1:] }
func (rm *RAMStorage) Append(v interface{})  { rm.mails = append(rm.mails, v) }

func New() *RAMStorage { return &RAMStorage{} }
