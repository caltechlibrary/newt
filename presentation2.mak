
# where FORMAT is either s5, slidy, slideous, dzslides, or revealjs.
WEB_FORMAT = slidy

build: .FORCE clean presentation_dir html pdf pptx

presentation_dir: .FORCE
	mkdir -p presentation2

html: .FORCE
	pandoc -V lang=en -s -t $(WEB_FORMAT) newt-presentation2.md -o presentation2/index.html
	git add presentation2/index.html

pdf: .FORCE
	pandoc -V lang=en -s -t beamer newt-presentation2.md -o presentation2/newt-p2.pdf

pptx: .FORCE
	pandoc -V lang=en -s newt-presentation2.md -o presentation2/newt-p2.pptx

clean: .FORCE
	if [ -f presentation2/index.html ]; then rm presentation2/*.html; fi
	if [ -f presentation2/newt-p*.pdf ]; then rm presentation2/*.pdf; fi
	if [ -f presentation2/newt-p*.pptx ]; then rm presentation2/*.pptx; fi

.FORCE:
