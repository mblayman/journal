.PHONY: local

local:
	uv run honcho start -f Procfile

shell:
	uv run manage.py shell

coverage:
	uv run pytest --cov=journal --migrations -n 2 --dist loadfile

test: fcov

# fcov == "fast coverage" by skipping migrations checking. Save that for CI.
fcov:
	@echo "Running fast coverage check"
	@uv run pytest --cov=journal -n 4 --dist loadfile -q

build:
	go build -o app

run: build
	./app
