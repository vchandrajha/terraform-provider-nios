#!/usr/bin/env python3
"""
Adds a PutExpand function to all model files that have:
1. A *Model struct with an Expand method
2. A *ResourceSchemaAttributes map variable

The PutExpand function uses reflection to iterate over schema attributes
and removes fields with empty/zero values from the API struct before PUT requests,
respecting schema defaults.

Usage:
    python3 scripts/add_put_expand.py <repo_path> [--dry-run]

Example:
    python3 scripts/add_put_expand.py /Users/vchandra/wspace/terraform-provider-nios
    python3 scripts/add_put_expand.py /Users/vchandra/wspace/terraform-provider-nios --dry-run
"""

import argparse
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path


BASE_DIR: Path = None  # Set in main() from CLI argument


def run_gofmt(filepath: Path):
    """Run goimports (preferred) or gofmt on a file to ensure proper formatting."""
    formatter = shutil.which("goimports") or shutil.which("gofmt")
    if not formatter:
        return
    try:
        subprocess.run([formatter, "-w", str(filepath)],
                       capture_output=True, text=True, timeout=10)
    except Exception:
        pass


def find_model_info(content: str) -> dict | None:
    """Extract model name, API type, schema attributes var, and package alias from a model file."""

    # Find Expand method: func (m *FooModel) Expand(...) *pkg.Foo {
    expand_match = re.search(
        r'func \(m \*(\w+Model)\) Expand\(.*?\) \*(\w+)\.(\w+)\s*\{',
        content
    )
    if not expand_match:
        return None

    model_name = expand_match.group(1)      # e.g. ZoneAuthModel
    pkg_alias = expand_match.group(2)        # e.g. dns
    api_type = expand_match.group(3)         # e.g. ZoneAuth

    # Derive the base name (model name without "Model" suffix)
    base_name = model_name.replace("Model", "")  # e.g. ZoneAuth

    # Find ResourceSchemaAttributes variable
    schema_var = f"{base_name}ResourceSchemaAttributes"
    if schema_var not in content:
        return None

    return {
        "model_name": model_name,
        "base_name": base_name,
        "pkg_alias": pkg_alias,
        "api_type": api_type,
        "schema_var": schema_var,
    }


def has_put_expand(content: str) -> bool:
    """Check if PutExpand already exists."""
    return "PutExpand" in content


def needs_import(content: str, import_str: str) -> bool:
    """Check if an import is missing."""
    return import_str not in content


def add_imports(content: str) -> str:
    """Add missing imports for reflect, strings, defaults, stringdefault."""
    imports_to_add = []

    # if needs_import(content, '"reflect"'):
    #     imports_to_add.append('"reflect"')
    # if needs_import(content, '"strings"'):
    #     imports_to_add.append('"strings"')

    # These are framework imports
    if needs_import(content, '"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"'):
        imports_to_add.append('"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"')
    if needs_import(content, '"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"'):
        imports_to_add.append('"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"')

    # utils import
    if needs_import(content, '"github.com/infobloxopen/terraform-provider-nios/internal/utils"'):
        imports_to_add.append('"github.com/infobloxopen/terraform-provider-nios/internal/utils"')

    if not imports_to_add:
        return content

    # Find the import block and add missing imports
    # Look for the closing ) of the import block
    import_end = content.find('\n)\n')
    if import_end == -1:
        return content

    insert_lines = "\n".join(f"\t{imp}" for imp in imports_to_add)
    content = content[:import_end] + "\n" + insert_lines + content[import_end:]

    return content


def generate_put_expand(info: dict) -> str:
    """Generate the PutExpand function."""
    return f'''
func (m *{info["model_name"]}) PutExpand(to *{info["pkg_alias"]}.{info["api_type"]}) *{info["pkg_alias"]}.{info["api_type"]} {{
	if m == nil {{
		return nil
	}}
	toType := reflect.TypeOf(to)
	if toType.Kind() == reflect.Ptr {{
		toType = toType.Elem()
	}}
	toVal := reflect.ValueOf(to).Elem()
	for field, attr := range {info["schema_var"]} {{
		attrVal := reflect.ValueOf(attr)
		attrType := attrVal.Type()
		if toType.Kind() == reflect.Struct {{
			for i := 0; i < toType.NumField(); i++ {{
				fieldValue := toVal.Field(i).Interface()
				tField := toType.Field(i)
				cleanTag := strings.Split(tField.Tag.Get("json"), ",")[0]
				cleanTag = strings.Trim(cleanTag, "_")
				txtFieldValue := utils.ToString(field, fieldValue)
				if field == cleanTag {{
					_, ok := attrType.FieldByName("Default")
					if ok {{
						defaultVal := attrVal.FieldByName("Default")
						if defaultVal.IsValid() && defaultVal.CanInterface() {{
							strDef, ok := defaultVal.Interface().(defaults.String)
							if ok {{
								if strDef == stringdefault.StaticString("") {{
									continue
								}} else if txtFieldValue == "" {{
									utils.DeleteBy(to, tField.Name)
								}}
							}}
							if !ok && txtFieldValue == "" {{
								utils.DeleteBy(to, tField.Name)
							}}
						}}
					}} else if txtFieldValue == "" {{
						utils.DeleteBy(to, tField.Name)
					}}
				}}
			}}
		}}
	}}
	return to
}}
'''


def process_file(filepath: Path, dry_run: bool = False) -> bool:
    """Process a single model file. Returns True if modified."""
    content = filepath.read_text()

    if has_put_expand(content):
        return False

    info = find_model_info(content)
    if info is None:
        return False

    # Add imports
    new_content = add_imports(content)

    # Add PutExpand function at the end of the file
    put_expand_func = generate_put_expand(info)
    new_content = new_content.rstrip() + "\n" + put_expand_func

    if dry_run:
        print(f"  [DRY-RUN] Would add PutExpand to: {filepath.relative_to(BASE_DIR.parent.parent)}")
        print(f"            Model: {info['model_name']}, API: {info['pkg_alias']}.{info['api_type']}, Schema: {info['schema_var']}")
        return True

    filepath.write_text(new_content)
    run_gofmt(filepath)
    print(f"  ✓ Added PutExpand to: {filepath.relative_to(BASE_DIR.parent.parent)}")
    print(f"    Model: {info['model_name']}, API: {info['pkg_alias']}.{info['api_type']}")
    return True


def process_resource_file(filepath: Path, dry_run: bool = False) -> bool:
    """
    Process a _resource.go file to wrap the Expand call in Update with PutExpand.
    
    Transforms:
        SomeType(*data.Expand(ctx, &resp.Diagnostics, false)).
    into:
        SomeType(*data.PutExpand(data.Expand(ctx, &resp.Diagnostics, false))).
    
    Only modifies the call inside the Update function (after an .Update( line).
    """
    content = filepath.read_text()

    if "PutExpand" in content:
        return False

    # We need to find the Expand call that is inside the Update function.
    # The pattern is: after a line containing ".Update(ctx," there will be a line
    # with SomeType(*data.Expand(ctx, &resp.Diagnostics, false)).
    # or SomeType(*data.Expand(ctx, &resp.Diagnostics)).

    lines = content.split('\n')
    modified = False
    in_update_block = False

    for i, line in enumerate(lines):
        # Detect we're in the Update API call block
        if re.search(r'Update\(ctx[,)]', line):
            in_update_block = True
            continue

        if in_update_block:
            # Look for the Expand pattern within the next few lines
            # Pattern: SomeType(*data.Expand(ctx, &resp.Diagnostics, false)).
            # or:      SomeType(*data.Expand(ctx, &resp.Diagnostics)).
            if '*data.Expand(ctx, &resp.Diagnostics' in line and 'PutExpand' not in line:
                new_line = line.replace(
                    '*data.Expand(ctx, &resp.Diagnostics, false)',
                    '*data.PutExpand(data.Expand(ctx, &resp.Diagnostics, false))'
                )
                if new_line == line:
                    # Try without the isCreate param
                    new_line = line.replace(
                        '*data.Expand(ctx, &resp.Diagnostics)',
                        '*data.PutExpand(data.Expand(ctx, &resp.Diagnostics))'
                    )
                if new_line != line:
                    lines[i] = new_line
                    modified = True
                    in_update_block = False
                    continue

            # If we hit Execute or another API method, stop looking
            if 'Execute()' in line or 'ReturnFields' in line:
                in_update_block = False

    if not modified:
        return False

    if dry_run:
        print(f"  [DRY-RUN] Would wrap Expand with PutExpand in: {filepath.relative_to(BASE_DIR.parent.parent)}")
        return True

    filepath.write_text('\n'.join(lines))
    run_gofmt(filepath)
    print(f"  ✓ Wrapped Expand with PutExpand in: {filepath.relative_to(BASE_DIR.parent.parent)}")
    return True


def main():
    global BASE_DIR

    parser = argparse.ArgumentParser(description="Add PutExpand to model and resource files.")
    parser.add_argument("--repo-path", "--repo_path", default=".", help="Path to the terraform-provider-nios repository root (default: current directory)")
    parser.add_argument("--dry-run", action="store_true", help="Preview changes without modifying files")
    args = parser.parse_args()

    repo_path = Path(args.repo_path).resolve()
    BASE_DIR = repo_path / "internal" / "service"
    dry_run = args.dry_run

    if not BASE_DIR.exists():
        print(f"Error: {BASE_DIR} does not exist. Check the repo path.")
        sys.exit(1)

    print(f"Repo: {repo_path}")
    if dry_run:
        print("=== DRY RUN MODE ===\n")

    # --- Step 1: Add PutExpand to model files ---
    print(f"=== Step 1: Adding PutExpand to model files ===\n")

    model_modified = 0
    model_skipped = 0
    model_errors = 0

    for filepath in sorted(BASE_DIR.rglob("model_*.go")):
        name = filepath.name
        if any(x in name for x in ["_test", "_old", "_new", "_edited"]):
            continue

        try:
            if process_file(filepath, dry_run):
                model_modified += 1
            else:
                model_skipped += 1
        except Exception as e:
            print(f"  ✗ Error processing {filepath}: {e}")
            model_errors += 1

    print(f"\n  Model files modified: {model_modified}")
    print(f"  Model files skipped: {model_skipped}")
    print(f"  Model file errors: {model_errors}")

    # --- Step 2: Wrap Expand with PutExpand in resource Update functions ---
    print(f"\n=== Step 2: Wrapping Expand with PutExpand in resource Update functions ===\n")

    res_modified = 0
    res_skipped = 0
    res_errors = 0

    for filepath in sorted(BASE_DIR.rglob("*_resource.go")):
        name = filepath.name
        if any(x in name for x in ["_test", "_old", "_new", "_edited"]):
            continue

        try:
            if process_resource_file(filepath, dry_run):
                res_modified += 1
            else:
                res_skipped += 1
        except Exception as e:
            print(f"  ✗ Error processing {filepath}: {e}")
            res_errors += 1

    print(f"\n  Resource files modified: {res_modified}")
    print(f"  Resource files skipped: {res_skipped}")
    print(f"  Resource file errors: {res_errors}")

    # --- Summary ---
    total_modified = model_modified + res_modified
    total_errors = model_errors + res_errors
    print(f"\n=== Total Summary ===")
    print(f"  Total modified: {total_modified}")
    print(f"  Total errors: {total_errors}")

    if not dry_run and total_modified > 0:
        print(f"\nAll modified files have been auto-formatted with gofmt/goimports.")
        print(f"Run 'go build ./...' to verify the changes compile.")


if __name__ == "__main__":
    main()
