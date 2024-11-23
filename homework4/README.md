# Объяснение скрипта

1. E. Coli SRR31294328

Скачиваем и распаковываем
```bash
curl -o SRR31294328.fastq.gz "https://trace.ncbi.nlm.nih.gov/Traces/sra-reads-be/fastq?acc=SRR31294328"

gzip -d SRR31294328.fastq.gz
```

2. Референсный геном

Скачиваем и распаковываем
```bash
curl -o GCF_000005845.2_ASM584v2_genomic.fna.gz https://ftp.ncbi.nlm.nih.gov/genomes/all/GCF/000/005/845/GCF_000005845.2_ASM584v2/GCF_000005845.2_ASM584v2_genomic.fna.gz

gzip -d GCF_000005845.2_ASM584v2_genomic.fna.gz

mv GCF_000005845.2_ASM584v2_genomic.fna reference.fna
```

3. Устанавливаем lh3/bwa, fastqc, samtools.

lh3/bwa
```bash
git clone https://github.com/lh3/bwa.git
cd bwa || exit
make
```

fastqc
```bash
https://www.bioinformatics.babraham.ac.uk/projects/fastqc/
```

samtools
```bash
https://www.htslib.org/download/
```

4. Запускаем анализ с помощью fastqc

```bash
mkdir fastqc_out

fastqc -o fastqc_out SRR31294328.fastq
```

5. Индексируем референсный геном с помощью bwa

```bash
mkdir bwa_out

cd bwa_out

mv ../reference.fna .

bwa index reference.fna
```

6. Теперь выравниваем риды из файлов fastqc на референсный геном

```bash
bwa mem reference.fna ../SRR31294328.fastq  | gzip -3 > aln-se.sam.gz
```

7. Распакуем

```bash
gunzip aln-se.sam.gz
```

8. Преобразование файла формата SAM в формат BAM
```bash
samtools view -bS aln-se.sam > aln-se.bam 
```

9. Анализируем и сохраняем в samtools_result.txt
```bash
samtools flagstat aln-se.bam > samtools_result.txt
```

10. Сохраняем скрипт для получения процента в parse.sh

```bash
grep "[0-9] mapped (" $1 | sed 's/^.*mapped/mapped/' | tr -d -c "0-9."
```

11. Даем права на исполнение
```bash
chmod +x parse.sh
```

12. Выполняем скрипт
```bash
./parse.sh samtools_result.txt > parse_result.txt
```

13. Проверяем результат
```bash
if (( $(echo "$(cat parse_result.txt) < 90" | bc -l) )); then
    echo "Not OK"
else
    samtools sort aln-se.bam > alignment.sorted.bam
    freebayes -f reference.fna -b alignment.sorted.bam > result.vcf
    echo "OK"
fi
```

# Устанавливаем фреймворк scipipe
#### https://scipipe.org/

1. Скачиваем и устанавливаем Go
```
https://go.dev/doc/install
```

2. Устанавливаем scipipe с помощью go
```bash
go install github.com/scipipe/scipipe/...@latest
```

3. Создаем папку, в котором будет код
```bash
mkdir workflow
cd workflow
```

4. Инициализируем go-модуль
```bash
go mod init workflow
```