#!/bin/bash
# Prep Data Menu Script

# trap ctrl-c and call ctrl_c()
trap ctrl_c INT
function ctrl_c() {
	tput rmcup
	exit 0
}

function print_menu(){
	clear
	echo "  ___                ___       _
 | _ \_ _ ___ _ __  |   \ __ _| |_ __ _
 |  _/ '_/ -_) '_ \ | |) / _\` |  _/ _\` |
 |_| |_| \___| .__/ |___/\__,_|\__\__,_|
             |_|
"
	echo "Folder Creation Toggles:
1) 1 minute
2) 5 minute
3) hourly
4) daily

File Creation Window:
5) Current and Past
6) Current, Past and Future
7) Current and Future

8) Number of files per folder

9) Execute
0) Quit"
	# write options status
	list=""
	for a in "${MODE[@]}"
	do
		list="$list $a"
	done
	printf '\n Prefix folder:%s\n Selected folders:%s\n        Time Mode: %s\n  Number of files: %d\n' "$PREFIX" "$list" "$TIME" "$FILES"
	# check if any aditional message to write
	if [ ! -z "$ERROR" ]
	then
		echo -e "\n"$ERROR
	fi
	# printf "\n> "
}

function toggle_action(){
	if printf '%s\0' "${MODE[@]}" | grep -Fxqz -- "$1"
	then
		# exists -> remove from array
		aux=()
		for item in ${MODE[@]}
		do
			if [ "$item" != "$1" ]; then aux+=("$item"); fi
		done
		MODE=("${aux[@]}")
	else
		# don't exists -> add to array
		MODE+=("$1")
	fi
}

function write_files() {
	start=${1:-0}
	end=${2:-5}
	step=${3:-minutes}
	prefix=${4:-file_mod_}
	extension=${5:-txt}
	increment=${6:-1}

	for i in $(seq $start $increment $end)
	do
		# stamp=$(date -ud "${i} ${step}" '+%Y%m%d%H%M')
		stamp=$(date -d "${i} ${step}" '+%Y%m%d%H%M')
		echo "Writen ${prefix}${stamp}.${extension}"
		# write file with data
		head -c 10K /dev/urandom > "${prefix}${stamp}.${extension}"
		# change modification time
		touch -t $stamp "${prefix}${stamp}.${extension}"
	done
}

function execute_actions() {
	# time mode
	if [ "$TIME" == "past" ]; then
		start=-$(( $FILES-1 ))
		end=0
	elif [ "$TIME" = "future" ]; then
		start=0
		end=$(( $FILES-1 ))
	elif [ "$TIME" = "all" ]; then
		start=-$(( $FILES/2 ))
		end=$(( $FILES/2 - (1-$FILES%2) ))
	else
		ERROR="${RED}ERROR${NC}: Unknow Time Mode ${TIME}"
		return
	fi

	# execute actions
	echo ""
	for act in ${MODE[@]}; do
		if [ $act = "1minute" ]; then
			mkdir -p $PREFIX/files_1min
			write_files $start $end minutes "$PREFIX/files_1min/min1_" log
		elif [ $act = "5minute" ]; then
			mkdir -p $PREFIX/files_5min
			write_files $(( $start*5 )) $(( $end*5 )) minutes "$PREFIX/files_5min/min5_" log 5
		elif [ $act = "hourly" ]; then
			mkdir -p $PREFIX/files_hour
			write_files $start $end hours "$PREFIX/files_hour/hourly_" log
		elif [ $act = "daily" ]; then
			mkdir -p $PREFIX/files_day
			write_files $start $end days "$PREFIX/files_day/daily_" log
		else
			ERROR="${RED}ERROR${NC}: Unknow mode"
		fi
	done
}

# terminal colors
RED='\033[0;31m'
NC='\033[0m'
# Global Variables
MODE=("1minute")
TIME="past"
FILES=5
ERROR=""
PREFIX="${1:-ftpdata}"

# create perfix folder
mkdir -p $PREFIX

# save terminal state
tput smcup

print_menu
while :
do
	read -p "> " opt
	case $opt in
		1)
			toggle_action "1minute";;
		2)
			toggle_action "5minute";;
		3)
			toggle_action "hourly";;
		4)
			toggle_action "daily";;
		5)
			TIME="past";;
		6)
			TIME="all";;
		7)
			TIME="future";;
		8)
			read -p "Select number of files: " num
			if [[ $num =~ ^[0-9]+$  && $num > 0 ]]; then
				FILES=$num
			else
				ERROR="${RED}ERROR${NC}: Number of files must be a integer, greather than zero"
			fi
			;;
		9)
			execute_actions
			echo -e "\nPress any key to exit"
			read -n1
			break 2
			;;
		0)
			break 2;;
		*)
			ERROR="${RED}ERROR${NC}: Invalid option ${opt}";;
	esac
	print_menu
	# clean error variable
	ERROR=""
done
# restore terminal
tput rmcup
