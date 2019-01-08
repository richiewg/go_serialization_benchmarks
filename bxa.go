package goserbench

import (
	"bytes"
)

type BXAA struct {
	Name     string
	BirthDay uint64
	Phone    string
	Siblings uint32
	Spouse   bool
	Money    uint64
}

func (o *BXAA) MarshalTo() ([]byte,error) {
	temp := bytes.NewBuffer(nil)
	err:=WriteString(temp,o.Name)
	if err != nil {
		return nil,err
	}
	err= WriteUint64(temp,o.BirthDay)
	if err != nil {
		return nil,err
	}
	err= WriteString(temp,o.Phone)
	if err != nil {
		return nil,err
	}
	err= WriteUint32(temp,o.Siblings)
	if err != nil {
		return nil,err
	}
	err=WriteBool(temp,o.Spouse)
	if err != nil {
		return nil,err
	}
	err= WriteUint64(temp,o.Money)
	if err != nil {
		return nil,err
	}
	return temp.Bytes(),err
}

func (o *BXAA) UnMarshal(data []byte) error{
	r:=bytes.NewBuffer(data)

	name, err := ReadString(r)
	if err != nil {
		return err
	}
	birthday, err := ReadUint64(r)
	if err != nil {
		return err
	}
	phone, err := ReadString(r)
	if err != nil {
		return err
	}
	siblings, err := ReadUint32(r)
	if err != nil {
		return err
	}
	spouse, err := ReadBool(r)
	if err != nil {
		return err
	}
	soney, err := ReadUint64(r)
	if err != nil {
		return err
	}
	o.Name= name
	o.BirthDay= birthday
	o.Phone= phone
	o.Siblings= siblings
	o.Spouse= spouse
	o.Money= soney
	return nil
}
