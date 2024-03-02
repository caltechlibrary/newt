#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = birds4_demo

MD_PAGES = $(shell ls -1 *.md | grep -v 'nav.md')

HTML_PAGES = $(shell ls -1 *.md | grep -v 'nav.md' | sed -E 's/.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES)

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html --lua-filter=../../links-to-html.lua --template=../../page.tmpl
	@if [ "$(basename $@)" = "README" ]; then mv README.html index.html; git add index.html; else git add "$(basename $@).html"; fi

clean:
	@if [ -f index.html ]; then rm *.html; fi

.FORCE:
