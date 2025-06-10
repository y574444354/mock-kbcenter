#!/bin/bash

# Check for untranslated i18n strings in Go files
# Usage: ./scripts/check_untranslated_i18n.sh

set -euo pipefail

# Find all Go files excluding vendor directory
GO_FILES=$(find . -name "*.go" -not -path "./vendor/*")

# Extract all i18n.Translate("...") strings
# Extract all i18n.Translate calls with file paths and line numbers
TRANSLATE_CALLS=$(grep -r -n 'i18n\.Translate(".*",' $GO_FILES)
# Extract just the translation keys
TRANSLATE_STRINGS=$(echo "$TRANSLATE_CALLS" | sed -E 's/.*i18n\.Translate\("([^"]*)",.*/\1/')

# Check if string exists in translation files
MISSING_TRANSLATIONS=0

for str in $TRANSLATE_STRINGS; do
    # Check zh-CN translation
    # Get the file and line info for this string (first occurrence)
    FILE_LINE=$(echo "$TRANSLATE_CALLS" | grep -m1 "i18n\.Translate(\"$str\"," | cut -d: -f1-2)
    
    # Check zh-CN translation
    if ! grep -q "$str:" i18n/locales/zh-CN.yaml; then
        echo "[zh-CN] Missing translation for: \"$str\" (at $FILE_LINE)"
        MISSING_TRANSLATIONS=$((MISSING_TRANSLATIONS + 1))
    fi
    
    # Check en translation
    if ! grep -q "$str:" i18n/locales/en.yaml; then
        echo "[en] Missing translation for: \"$str\" (at $FILE_LINE)"
        MISSING_TRANSLATIONS=$((MISSING_TRANSLATIONS + 1))
    fi
done

if [ $MISSING_TRANSLATIONS -gt 0 ]; then
    echo "Found $MISSING_TRANSLATIONS missing translations"
    exit 1
else
    echo "All translations are present"
    exit 0
fi