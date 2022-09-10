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
	Skin     Skin
	PostList []Document_Meta
}

func MakeNewBuilder() Builder {
	return Builder{Skin: MakeSkin()}
}

func (b *Builder) build_page(postpath string) (string, Document_Meta, bool) {
	markup := nemomark.MakeNemomark()

	ctx, perr := os.ReadFile(postpath)

	if perr != nil {
		panic(perr)
	}

	post_rawctx := string(ctx)

	if !strings.ContainsAny(post_rawctx, Spliter) {
		return "", Document_Meta{}, false
	}
	post_ctx := strings.Split(post_rawctx, Spliter)

	document := NewDocument()

	document.Content = markup.Mark(post_ctx[1])
	document.Content = strings.ReplaceAll(document.Content, "\n", "\n<br />")

	document.ParseMeta(post_ctx[0])

	headd := Header{IsNotIndex: true}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: true, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	document.Foot = foot

	navd := Nav{IsNotIndex: true}
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
	t.Execute(&writer, document)

	return writer.String(), document.Meta, true

}

func (b *Builder) build_index() string {
	indexs := NewIndexData()
	indexs.Indexs = b.PostList

	headd := Header{IsNotIndex: false}
	head := BuildHeader(b.Skin, headd)
	indexs.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: false, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	indexs.Foot = foot

	navd := Nav{IsNotIndex: false}
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
	t.Execute(&writer, indexs)

	return writer.String()
}

func (b *Builder) build_about_page() string {
	wd, _ := os.Getwd()
	markup := nemomark.MakeNemomark()

	ctx, perr := os.ReadFile(wd + "/post/about.ps")

	if perr != nil {
		panic(perr)
	}

	post_rawctx := string(ctx)

	document := NewDocument()

	document.Content = markup.Mark(post_rawctx)

	document.Content = strings.ReplaceAll(document.Content, "\n", "\n<br />")

	headd := Header{IsNotIndex: false}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	vinfo := "Skin with - " + b.Skin.Info.Name + " / Made by - " + b.Skin.Info.Author + " / " + b.Skin.Info.Summary
	footd := Footer{IsNotIndex: false, VInfo: vinfo}
	foot := BuildFooter(b.Skin, footd)
	document.Foot = foot

	navd := Nav{IsNotIndex: false}
	nav := BuildNav(b.Skin, navd)
	document.Nav = nav

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
	t.Execute(&writer, document)

	return writer.String()

}

func (b *Builder) pack_res() {

	wd, werr := os.Getwd()
	if werr != nil {
		panic(werr)
	}

	_, ex := os.Stat("dist/res")
	if os.IsNotExist(ex) {
		os.Chdir("dist")
		os.Mkdir("res", os.ModePerm)
		os.Chdir("res")
		os.Mkdir("skin", os.ModePerm)
		os.Chdir("../..")
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

	skinsrc := wd + "/skin/static"
	skindet := wd + "/dist/res/skin"

	cerr := DirCopy(skinsrc, skindet)
	if cerr != nil {
		panic(cerr)
	}

}

func make_file_name(title string, smp TimeStamp) string {
	timesp := smp.StampSize()

	var file_title string

	if len(title) > 10 {
		file_title = title[:10]
	} else {
		file_title = title
	}

	file_title = base64.StdEncoding.EncodeToString([]byte(file_title))

	fname := strconv.Itoa(int(timesp)) + "-" + file_title + ".html"
	return fname
}

func (b *Builder) Build() {
	b.PostList = make([]Document_Meta, 0)

	b.Skin.Get_skin()

	wd, derr := os.Getwd()
	workd := wd + "/post/"

	if derr != nil {
		panic(derr)
	}

	dir, rderr := os.ReadDir(workd)

	if rderr != nil {
		panic(rderr)
	}

	var BuildTargets []string = make([]string, 0)

	for _, ctx := range dir {
		name := ctx.Name()
		if strings.ContainsAny(name, ".ps") && (name != "about.ps") && (!ctx.Type().IsDir()) {
			BuildTargets = append(BuildTargets, (workd + name))
		}
	}

	fmt.Print("Building...\n")

	_, exerr := os.Stat("dist")
	if os.IsNotExist(exerr) {
		os.Mkdir("dist", os.ModePerm)
		os.Chdir("dist")
		os.Mkdir("page", os.ModePerm)
		os.Chdir("..")
	}

	for _, ctx := range BuildTargets {
		content, dmeta, iscom := b.build_page(ctx)
		if !iscom {
			continue
		}

		name := make_file_name(dmeta.Title, dmeta.Timestamp)
		fdir := wd + "/dist/page/" + name

		fmt.Println(" => ", name)

		os.WriteFile(fdir, []byte(content), 0777)

		dmeta.Path = name
		b.PostList = append(b.PostList, dmeta)
	}

	sort.Slice(b.PostList, func(i, j int) bool {
		return b.PostList[i].Timestamp.StampSize() < b.PostList[j].Timestamp.StampSize()
	})

	indext := b.build_index()
	indexpath := wd + "/dist/index.html"

	os.WriteFile(indexpath, []byte(indext), 0777)

	_, ex := os.Stat(wd + "/post/about.ps")
	if !os.IsNotExist(ex) {
		about_ctx := b.build_about_page()
		about_path := wd + "/dist/about.html"
		os.WriteFile(about_path, []byte(about_ctx), 0777)
	}

	b.pack_res()
}
