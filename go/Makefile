# export GOPATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# SRCS = step0_repl.go step1_read_print.go step2_eval.go step3_env.go \
#        step4_if_fn_do.go step5_tco.go step6_file.go step7_quote.go \
#        step8_macros.go step9_try.go stepA_mal.go
SRCS = step0_repl.go step1_read_print.go step2_eval.go step3_env.go \
	   step4_if_fn_do.go
BINS = $(SRCS:%.go=%)

#####################

all: $(BINS)

dist: mal

mal: $(word $(words $(BINS)),$(BINS))
	cp $< $@

define dep_template
$(1): src/$(1)/$(1).go
	cd src; go build -o ../mal/$(1) $(1)/$(1).go
endef

$(foreach b,$(BINS),$(eval $(call dep_template,$(b))))

clean:
	rm -f $(BINS) mal
