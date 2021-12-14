script   := "./.run-script.sh"
get-id   := $(filter-out $@,$(MAKECMDGOALS))
entity   := broker leia fulcrum informant

default: noinput

noinput:
	@$(script) fulcrum 0
	@$(script) fulcrum 1
	@$(script) fulcrum 2
	@$(script) broker

dist-181:
	@$(script) fulcrum 0
	@$(script) informant 0

dist-182:
	@$(script) fulcrum 1
	@$(script) informant 1

dist-183:
	@$(script) fulcrum 2
	@$(script) leia

dist-184:
	@$(script) broker

$(entity):
	@$(script) $@ $(call get-id)

stop:
	@-pkill -x go ||:
	@-pkill -x main ||:

help:
	@echo "Running entities:"
	@echo "    make <entity> [id]" # run entity with id 0 by default
	@echo "    make dits-[#vm]     # run specific entities for a given VM"
	@echo " *  make noinput        # run all entities that don't read from stdin"
	@echo ""
	@echo "Other:"
	@echo "    make clean          # remove .log and .txt files"
	@echo "    make stop           # kill all \"go\" and \"main\" processes"
	@echo "    make reset          # both stop and clean"
	@echo ""
	@echo "[ * default ]"

reset: stop clean

clean:
	@rm -fv .logs/*.log
	@rm -fv *.txt

%:
	@:

.PHONY: $(entity) dist-0 dist-1 dist-2 dist-3 help stop noinput