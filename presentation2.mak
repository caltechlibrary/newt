
# where FORMAT is either s5, slidy, slideous, dzslides, or revealjs.
WEB_FORMAT = slidy

build: .FORCE clean presentation_dir html

presentation_dir: .FORCE
	mkdir -p presentation2

html: .FORCE
	pandoc -s -t $(WEB_FORMAT) newt-presentation2.md -o presentation2/index.html
	git add presentation2/index.html

pdf: .FORCE
	pandoc -s -t beamer newt-presentation2.md -o presentation2/newt-presentation.pdf

pptx: .FORCE
	pandoc -s newt-presentation2.md -o presentation2/newt-presentation.pptx

clean: .FORCE
	if [ -f presentation2/index.html ]; then rm presentation2/*.html; fi
	if [ -f presentation2/newt-presentation.pdf ]; then rm presentation2/*.pdf; fi
	if [ -f presentation2/newt-presentation.pptx ]; then rm presentation2/*.pptx; fi

.FORCE:
