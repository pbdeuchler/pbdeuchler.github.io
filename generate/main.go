package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Post struct {
	Title          string
	Subtitle       string
	Date           time.Time
	DateFmt        string
	LastModified   time.Time
	LastModifiedFmt string
	Slug           string
	Draft          bool
	Content        template.HTML
}

type Page struct {
	Title   string
	Slug    string
	Content template.HTML
}

type SiteData struct {
	Title        string
	Path         string
	HeaderImages []string
	Posts        []Post
	Post         *Post
	Page         *Page
}

func main() {
	root := ".."
	contentDir := filepath.Join(root, "contents")
	outDir := filepath.Join(root, "site")
	tmplDir := filepath.Join(contentDir, "templates")

	// Clean output
	os.RemoveAll(outDir)

	// Load base template, then clone + overlay each child
	base, err := template.ParseFiles(filepath.Join(tmplDir, "base.html"))
	must(err, "parsing base template")

	loadTemplate := func(name string) *template.Template {
		t, err := template.Must(base.Clone()).ParseFiles(filepath.Join(tmplDir, name))
		must(err, "parsing template "+name)
		return t
	}

	indexTmpl := loadTemplate("index.html")
	postTmpl := loadTemplate("post.html")
	pageTmpl := loadTemplate("page.html")
	fourOhFourTmpl := loadTemplate("404.html")

	// Scan header images
	headerImages := scanHeaderImages(filepath.Join(contentDir, "static", "img", "headers"))

	// Parse posts
	posts := parsePosts(filepath.Join(contentDir, "posts"))
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// Filter drafts
	var published []Post
	for _, p := range posts {
		if !p.Draft {
			published = append(published, p)
		}
	}

	// Parse pages
	pages := parsePages(filepath.Join(contentDir, "pages"))

	// Render index
	renderTemplate(indexTmpl, filepath.Join(outDir, "index.html"), SiteData{
		Title:        "Philip Deuchler",
		Path:         "/",
		HeaderImages: headerImages,
		Posts:        published,
	})

	// Render posts
	for i := range published {
		p := &published[i]
		dest := filepath.Join(outDir, "posts", p.Slug, "index.html")
		renderTemplate(postTmpl, dest, SiteData{
			Title:        p.Title + " \u2013 Philip Deuchler",
			Path:         "/posts/" + p.Slug + "/",
			HeaderImages: headerImages,
			Post:         p,
		})
	}

	// Render pages
	for i := range pages {
		pg := &pages[i]
		dest := filepath.Join(outDir, pg.Slug, "index.html")
		renderTemplate(pageTmpl, dest, SiteData{
			Title:        pg.Title + " \u2013 Philip Deuchler",
			Path:         "/" + pg.Slug + "/",
			HeaderImages: headerImages,
			Page:         pg,
		})
	}

	// Render 404
	renderTemplate(fourOhFourTmpl, filepath.Join(outDir, "404.html"), SiteData{
		Title:        "404 \u2013 Philip Deuchler",
		Path:         "/404.html",
		HeaderImages: headerImages,
	})

	// Copy static files
	staticDir := filepath.Join(contentDir, "static")
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		copyDir(staticDir, outDir)
	}

	fmt.Printf("Generated %d posts, %d pages → %s\n", len(published), len(pages), outDir)
}

func parseFrontmatter(raw []byte) (map[string]string, string) {
	content := string(raw)
	if !strings.HasPrefix(content, "---\n") {
		return nil, content
	}
	end := strings.Index(content[4:], "\n---\n")
	if end < 0 {
		return nil, content
	}
	fmBlock := content[4 : 4+end]
	body := content[4+end+5:]

	fm := make(map[string]string)
	for _, line := range strings.Split(fmBlock, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		fm[key] = val
	}
	return fm, body
}

var md = goldmark.New(
	goldmark.WithExtensions(extension.Footnote),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func renderMarkdown(src string) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		fmt.Fprintf(os.Stderr, "warning: markdown render error: %v\n", err)
	}
	return template.HTML(buf.String())
}

func parsePosts(dir string) []Post {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		must(err, "reading posts dir")
	}
	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		must(err, "reading post "+e.Name())

		fm, body := parseFrontmatter(raw)
		slug := strings.TrimSuffix(e.Name(), ".md")

		title := fm["title"]
		if title == "" {
			title = slug
		}

		var date time.Time
		if d, ok := fm["date"]; ok {
			date, _ = time.Parse("2006-01-02", d)
		}

		draft := fm["draft"] == "true"

		var lastModified time.Time
		if lm, ok := fm["last-modified"]; ok {
			lastModified, _ = time.Parse("2006-01-02", lm)
		}

		var lastModifiedFmt string
		if !lastModified.IsZero() {
			lastModifiedFmt = lastModified.Format("Jan 2, 2006")
		}

		posts = append(posts, Post{
			Title:           title,
			Subtitle:        fm["subtitle"],
			Date:            date,
			DateFmt:         date.Format("Jan 2, 2006"),
			LastModified:    lastModified,
			LastModifiedFmt: lastModifiedFmt,
			Slug:            slug,
			Draft:           draft,
			Content:         renderMarkdown(body),
		})
	}
	return posts
}

func parsePages(dir string) []Page {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		must(err, "reading pages dir")
	}
	var pages []Page
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		must(err, "reading page "+e.Name())

		fm, body := parseFrontmatter(raw)
		slug := strings.TrimSuffix(e.Name(), ".md")

		title := fm["title"]
		if title == "" {
			title = slug
		}

		pages = append(pages, Page{
			Title:   title,
			Slug:    slug,
			Content: renderMarkdown(body),
		})
	}
	return pages
}

func scanHeaderImages(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		must(err, "reading header images dir")
	}
	var images []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".webp") {
			images = append(images, e.Name())
		}
	}
	return images
}

func renderTemplate(tmpl *template.Template, dest string, data SiteData) {
	must(os.MkdirAll(filepath.Dir(dest), 0o755), "creating dir for "+dest)
	var buf bytes.Buffer
	must(tmpl.ExecuteTemplate(&buf, "base", data), "executing template for "+dest)
	must(os.WriteFile(dest, buf.Bytes(), 0o644), "writing "+dest)
}

func copyDir(src, dst string) {
	filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func must(err error, context string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s: %v\n", context, err)
		os.Exit(1)
	}
}
