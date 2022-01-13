from tyk.decorators import *
from gateway import TykGateway as tyk

def dump(obj):
	for attr in dir(obj):
		tyk.log("obj.%s = %r" % (attr, getattr(obj, attr)), "info")

@Hook
def Headers(request, session, spec):
	logLevel = "info"
	tyk.log("Headers START", logLevel)
	for key in request.object.headers:
		tyk.log("Header: " + key + " is " + request.object.headers[key], logLevel)
	tyk.log("Headers END", logLevel)
	return request, session
