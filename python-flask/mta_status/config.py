import requests
from cachecontrol import CacheControl

SESSION = CacheControl(requests.Session())
MTA_FEED_URL = 'http://web.mta.info/status/serviceStatus.txt'
FEED_CASE = {
	'subway': 'Subway',
	'bus': 'Bus',
	# 'BT': 'BT',
	# 'LIRR': 'LIRR',
	# 'MetroNorth': 'MetroNorth',
}
