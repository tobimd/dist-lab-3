script   := "./.run-script.sh"
get-id   := $(filter-out $@,$(MAKECMDGOALS))
entity   := broker leia fulcrum informant

default: noinput

noinput:
	@$(script) fulcrum 0
	@$(script) fulcrum 1
	@$(script) fulcrum 2
	@$(script) broker

dist-0:
	@$(script) fulcrum 0
	@$(script) informant 0

dist-1:
	@$(script) fulcrum 1
	@$(script) informant 1

dist-2:
	@$(script) fulcrum 2
	@$(script) leia

dist-3:
	@$(script) broker

$(entity):
	@$(script) $@ $(call get-id)

stop:
	@-pkill -ex go ||:
	@-pkill -ex main ||:

help:
	@echo "Running entities:"
	@echo "    make <entity> [id]"
	@echo "    make dist-[0..3]    # run specific entities"
	@echo " *  make noinput        # run all entities that don't read from stdin"
	@echo ""
	@echo "Other:"
	@echo "    make clean          # remove .log and .txt files"
	@echo "    make stop           # kill all \"go\" and \"main\" processes"
	@echo ""
	@echo "[ * default ]"


clean:
	@rm -fv .logs/*.log
	@rm -fv *.txt

%:
	@:

.PHONY: $(entity) dist-0 dist-1 dist-2 dist-3 help stop noinput