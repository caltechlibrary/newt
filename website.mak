#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = newt

BASE_URL =
ifneq ($(base_url),)
  BASE_URL = $(base_url)
endif

MD_PAGES = $(shell ls -1 *.md | grep -v 'nav.md')

HTML_PAGES = $(shell ls -1 *.md | grep -v 'nav.md' | sed -E 's/.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES) pagefind

$(HTML_PAGES): $(MD_PAGES) .FORCE
	echo "" | pandoc --metadata base_url=$(base_url) -s --to markdown -o search.md --template=search.md.tmpl; git add search.md
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html --lua-filter=links-to-html.lua --template=page.tmpl
	@if [ "$(basename $@)" = "README" ]; then mv README.html index.html; git add index.html; else git add "$(basename $@).html"; fi

pagefind: .FORCE
	pagefind --verbose --exclude-selectors="nav,header,footer" --site .
	git add pagefind

clean:
	@if [ -f index.html ]; then rm *.html; fi

.FORCE:
