# journal

A SaaS Journal - future source home of journeyinbox.com

## setup

```
python3 -m venv venv
pip install -r requirements-dev.txt -r requirements.txt
```

Local superuser:

* username: matt
* email: matt@testing.com
* password: createhorse

## local development tools

```
CFLAGS="-I/opt/homebrew/Cellar/graphviz/7.0.5/include" LDFLAGS="-L/opt/homebrew/Cellar/graphviz/7.0.5/lib" pip install -r requirements-dev.txt -r requirements.txt
```
