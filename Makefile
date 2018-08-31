all: plots/ASTBT.gif plots/JLNVT2.gif plots/SELVPNS.gif

TOTAL=30
BANDWIDTH=1000
THIN=20
HEAD=101
DELAY=25
LOOP=1

GOS=\
	./cmd/plot/main.go \
	./cmd/plot/main_test.go \
	./cmd/scan/main.go \
	./cmd/scan/main_test.go \
	./cmd/tune/main.go \
	./cmd/tune/main_test.go \
	./model/model_test.go \
	./model/beta_test.go \
	./model/beta.go \
	./model/model.go \
	./model/query/query.go \
	./model/query/query_test.go \
	./model/query/model_test.go \
	./model/query/model.go \
	./infer/infer.go \
	./infer/infer_test.go

pps-%.csv: data/%.csv scan
	./scan -total $(TOTAL) -bandwidth `./tune -bandwidth $(BANDWIDTH) -thin $(THIN) < $< | head -$(HEAD) > $@

plots/%.gif: pps-%.csv plot
	./plot < $<
	convert -delay $(DELAY) -loop $(LOOP) pps-*.gif $@ ; rm -f pps-*.gif

scan: $(GOS)
	go test ./...
	go build ./cmd/scan

plot: $(GOS)
	go test ./...
	go build ./cmd/plot

clean:
	rm -f pps-*.csv

distclean: clean
	rm -f plot scan plots/*.gif

