
# where FORMAT is either s5, slidy, slideous, dzslides, or revealjs.
WEB_FORMAT = slidy

build: .FORCE clean presentation_dir html pdf pptx

presentation_dir: .FORCE
	mkdir -p presentation3

html: .FORCE
	pandoc -V lang=en -s -t $(WEB_FORMAT) presentation3.md -o presentation3/index.html
	git add presentation3/index.html

pdf: .FORCE
	pandoc -V lang=en -s -t beamer presentation3.md -o presentation3/newt-p3.pdf

pptx: .FORCE
	pandoc -V lang=en -s presentation3.md -o presentation3/newt-p3.pptx

clean: .FORCE
	if [ -f presentation3/index.html ]; then rm presentation3/*.html; fi
	if [ -f presentation3/newt-p*.pdf ]; then rm presentation3/*.pdf; fi
	if [ -f presentation3/newt-p*.pptx ]; then rm presentation3/*.pptx; fi

.FORCE:
