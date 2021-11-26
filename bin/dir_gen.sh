#!/bin/bash
# cmd: /bin/bash bin/dir_gen.sh

# define
file_name_list=("log")

# operate
if [ ! -d "runtime/" ];then
	mkdir runtime
	echo "generated runtime!"
fi
chmod -R 777 runtime

cd runtime
for item in ${file_name_list[*]}
do 
	if [ ! -d "${item}/" ];then
		mkdir ${item}
		echo "generated ${item}!"
	fi

	chmod -R 777 ${item}
done
