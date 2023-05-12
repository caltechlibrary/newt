#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = newt

MD_PAGES = $(shell ls -1 *.md)

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/\.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES) pagefind
	@for FNAME in $(HTML_PAGES); do git add "$$FNAME"; done
	@if [ -f index.html ]; then git add index.html; fi
	@git commit -am 'website build process'

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=page.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; git add index.html; fi

pagefind: .FORCE
	pagefind --verbose --exclude-selectors="nav,header,footer" --bundle-dir ./pagefind --source .
	git add pagefind

clean:
	@if [ -f index.html ]; then rm index.html; fi
	@for FNAME in $(HTML_PAGES); do if [ -f "$$FNAME" ]; then rm "$$FNAME"; fi; done
	@if [ -f libdataset/index.html ]; then rm libdataset/index.html; fi


.FORCE:
