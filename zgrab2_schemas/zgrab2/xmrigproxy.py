# zschema sub-schema for zgrab2's xmrigproxy module
# Registers zgrab2-xmrigproxy globally, and xmrigproxy with the main zgrab2 schema.
from zschema.leaves import *
from zschema.compounds import *
import zschema.registry

import zcrypto_schemas.zcrypto as zcrypto
import zgrab2

xmrigproxy_scan_response = SubRecord({
    "result": SubRecord({
        # TODO FIXME IMPLEMENT SCHEMA
    })
}, extends=zgrab2.base_scan_response)

zschema.registry.register_schema("zgrab2-xmrigproxy", xmrigproxy_scan_response)

zgrab2.register_scan_response_type("xmrigproxy", xmrigproxy_scan_response)
