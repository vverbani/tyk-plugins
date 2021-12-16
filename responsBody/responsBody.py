from tyk.decorators import *
from gateway import TykGateway as tyk

@Hook
def ChangeResponseBody(request, response, session, metadata, spec):
	tyk.log("ResponseHook is called", "info")
	tyk.log("ResponseHook: upstream returned {0}".format(response.status_code), "info")
	tyk.log(f"response object is: {response}", "info")
	response.headers["injectedkey"] = "injectedvalue"
	response.raw_body = b"{Pete Woz Here}"
	response.headers["Content-Length"] = str(len(response.raw_body))
	return response
