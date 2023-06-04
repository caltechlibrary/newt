
# where FORMAT is either s5, slidy, slideous, dzslides, or revealjs.
WEB_FORMAT = slidy

build: clean html

html: presentation_dir presentation/index.html

presentation_dir: presentation/
	mkdir -p presentation

presentation/index.html: newt-presentation.md
	pandoc -s -t $(WEB_FORMAT) newt-presentation.md -o presentation/index.html
	pandoc -s newt-presentation.md -o newt-presentation.pptx

clean: .FORCE
	if [ -f presentation/index.html ]; then rm presentation/*.html; fi

pdf: .FORCE
	pandoc -s -t beamer newt-presentation.md -o presentation/newt-presentation.pdf

.FORCE:
