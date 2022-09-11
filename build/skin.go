package build

import (
	"encoding/json"
	"os"
	"reflect"
)

type SkinPath struct {
	Index  string `json:"index"`
	Post   string `json:"post"`
	About  string `json:"about"`
	Header string `json:"header"`
	Footer string `json:"footer"`
	Nav    string `json:"nav"`
}

type SkinInfo struct {
	Name    string   `json:"name"`
	Author  string   `json:"author"`
	Version string   `json:"version"`
	Summary string   `json:"summary"`
	Paths   SkinPath `json:"paths"`
}

type Skin struct {
	Info SkinInfo
}

func MakeSkin() Skin {
	return Skin{}
}

func (s *Skin) GetSkin() {
	skinpath, perr := os.Getwd()

	if perr != nil {
		panic(perr)
	}

	ctx, ferr := os.ReadFile(skinpath + "/skin/skin.json")

	if ferr != nil {
		panic(ferr)
	}

	var sinfo = SkinInfo{}
	jerr := json.Unmarshal(ctx, &sinfo)

	if jerr != nil {
		panic(jerr)
	}

	s.Info = sinfo

	var infos = SkinInfo{}

	ref := reflect.ValueOf(&sinfo.Paths).Elem()
	sref := reflect.ValueOf(&infos.Paths).Elem()

	for i := 0; i < ref.NumField(); i++ {
		pelmVal := ref.Field(i)
		pelmTyp := ref.Type().Field(i)

		srefval := sref.FieldByName(pelmTyp.Name)
		pval := skinpath + (pelmVal.Interface().(string))
		srefval.SetString(pval)
	}

	s.Info.Paths = infos.Paths

}
