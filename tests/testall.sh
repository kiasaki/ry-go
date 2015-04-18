#!/bin/sh

red='\e[0;31m' # Red
grn='\e[0;32m' # Green
rst='\e[0m'    # Text Reset

exitstatus=0
failed=0
passed=0

for file in tests/*.syp
do
  ./ry-repl -exitonfail "${file}"
  status=$?

  if [ $status -ne 0 ]; then
    exitstatus=1
    failed=$((failed + 1))
    printf "$red${file} failed$rst\n"
  else
    passed=$((passed + 1))
    printf "$grn${file} passed$rst\n"
  fi
done

total=$((failed + passed))

if [ $exitstatus -ne 0 ]; then
  printf "\nRan $total test files, $failed failed\n"
  printf "$red\nFAILED$rst\n"
else
  printf "\nRan $total test files\n"
  printf "$grn\nOK$rst\n"
fi

exit $exitstatus
