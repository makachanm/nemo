package build

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"nemo/nemomark"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const Spliter = "==========\n"

type Builder struct {
	Manifest Manifest
	Skin     Skin
	PostList []DocumentMeta
}

func MakeNewBuilder() Builder {
	mfest := GetManifest()
	return Builder{Skin: MakeSkin(), Manifest: mfest}
}

func (b *Builder) buildPage(postpath string) (string, DocumentMeta, bool) {
	markup := nemomark.MakeNemomark()

	ctx, perr := os.ReadFile(postpath)

	if perr != nil {
		panic(perr)
	}

	postRawctx := string(ctx)

	if !strings.ContainsAny(postRawctx, Spliter) {
		return "", DocumentMeta{}, false
	}
	postCtx := strings.Split(postRawctx, Spliter)

	document := NewDocument()

	document.Content = markup.Mark(strings.Join(postCtx[1:], ""))
	document.Content = strings.ReplaceAll(document.Content, "\n", "\n<br />")

	document.ParseMeta(postCtx[0])

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: true, BlogName: bname, PostName: document.Meta.Title}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: true, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	document.Foot = foot

	navd := Nav{IsNotIndex: true, BlogName: bname}
	nav := BuildNav(b.Skin, navd)
	document.Nav = nav

	file, fserr := os.ReadFile(b.Skin.Info.Paths.Post)

	if fserr != nil {
		panic(fserr)
	}

	var builder template.Template
	t, err := builder.Parse(string(file))

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, document)

	return writer.String(), document.Meta, true

}

func (b *Builder) buildIndex() string {
	indexs := NewIndexData()
	indexs.Indexs = b.PostList

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: false, BlogName: bname}
	head := BuildHeader(b.Skin, headd)
	indexs.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: false, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	indexs.Foot = foot

	navd := Nav{IsNotIndex: false, BlogName: bname}
	nav := BuildNav(b.Skin, navd)
	indexs.Nav = nav

	file, fserr := os.ReadFile(b.Skin.Info.Paths.Index)

	if fserr != nil {
		panic(fserr)
	}

	var builder template.Template
	t, err := builder.Parse(string(file))

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, indexs)

	return writer.String()
}

func (b *Builder) buildAboutPage() string {
	wd, _ := os.Getwd()
	markup := nemomark.MakeNemomark()

	ctx, perr := os.ReadFile(wd + "/post/about.ps")

	if perr != nil {
		panic(perr)
	}

	postRawctx := string(ctx)

	document := NewAboutPage()

	document.Content = markup.Mark(postRawctx)

	document.Content = strings.ReplaceAll(document.Content, "\n", "\n<br />")

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: false, BlogName: bname}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: false, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	document.Foot = foot

	navd := Nav{IsNotIndex: false, BlogName: bname}
	nav := BuildNav(b.Skin, navd)
	document.Nav = nav

	document.BuildInfo = "" //WIP

	file, fserr := os.ReadFile(b.Skin.Info.Paths.About)

	if fserr != nil {
		panic(fserr)
	}

	var builder template.Template
	t, err := builder.Parse(string(file))

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, document)

	return writer.String()

}

func (b *Builder) packRes() {

	wd, werr := os.Getwd()
	if werr != nil {
		panic(werr)
	}

	_, ex := os.Stat("dist/res")
	if os.IsNotExist(ex) {
		_ = os.Chdir("dist")
		_ = os.Mkdir("res", os.ModePerm)
		_ = os.Chdir("res")
		_ = os.Mkdir("skin", os.ModePerm)
		_ = os.Chdir("../..")
	}

	_, roex := os.Stat("skin/static")
	if os.IsNotExist(roex) {
		fmt.Println("skin static resource not found.")
		return
	}

	_, svex := os.Stat("post/res")
	if os.IsNotExist(svex) {
		fmt.Println("post/res folder not found.")
		return
	}

	skinsrc := "./skin/static"
	skindst := "./dist/res/skin"

	_ = os.Chdir(wd)

	skinfiles, _ := os.ReadDir(skinsrc)
	if len(skinfiles) == 0 {
		panic("skin/static resource not found.")
	}
	for _, file := range skinfiles {
		skinfile, status := os.ReadFile(skinsrc + "/" + file.Name())
		if status != nil {
			panic(status)
		}
		cerr := os.WriteFile(skindst+"/"+file.Name(), skinfile, os.ModePerm)
		if cerr != nil {
			panic(cerr)
		}
	}
	check, _ := os.ReadDir(skindst)
	if len(check) == 0 {
		panic("dist/res/skin resource not found.")
	}

	resrc := "./post/res"
	resdst := "./dist/res"

	_ = os.Chdir(wd)

	resfiles, _ := os.ReadDir(resrc)
	for _, file := range resfiles {
		resfile, status := os.ReadFile(resrc + "/" + file.Name())
		if status != nil {
			panic(status)
		}
		rerr := os.WriteFile(resdst+"/"+file.Name(), resfile, os.ModePerm)
		if rerr != nil {
			panic(rerr)
		}
	}
	if len(check) == 0 {
		panic("dist/res folder not found.")
	}

}

func makeFileName(title string, smp TimeStamp) string {
	timesp := smp.StampSize()

	var fileTitle string

	if len(title) > 10 {
		fileTitle = title[:10]
	} else {
		fileTitle = title
	}

	fileTitle = base64.StdEncoding.EncodeToString([]byte(fileTitle))

	fname := strconv.Itoa(timesp) + "-" + fileTitle + ".html"
	return fname
}

func (b *Builder) Build() {
	b.PostList = make([]DocumentMeta, 0)

	b.Skin.GetSkin()

	wd, derr := os.Getwd()
	workd := wd + "/post/"

	if derr != nil {
		panic(derr)
	}

	dir, rderr := os.ReadDir(workd)

	if rderr != nil {
		panic(rderr)
	}

	var BuildTargets = make([]string, 0)

	for _, ctx := range dir {
		name := ctx.Name()
		if strings.ContainsAny(name, ".ps") && (name != "about.ps") && (!ctx.Type().IsDir()) {
			BuildTargets = append(BuildTargets, workd+name)
		}
	}

	fmt.Print("Building...\n")

	_, exerr := os.Stat("dist")
	if os.IsNotExist(exerr) {
		_ = os.Mkdir("dist", os.ModePerm)
		_ = os.Chdir("dist")
		_ = os.Mkdir("page", os.ModePerm)
		_ = os.Chdir("..")
	}

	for _, ctx := range BuildTargets {
		content, dmeta, iscom := b.buildPage(ctx)
		if !iscom {
			continue
		}

		name := makeFileName(dmeta.Title, dmeta.Timestamp)
		fdir := wd + "/dist/page/" + name

		fmt.Println(" => ", name)

		_ = os.WriteFile(fdir, []byte(content), 0777)

		dmeta.Path = name
		b.PostList = append(b.PostList, dmeta)
	}

	sort.Slice(b.PostList, func(i, j int) bool {
		return b.PostList[i].Timestamp.StampSize() < b.PostList[j].Timestamp.StampSize()
	})

	indext := b.buildIndex()
	indexpath := wd + "/dist/index.html"

	_ = os.WriteFile(indexpath, []byte(indext), 0777)

	_, ex := os.Stat(wd + "/post/about.ps")
	if !os.IsNotExist(ex) {
		aboutCtx := b.buildAboutPage()
		aboutPath := wd + "/dist/about.html"
		_ = os.WriteFile(aboutPath, []byte(aboutCtx), 0777)
	}

	b.packRes()
}
