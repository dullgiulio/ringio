#!/bin/sh

export RINGIO_HOME=.
export RINGIO=../ringio
export RINGIO_TEST_FAILURE=0

ringio_test_error() {
	echo -e "err\t${@:1}" >&2
	RINGIO_TEST_FAILURE=1
}

nohup $RINGIO test open >/dev/null 2>&1 &

# TODO: Implement --quiet option.
RINGIO_DEBUG_VAR=debugval $RINGIO test input print-env.sh >/dev/null
if [ "$?" != '0' ]; then
	ringio_test_error "Adding agent failed"
fi

if [ "$($RINGIO test output --no-wait | grep debugval)" == '' ]; then
	ringio_test_error "Environment variable not passed to subprocess"
fi

$RINGIO test close
if [ "$?" != '0' ]; then
	ringio_test_error "Ending session failed"
fi

if [ "$(ls -1 .ringio/ | wc -l)" != '0' ]; then
	ringio_test_error "Expected ringio home to be empty"
fi

rm -rf .ringio

if [ "$RINGIO_TEST_FAILURE" == '0' ]; then
	echo "OK"
fi

exit $RINGIO_TEST_FAILURE

