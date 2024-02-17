#
# Makefile to assemble and render Book of Newt as PDF
#

BOOK_SECTIONS = preface.md \
    introduction.md \
	design_choices.md \
	newt_yaml_syntax.md \
	newtgenerator_explained.md \
	newtmustache_explained.md \
	newtrouter_explained.md \
	reference_material.md

build: setup markdown book

# Setup everything for making the book
setup: .FORCE
	if [ ! -d pdf ]; then mkdir -p pdf; fi
	if [ -f pdf/book_of_newt.pdf ]; then rm pdf/book_of_newt.pdf; fi
	
# Assemble the text as Markdown
markdown: $(BOOK_SECTIONS) 
	pandoc -s \
	    -f markdown -t markdown --toc -N \
		cover.md $(BOOK_SECTIONS)\
		>book.md

# Assemble the text as HTML
html: $(BOOK_SECTIONS) 
	pandoc -s \
	    -f markdown -t HTML5 --toc -N \
		cover.md $(BOOK_SECTIONS)\
		>book.html

# Make the book as a PDF
book: .FORCE
	pandoc -s -f markdown -t pdf --toc -N \
		--pdf-engine=xelatex \
		-V documentclass=book \
		cover.md \
	    $(BOOK_SECTIONS) \
		-o pdf/book_of_newt.pdf

.FORCE:
