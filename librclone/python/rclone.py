"""
Python interface to librclone.so using ctypes

Create an rclone object

    rclone = SecureCloudEngine(shared_object="/path/to/librclone.so")

Then call rpc calls on it

    rclone.rpc("rc/noop", a=42, b="string", c=[1234])

When finished, close it

    rclone.close()
"""

__all__ = ('SecureCloudEngine', 'SecureCloudEngineException')

import os
import json
import subprocess
from ctypes import *

class SecureCloudEngineRPCString(c_char_p):
    """
    This is a raw string from the C API

    With a plain c_char_p type, ctypes will replace it with a
    regular Python string object that cannot be used with
    SecureCloudEngineFreeString. Subclassing prevents it, while the string
    can still be retrieved from attribute value.
    """
    pass

class SecureCloudEngineRPCResult(Structure):
    """
    This is returned from the C API when calling SecureCloudEngineRPC
    """
    _fields_ = [("Output", SecureCloudEngineRPCString),
                ("Status", c_int)]

class SecureCloudEngineException(Exception):
    """
    Exception raised from rclone

    This will have the attributes:

    output - a dictionary from the call
    status - a status number
    """
    def __init__(self, output, status):
        self.output = output
        self.status = status
        message = self.output.get('error', 'Unknown rclone error')
        super().__init__(message)

class SecureCloudEngine():
    """
    Interface to SecureCloudEngine via librclone.so

    Initialise with shared_object as the file path of librclone.so
    """
    def __init__(self, shared_object=f"./librclone{'.dll' if os.name == 'nt' else '.so'}"):
        self.rclone = CDLL(shared_object)
        self.rclone.SecureCloudEngineRPC.restype = SecureCloudEngineRPCResult
        self.rclone.SecureCloudEngineRPC.argtypes = (c_char_p, c_char_p)
        self.rclone.SecureCloudEngineFreeString.restype = None
        self.rclone.SecureCloudEngineFreeString.argtypes = (c_char_p,)
        self.rclone.SecureCloudEngineInitialize.restype = None
        self.rclone.SecureCloudEngineInitialize.argtypes = ()
        self.rclone.SecureCloudEngineFinalize.restype = None
        self.rclone.SecureCloudEngineFinalize.argtypes = ()
        self.rclone.SecureCloudEngineInitialize()
    def rpc(self, method, **kwargs):
        """
        Call an rclone RC API call with the kwargs given.

        The result will be a dictionary.

        If an exception is raised from rclone it will of type
        SecureCloudEngineException.
        """
        method = method.encode("utf-8")
        parameters = json.dumps(kwargs).encode("utf-8")
        resp = self.rclone.SecureCloudEngineRPC(method, parameters)
        output = json.loads(resp.Output.value.decode("utf-8"))
        self.rclone.SecureCloudEngineFreeString(resp.Output)
        status = resp.Status
        if status != 200:
            raise SecureCloudEngineException(output, status)
        return output
    def close(self):
        """
        Call to finish with the rclone connection
        """
        self.rclone.SecureCloudEngineFinalize()
        self.rclone = None
    @classmethod
    def build(cls, shared_object):
        """
        Builds rclone to shared_object if it doesn't already exist

        Requires go to be installed
        """
        if os.path.exists(shared_object):
            return
        print("Building "+shared_object)
        subprocess.check_call(["go", "build", "--buildmode=c-shared", "-o", shared_object, "github.com/wizdom13/SecureCloudEngine/librclone"])
