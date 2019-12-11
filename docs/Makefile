
pages: html fix

html: head
	grep -v ':::' README.0.md > README.md
	grep -v ':::' DETAILS.0.md > DETAILS.md
	grep -v ':::' DIAGRAM.0.md > DIAGRAM.md
	grep -v ':::' SERVICE.0.md > SERVICE.md
	grep -v ':::' SIDE_BY_SIDE_LEFT.0.md > SIDE_BY_SIDE_LEFT.md
	grep -v ':::' SIDE_BY_SIDE_RIGHT.0.md > SIDE_BY_SIDE_RIGHT.md
	pandoc ../README.0.md -o - >> USAGE.html
	pandoc ./README.0.md -o - >> index.html
	pandoc ./DETAILS.0.md -o - >> DETAILS.html
	pandoc ./DIAGRAM.0.md -o - >> DIAGRAM.html
	pandoc ./SERVICE.0.md -o - >> SERVICE.html
	pandoc ./SIDE_BY_SIDE_LEFT.0.md -o - >> SIDE_BY_SIDE_LEFT.html
	pandoc ./SIDE_BY_SIDE_RIGHT.0.md -o - >> SIDE_BY_SIDE_RIGHT.html
	pandoc --highlight-style=tango index.html -t html5 -o index.pdf
	pandoc --highlight-style=tango DETAILS.html -t html5 -o details.pdf
	pandoc --highlight-style=tango SERVICE.html -t html5 -o service.pdf


fix:
	find . -name '*.html' -exec sed -i 's|.md|.html|g' {} \;
	find . -name '*.html' -exec tidy -i -m -w 80 -ashtml -utf8 {} \;

head:
	@echo '<!DOCTYPE html>' > head.html
	@echo '<html>' >> head.html
	@echo '<head>' >> head.html
	@echo '  <link rel="stylesheet" type="text/css" href="style.css">' >> head.html
	@echo "  <title>AnonVPN</title>" >> head.html
	@echo '</head>' >> head.html
	@echo '' >> head.html
	cp head.html USAGE.html
	cp head.html index.html
	cp head.html DETAILS.html
	cp head.html DIAGRAM.html
	cp head.html SERVICE.html
	@echo "" > SIDE_BY_SIDE_LEFT.html
	@echo "" > SIDE_BY_SIDE_RIGHT.html

host:
	cd ../ && eephttpd -d ./docs
