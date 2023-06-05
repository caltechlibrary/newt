
# where FORMAT is either s5, slidy, slideous, dzslides, or revealjs.
WEB_FORMAT = slidy

build: .FORCE clean presentation_dir html pptx

presentation_dir: presentation/
	mkdir -p presentation

html: .FORCE
	pandoc -s -t $(WEB_FORMAT) newt-presentation.md -o presentation/index.html

pdf: .FORCE
	pandoc -s -t beamer newt-presentation.md -o presentation/newt-presentation.pdf

pptx: .FORCE
	pandoc -s newt-presentation.md -o presentation/newt-presentation.pptx

clean: .FORCE
	if [ -f presentation/index.html ]; then rm presentation/*.html; fi
	if [ -f presentation/newt-presentation.pdf ]; then rm presentation/*.pdf; fi
	if [ -f presentation/newt-presentation.pptx ]; then rm presentation/*.pptx; fi

.FORCE:
