import xml.etree.ElementTree as ET
import requests
from . import app

FEED_META_TAGS = ['responsecode', 'timestamp']

def get_xml(url):
	r = app.config['SESSION'].get(url)
	r.raise_for_status()
	return ET.fromstring(r.text)

def normalize_case(service):
	"""Normalize service name to have proper casing.

	Some services in the feed are properly cased (e.g. "LIRR", "MetroNorth"),
	while others are lowercase (e.g. "subway", "bus").
	:param service: str service name
	:return: str service name properly cased
	"""
	# if service[0].isupper():
	# 	return service
	# return service.title()
	proper_name = app.config['FEED_CASE'].get(service)
	if proper_name is None:
		return service
	return proper_name

def get_services(url=None):
	if url is None:
		url = app.config['MTA_FEED_URL']
	try:
		x = get_xml(url)
	except requests.exceptions.HTTPError as e:
		# return e.response.status_code, {'status': 'error', 'message': e.message}
		return {'status': 'error', 'message': e.message, 'data': []}

	services = []
	for child in x:
		tag = child.tag
		if tag in FEED_META_TAGS:
			continue
		services.append(normalize_case(tag))

	# return 200, {'status': 'success', 'data': services}
	return {'status': 'success', 'data': services}

def get_service_status(service, url=None):
	if url is None:
		url = app.config['MTA_FEED_URL']
	try:
		x = get_xml(url)
	except requests.exceptions.HTTPError:
		return {'status': 'error', 'message': e.message, 'data': []}

	statuses = []
	for line in x.findall(service+'/line'):
		st = {'name': line.find('name').text,
				'status': line.find('status').text,
				'message': line.find('text').text}
		statuses.append(st)

	return {'status': 'success', 'data': statuses}
