#!/usr/bin/env python3
"""
WAPI Schema Extractor — Static AST-based extraction of WAPI object schemas
from the NIOS codebase.

This script implements "Approach B" from the ADR:
  wapi-terraform-schema-generation.md

It parses Python source files under products/*/server/src/wapi/ using AST
analysis (no import of NIOS modules required), extracts WAPIObject class
definitions, WF field metadata, and cross-references with RTXML XML files
to produce a comprehensive JSON schema suitable for Terraform provider
code generation.

Usage:
    python extract_schema.py --nios-root /path/to/nios --output-dir ./schemas

    # Extract only specific object types
    python extract_schema.py --nios-root /path/to/nios --output-dir ./schemas \
        --objects record:a,network,zone_auth

    # Show summary without writing files
    python extract_schema.py --nios-root /path/to/nios --summary
"""

import argparse
import ast
import glob
import json
import os
import re
import sys
import xml.etree.ElementTree as ET
from collections import defaultdict
from dataclasses import dataclass, field, asdict
from pathlib import Path
from typing import Optional


# ---------------------------------------------------------------------------
# Data models
# ---------------------------------------------------------------------------

@dataclass
class RTXMLMember:
    """A <member> element from an RTXML XML file."""
    name: str
    type: str = "rtxml.string"
    key_type: Optional[str] = None
    ref_type: Optional[str] = None
    unique_id: bool = False
    version_id: bool = False
    internal_ro: bool = False
    auto_uid: bool = False
    synthetic_direct: bool = False
    synthetic_direct_func: Optional[str] = None
    metadata_of: Optional[str] = None
    default_value: Optional[str] = None
    index: bool = False
    syntax_name: Optional[str] = None
    syntax_param: Optional[str] = None
    description: Optional[str] = None


@dataclass
class RTXMLStructure:
    """A <structure> element from an RTXML XML file."""
    name: str
    package: str
    full_type: str  # e.g. ".com.infoblox.dns.bind_a"
    members: dict = field(default_factory=dict)  # name -> RTXMLMember


@dataclass
class WAPIField:
    """A WF* field extracted from a WAPIObject class."""
    wname: str
    wf_class: str  # e.g. "WF", "WFBool", "WFUInt", etc.
    doc: Optional[str] = None
    iname: Optional[str] = None
    purpose: str = "rwu"
    create: bool = False
    create_only: bool = False
    create_default: Optional[str] = None
    use_flag: Optional[str] = None
    search: Optional[str] = None
    std_field: bool = False
    nillable: bool = False
    default_value: Optional[object] = None
    default_docstring: Optional[str] = None
    enumlist: Optional[list] = None
    strip: bool = False
    i2w_func: Optional[str] = None
    w2i_func: Optional[str] = None
    in_struct: Optional[str] = None
    typecls: Optional[str] = None
    scope_of_delegation: bool = False
    inheritance: bool = False
    # Derived from WF class
    is_array: bool = False
    is_struct: bool = False
    is_subobj: bool = False
    is_func_call: bool = False
    is_ext_attrs: bool = False
    is_obsolete: bool = False


@dataclass
class WAPIObjectDef:
    """A WAPIObject class definition."""
    class_name: str
    wapitype: str
    version: str
    source_file: str
    parent_class: Optional[str] = None
    cmdclass: Optional[str] = None
    fields: list = field(default_factory=list)  # list of WAPIField
    is_valid_version: bool = True
    restrictions: Optional[tuple] = None
    schedulable: bool = False
    globalsearch: bool = False
    csvexport: bool = False
    has_modify_fields: bool = False
    doc: Optional[str] = None


@dataclass
class SchemaField:
    """Final field schema for JSON output."""
    name: str
    type: str
    wf_class: str
    is_array: bool = False
    is_struct: bool = False
    standard_field: bool = False
    supports: str = "rwu"
    searchable_by: Optional[str] = None
    create_required: bool = False
    create_only: bool = False
    nillable: bool = False
    default_value: Optional[object] = None
    has_server_default: bool = False
    enum_values: Optional[list] = None
    use_flag: Optional[str] = None
    iname: Optional[str] = None
    doc: Optional[str] = None
    # Mutability classification
    mutability: str = "mutable"  # mutable, immutable, read_only
    computation: str = "stored"  # stored, server_computed, derived, client_derivable, inherited
    derived_from: Optional[list] = None
    derivation_function: Optional[str] = None
    stable_after_create: bool = False
    reconciliation_action: Optional[str] = None
    i2w_function: Optional[str] = None
    w2i_function: Optional[str] = None
    # RTXML cross-reference
    rtxml_key_type: Optional[str] = None
    rtxml_syntax: Optional[str] = None
    rtxml_unique_id: bool = False
    rtxml_synthetic_direct: bool = False


# ---------------------------------------------------------------------------
# RTXML Parser
# ---------------------------------------------------------------------------

class RTXMLParser:
    """Parse RTXML XML files to extract structure and member definitions."""

    def __init__(self, nios_root: str):
        self.nios_root = nios_root
        self.structures: dict[str, RTXMLStructure] = {}  # full_type -> structure

    def parse_all(self):
        """Parse all XML files under products/*/xml/."""
        pattern = os.path.join(self.nios_root, "products", "*", "xml", "*.xml")
        xml_files = glob.glob(pattern)
        parsed = 0
        errors = 0
        for xml_file in sorted(xml_files):
            try:
                self._parse_file(xml_file)
                parsed += 1
            except ET.ParseError as e:
                errors += 1
                # Some XML files may be malformed or use DTD features
                # that ElementTree doesn't handle — skip them
                pass
            except Exception as e:
                errors += 1
                print(f"  Warning: Error parsing {xml_file}: {e}",
                      file=sys.stderr)
        print(f"  Parsed {parsed} RTXML files ({errors} skipped)",
              file=sys.stderr)

    def _parse_file(self, xml_file: str):
        """Parse a single RTXML XML file."""
        # RTXML files reference a DTD we don't have locally in a parseable
        # form. Strip the DOCTYPE and parse as plain XML.
        with open(xml_file, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()

        # Remove DOCTYPE declaration
        content = re.sub(
            r'<!DOCTYPE[^>]*>',
            '',
            content,
            count=1
        )

        try:
            root = ET.fromstring(content)
        except ET.ParseError:
            return

        # Get package name
        package_name = root.get('name', '')

        for struct_elem in root.iter('structure'):
            struct_name = struct_elem.get('name')
            if not struct_name:
                continue

            full_type = f"{package_name}.{struct_name}"
            structure = RTXMLStructure(
                name=struct_name,
                package=package_name,
                full_type=full_type,
                members={}
            )

            for member_elem in struct_elem.findall('member'):
                member = self._parse_member(member_elem)
                if member:
                    structure.members[member.name] = member

            self.structures[full_type] = structure
            # Also index by short name for easier lookup
            self.structures[struct_name] = structure

    def _parse_member(self, elem) -> Optional[RTXMLMember]:
        """Parse a <member> element."""
        name = elem.get('name')
        if not name:
            return None

        member = RTXMLMember(name=name)
        member.type = elem.get('type', 'rtxml.string')
        member.key_type = elem.get('key-type')
        member.ref_type = elem.get('ref-type')
        member.unique_id = elem.get('unique-id', 'false').lower() == 'true'
        member.version_id = elem.get('version-id', 'false').lower() == 'true'
        member.internal_ro = elem.get('internal-ro', 'false').lower() == 'true'
        member.auto_uid = elem.get('auto-uid', 'false').lower() == 'true'
        member.synthetic_direct = elem.get('synthetic-direct', 'false').lower() == 'true'
        member.synthetic_direct_func = elem.get('synthetic-direct-func')
        member.metadata_of = elem.get('metadata-of')
        member.default_value = elem.get('default-value')
        member.index = elem.get('index', 'false').lower() == 'true'

        # Get syntax info
        syntax_elem = elem.find('syntax')
        if syntax_elem is not None:
            member.syntax_name = syntax_elem.get('name')
            member.syntax_param = syntax_elem.get('param')

        desc_elem = elem.find('description')
        if desc_elem is not None and desc_elem.text:
            member.description = desc_elem.text.strip()

        return member

    def lookup(self, rtxml_type: str) -> Optional[RTXMLStructure]:
        """Look up an RTXML structure by type name."""
        # Try full type first
        if rtxml_type in self.structures:
            return self.structures[rtxml_type]
        # Try just the structure name
        short = rtxml_type.rsplit('.', 1)[-1] if '.' in rtxml_type else rtxml_type
        return self.structures.get(short)

    def lookup_member(self, rtxml_type: str, member_name: str) -> Optional[RTXMLMember]:
        """Look up a specific member in an RTXML structure."""
        structure = self.lookup(rtxml_type)
        if structure:
            return structure.members.get(member_name)
        return None


# ---------------------------------------------------------------------------
# Python AST-based WAPI Parser
# ---------------------------------------------------------------------------

# Maps WF class names to inferred type information.
# Base types from wapibaseobj.py:
WF_CLASS_TYPE_MAP = {
    'WF': 'string',
    'WFBool': 'bool',
    'WFInt': 'int',
    'WFUInt': 'uint',
    'WFLong': 'int',
    'WFULong': 'uint',
    'WFEnum': 'enum',
    'WFTimeStamp': 'timestamp',
    'WFIPAddress': 'string',
    'WFIPv4Address': 'string',
    'WFIPv6Address': 'string',
    'WFIPv4Cidr': 'string',
    'WFIPv6Cidr': 'string',
    'WFFuncResult': 'string',
    'WFFuncResultSubObj': 'object_ref',
    'WFArray': 'string[]',
    'WFEnumArray': 'enum[]',
    'WFEnumArrayIString': 'enum[]',
    'WFStringFieldArray': 'string[]',
    'WFStruct': 'struct',
    'WFStructOrNone': 'struct',
    'WFStructArray': 'struct[]',
    'WFFlatStruct': 'struct',
    'WFSubObj': 'object_ref',
    'WFSubObjArray': 'object_ref[]',
    'WFExtensibleAttributes': 'ea',
    'WFNotInheritedExtensibleAttributes': 'ea',
    'WFFuncCall': 'funccall',
    'WFDBObjectByString': 'string',
    'WFDBObjectByStringArray': 'string[]',
    'WFRRSharedRecordGroup': 'string',
    'WFSchedule': 'struct',
    'WFObsolete': 'obsolete',
}


# Domain-specific WF subclasses that hardcode purpose='r' in their __init__.
# The AST parser can't see these overrides because they happen inside the
# subclass constructor, not in the call-site keyword arguments. Without this
# map, these fields are incorrectly classified as 'rwu' (the WF default)
# instead of 'r' (read-only), causing the Terraform provider to generate
# Optional instead of Computed schema attributes.
#
# Discovered by diffing AST output against v2 runtime schema: these fields
# showed supports='rwu' in extracted but supports='r' in runtime.
WF_CLASS_FORCED_PURPOSE = {
    'WFMSADUserData': 'r',        # Always read-only — server populates from AD
    'WFDiscoveryData': 'r',       # Always read-only — server populates from discovery
    'WFCloudInfo': 'r',           # Default purpose='r' in __init__ signature
    'WFAWSRte53RecordInfo': 'r',  # Read-only AWS Route53 integration data
}


def _infer_wf_type(wf_class: str) -> str:
    """Infer field type from WF class name using the explicit map first,
    then falling back to suffix-based heuristics for the ~250 domain-specific
    WF subclasses (e.g. WFDhcpMembersArray -> struct[], WFSNMPCredential -> struct).
    """
    if wf_class in WF_CLASS_TYPE_MAP:
        return WF_CLASS_TYPE_MAP[wf_class]

    name = wf_class
    # Suffix-based inference: order matters (most specific first)
    if name.endswith('Array'):
        return 'struct[]'
    if name.endswith('Struct'):
        return 'struct'
    if name.endswith('SubObj'):
        return 'object_ref'
    if name.endswith('Bool'):
        return 'bool'
    if name.endswith('Enum'):
        return 'enum'
    if name.endswith('Int') or name.endswith('UInt'):
        return 'uint'
    if name.endswith('Address') or name.endswith('Cidr'):
        return 'string'
    # *Name classes are typically string wrappers (WFViewName, WFNetViewName, etc.)
    if name.endswith('Name') or name.endswith('NameNone'):
        return 'string'
    # *FuncCall classes are function invocations, not data
    if name.endswith('FuncCall') or name.endswith('Function'):
        return 'funccall'
    # *Token and *RO classes are typically string values
    if name.endswith('Token') or name.endswith('RO'):
        return 'string'
    # *Data, *Info, *Setting(s), *Config, *Credential classes are structs
    if name.endswith('Data') or name.endswith('Info'):
        return 'struct'
    if 'Setting' in name or 'Config' in name:
        return 'struct'
    # Default: treat as struct (most domain WF subclasses wrap structs)
    return 'struct'

# Well-known i2w functions that are client-derivable
CLIENT_DERIVABLE_TRANSFORMS = {
    'i2w_wrap(str.lower)': 'lowercase',
    'i2w_invert_boolean': 'boolean_invert',
    'str.lower': 'lowercase',
}

# Well-known iname patterns for client-derivable fields
CLIENT_DERIVABLE_INAMES = {
    'dns_fqdn': 'punycode',
    'dns_rdata': 'punycode',
}


class WAPIASTParser:
    """Parse WAPI Python files using AST to extract object definitions."""

    def __init__(self, nios_root: str):
        self.nios_root = nios_root
        self.objects: dict[str, list[WAPIObjectDef]] = defaultdict(list)
        # wapitype -> [WAPIObjectDef sorted by version desc]

    def parse_all(self):
        """Parse all WAPI Python files."""
        pattern = os.path.join(
            self.nios_root, "products", "*", "server", "src", "wapi", "*.py"
        )
        wapi_files = glob.glob(pattern)
        parsed = 0
        errors = 0
        total_objects = 0

        for wapi_file in sorted(wapi_files):
            # Skip test files
            if '/test/' in wapi_file or wapi_file.endswith('_test.py'):
                continue
            try:
                objects = self._parse_file(wapi_file)
                total_objects += len(objects)
                parsed += 1
            except SyntaxError as e:
                errors += 1
            except Exception as e:
                errors += 1
                print(f"  Warning: Error parsing {wapi_file}: {e}",
                      file=sys.stderr)

        print(f"  Parsed {parsed} WAPI files ({errors} skipped), "
              f"found {total_objects} object definitions across "
              f"{len(self.objects)} WAPI types",
              file=sys.stderr)

    def _parse_file(self, filepath: str) -> list[WAPIObjectDef]:
        """Parse a single Python file and extract WAPIObject definitions.

        Uses a sequential walk through top-level statements (not ast.walk)
        so we can:
        1. Track abstract base classes with base_fields
        2. Associate module-level fields.extend() calls with the most
           recently defined WAPIObject class
        """
        with open(filepath, 'r', encoding='utf-8', errors='replace') as f:
            source = f.read()

        try:
            tree = ast.parse(source, filename=filepath)
        except SyntaxError:
            return []

        results = []
        last_wapi_obj = None  # Track most recently parsed WAPIObject

        # First pass: collect base_fields from abstract classes in this file
        self._current_file_base_fields = {}  # class_name -> list of WAPIField
        base_fields_map = self._current_file_base_fields
        for node in tree.body:
            if isinstance(node, ast.ClassDef):
                for stmt in node.body:
                    if (isinstance(stmt, ast.Assign) and stmt.targets
                            and isinstance(stmt.targets[0], ast.Name)
                            and stmt.targets[0].id == 'base_fields'
                            and isinstance(stmt.value, ast.List)):
                        fields = []
                        for elem in stmt.value.elts:
                            wf = self._parse_wf_field(elem)
                            if wf:
                                fields.append(wf)
                        if fields:
                            base_fields_map[node.name] = fields

        # Second pass: parse WAPIObject classes and module-level extends
        for node in tree.body:
            if isinstance(node, ast.ClassDef):
                obj = self._try_parse_wapi_class(node, filepath, source)
                if obj:
                    results.append(obj)
                    self.objects[obj.wapitype].append(obj)
                    last_wapi_obj = obj

            # Detect module-level fields.extend(...) after a class definition
            # Pattern: fields.extend(BaseClass.base_fields)
            # Pattern: fields.extend((WF(...), WF(...), ...))
            elif isinstance(node, ast.Expr) and isinstance(node.value, ast.Call):
                call = node.value
                if (isinstance(call.func, ast.Attribute)
                        and call.func.attr == 'extend'
                        and isinstance(call.func.value, ast.Name)
                        and call.func.value.id == 'fields'
                        and last_wapi_obj is not None
                        and call.args):
                    arg = call.args[0]

                    # Case 1: fields.extend(ClassName.base_fields)
                    if (isinstance(arg, ast.Attribute)
                            and arg.attr == 'base_fields'
                            and isinstance(arg.value, ast.Name)):
                        base_cls = arg.value.id
                        if base_cls in base_fields_map:
                            last_wapi_obj.fields.extend(
                                base_fields_map[base_cls]
                            )

                    # Case 2: fields.extend((WF(...), ...)) or fields.extend([...])
                    elif isinstance(arg, (ast.Tuple, ast.List)):
                        for elem in arg.elts:
                            wf = self._parse_wf_field(elem)
                            if wf:
                                last_wapi_obj.fields.append(wf)

        return results

    def _try_parse_wapi_class(self, node: ast.ClassDef, filepath: str,
                               source: str) -> Optional[WAPIObjectDef]:
        """Try to parse an AST ClassDef as a WAPIObject definition."""
        name = node.name

        # Match WAPIObject_xxx_N_N pattern
        m = re.match(r'^WAPIObject_([a-z\d]+?)_([\d_]+)$', name)
        if not m:
            return None

        wapitype_hint = m.group(1)
        version = m.group(2).replace('_', '.')

        # Extract class-level attributes
        obj = WAPIObjectDef(
            class_name=name,
            wapitype=wapitype_hint,  # May be overridden by explicit wapitype
            version=version,
            source_file=os.path.relpath(filepath, self.nios_root),
        )

        # Get parent class
        if node.bases:
            base = node.bases[0]
            if isinstance(base, ast.Name):
                obj.parent_class = base.id
            elif isinstance(base, ast.Attribute):
                obj.parent_class = ast.dump(base)

        # Get docstring
        if (node.body and isinstance(node.body[0], ast.Expr)
                and isinstance(node.body[0].value, ast.Constant) and isinstance(node.body[0].value.value, str)):
            val = node.body[0].value
            obj.doc = val.value
            # Truncate doc to first paragraph (before %CUT%)
            if obj.doc and '%CUT%' in obj.doc:
                obj.doc = obj.doc.split('%CUT%')[0].strip()

        # Parse class body for attributes, fields, and fields.extend() calls
        for stmt in node.body:
            if isinstance(stmt, ast.Assign):
                self._parse_class_assign(stmt, obj)
            elif isinstance(stmt, ast.FunctionDef):
                if stmt.name == 'modify_fields':
                    obj.has_modify_fields = True
            elif isinstance(stmt, ast.Expr) and isinstance(stmt.value, ast.Call):
                # Handle fields.extend(...) inside the class body
                self._parse_fields_extend(stmt.value, obj)

        return obj

    def _parse_class_assign(self, stmt: ast.Assign, obj: WAPIObjectDef):
        """Parse a class-level assignment like `wapitype = 'record:a'`."""
        if not stmt.targets or not isinstance(stmt.targets[0], ast.Name):
            return

        attr_name = stmt.targets[0].id
        value = stmt.value

        if attr_name == 'wapitype' and isinstance(value, ast.Constant):
            obj.wapitype = value.value

        elif attr_name == 'cmdclass' and isinstance(value, ast.Attribute):
            obj.cmdclass = self._get_dotted_name(value)

        elif attr_name == 'is_valid_version':
            if isinstance(value, ast.Constant):
                obj.is_valid_version = bool(value.value)
            elif isinstance(value, ast.NameConstant):
                obj.is_valid_version = bool(value.value)

        elif attr_name == 'schedulable':
            obj.schedulable = self._get_bool(value, False)

        elif attr_name == 'globalsearch':
            obj.globalsearch = self._get_bool(value, False)

        elif attr_name == 'csvexport':
            obj.csvexport = self._get_bool(value, False)

        elif attr_name == 'restrictions' and isinstance(value, ast.Tuple):
            obj.restrictions = tuple(
                e.value for e in value.elts
                if isinstance(e, ast.Constant) and isinstance(e.value, str)
            )

        elif attr_name == 'fields' and isinstance(value, ast.List):
            for elem in value.elts:
                wf = self._parse_wf_field(elem)
                if wf:
                    obj.fields.append(wf)

        elif attr_name == 'fields' and isinstance(value, ast.BinOp):
            # Handle: fields = BaseClass.base_fields + [WF(...), ...]
            self._parse_fields_binop(value, obj)

    def _parse_fields_extend(self, call: ast.Call, obj: WAPIObjectDef):
        """Parse a fields.extend(...) call inside or after a class definition."""
        if not (isinstance(call.func, ast.Attribute)
                and call.func.attr == 'extend'
                and isinstance(call.func.value, ast.Name)
                and call.func.value.id == 'fields'
                and call.args):
            return

        arg = call.args[0]

        # Case 1: fields.extend(ClassName.base_fields)
        if (isinstance(arg, ast.Attribute)
                and arg.attr == 'base_fields'
                and isinstance(arg.value, ast.Name)):
            base_cls = arg.value.id
            base_fields = self._current_file_base_fields.get(base_cls, [])
            if base_fields:
                obj.fields.extend(base_fields)

        # Case 2: fields.extend((WF(...), ...)) or fields.extend([...])
        elif isinstance(arg, (ast.Tuple, ast.List)):
            for elem in arg.elts:
                wf = self._parse_wf_field(elem)
                if wf:
                    obj.fields.append(wf)

    def _parse_fields_binop(self, node: ast.BinOp, obj: WAPIObjectDef):
        """Parse fields = X.base_fields + [...] pattern."""
        if not isinstance(node.op, ast.Add):
            return

        # Left side: ClassName.base_fields
        if (isinstance(node.left, ast.Attribute)
                and node.left.attr == 'base_fields'
                and isinstance(node.left.value, ast.Name)):
            base_cls = node.left.value.id
            base_fields = self._current_file_base_fields.get(base_cls, [])
            if base_fields:
                obj.fields.extend(base_fields)

        # Right side: [WF(...), ...]
        if isinstance(node.right, ast.List):
            for elem in node.right.elts:
                wf = self._parse_wf_field(elem)
                if wf:
                    obj.fields.append(wf)

    def _parse_wf_field(self, node: ast.expr) -> Optional[WAPIField]:
        """Parse a WF*(...) call expression into a WAPIField."""
        if not isinstance(node, ast.Call):
            return None

        # Get the WF class name
        func = node.func
        if isinstance(func, ast.Name):
            wf_class = func.id
        elif isinstance(func, ast.Attribute):
            wf_class = func.attr
        else:
            return None

        # Only process known WF classes
        if not wf_class.startswith('WF'):
            return None

        # WFFuncCall fields are function calls, not data fields
        is_func_call = wf_class == 'WFFuncCall'

        # Special case: WFExtensibleAttributes() has no positional args —
        # its wname is hardcoded to 'extattrs' in __init__
        if wf_class in ('WFExtensibleAttributes',
                         'WFNotInheritedExtensibleAttributes'):
            wf = WAPIField(wname='extattrs', wf_class=wf_class)
            wf.is_ext_attrs = True
            wf.iname = 'extensible_attributes'
            wf.purpose = 'rwu'
            wf.default_value = {}
            # Parse any keyword overrides
            for kw in node.keywords:
                if kw.arg is not None:
                    self._parse_wf_keyword(kw, wf)
            return wf

        # Get the wname (first positional argument)
        if not node.args:
            return None

        wname_node = node.args[0]
        if isinstance(wname_node, ast.Constant) and isinstance(wname_node.value, str):
            wname = wname_node.value
        else:
            return None

        wf = WAPIField(wname=wname, wf_class=wf_class)
        wf.is_func_call = is_func_call

        # Get doc (second positional argument)
        if len(node.args) >= 2:
            doc_node = node.args[1]
            if isinstance(doc_node, ast.Constant) and isinstance(doc_node.value, str):
                wf.doc = doc_node.value.strip()

        # Infer type from WF class (with suffix-based fallback)
        wf_type = _infer_wf_type(wf_class)
        wf.is_array = '[]' in wf_type
        wf.is_struct = 'struct' in wf_type
        wf.is_subobj = 'object_ref' in wf_type
        wf.is_ext_attrs = wf_class in (
            'WFExtensibleAttributes', 'WFNotInheritedExtensibleAttributes'
        )
        wf.is_obsolete = wf_type == 'obsolete'

        # Parse keyword arguments
        for kw in node.keywords:
            if kw.arg is None:
                continue  # **kwargs
            self._parse_wf_keyword(kw, wf)

        # Default iname = wname
        if wf.iname is None:
            wf.iname = wf.wname

        # Apply forced purpose overrides for domain-specific WF subclasses
        # that hardcode purpose in their __init__ (not visible to AST parser).
        if wf_class in WF_CLASS_FORCED_PURPOSE:
            wf.purpose = WF_CLASS_FORCED_PURPOSE[wf_class]

        # The WAPI runtime automatically adds 's' to purpose when search
        # is set on a field. Replicate that here so supports matches.
        if wf.search and 's' not in wf.purpose:
            wf.purpose += 's'

        return wf

    def _parse_wf_keyword(self, kw: ast.keyword, wf: WAPIField):
        """Parse a keyword argument to a WF constructor."""
        name = kw.arg
        value = kw.value

        if name == 'iname':
            wf.iname = self._get_str(value)

        elif name == 'purpose':
            wf.purpose = self._get_str(value) or 'rwu'

        elif name == 'create':
            wf.create = self._get_bool(value, False)

        elif name == 'create_only':
            wf.create_only = self._get_bool(value, False)

        elif name == 'create_default':
            wf.create_default = self._get_str(value)

        elif name == 'use_flag':
            wf.use_flag = self._get_str(value)

        elif name == 'search':
            wf.search = self._get_str(value)

        elif name == 'std_field':
            wf.std_field = self._get_bool(value, False)

        elif name == 'nillable':
            wf.nillable = self._get_bool(value, False)

        elif name == 'strip':
            wf.strip = self._get_bool(value, False)

        elif name == 'scope_of_delegation':
            wf.scope_of_delegation = self._get_bool(value, False)

        elif name == 'inheritance':
            wf.inheritance = self._get_bool(value, False)

        elif name == 'typecls':
            wf.typecls = self._get_name(value)

        elif name == 'enumlist':
            wf.enumlist = self._get_list_of_strings(value)

        elif name == 'default':
            self._parse_default(value, wf)

        elif name == 'i2w':
            wf.i2w_func = self._get_func_name(value)

        elif name == 'w2i':
            wf.w2i_func = self._get_func_name(value)

        elif name == 'in_struct':
            wf.in_struct = self._get_str(value)

    def _parse_default(self, node: ast.expr, wf: WAPIField):
        """Parse a default={'value': X} or default={'docstring': X} dict."""
        if not isinstance(node, ast.Dict):
            return

        for key, val in zip(node.keys, node.values):
            key_str = self._get_str(key)
            if key_str == 'value':
                if isinstance(val, ast.Constant):
                    wf.default_value = val.value
                elif isinstance(val, ast.List):
                    wf.default_value = []
                elif isinstance(val, ast.Dict):
                    wf.default_value = {}
            elif key_str == 'docstring':
                wf.default_docstring = self._get_str(val)

    # --- AST helper methods ---

    @staticmethod
    def _get_str(node) -> Optional[str]:
        if isinstance(node, ast.Constant) and isinstance(node.value, str):
            return node.value
        return None

    @staticmethod
    def _get_bool(node, default=False) -> bool:
        if isinstance(node, ast.Constant):
            return bool(node.value)
        if isinstance(node, ast.NameConstant):
            return bool(node.value)
        if isinstance(node, ast.Name):
            return node.id == 'True'
        return default

    @staticmethod
    def _get_name(node) -> Optional[str]:
        if isinstance(node, ast.Name):
            return node.id
        if isinstance(node, ast.Attribute):
            return node.attr
        return None

    @staticmethod
    def _get_dotted_name(node) -> Optional[str]:
        """Get a dotted name like 'a_record.ARecordCmd'."""
        parts = []
        while isinstance(node, ast.Attribute):
            parts.append(node.attr)
            node = node.value
        if isinstance(node, ast.Name):
            parts.append(node.id)
        return '.'.join(reversed(parts)) if parts else None

    @staticmethod
    def _get_list_of_strings(node) -> Optional[list]:
        if isinstance(node, ast.List):
            result = []
            for elem in node.elts:
                if isinstance(elem, ast.Constant) and isinstance(elem.value, str):
                    result.append(elem.value)
            return result if result else None
        return None

    @staticmethod
    def _get_func_name(node) -> Optional[str]:
        """Get a function reference name from AST."""
        if isinstance(node, ast.Name):
            return node.id
        if isinstance(node, ast.Attribute):
            # e.g. self.some_func or module.func
            parts = []
            while isinstance(node, ast.Attribute):
                parts.append(node.attr)
                node = node.value
            if isinstance(node, ast.Name):
                parts.append(node.id)
            return '.'.join(reversed(parts))
        if isinstance(node, ast.Call):
            # e.g. i2w_wrap(str.lower) — get the call representation
            func_name = WAPIASTParser._get_func_name(node.func)
            if func_name and node.args:
                arg = WAPIASTParser._get_func_name(node.args[0])
                if arg:
                    return f"{func_name}({arg})"
            return func_name
        return None

    def get_latest_version(self, wapitype: str) -> Optional[WAPIObjectDef]:
        """Get the latest (highest) valid version of a WAPI object type."""
        versions = self.objects.get(wapitype, [])
        if not versions:
            return None

        # Sort by parsed version (descending), prefer is_valid_version=True
        def version_key(obj):
            parts = obj.version.split('.')
            try:
                nums = tuple(int(p) for p in parts)
            except ValueError:
                nums = (0,)
            return (obj.is_valid_version, nums)

        return max(versions, key=version_key)


# ---------------------------------------------------------------------------
# Schema Classifier
# ---------------------------------------------------------------------------

class FieldClassifier:
    """
    Classify WAPI fields into mutability/computation categories by
    cross-referencing WF field attributes, RTXML member attributes,
    and i2w function names.
    """

    def classify(self, wf: WAPIField,
                 rtxml_member: Optional[RTXMLMember] = None) -> SchemaField:
        """Classify a WAPIField into a SchemaField with mutability metadata."""

        # Determine base type
        field_type = self._infer_type(wf)

        sf = SchemaField(
            name=wf.wname,
            type=field_type,
            wf_class=wf.wf_class,
            is_array=wf.is_array,
            is_struct=wf.is_struct,
            standard_field=wf.std_field,
            supports=wf.purpose,
            searchable_by=wf.search,
            create_required=wf.create,
            create_only=wf.create_only,
            nillable=wf.nillable,
            default_value=wf.default_value,
            has_server_default=wf.default_value is not None,
            enum_values=wf.enumlist,
            use_flag=wf.use_flag,
            iname=wf.iname if wf.iname != wf.wname else None,
            doc=wf.doc,
            i2w_function=wf.i2w_func,
            w2i_function=wf.w2i_func,
        )

        # Cross-reference RTXML
        if rtxml_member:
            sf.rtxml_key_type = rtxml_member.key_type
            sf.rtxml_syntax = rtxml_member.syntax_name
            sf.rtxml_unique_id = rtxml_member.unique_id
            sf.rtxml_synthetic_direct = rtxml_member.synthetic_direct
            # RTXML default-value is another source of server defaults
            if rtxml_member.default_value is not None:
                sf.has_server_default = True

        # --- Classification logic (mirrors the ADR algorithm) ---

        # 1. RTXML-level signals
        if rtxml_member:
            if rtxml_member.unique_id:
                sf.mutability = 'read_only'
                sf.computation = 'server_computed'
                sf.stable_after_create = True
                self._set_plan_modifier(sf)
                return sf

            if rtxml_member.key_type == 'key':
                sf.mutability = 'immutable'
                sf.computation = 'stored'
                sf.stable_after_create = True

            if rtxml_member.version_id:
                sf.mutability = 'read_only'
                sf.computation = 'server_computed'
                sf.stable_after_create = False
                self._set_plan_modifier(sf)
                return sf

            if rtxml_member.internal_ro:
                sf.mutability = 'read_only'
                sf.computation = 'server_computed'
                sf.stable_after_create = False
                self._set_plan_modifier(sf)
                return sf

            if rtxml_member.synthetic_direct:
                sf.mutability = 'read_only'
                sf.computation = 'derived'
                sf.stable_after_create = False
                self._set_plan_modifier(sf)
                return sf

        # 2. WAPI WF-level signals

        # Fully read-only fields
        purpose = wf.purpose or 'rwu'
        is_read_only = (purpose in ('r', 'rs'))
        is_writable = 'w' in purpose
        is_updatable = 'u' in purpose

        if is_read_only:
            sf.mutability = 'read_only'

            if wf.i2w_func:
                sf.computation = 'derived'
                sf.stable_after_create = False
                # Check client-derivable
                self._check_client_derivable(wf, sf)
            elif wf.wname in ('_ref', 'creation_time'):
                sf.computation = 'server_computed'
                sf.stable_after_create = True
            else:
                sf.computation = 'server_computed'
                # Some read-only fields are immutable after creation
                if wf.wname in ('zone', 'network_view', 'creator',
                                'creation_time'):
                    sf.stable_after_create = True
                else:
                    sf.stable_after_create = False

            # Check client-derivable iname patterns
            if wf.iname and wf.iname in CLIENT_DERIVABLE_INAMES:
                sf.computation = 'client_derivable'
                sf.derivation_function = CLIENT_DERIVABLE_INAMES[wf.iname]
                sf.stable_after_create = True
                # Infer derived_from
                if sf.derivation_function == 'punycode':
                    # dns_fqdn derives from fqdn/name
                    sf.derived_from = ['name']
                    if wf.wname == 'dns_fqdn':# dns_name
                        print("  Warning: inferring dns_fqdn derives from name (could be fqdn) {0}".format(wf.wname),
                              file=sys.stderr)
                        sf.derived_from = ['fqdn']

            # Check dns_* naming pattern: a read-only field named dns_X
            # where X is another field on the same object is the punycode
            # form of that field. E.g. dns_target_name → punycode(target_name),
            # dns_canonical → punycode(canonical), dns_name → punycode(name).
            if (sf.computation != 'client_derivable'
                    and wf.wname.startswith('dns_')):
                source_candidate = wf.wname[4:]  # strip 'dns_' prefix
                sf.computation = 'client_derivable'
                sf.derivation_function = 'punycode'
                sf.derived_from = [source_candidate]
                sf.stable_after_create = True

            self._set_plan_modifier(sf)
            return sf

        # Create-only fields
        if wf.create_only:
            sf.mutability = 'immutable'
            sf.computation = 'stored'
            sf.stable_after_create = True
            self._set_plan_modifier(sf)
            return sf

        # Writable but not updatable → immutable after creation
        if is_writable and not is_updatable:
            sf.mutability = 'immutable'
            sf.computation = 'stored'
            sf.stable_after_create = True
            self._set_plan_modifier(sf)
            return sf

        # Inherited fields (use_flag pattern)
        if wf.use_flag:
            sf.computation = 'inherited'
            sf.stable_after_create = False
            self._set_plan_modifier(sf)
            return sf

        # Inheritance flag
        if wf.inheritance:
            sf.computation = 'inherited'
            sf.stable_after_create = False
            self._set_plan_modifier(sf)
            return sf

        # SubObjArray fields (joined from child objects)
        if wf.is_subobj and wf.is_array:
            sf.computation = 'derived'
            sf.derived_from = [wf.iname or wf.wname]
            sf.stable_after_create = False
            self._set_plan_modifier(sf)
            return sf

        # ExtensibleAttributes
        if wf.is_ext_attrs:
            sf.computation = 'derived'
            sf.derived_from = ['_ea_store']
            sf.stable_after_create = False
            self._set_plan_modifier(sf)
            return sf

        # Fields with i2w converters (derived even if writable)
        if wf.i2w_func:
            self._check_client_derivable(wf, sf)
            if sf.computation != 'client_derivable':
                sf.computation = 'derived'
                sf.stable_after_create = False

        # Default: mutable/stored
        self._set_plan_modifier(sf)
        return sf

    def _check_client_derivable(self, wf: WAPIField, sf: SchemaField):
        """Check if an i2w function is a client-derivable transformation."""
        if not wf.i2w_func:
            return

        if wf.i2w_func in CLIENT_DERIVABLE_TRANSFORMS:
            sf.computation = 'client_derivable'
            sf.derivation_function = CLIENT_DERIVABLE_TRANSFORMS[wf.i2w_func]
            sf.stable_after_create = True

            # Infer derived_from based on the transform type
            if sf.derivation_function == 'boolean_invert':
                # The source is typically the iname with inverted naming
                sf.derived_from = [wf.iname or wf.wname]
            elif sf.derivation_function == 'lowercase':
                sf.derived_from = [wf.wname]

    def _set_plan_modifier(self, sf: SchemaField):
        """Set the reconciliation_action based on classification.

        Generic reconciliation actions (tool-agnostic):
          suppress_diff      — field is stable after create; tools can cache/suppress
          compute_on_diff(X) — field is client-derivable via transform X; tools can compute
          always_refresh     — field may change server-side; tools must re-read every cycle
        """
        if sf.computation == 'client_derivable' and sf.stable_after_create:
            func = sf.derivation_function or 'unknown'
            sf.reconciliation_action = f"compute_on_diff({func})"
        elif sf.stable_after_create and sf.mutability in ('immutable', 'read_only'):
            sf.reconciliation_action = "suppress_diff"
        else:
            sf.reconciliation_action = "always_refresh"

    def _infer_type(self, wf: WAPIField) -> str:
        """Infer the schema type from WF class and typecls."""
        # Use the suffix-based fallback for domain-specific WF subclasses
        wf_type = _infer_wf_type(wf.wf_class)

        # Refine based on typecls
        if wf.typecls:
            typecls_map = {
                'WTypeBool': 'bool',
                'WTypeUInt': 'uint',
                'WTypeInt': 'int',
                'WTypeTS': 'timestamp',
                'WTypeIPv4Address': 'string',
                'WTypeIPv6Address': 'string',
                'WTypeIPAddress': 'string',
                'WTypeIPv4Cidr': 'string',
                'WTypeIPv6Cidr': 'string',
            }
            if wf.typecls in typecls_map:
                base_type = typecls_map[wf.typecls]
                if wf.is_array:
                    return f"{base_type}[]"
                return base_type

        return wf_type


# ---------------------------------------------------------------------------
# pyabs → RTXML type resolver
# ---------------------------------------------------------------------------

class PyabsResolver:
    """
    Resolve the RTXML structure type from a WAPI object's cmdclass.

    Scans pyabs files to find:
      class ARecord(BaseAddressRecord):
          type = 'dns.bind_a'

    Then resolves dns.bind_a → .com.infoblox.dns.bind_a
    """

    def __init__(self, nios_root: str):
        self.nios_root = nios_root
        # Map: pyabs class name → RTXML type (e.g., 'ARecord' → 'dns.bind_a')
        self.class_to_rtxml: dict[str, str] = {}
        # Map: cmdclass ref → pyabs class (e.g. 'a_record.ARecordCmd' → 'ARecord')
        self.cmdclass_to_pyabs: dict[str, str] = {}

    def scan(self):
        """Scan pyabs files to build the class → RTXML mapping."""
        pattern = os.path.join(
            self.nios_root, "products", "*", "server", "src", "pyabs", "*.py"
        )
        pyabs_files = glob.glob(pattern)
        count = 0

        for filepath in sorted(pyabs_files):
            try:
                self._scan_file(filepath)
                count += 1
            except Exception:
                pass

        print(f"  Scanned {count} pyabs files, found {len(self.class_to_rtxml)} "
              f"type mappings", file=sys.stderr)

    def _scan_file(self, filepath: str):
        """Scan a pyabs file for type = 'xxx.yyy' assignments."""
        with open(filepath, 'r', encoding='utf-8', errors='replace') as f:
            source = f.read()

        try:
            tree = ast.parse(source, filename=filepath)
        except SyntaxError:
            return

        for node in ast.walk(tree):
            if not isinstance(node, ast.ClassDef):
                continue

            class_name = node.name
            for stmt in node.body:
                if not isinstance(stmt, ast.Assign):
                    continue
                for target in stmt.targets:
                    if not isinstance(target, ast.Name):
                        continue
                    if target.id == 'type' and isinstance(stmt.value, ast.Constant):
                        val = stmt.value.value
                        if isinstance(val, str) and '.' in val:
                            self.class_to_rtxml[class_name] = val

    def resolve_rtxml_type(self, cmdclass: Optional[str]) -> Optional[str]:
        """
        Resolve a cmdclass reference to an RTXML full type.

        cmdclass is like 'a_record.ARecordCmd' → look up ARecordCmd's parent
        class which would be ARecord, find its type = 'dns.bind_a',
        then expand to '.com.infoblox.dns.bind_a'.
        """
        if not cmdclass:
            return None

        # The cmdclass is usually like 'a_record.ARecordCmd'
        # The base pyabs class is typically named without 'Cmd' suffix
        parts = cmdclass.split('.')
        cmd_name = parts[-1]  # e.g. 'ARecordCmd'

        # Try removing 'Cmd' suffix to get the pyabs class name
        base_name = cmd_name.replace('Cmd', '')

        if base_name in self.class_to_rtxml:
            short_type = self.class_to_rtxml[base_name]
            # Expand: dns.bind_a → .com.infoblox.dns.bind_a
            return f".com.infoblox.{short_type}"

        return None


# ---------------------------------------------------------------------------
# Schema Generator
# ---------------------------------------------------------------------------

class SchemaGenerator:
    """Combine RTXML, WAPI AST, and pyabs data into final JSON schemas."""

    def __init__(self, rtxml: RTXMLParser, wapi: WAPIASTParser,
                 pyabs: PyabsResolver, classifier: FieldClassifier):
        self.rtxml = rtxml
        self.wapi = wapi
        self.pyabs = pyabs
        self.classifier = classifier

    def generate(self, wapitype: str) -> Optional[dict]:
        """Generate schema JSON for a single WAPI object type."""
        obj = self.wapi.get_latest_version(wapitype)
        if not obj:
            return None

        # Resolve RTXML type
        rtxml_type = self.pyabs.resolve_rtxml_type(obj.cmdclass)
        rtxml_struct = self.rtxml.lookup(rtxml_type) if rtxml_type else None

        # Build field schemas
        fields = {}
        use_flag_companions = []  # Track use_* companion fields to generate

        for wf in obj.fields:
            # Skip function call fields — they're not data attributes
            if wf.is_func_call:
                continue

            # Skip WFFlatStruct fields — these are internal containers whose
            # sub-fields are flattened into the parent object by the runtime.
            # They don't appear as top-level WAPI fields.
            if wf.wf_class == 'WFFlatStruct':
                continue

            # Look up RTXML member
            rtxml_member = None
            if rtxml_struct:
                # Try iname first (iname maps to RTXML member name)
                lookup_name = wf.iname or wf.wname
                rtxml_member = rtxml_struct.members.get(lookup_name)

            sf = self.classifier.classify(wf, rtxml_member)
            fields[sf.name] = self._schema_field_to_dict(sf)

            # Track use_flag companions — the runtime auto-generates a
            # boolean use_* field for each field that has use_flag set.
            if wf.use_flag:
                flag_name = wf.use_flag.split('=')[0]
                if flag_name not in fields:
                    use_flag_companions.append(flag_name)

        # Generate use_* companion fields that weren't explicitly declared
        for flag_name in use_flag_companions:
            if flag_name not in fields:
                fields[flag_name] = {
                    "name": flag_name,
                    "type": "bool",
                    "wf_class": "WFBool",
                    "supports": "rwu",
                    "mutability": "mutable",
                    "computation": "stored",
                }

        # Determine CRUD restrictions
        restrictions = self._get_restrictions(obj)

        # Collect all versions
        all_versions = []
        for ver_obj in self.wapi.objects.get(wapitype, []):
            if ver_obj.is_valid_version:
                all_versions.append(ver_obj.version)
        all_versions.sort(key=lambda v: tuple(int(p) for p in v.split('.')
                                                if p.isdigit()),
                          reverse=True)

        return {
            "type": obj.wapitype,
            "version": obj.version,
            "all_versions": all_versions,
            "source_file": obj.source_file,
            "cmdclass": obj.cmdclass,
            "rtxml_type": rtxml_type,
            "restrictions": restrictions,
            "schedulable": obj.schedulable,
            "globalsearch": obj.globalsearch,
            "csvexport": obj.csvexport,
            "doc": obj.doc,
            "field_count": len(fields),
            "fields": fields,
        }

    def generate_all(self, object_filter: Optional[list] = None) -> dict:
        """Generate schemas for all (or filtered) WAPI object types."""
        result = {}
        types = object_filter if object_filter else sorted(self.wapi.objects.keys())

        for wapitype in types:
            schema = self.generate(wapitype)
            if schema:
                result[wapitype] = schema

        return result

    def _get_restrictions(self, obj: WAPIObjectDef) -> list:
        """Determine which CRUD operations are supported."""
        supported = ['create', 'read', 'update', 'delete']
        if obj.restrictions:
            # restrictions list items that are NOT allowed
            supported = [op for op in supported if op not in obj.restrictions]
        ops = []
        if 'read' in supported:
            ops.append('read')
        if 'create' in supported:
            ops.append('create')
        if 'update' in supported:
            ops.append('update')
        if 'delete' in supported:
            ops.append('delete')
        if obj.schedulable:
            ops.append('scheduling')
        if obj.csvexport:
            ops.append('csv')
        return ops

    @staticmethod
    def _schema_field_to_dict(sf: SchemaField) -> dict:
        """Convert a SchemaField to a dictionary, omitting None values."""
        d = asdict(sf)
        # Remove None and False values for cleaner output
        return {k: v for k, v in d.items()
                if v is not None and v is not False and v != ''}


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def parse_args():
    parser = argparse.ArgumentParser(
        description='Extract WAPI schema from NIOS codebase (static AST analysis)',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Full extraction
  %(prog)s --nios-root /path/to/nios --output-dir ./schemas

  # Specific objects only
  %(prog)s --nios-root /path/to/nios --output-dir ./schemas \\
      --objects record:a,network,zone_auth

  # Summary mode (no file output)
  %(prog)s --nios-root /path/to/nios --summary

  # Single object to stdout
  %(prog)s --nios-root /path/to/nios --objects record:a --stdout
        """
    )
    parser.add_argument('--nios-root', required=True,
                        help='Path to the NIOS source root')
    parser.add_argument('--output-dir', default=None,
                        help='Directory to write JSON schema files')
    parser.add_argument('--objects', default=None,
                        help='Comma-separated list of WAPI object types to extract')
    parser.add_argument('--summary', action='store_true',
                        help='Print summary statistics only')
    parser.add_argument('--stdout', action='store_true',
                        help='Print JSON to stdout instead of files')
    parser.add_argument('--pretty', action='store_true', default=True,
                        help='Pretty-print JSON output (default: true)')
    parser.add_argument('--compact', action='store_true',
                        help='Compact JSON output')
    return parser.parse_args()


def main():
    args = parse_args()

    nios_root = os.path.abspath(args.nios_root)
    if not os.path.isdir(nios_root):
        print(f"Error: NIOS root directory not found: {nios_root}",
              file=sys.stderr)
        sys.exit(1)

    print(f"WAPI Schema Extractor", file=sys.stderr)
    print(f"  NIOS root: {nios_root}", file=sys.stderr)
    print(file=sys.stderr)

    # Phase 1: Parse RTXML
    print("Phase 1: Parsing RTXML XML files...", file=sys.stderr)
    rtxml = RTXMLParser(nios_root)
    rtxml.parse_all()
    print(f"  Found {len(rtxml.structures)} RTXML structures", file=sys.stderr)
    print(file=sys.stderr)

    # Phase 2: Scan pyabs for type mappings
    print("Phase 2: Scanning pyabs type mappings...", file=sys.stderr)
    pyabs = PyabsResolver(nios_root)
    pyabs.scan()
    print(file=sys.stderr)

    # Phase 3: Parse WAPI definitions
    print("Phase 3: Parsing WAPI object definitions...", file=sys.stderr)
    wapi = WAPIASTParser(nios_root)
    wapi.parse_all()
    print(file=sys.stderr)

    # Phase 4: Generate schemas
    print("Phase 4: Generating schemas...", file=sys.stderr)
    classifier = FieldClassifier()
    generator = SchemaGenerator(rtxml, wapi, pyabs, classifier)

    object_filter = None
    if args.objects:
        object_filter = [o.strip() for o in args.objects.split(',')]

    schemas = generator.generate_all(object_filter)
    print(f"  Generated schemas for {len(schemas)} WAPI object types",
          file=sys.stderr)
    print(file=sys.stderr)

    # Summary mode
    if args.summary:
        _print_summary(wapi, schemas)
        return

    # Build the full output document
    output = {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "generator": "wapi_schema_extractor",
        "generator_version": "1.0.0",
        "extraction_method": "static_ast_analysis",
        "total_objects": len(schemas),
        "objects": schemas,
    }

    indent = 2 if (args.pretty and not args.compact) else None

    if args.stdout:
        print(json.dumps(output, indent=indent, default=str))
        return

    if not args.output_dir:
        print("Error: --output-dir or --stdout required", file=sys.stderr)
        sys.exit(1)

    output_dir = os.path.abspath(args.output_dir)
    os.makedirs(output_dir, exist_ok=True)

    # Write combined schema
    combined_path = os.path.join(output_dir, "wapi_schema_all.json")
    with open(combined_path, 'w') as f:
        json.dump(output, f, indent=indent, default=str)
    print(f"  Wrote combined schema: {combined_path}", file=sys.stderr)

    # Write per-object schemas
    per_obj_dir = os.path.join(output_dir, "objects")
    os.makedirs(per_obj_dir, exist_ok=True)

    for wapitype, schema in schemas.items():
        # Sanitize filename: record:a → record_a.json
        filename = wapitype.replace(':', '_').replace('.', '_') + '.json'
        filepath = os.path.join(per_obj_dir, filename)
        with open(filepath, 'w') as f:
            json.dump(schema, f, indent=indent, default=str)

    print(f"  Wrote {len(schemas)} per-object schemas to: {per_obj_dir}",
          file=sys.stderr)

    # Write type index
    index = {
        "total_types": len(schemas),
        "types": {}
    }
    for wapitype, schema in sorted(schemas.items()):
        index["types"][wapitype] = {
            "version": schema["version"],
            "field_count": schema["field_count"],
            "restrictions": schema["restrictions"],
            "source_file": schema["source_file"],
        }

    index_path = os.path.join(output_dir, "wapi_type_index.json")
    with open(index_path, 'w') as f:
        json.dump(index, f, indent=indent, default=str)
    print(f"  Wrote type index: {index_path}", file=sys.stderr)

    print(file=sys.stderr)
    print("Done.", file=sys.stderr)


def _print_summary(wapi: WAPIASTParser, schemas: dict):
    """Print a summary of extracted data."""
    print("=" * 70)
    print("WAPI SCHEMA EXTRACTION SUMMARY")
    print("=" * 70)

    # Overall stats
    total_versions = sum(len(v) for v in wapi.objects.values())
    print(f"\nTotal WAPI types:                {len(wapi.objects)}")
    print(f"Total versioned definitions:     {total_versions}")
    print(f"Schemas generated (latest ver):  {len(schemas)}")

    # Field stats
    total_fields = 0
    type_counts = defaultdict(int)
    mutability_counts = defaultdict(int)
    computation_counts = defaultdict(int)
    use_state_safe_count = 0
    use_state_unsafe_count = 0
    has_use_flag = 0
    has_enum = 0

    for schema in schemas.values():
        for fname, fdata in schema.get('fields', {}).items():
            total_fields += 1
            type_counts[fdata.get('type', 'unknown')] += 1
            mutability_counts[fdata.get('mutability', 'unknown')] += 1
            computation_counts[fdata.get('computation', 'unknown')] += 1
            if fdata.get('stable_after_create'):
                use_state_safe_count += 1
            if fdata.get('use_flag'):
                has_use_flag += 1
            if fdata.get('enum_values'):
                has_enum += 1

    print(f"\nTotal fields across all types:   {total_fields}")
    print(f"Fields with use_flag:            {has_use_flag}")
    print(f"Fields with enum_values:         {has_enum}")

    print(f"\n--- Field Types ---")
    for t, c in sorted(type_counts.items(), key=lambda x: -x[1]):
        print(f"  {t:25s}  {c}")

    print(f"\n--- Mutability ---")
    for m, c in sorted(mutability_counts.items(), key=lambda x: -x[1]):
        print(f"  {m:25s}  {c}")

    print(f"\n--- Computation ---")
    for comp, c in sorted(computation_counts.items(), key=lambda x: -x[1]):
        print(f"  {comp:25s}  {c}")

    print(f"\n--- UseStateForUnknown Safety ---")
    print(f"  Safe (can use):                {use_state_safe_count}")
    print(f"  Unsafe (must NOT use):         {total_fields - use_state_safe_count}")

    # Terraform-relevant objects (have CRUD)
    tf_candidates = 0
    for schema in schemas.values():
        ops = schema.get('restrictions', [])
        if 'create' in ops and 'read' in ops:
            tf_candidates += 1

    print(f"\n--- Terraform Relevance ---")
    print(f"  Types with CRUD (TF candidates): {tf_candidates}")
    print(f"  Read-only types:                 {len(schemas) - tf_candidates}")

    # Show top 10 types by field count
    print(f"\n--- Top 10 Types by Field Count ---")
    by_fields = sorted(schemas.items(),
                       key=lambda x: x[1].get('field_count', 0),
                       reverse=True)[:10]
    for wapitype, schema in by_fields:
        ops = ','.join(schema.get('restrictions', []))
        print(f"  {wapitype:35s}  {schema['field_count']:3d} fields  "
              f"v{schema['version']}  [{ops}]")


if __name__ == '__main__':
    main()
