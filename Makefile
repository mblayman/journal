.PHONY: local

local:
	heroku local

deploy:
	git push heroku main

# Generate an image of the models in the system.
graph:
	./manage.py graph_models

coverage:
	pytest --cov=journal --migrations -n 2 --dist loadfile

# fcov == "fast coverage" by skipping migrations checking. Save that for CI.
fcov:
	@echo "Running fast coverage check"
	@pytest --cov=journal -n 4 --dist loadfile -q

# mypy:
# 	mypy homeschool project manage.py
