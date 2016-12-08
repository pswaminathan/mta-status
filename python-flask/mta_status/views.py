from flask import flash, render_template
from . import app
from .controllers import normalize_case, get_services, get_service_status

@app.route('/')
def index():
	services = get_services()
	if services.get('status') != 'success':
		flash('Error retrieving services. Message: '
			  '{}'.format(services.get('message')))
	return render_template("index.html")

@app.route('/status/<service>')
def status(service):
	pass
