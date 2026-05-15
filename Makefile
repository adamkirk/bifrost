.PHONY: npm

npm:
	@[ $$(node -v | tr -d v | cut -d. -f1) -ge 25 ] || { echo "Error: node $$(node -v) < required v25"; exit 1; }
	@[ $$(npm -v | cut -d. -f1) -ge 11 ] || { echo "Error: npm $$(npm -v) < required v11"; exit 1; }
	npm install


