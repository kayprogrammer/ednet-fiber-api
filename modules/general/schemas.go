package general

import (
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

type SiteDetailSchema struct {
	Name 		string 		`json:"name" example:"SocialNet"`
	Email 		string 		`json:"email" example:"johndoe@email.com"`
	Phone 		string		`json:"phone" example:"+2348133831036"`
	Address 	string		`json:"address" example:"234, Lagos, Nigeria"`
	Fb 			string		`json:"fb" example:"https://facebook.com"`
	Tw 			string		`json:"tw" example:"https://twitter.com"`
	Wh 			string		`json:"wh" example:"https://wa.me/2348133831036"`
	Ig 			string		`json:"ig" example:"https://instagram.com"`
}

func (s SiteDetailSchema) Init(obj *ent.SiteDetail) SiteDetailSchema {
	s.Name = obj.Name
	s.Email = obj.Email
	s.Phone = obj.Phone
	s.Address = obj.Address
	s.Fb = obj.Fb
	s.Tw = obj.Tw
	s.Wh = obj.Wh
	s.Ig = obj.Ig
	return s
}

type SiteDetailResponseSchema struct {
	base.ResponseSchema
	Data			SiteDetailSchema		`json:"data"`
}