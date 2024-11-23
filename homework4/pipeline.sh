#!/bin/sh

curl -o SRR31294328.fastq.gz "https://trace.ncbi.nlm.nih.gov/Traces/sra-reads-be/fastq?acc=SRR31294328"

gzip -d SRR31294328.fastq.gz

curl -o GCF_000005845.2_ASM584v2_genomic.fna.gz https://ftp.ncbi.nlm.nih.gov/genomes/all/GCF/000/005/845/GCF_000005845.2_ASM584v2/GCF_000005845.2_ASM584v2_genomic.fna.gz

gzip -d GCF_000005845.2_ASM584v2_genomic.fna.gz

git clone https://github.com/lh3/bwa.git
cd bwa || exit
make