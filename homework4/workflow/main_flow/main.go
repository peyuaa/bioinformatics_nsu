package main

import (
	sp "github.com/scipipe/scipipe"
)

func main() {
	// Initialize the workflow with a concurrency level of 4
	wf := sp.NewWorkflow("bioinformatics_pipeline", 4)

	// Step 1: Download SRR31294328.fastq.gz
	downloadFastq := wf.NewProc("download_fastq", "curl -o {o:fastq_gz} https://trace.ncbi.nlm.nih.gov/Traces/sra-reads-be/fastq?acc=SRR31294328")
	downloadFastq.SetOut("fastq_gz", "SRR31294328.fastq.gz")

	// Step 2: Decompress SRR31294328.fastq.gz
	decompressFastq := wf.NewProc("decompress_fastq", "gzip -d -c {i:fastq_gz} > {o:fastq}")
	decompressFastq.In("fastq_gz").From(downloadFastq.Out("fastq_gz"))
	decompressFastq.SetOut("fastq", "SRR31294328.fastq")

	// Step 3: Download reference genome
	downloadGenome := wf.NewProc("download_genome", "curl -o {o:fna_gz} https://ftp.ncbi.nlm.nih.gov/genomes/all/GCF/000/005/845/GCF_000005845.2_ASM584v2/GCF_000005845.2_ASM584v2_genomic.fna.gz")
	downloadGenome.SetOut("fna_gz", "GCF_000005845.2_ASM584v2_genomic.fna.gz")

	// Step 4: Decompress reference genome
	decompressGenome := wf.NewProc("decompress_genome", "gzip -d -c {i:fna_gz} > {o:fna}")
	decompressGenome.In("fna_gz").From(downloadGenome.Out("fna_gz"))
	decompressGenome.SetOut("fna", "GCF_000005845.2_ASM584v2_genomic.fna")

	//// Step 5: Run FastQC and bwa index
	//fastqc_bwa := wf.NewProc("fastqc_bwa", "fastqc {i:fastq} && bwa index {i:fna}")
	//fastqc_bwa.In("fastq").From(decompressFastq.Out("fastq"))
	//fastqc_bwa.In("fna").From(decompressGenome.Out("fna"))

	// Step 7: BWA mem
	bwaMem := wf.NewProc("bwa_mem", "bwa mem {i:fna} {i:fastq} | gzip -3 > {o:sam_gz}")
	bwaMem.In("fna").From(decompressGenome.Out("fna"))
	bwaMem.In("fastq").From(decompressFastq.Out("fastq"))
	bwaMem.SetOut("sam_gz", "aln-se.sam.gz")

	// Step 8: Decompress SAM file
	decompressSam := wf.NewProc("decompress_sam", "gzip -d -c {i:sam_gz} > {o:sam}")
	decompressSam.In("sam_gz").From(bwaMem.Out("sam_gz"))
	decompressSam.SetOut("sam", "aln-se.sam")

	// Step 9: Convert SAM to BAM
	samToBam := wf.NewProc("sam_to_bam", "samtools view -bS {i:sam} > {o:bam}")
	samToBam.In("sam").From(decompressSam.Out("sam"))
	samToBam.SetOut("bam", "aln-se.bam")

	// Step 10: Samtools flagstat
	flagstat := wf.NewProc("flagstat", "samtools flagstat {i:bam} > {o:stats}")
	flagstat.In("bam").From(samToBam.Out("bam"))
	flagstat.SetOut("stats", "samtools_result.txt")

	// Step 11: Parse and evaluate results
	parseResult := wf.NewProc("parse_result", `grep "[0-9] mapped (" {i:stats} | sed 's/^.*mapped/mapped/' | tr -d -c "0-9." > {o:parse_out}`)
	parseResult.In("stats").From(flagstat.Out("stats"))
	parseResult.SetOut("parse_out", "parse_result.txt")

	sortBam := wf.NewProc("final", `
if (( $(echo "$(cat {i:parse_out}) < 90" | bc -l) )); then
    echo "Not OK"
else
    samtools sort {i:bam} > {o:sorted_bam}
fi`)
	sortBam.In("parse_out").From(parseResult.Out("parse_out"))
	sortBam.In("bam").From(samToBam.Out("bam"))
	sortBam.SetOut("sorted_bam", "alignment.sorted.bam")

	freebayes := wf.NewProc("freebayes", `
if (( $(echo "$(cat {i:parse_out}) >= 90" | bc -l) )); then
    freebayes -f {i:fna} -b {i:sorted_bam} > {o:vcf}
fi`)
	freebayes.In("parse_out").From(parseResult.Out("parse_out"))
	freebayes.In("fna").From(decompressGenome.Out("fna"))
	freebayes.In("sorted_bam").From(sortBam.Out("sorted_bam"))
	freebayes.SetOut("vcf", "result.vcf")

	// Run the workflow
	wf.Run()
}
