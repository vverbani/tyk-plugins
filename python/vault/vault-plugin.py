from tyk.decorators import *
from gateway import TykGateway as tyk

sys.path.append(os.path.join(os.path.abspath(os.path.dirname(__file__)), 'vendor/lib/python3.6/site-packages'))

import six
import hvac

@Hook
def PreHook(request, session, spec):
  vault_url='http://10.0.0.21:8200/'
  vault_token='myroot'
  client = hvac.Client(url=vault_url, token=vault_token)
  is_authenticated = client.is_authenticated()

  print("Vault authenticated: " + is_authenticated.__str__())

  if is_authenticated:
    key = "fred"
    secret = client.secrets.kv.v2.read_secret(mount_point='secret', path='data')['data']['data'][key]
    tyk.log(f"My secret:{key} -> " + secret, "info")
  return request, session

