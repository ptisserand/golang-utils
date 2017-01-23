SUBDIRS = finddup elfenstein pushme
SUBCLEAN = $(addsuffix .clean, $(SUBDIRS))

.PHONY: clean $(SUBCLEAN) subdirs $(SUBDIRS)

subdirs: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

clean: $(SUBCLEAN)

$(SUBCLEAN):
	$(MAKE) -C $(@:.clean=) clean

