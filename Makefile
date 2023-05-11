
ip2asn-combined.tsv:
	curl -LO https://iptoasn.com/data/ip2asn-combined.tsv.gz
	gunzip $@.gz
