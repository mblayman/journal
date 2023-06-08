.PHONY: local

local:
	heroku local

deploy:
	git push heroku main

# graph:
# 	./manage.py graph_models \
# 		--rankdir BT \
# 		accounts \
# 		core \
# 		courses \
# 		notifications \
# 		referrals \
# 		reports \
# 		schools \
# 		students \
# 		teachers \
# 		users \
# 		-o models.png

coverage:
	pytest --cov=journal --migrations -n 2 --dist loadfile

# fcov == "fast coverage" by skipping migrations checking. Save that for CI.
fcov:
	@echo "Running fast coverage check"
	@pytest --cov=journal -n 4 --dist loadfile -q

# mypy:
# 	mypy homeschool project manage.py
