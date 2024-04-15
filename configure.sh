#!/bin/bash

#
# This script check see the needed programs are installed
#
#!/bin/bash

#
# Check software needed and their versions
#
function check_version() {
    OPT="$1"
	VERSION="$2"
	CMD="$3"
	if command -v "${CMD}" >/dev/null; then
		HAS_VERSION=$("${CMD}" "${OPT}" | grep "${VERSION}")
		if [ "$HAS_VERSION" = "" ]; then
			echo "${CMD} version check: expected ${VERSION}, got $("${CMD}" "${OPT}")"
		fi
	else
		echo "${CMD} is missing, aborting"
		exit 1
	fi
}

#
# Check all of irdmtools is installed
#
echo "Checking for Go 1.22, Pandoc 3.1, PageFind 1.1.0, Postgres 16 and PostgREST 12"
## - Pandoc >= 3.1
check_version "--version" "3.1" "pandoc"
## - PageFind >= 1.1.0
check_version "--version" "1.1.0" "pagefind"
## - Check for Go version 1.22
check_version "version" "1.22" "go"
## - Check for Postgres 16
check_version "--version" "16" "psql"
## - Check for PostgREST 121
check_version "--version" "12" "postgrest"
echo "Software installed. Review version info if any is displayed."

