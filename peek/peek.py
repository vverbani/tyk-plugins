from tyk.decorators import *
from gateway import TykGateway as tyk

import os
import json

@Hook
def Peek(request, session, spec):
  logLevel = "info"
  tyk.log("PreHook START", logLevel)
  # load the config data from the API definition
  api_config_data = json.loads(spec['config_data'])
  #peekPath = "/opt/tyk-gateway/tyk.conf"
  if "peek" in api_config_data:
    peekPath = api_config_data["peek"]
    request.object.return_overrides.headers['content-type'] = 'text/plain'
    tyk.log(f"PreHook Checking {peekPath}", logLevel)
    if peekPath == 'ENV':
      request.object.return_overrides.response_code = 200
      env = ""
      for k, v in os.environ.items():
        env = env + f'{k}={v}\n'
      request.object.return_overrides.response_body = env
    else:
      if os.path.exists(peekPath):
        if os.path.isdir(peekPath):
          # directory, return a listing
          tyk.log(f"PreHook {peekPath} is a directory", logLevel)
          request.object.return_overrides.response_code = 200
          request.object.return_overrides.response_body = "\n".join(os.listdir(peekPath))
        elif os.path.isfile(peekPath):
          # a file open and return it
          tyk.log(f"PreHook {peekPath} is a file", logLevel)
          with open(peekPath) as f:
            request.object.return_overrides.response_code = 200
            request.object.return_overrides.response_body = " ".join(f.readlines())

      else:
        # missing, give error
        request.object.return_overrides.response_code = 404
        request.object.return_overrides.response_error = f"PreHook {peekPath} not found or not a file or directory"
        tyk.log(f"PreHook {peekPath} not found or not a file or directory", logLevel)
      tyk.log("PreHook END", logLevel)
  return request, session
