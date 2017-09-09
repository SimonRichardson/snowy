all: scripts/snowy.js

scripts/snowy.js:
	elm make --yes src/Main.elm --output scripts/snowy.js

clean: FORCE
	rm -f scripts/snowy.js

FORCE: