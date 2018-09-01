# Code for page-per-session forecasting case study

The code in this repository accompanies a case
study on page-per-session forecasting.

## Code layout

* model — deterministic model.
* model/query — probabilistic query on top of the probabilistic model.
* infer — Monte Carlo inference on the probabilistic query.
* cmd/scan — a command-line utility which takes a PPS log and computes
	posterior beliefs.
* cmd/plot — a command-line utility which transforms beliefs into
	graphs, each graph is frame for an animated GIF file.
* cmd/tune — a command-line utility for inferring the bandwidth.

## Usage

Go 1.11 is required because this package uses [Go Modules](https://github.com/golang/go/wiki/Modules).

You will need `convert` from `imagemagick` in the path.
Run `make` on a Unix-like system to generate animated GIF.

Build and run `tune` on a file from `data/` to infer the bandwidth 
for each campaign.

# License

See the LICENSE file.
