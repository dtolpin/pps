all: plots/ASTBT.gif plots/JLNVT2.gif plots/SELVPNS.gif

TOTAL=30
BANDWIDTH=1000
THIN=20
HEAD=101
DELAY=25
LOOP=1

pps-%.csv: data/%.csv scan
	./scan -total $(TOTAL) -bandwidth $(BANDWIDTH) -thin $(THIN) < $< | head -$(HEAD) > $@

plots/%.gif: pps-%.csv plot
	./plot < $<
	convert -delay $(DELAY) -loop $(LOOP) pps-*.gif $@ ; rm -f pps-*.gif


scan: model/model.go model/model_test.go cmd/scan/main.go cmd/scan/main_test.go
	go test ./model ./cmd/scan
	go build ./cmd/scan

plot: model/model.go model/model_test.go cmd/plot/main.go cmd/plot/main_test.go
	go test ./model ./cmd/plot
	go build ./cmd/plot

clean:
	rm -f pps-*.csv

distclean: clean
	rm -f plot scan plots/*.gif

