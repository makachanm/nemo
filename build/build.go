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
	Manifest     Manifest
	Skin         Skin
	PostList     []DocumentMeta
	TagList      map[string][]DocumentMeta
	IndexPageNum int
	wd           string
}

func MakeNewBuilder() Builder {
	mfest := GetManifest()
	return Builder{Skin: MakeSkin(), Manifest: mfest, TagList: make(map[string][]DocumentMeta)}
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

	document.ParseMeta(postCtx[0])

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: true, BlogName: bname, PostName: document.Meta.Title}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	footd := Footer{IsNotIndex: true}
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

func (b *Builder) buildIndex(indexnum int, isFirstIndexBuild bool) {
	var plist []DocumentMeta
	var iname string
	var prevurl, nexturl string

	if len(b.PostList) <= indexnum && isFirstIndexBuild {
		plist = append(plist, b.PostList...)
		b.PostList = nil
		iname = "index.html"
		prevurl = ""
		nexturl = ""

	} else if b.IndexPageNum == 0 {
		plist = append(plist, b.PostList[:indexnum]...)
		b.PostList = b.PostList[indexnum:]

		iname = "index.html"
		b.IndexPageNum++
		prevurl = ""
		nexturl = fmt.Sprintf("./index-%v.html", b.IndexPageNum)

	} else if len(b.PostList) <= indexnum && !(isFirstIndexBuild) {
		plist = append(plist, b.PostList...)
		b.PostList = nil

		iname = fmt.Sprintf("index-%v.html", b.IndexPageNum)
		b.IndexPageNum++
		if b.IndexPageNum > 2 {
			prevurl = "./index.html"
		} else {
			prevurl = fmt.Sprintf("./index-%v.html", (b.IndexPageNum - 1))
		}

		nexturl = ""
	} else {
		plist = append(plist, b.PostList[:indexnum]...)
		b.PostList = b.PostList[indexnum:]

		iname = fmt.Sprintf("index-%v.html", b.IndexPageNum)
		b.IndexPageNum++
		if b.IndexPageNum > 2 {
			prevurl = "./index.html"
		} else {
			prevurl = fmt.Sprintf("./index-%v.html", (b.IndexPageNum - 1))
		}
		nexturl = fmt.Sprintf("./index-%v.html", b.IndexPageNum)
	}

	indexs := NewIndexData()
	indexs.Indexs = plist
	indexs.PrevPage = prevurl
	indexs.NextPage = nexturl

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: false, BlogName: bname}
	head := BuildHeader(b.Skin, headd)
	indexs.Head = head

	footd := Footer{IsNotIndex: false}
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

	indexpath := b.wd + "/dist/" + iname

	_ = os.WriteFile(indexpath, writer.Bytes(), 0777)

}

func (b *Builder) buildTagsPage() {
   fmt.Println("Tag Page Building...");
	tags := NewTagsData()
	tags.Tags = b.TagList

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: false, BlogName: bname}
	head := BuildHeader(b.Skin, headd)
	tags.Head = head
   
	footd := Footer{IsNotIndex: false}
	foot := BuildFooter(b.Skin, footd)
	tags.Foot = foot

	navd := Nav{IsNotIndex: false, BlogName: bname}
	nav := BuildNav(b.Skin, navd)
	tags.Nav = nav

	file, fserr := os.ReadFile(b.Skin.Info.Paths.Tags)

	if fserr != nil {
		panic(fserr)
	}

	var builder template.Template
	t, err := builder.Parse(string(file))

	if err != nil {
		panic(err)
	}

	var writer bytes.Buffer
	_ = t.Execute(&writer, tags)

	indexpath := b.wd + "/dist/tags.html"

	_ = os.WriteFile(indexpath, writer.Bytes(), 0777)

}

func (b *Builder) buildAboutPage() {
	markup := nemomark.MakeNemomark()

	ctx, perr := os.ReadFile(b.wd + "/post/about.ps")

	if perr != nil {
		panic(perr)
	}

	postRawctx := string(ctx)

	document := NewAboutPage()

	document.Content = markup.Mark(postRawctx)

	bname := b.Manifest.Name

	headd := Header{IsNotIndex: false, BlogName: bname}
	head := BuildHeader(b.Skin, headd)
	document.Head = head

	footd := Footer{IsNotIndex: false}
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

	aboutPath := b.wd + "/dist/about.html"
	os.WriteFile(aboutPath, writer.Bytes(), 0777)

}

func (b *Builder) packRes() {
	//TODO: REMOVE
	_, ex := os.Stat("dist")
	if os.IsNotExist(ex) {
		_ = os.Mkdir("dist", os.ModePerm)
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

	skinsrc := "./skin/static/"
	skindet := "./dist/static/"

	cerr := DirCopy(skinsrc, skindet)
	if cerr != nil {
		panic(cerr)
	}

	resrc := "./post/res"
	resdet := "./dist/page/res"

	rerr := DirCopy(resrc, resdet)
	if rerr != nil {
		panic(rerr)
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
	b.wd = wd
	workd := b.wd + "/post/"

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
		fdir := b.wd + "/dist/page/" + name

		fmt.Println(" => ", name)

		_ = os.WriteFile(fdir, []byte(content), 0777)

		dmeta.Path = name
		b.PostList = append(b.PostList, dmeta)
		if dmeta.Tags != "" {
			b.TagList[dmeta.Tags] = append(b.TagList[dmeta.Tags], dmeta)
		}
	}

	fmt.Println("TGL :", b.TagList)

	sort.Slice(b.PostList, func(i, j int) bool {
		return b.PostList[i].Timestamp.StampSize() < b.PostList[j].Timestamp.StampSize()
	})

	isFirst := true
	for !(len(b.PostList) == 0) {
		b.buildIndex(b.Skin.Info.Conf.IndexNum, isFirst)
		isFirst = false
	}

	_, ex := os.Stat(b.wd + "/post/about.ps")
	if !os.IsNotExist(ex) {
		b.buildAboutPage()
	}

	b.buildTagsPage()

	b.packRes()
}
