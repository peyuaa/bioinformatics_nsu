#!/bin/sh

curl -o SRR31294328.fastq.gz "https://trace.ncbi.nlm.nih.gov/Traces/sra-reads-be/fastq?acc=SRR31294328"

gzip -d SRR31294328.fastq.gz

curl -o GCF_000005845.2_ASM584v2_genomic.fna.gz https://ftp.ncbi.nlm.nih.gov/genomes/all/GCF/000/005/845/GCF_000005845.2_ASM584v2/GCF_000005845.2_ASM584v2_genomic.fna.gz

gzip -d GCF_000005845.2_ASM584v2_genomic.fna.gz

git clone https://github.com/lh3/bwa.git
cd bwa || exit
make

mkdir fastqc_out

fastqc -o fastqc_out SRR31294328.fastq

mkdir bwa_out

cd bwa_out

mv ../reference.fna .

bwa index reference.fna

bwa mem reference.fna ../SRR31294328.fastq  | gzip -3 > aln-se.sam.gz

gunzip aln-se.sam.gz

samtools view -bS aln-se.sam > aln-se.bam

samtools flagstat aln-se.bam > samtools_result.txt

chmod +x parse.sh

./parse.sh samtools_result.txt > parse_result.txt

if (( $(echo "$(cat parse_result.txt) < 90" | bc -l) )); then
    echo "Not OK"
else
    samtools sort aln-se.bam > alignment.sorted.bam
    freebayes -f reference.fna -b alignment.sorted.bam > result.vcf
    echo "OK"
fi