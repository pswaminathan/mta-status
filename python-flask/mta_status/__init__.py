import os.path
from flask import Flask

app = Flask(__name__)
app.config.from_pyfile(os.path.join(os.path.abspath(os.path.dirname(__file__)),
									'config.py'))
