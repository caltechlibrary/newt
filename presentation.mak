
build: html

html: presentation_dir presentation/index.html

presentation_dir: presentation/
	mkdir -p presentation

presentation/index.html: newt-presentation.md
	pandoc -s -t slidy newt-presentation.md -o presentation/index.html

clean: .FORCE
	if [ -f presentation/index.html ]; then rm presentation/*.html; fi

pdf: .FORCE
	pandoc -s -t beamer newt-presentation.md -o presentation/newt-presentation.pdf

.FORCE:
